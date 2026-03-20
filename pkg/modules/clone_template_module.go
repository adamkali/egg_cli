/*
Copyright © 2025 Adam Kalinowski <adam.kalilarosa@proton.me>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package modules

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/adamkali/egg_cli/pkg/configuration"
	"github.com/adamkali/egg_cli/pkg/models"
	"github.com/adamkali/egg_cli/styles"
)

const (
	jwtTemplateRepo      = "https://github.com/adamkali/egg-template-jwt"
	auth0TemplateRepo    = "https://github.com/adamkali/egg-template-auth0"
	tursoJwtTemplateRepo = "https://github.com/adamkali/egg-template-turso-jwt"
)

var binaryExtensions = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".ico": true,
	".woff": true, ".woff2": true, ".gif": true, ".webp": true, ".svg": true,
}

// CloneTemplateModule clones the appropriate egg template repo and substitutes
// all __EGG_*__ placeholders with values from the configuration.
type CloneTemplateModule struct {
	config *configuration.Configuration
	eggl   *models.EggLog
	err    error
}

func (m *CloneTemplateModule) Name() string { return "egg::clone_template" }
func (m *CloneTemplateModule) IsError() error { return m.err }

func (m *CloneTemplateModule) LoadFromConfig(config *configuration.Configuration, eggl *models.EggLog) {
	m.config = config
	m.eggl = eggl
}

func (m *CloneTemplateModule) Run() {
	// Pick template repo based on auth provider
	repoURL := jwtTemplateRepo
	switch m.config.Auth.Provider {
	case "auth0":
		repoURL = auth0TemplateRepo
	case "turso-jwt":
		repoURL = tursoJwtTemplateRepo
	}

	fmt.Println(styles.EggProgressInfo.Render(fmt.Sprintf("Cloning template from %s", repoURL)))

	// Clone into a temp directory
	tmpDir, err := os.MkdirTemp("", "egg-template-*")
	if err != nil {
		m.err = fmt.Errorf("failed to create temp dir: %w", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("git", "clone", "--depth=1", repoURL, tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		m.err = fmt.Errorf("git clone failed: %w", err)
		return
	}

	// Build the placeholder map
	placeholders := m.buildPlaceholders()

	// Walk and copy all files (excluding .git/)
	destDir, err := os.Getwd()
	if err != nil {
		m.err = fmt.Errorf("failed to get working directory: %w", err)
		return
	}

	err = filepath.Walk(tmpDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get path relative to tmpDir
		relPath, err := filepath.Rel(tmpDir, srcPath)
		if err != nil {
			return err
		}

		// Skip .git directory
		if relPath == ".git" || strings.HasPrefix(relPath, ".git"+string(os.PathSeparator)) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Apply placeholder substitution to the path itself
		destRelPath := applyPlaceholders(relPath, placeholders)
		destPath := filepath.Join(destDir, destRelPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return m.copyFile(srcPath, destPath, info.Mode(), placeholders)
	})

	if err != nil {
		m.err = fmt.Errorf("failed to copy template files: %w", err)
		return
	}

	fmt.Println(styles.EggProgressInfo.Render("Template files copied, downloading dependencies..."))

	// Use go mod download rather than go mod tidy here — db/repository/ is
	// still empty at this point (sqlc generates it in a later module).
	// go mod download fetches all modules listed in go.mod without needing
	// to resolve every internal package.
	download := exec.Command("go", "mod", "download")
	download.Dir = destDir
	download.Stdout = os.Stdout
	download.Stderr = os.Stderr
	if err := download.Run(); err != nil {
		m.err = fmt.Errorf("go mod download failed: %w", err)
		return
	}

	fmt.Println(styles.EggProgressComplete.Render("Template cloned and dependencies downloaded"))
}

func (m *CloneTemplateModule) buildPlaceholders() map[string]string {
	c := m.config

	// Extract namespace owner (e.g. "adamkali" from "github.com/adamkali/myapp")
	owner := namespaceOwner(c.Namespace)

	p := map[string]string{
		"__EGG_NAMESPACE__":       c.Namespace,
		"__EGG_APP_NAME__":        c.Name,
		"__EGG_VERSION__":         c.Semver,
		"__EGG_LICENSE__":         c.License,
		"__EGG_AUTHOR__":          c.Copyright.Author,
		"__EGG_YEAR__":            strconv.Itoa(c.Copyright.Year),
		"__EGG_PORT__":            strconv.Itoa(c.Server.Port),
		"__EGG_NAMESPACE_OWNER__": owner,
		"__EGG_DB_URL__":          c.Database.URL,
		"__EGG_SQLC_REPO__":       c.Database.Sqlc.RepositoryLocation,
		"__EGG_SQLC_SCHEMA__":     c.Database.Sqlc.Schema,
		"__EGG_SQLC_TYPE__":       c.Database.Sqlc.SqlOrGo,
		"__EGG_DB_QUERIES__":      c.Database.QueriesLocation,
		"__EGG_MIGRATION_PROTO__": c.Database.Migration.Protocol,
		"__EGG_MIGRATION_DEST__":  c.Database.Migration.Destination,
		"__EGG_FRONTEND_DIR__":    c.Server.Frontend.Dir,
		"__EGG_FRONTEND_API__":    c.Server.Frontend.Api,
		"__EGG_CACHE_URL__":       c.Cache.URL,
		"__EGG_S3_URL__":          c.S3.URL,
		"__EGG_S3_ACCESS__":       c.S3.Access,
		"__EGG_S3_SECRET__":       c.S3.Secret,
	}

	// Auth-provider-specific placeholders
	switch c.Auth.Provider {
	case "auth0":
		p["__EGG_AUTH0_DOMAIN__"]        = c.Auth.Auth0Domain
		p["__EGG_AUTH0_AUDIENCE__"]      = c.Auth.Auth0Audience
		p["__EGG_AUTH0_CLIENT_ID__"]     = c.Auth.Auth0ClientID
		p["__EGG_AUTH0_CLIENT_SECRET__"] = c.Auth.Auth0ClientSecret
	case "turso-jwt":
		p["__EGG_JWT_SECRET__"]          = c.Auth.JWT
		p["__EGG_TURSO_AUTH_TOKEN__"]    = c.Database.TursoAuthToken
	default:
		p["__EGG_JWT_SECRET__"] = c.Auth.JWT
	}

	return p
}

func (m *CloneTemplateModule) copyFile(srcPath, destPath string, mode os.FileMode, placeholders map[string]string) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(srcPath))
	if binaryExtensions[ext] {
		// Binary file: copy as-is
		return copyBinary(srcPath, destPath, mode)
	}

	// Text file: read, substitute, write
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	content := applyPlaceholders(string(data), placeholders)

	return os.WriteFile(destPath, []byte(content), mode)
}

func applyPlaceholders(s string, placeholders map[string]string) string {
	for k, v := range placeholders {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

func copyBinary(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// namespaceOwner extracts the owner/username from a Go module namespace.
// e.g. "github.com/adamkali/myapp" → "adamkali"
func namespaceOwner(namespace string) string {
	parts := strings.Split(namespace, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return namespace
}
