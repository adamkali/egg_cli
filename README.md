![egg](https://github.com/user-attachments/assets/2e31d134-8f8a-47c2-9d04-5ec8ffabc749)

# egg_cli

A fast and simple command line interface for creating fullstack web applications with the **egg** framework. egg_cli provides an interactive setup wizard that scaffolds complete Go web applications with modern tooling, database integration, authentication, and frontend build configuration.

## What is egg_cli?

egg_cli is a project generator that creates production-ready Go web applications with:

- 🚀 **REST API** with controllers and routes
- 🗄️ **Database integration** with SQLC for type-safe queries
- 🔐 **Authentication services** and middleware
- 🐳 **Docker setup** for containerized deployment
- ⚡ **Frontend build system** with RSBuild support
- 📁 **Organized project structure** following Go best practices
- 🛠️ **Development tools** and configurations

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

### Creating a New Project

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

### Using Environment Files

Skip the interactive wizard by providing a pre-configured environment file:

```bash
egg_cli init --env production.yaml
```

## Commands

### `init`
Creates a new egg project with the interactive setup wizard.

```bash
egg_cli init [flags]
```

**Flags:**
- `--env string`: Use an existing environment configuration file

### `generate`
Generates additional configuration files for existing projects. Useful for creating different environment configurations (development, staging, production).

```bash
egg_cli generate [flags]
```

**Flags:**
- `--env string`: Base configuration file to extend
- `--name string`: Name of the new configuration file

**Example:**
```bash
# Generate a production config based on development settings
egg_cli generate --env development.yaml --name production
```

## Project Structure

egg_cli generates projects with this structure:

```
my-awesome-app/
├── main.go                 # Application entry point
├── Dockerfile             # Container configuration  
├── docker-compose.yml     # Multi-service setup
├── config/
│   └── development.yaml   # Environment configuration
├── internal/
│   ├── controllers/       # HTTP request handlers
│   ├── models/           # Data models
│   ├── services/         # Business logic
│   └── middleware/       # HTTP middleware
├── migrations/           # Database migrations
├── queries/             # SQLC database queries
└── frontend/            # Frontend application (if selected)
    ├── src/
    └── rsbuild.config.js
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

MIT License. See [LICENSE](LICENSE) for details.

## Support

- 📚 [Documentation](https://github.com/adamkali/egg_cli/wiki)
- 🐛 [Issues](https://github.com/adamkali/egg_cli/issues)
- 💬 [Discussions](https://github.com/adamkali/egg_cli/discussions)





