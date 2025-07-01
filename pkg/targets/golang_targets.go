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
		"github.com/redis/go-redis/v9",
	}
)
