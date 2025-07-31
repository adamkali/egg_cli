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
go mod tidy        # Update dependencies
go build -o egg_cli # Build binary
go test ./...      # Run tests
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

## Current State
This project has just completed transitioning from hardcoded template strings to a embedded Go file template system. It's still in early development, so expect changes and improvements.

## Future Plans
- Add a way to create either a preferred React or Svelte frontend with RSBuild that will just be very simple to use and work with. 
- Add more customization options for the project structure and file generation.


