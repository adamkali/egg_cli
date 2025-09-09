<img width="200" height="200" alt="egg" src="https://github.com/user-attachments/assets/cb9da171-2b36-4a74-8e0b-629c38bab2bc" />

# egg_cli

A fast and simple command line interface for creating fullstack web applications with the **egg** opinionated framework. 

In truth the **egg** framework is a collection of tools and libraries that can be used to build fullstack web applications quickly and efficiently.

## What is in egg?

Egg is built off of the idea that for many developers are going to need similar tools to get up and started when building a backend in go. 
- **REST API** 
- **Database integration** 
- **Authentication services**
- **Cache integration**
- **S3 integration**
- **Docker setup** for containerized deployment
- **Organized project structure** following Go best practices

With this in mind **egg** takes the approach of builing onto existing tools and libraries.

- **SQLC** for type safe queries 
- **Goose** for database migrations
- **Echo** for REST API actions and services
- **Go/Echo JWT** for authentication
- **Minio for go** for S3 integration
- **Redis-sdk** for cache integration

as a result **egg** does not try to be a dogmatized "YOU ARE GOING TO USE THIS A PARTICULAR WAY" framework. Rather a guiding post for developers to get something started without the boilerplate nuances. 

## You could ask, "why is it opinionated then?"
Simple. Because there are certain aspects that i want to have that I like to have when builing a web api. 

### 1. Services: 
Egg uses services to abstract away the details of how a service works, while allowing test services for tdd can be achieved. As long as You can implement an interface for the service it makes it easy for setting up testing in the future. (coming soon set up by the framework)

### 2. Handlers: 
The Idea for these handlers in the `project/controllers/user_controller.go` also allows for developers to implement only a reciever. this allows you the developer to only match the type for the handle function. This in turn allows for greater flex abilty in ab testing, the Current Instantiated service becomes the implementatino and the Handler expects a result. Allowing for a clear path for a2b testing. ie as long as a service has two or more functions that call the same tyes and the return the same types the handler does not need to understand anything about the current implementation inside of the call, and if anything goes wrong the error is locked and nothing else called. 

### 3. Reliance on echo 
Echo is a fantastic framework for building web services, so using echo for routing and plugin echosystem allows `egg` to piggyback off their plugin ecosystem. 

`Egg` uses the routing features along with the `echo-jwt` and `echo-swagger` plugins to create a REST API with authentication and documentation.

## Installation

### From Go
```bash
go install github.com/adamkali/egg_cli@latest
```

### Build from Source
```bash
git clone https://github.com/adamkali/egg_cli.git
cd egg_cli
make build
```

## Quick Start

### Creating a New Project (`init`)

The `init` command launches an interactive wizard that guides you through setting up your new egg project:

```bash
egg_cli init
```

The wizard will ask you about:
- Project name and description
- Database configuration (PostgreSQL, MySQL, SQLite)
- Authentication setup
- Frontend framework choice (React/Svelte with RSBuild)
- Additional libraries and tools

### Creating a Project with a Configuration (`generate`)

The `generate` command allows you to create a new project with a yaml configuration file:

```bash
egg_cli generate  --config egg.yaml
```

This command will create a new project based on the configuration file provided. By default, it will use the `egg.yaml` configuration file. 

### Creating a Project Configuration File (`generate example`)

The `generate example` command allows you to create a new project with a yaml configuration file:

```bash
egg_cli generate example --config egg.yaml
```

This command will create a new project based on the configuration file provided. By default, it will use the `egg.yaml` configuration file. 

When using --config flag, it will create a project configuration at file path provided. So be sure to also use 

```bash
egg_cli generate --config <file>
```

when calling `generate`

### Generating Environment Files (`generate dotenv`)

The `generate dotenv` command creates a `.env` file from your configuration:

```bash
egg_cli generate dotenv
```

This command:
- Creates a `.env` file from `config/developent.yaml` by default
- Can accept a specific config file path as an argument:
  ```bash
  egg_cli generate dotenv path/to/config.yaml
  ```
- Can read configuration from stdin for pipeline usage
- Supports the `--output` flag to specify output file location

### Creating Configuration from Environment Variables (`env`)

The `env` command creates a configuration file from environment variables:

```bash
egg_cli env
```

This command:
- Reads environment variables with the `EGG_` prefix
- Validates required variables (marked with X in the command help)
- Generates a configuration file for the specified environment
- Supports `--env` flag to specify environment (default: "dev")
- Supports `--input` flag to load variables from a specific `.env` file

Key environment variables include:
- `EGG_NAME` (required) - Project name
- `EGG_SEMVER` (required) - Project version
- `EGG_SERVER_PORT` (required) - Server port
- `EGG_SERVER_JWT` (required) - JWT secret
- Database, cache, and S3 configuration variables

### Version Information (`version`)

The `version` command displays version information:

```bash
egg_cli version
```

Available options:
- `--verbose` or `-v` - Show detailed version information
- `--hash` or `-#` - Show git commit hash
- `--build-time` or `-t` - Show build timestamp
- `--oneline` or `-1` - Display all info on one line
- `--compact` or `-c` - Compact format (version+hash+buildtime)

Examples:
```bash
egg_cli version                    # Basic version
egg_cli version --verbose          # Full details
egg_cli version --oneline          # Single line format
egg_cli version --compact          # Compact format
```

## Project Structure

egg_cli generates projects with this structure:

```
le-epic-project
├── cmd
│  ├── bump.go
│  ├── configuration
│  │  └── configuration.go
│  ├── db.go
│  ├── down.go
│  ├── generate.go
│  ├── migrate.go
│  ├── root.go
│  ├── swag.go
│  ├── up.go
│  └── version.go
├── config
│  └── development.yaml
├── controllers
│  ├── controller.go
│  ├── routes.go
│  └── user_controller.go
├── db
│  ├── migrations
│  │  └── 20250220222610_init.sql
│  ├── queries
│  │  ├── token.sql
│  │  └── user.sql
│  └── repository
│     ├── db.go
│     ├── models.go
│     ├── token.sql.go
│     └── user.sql.go
├── Dockerfile
├── docs
│  ├── docs.go
│  ├── swagger.json
│  └── swagger.yaml
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── middlewares
│  └── configs
│     ├── AuthConfig.go
│     └── StaticConfig.go
├── models
│  ├── handlers
│  │  ├── DeleteUserHandler.go
│  │  ├── GetCurrentLoggedInUserHandler.go
│  │  ├── GetProfilePictureHandler.go
│  │  ├── GetUsersHandler.go
│  │  ├── LoginHandler.go
│  │  ├── RegisterHandler.go
│  │  └── UploadProfilePictureHandler.go
│  ├── requests
│  │  ├── LoginRequest.go
│  │  └── NewUserRequest.go
│  └── responses
│     ├── DeleteUserResponse.go
│     ├── LoginResponse.go
│     ├── StringResponse.go
│     ├── UserResponse.go
│     └── UsersResponse.go
├── openapitools.json
├── README.md
├── services
│  ├── AuthService.go
│  ├── IAuthService.go
│  ├── IMinioService.go
│  ├── IRedisService.go
│  ├── IUserService.go
│  ├── MinioService.go
│  ├── MockAuthService.go
│  ├── MockUserService.go
│  ├── RedisService.go
│  ├── UserService.go
│  └── ValidatorService.go
├── sqlc.yml
├── tmp
│  └── main
└── web
   └── dist
      ├── egg.png
      └── index.html
```

## Development

### Building
```bash
make build
```

### Running Tests
```bash
go test ./...
```

### Clean 
```bash
make clean
```

### Development Dependencies
- Go 1.21+
- Make (for build automation)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./...`
5. Submit a pull request

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.

## Support

- 📚 [Documentation](https://github.com/adamkali/egg_cli/wiki)
- 🐛 [Issues](https://github.com/adamkali/egg_cli/issues)
- 💬 [Discussions](https://github.com/adamkali/egg_cli/discussions)





