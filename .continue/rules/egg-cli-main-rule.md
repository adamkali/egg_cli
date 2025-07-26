---
name: egg_cli
description: egg_cli is a tool for creating fullstack frameworks with using golang, the echo framework, and goose for database migrations.
---

- Egg CLI is a tool for creating fullstack frameworks with using golang, the echo framework, and goose for database migrations.
- This is a description of the `init` command in the egg_cli tool.
- The base code for the `init` command is located in the following directory: `/home/adamkali/projects/egg_cli/cmd/init.go`
- It uses the bubbletea framework for creating the installtion wizard.
- The installation wizard code is located in the following directory: `/home/adamkali/projects/egg_cli/pkg/models/` and is called in `/home/adamkali/projects/egg_cli/pkg/models/root_model.go`.
- After the user is done with the installation wizard, the code is kicked off to generate the project.
- this generation is handled by the code in the following directory: `/home/adamkali/projects/egg_cli/pkg/modules/` and is called in `/home/adamkali/projects/egg_cli/pkg/module_runner.go`. 
- `/home/adamkali/projects/egg_cli/pkg/module_runner.go` is runnes each of the moduelse in the `/modules/` directory and then generates what is needed from the requirements in defined in struct that satisfies `/home/adamkali/projects/egg_cli/pkg/modules/i_module.go`.
- The generator code then runs the code, and if it does not succeed the generator should output a file to `.scrambled` file to see what happend and then use the scrambled file to reload the state of the generator.
