# CLAUDE.md - egg_cli Project Context

## Project Overview
egg_cli is a command line interface tool for creating fullstack web framework projects called "egg". It's written in Go and uses Cobra for CLI functionality and Bubbletea for TUI interactions.

## Key Components

### Architecture
- **CLI Framework**: Cobra-based CLI with subcommands
- **TUI**: Bubbletea for interactive configuration wizards
- **Module System**: Pluggable modules for different project setup tasks
- **Template System**: Go templates embedded in the binary for file generation
- **Configuration**: YAML-based configuration with interactive generation

### Main Commands
- `egg_cli init` - Initialize new egg project with TUI wizard
- `egg_cli generate` - Generate configuration files with TUI wizard

### Module System
Located in `pkg/modules/`, implements `IModule` interface:
- `egg::initialize` - Project initialization
- `egg::install_tools` - Tool installation 
- `egg::install_libraries` - Library installation
- `egg::bootstrap_directories` - Directory structure creation
- `egg::generate_configuration` - Configuration file generation
- `egg::bootstrap_framework` - Framework file generation from templates
- `egg::rsbuild_frontend` - Frontend setup with RSBuild

### Template System
- Templates are embedded Go files in `pkg/templates/`
- Each template generates specific project files (main.go, Dockerfile, etc.)
- Mapping defined in `pkg/targets/mapping.yaml`
- Generates fullstack Go web applications with:
  - REST API controllers and routes
  - Database models and migrations (SQLC)
  - Authentication services
  - Middleware configurations
  - Docker setup
  - Frontend build configuration

### Project Structure
```
cmd/           - CLI command definitions
pkg/
  configuration/ - Configuration management
  models/       - Data models and TUI models
  modules/      - Module implementations
  targets/      - Target mapping configuration
  templates/    - Embedded Go templates
  state/        - Application state management
styles/        - Bubbletea styling components
```

## Development Commands

### Build and Test
```bash
make build         # Build binary with version info
make test          # Run tests
make tidy          # Update and verify dependencies
make clean         # Clean build artifacts
make all           # Run tidy, test, and build
```

### Legacy Commands (still supported)
```bash
go mod tidy        # Update dependencies
go build -o egg_cli # Build binary
go test ./...      # Run tests
```

### Release Management
```bash
make tag VERSION=v1.0.0  # Create and push a new version tag
make push-tag            # Push all existing tags to remote
make version             # Show current version information
```

### Key Files to Understand
- `main.go` - Entry point with embedded templates
- `cmd/root.go` - Root command definition
- `pkg/modules/i_module.go` - Module interface and factory
- `pkg/targets/mapping.yaml` - Template to file mapping
- `pkg/configuration/configuration.go` - Configuration structure

## Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - TUI styling
- `gopkg.in/yaml.v3` - YAML parsing

## Current State (Updated July 31, 2025)
This project has completed several major improvements:

### Recently Completed
- ✅ **Template System**: Transitioned from hardcoded template strings to embedded Go file template system
- ✅ **Documentation**: Enhanced README.md with comprehensive user guides, examples, and project structure explanations
- ✅ **Build System**: Added Makefile with automated build, test, and release management targets
- ✅ **Developer Experience**: Improved development workflow with standardized commands and version tracking
- ✅ **CI/CD Pipeline**: Implemented comprehensive testing and automated release pipeline with branch-specific versioning

### Development Maturity
- The core CLI and TUI framework is stable and functional
- Template system is working with embedded Go templates
- Interactive project generation wizard is operational
- CI/CD pipeline with automated testing and releases is now operational
- Branch-based development workflow established with dev and main branches
- Ready for broader testing and usage

### CI/CD Implementation (July 31, 2025)
- **Testing Integration**: Added automated testing with `make test`, Go vet, and staticcheck
- **Branch Strategy**: 
  - `dev` branch: Development and testing, increments minor version (x.Y.z) on merge
  - `main` branch: Production releases, increments major version (X.y.z) on merge
- **Automated Releases**: Creates GitHub releases with cross-platform binaries automatically
- **Build Matrix**: Supports Linux, macOS, and Windows on both amd64 and arm64 architectures

## Future Plans
- Add a way to create either a preferred React or Svelte frontend with RSBuild that will just be very simple to use and work with
- Add more customization options for the project structure and file generation
- Implement automated testing for generated project templates
- Add integration tests for the full project generation workflow
- Consider adding support for additional database types and ORMs
- Improve error handling and user feedback in the TUI wizard
- Consider adding automated changelog generation
- Implement semantic versioning with patch releases for bug fixes

## For Other AI Agents
When working on this project:
1. Always use `make build` instead of `go build` for consistent builds with version information
2. Run `make test` before making changes to ensure nothing breaks
3. The project follows Go best practices and uses Cobra+Bubbletea for CLI/TUI
4. Templates are located in `pkg/templates/` and mapped via `pkg/targets/mapping.yaml`
5. The module system in `pkg/modules/` handles different aspects of project generation
6. Configuration is YAML-based and managed through the TUI wizard
7. **Branch Strategy**: Work on `dev` branch for development, merge to `main` for production releases
8. **CI/CD**: The pipeline automatically runs tests, builds cross-platform binaries, and creates releases on branch pushes
9. **Version Management**: Automated semantic versioning - dev branch increments minor, main branch increments major


