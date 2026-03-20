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

package targets


var (
	// the packagesPackages that are required for the fullstack_app (to the best of my knowledge)
	// this will be tested later to make sure that everything still works
	GolangPackages = []string{
		"github.com/labstack/echo",
		"github.com/spf13/cobra",
		"gopkg.in/yaml.v3",
		"github.com/go-openapi/jsonpointer",
		"github.com/go-openapi/jsonreference",
		"github.com/go-openapi/spec",
		"github.com/go-openapi/swag",
		"github.com/golang-jwt/jwt/v5",
		"github.com/labstack/echo-jwt/v4",
		"github.com/golang-jwt/jwt",
		"github.com/labstack/echo/v4",
		"github.com/joho/godotenv",
		"github.com/minio/crc64nvme",
		"github.com/minio/md5-simd",
		"github.com/minio/minio-go/v7",
		"github.com/pkg/errors",
		"github.com/swaggo/echo-swagger",
		"github.com/swaggo/files/v2",
		"github.com/swaggo/swag",
		"github.com/google/uuid",
		"github.com/jackc/pgx",
		"github.com/jackc/pgx/v5",
		"github.com/jackc/pgx/v5/pgxpool",
		"github.com/redis/go-redis/v9",
	}
)
