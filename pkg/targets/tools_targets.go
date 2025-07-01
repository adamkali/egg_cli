package targets


var (
	// the tools that are required for the fullstack_app
	// air   -> hot reloading
	// swag  -> generating docs and frontend/api endpoint connections 
	// goose -> migrations
	// sqlc  -> generating queries
	RequiredTools = []string{
		"github.com/air-verse/air@latest",
		"github.com/swaggo/swag/cmd/swag@latest",
		"github.com/pressly/goose/v3/cmd/goose@latest",
		"github.com/sqlc-dev/sqlc/cmd/sqlc@latest",
	}
)
