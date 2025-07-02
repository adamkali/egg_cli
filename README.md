![egg](https://github.com/user-attachments/assets/80dcb1e9-b7a5-486c-8ed9-d628ed265ee2)

# egg_cli

A Fast and Simple Command Line Interface for Creating the Fullstack web framework egg. 

## Installation

```bash
go install github.com/adamkali/egg_cli@latest
```

## Usage

### Init
This will spin up the TUI configuration wizard to guide you through the creation of a new egg project, and is the recommended way to get started.
```bash
egg_cli init
```

This can also be used with the `--env` flag to use a certain environment file to skip through the wizard.

### Generate
This will spin up the TUI configuration wizard to guide you through the creation of a new configuraion file,
this can be useful if you already have an existing project but need to test a new database that has many nodes 
or other configuration options that you want to change from the development.yaml (default).

```bash
egg_cli generate
```

This can also be used with the `--env` flag to use a certain environment file to skip through the wizard.

This can also be used with the `--name` flag to set the name of the configuration file to be created outside of the wizard.





