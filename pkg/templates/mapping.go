package templates


import (
	"text/template"
)

var mapping = map[string]*template.Template{
	"main.go": template.Must(template.New("main.go").Parse(MainGoTemplate)),
	"openapitools.json": template.Must(template.New("openapitools.yaml").Parse(OpenapitoolsJSONTemplate)),
	"sqlc.yaml": template.Must(template.New("sqlc.yaml").Parse(SQLCYamlTemplate)),
	"README.md": template.Must(template.New("README.md").Parse(READMETemplate)),
	"Makefile": template.Must(template.New("Makefile").Parse(MakefileTemplate)),
	"Dockerfile": template.Must(template.New("Dockerfile").Parse(DockerfileTemplate)),
	".gitignore": template.Must(template.New(".gitignore").Parse(GitignoreTemplate)),
	".dockerignore": template.Must(template.New(".dockerignore").Parse(DockerignoreTemplate)),
	".air.toml": template.Must(template.New(".air.toml").Parse(AirTomlTemplate)),

	
}
