# Todo CLI

## Overview
Todo CLI is a command-line application for managing a to-do list. It allows users to add, view, and manage tasks efficiently. The application is written in Go and uses JSON for data serialization.

This based on the tutorial https://codingwithpatrik.dev/posts/how-to-build-a-cli-todo-app-in-go/

## Features
- Add tasks to the to-do list
- View all tasks
- Mark tasks as completed
- Delete tasks
- Edit existing tasks
- Local and global storage options
- Multiple output formats (table, JSON, pretty JSON)

## Storage Options

### Local Storage (Default)
By default, todos are stored in `.todos.json` in the current working directory.

### Global Storage
Use the `--global` or `-g` flag to store todos in your user home directory at `~/.todo/todos.json`. This allows you to access your todos from anywhere on your system.

```bash
# Add a todo to global storage
.\todo.exe --global add "Global task"

# List todos from global storage
.\todo.exe --global list

# Edit a global todo
.\todo.exe --global edit 1 "Updated global task"
```

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/rocketcentral/todo-cli.git
   ```
2. Navigate to the project directory:
   ```bash
   cd todo-cli
   ```
3. Build the application:
   ```bash
   go build -o todo.exe
   ```

## Usage

### Basic Commands
Run the application using the following commands:

```bash
# Show all todos (default action)
.\todo.exe

# Add a new todo
.\todo.exe add "Buy groceries"

# List todos in different formats
.\todo.exe list                    # Table format (default)
.\todo.exe list --format json      # JSON format
.\todo.exe list --format pretty    # Pretty JSON format

# Edit a todo
.\todo.exe edit 1 "Updated task"

# Toggle completion status
.\todo.exe toggle 1

# Delete a todo
.\todo.exe delete 1

# Show version
.\todo.exe version

# Show help
.\todo.exe help
```

### Global Storage
Use the `--global` or `-g` flag with any command to work with global storage:

```bash
# Add to global storage
.\todo.exe --global add "Global task"

# List from global storage
.\todo.exe --global list

# Default action with global storage
.\todo.exe --global
```

### Installation to PATH
Add the application to your PATH, and then call it with:
```bash
todo
```

I created a symlink in a folder that is in my PATH with the following powershell:
```pwsh
New-Item -ItemType SymbolicLink -Path <LINK> -Target <ACTUAL SOURCE>
```

## Examples

### Local Storage Example
```bash
# Add and list todos locally
.\todo.exe add "Buy groceries"
.\todo.exe add "Walk the dog"
.\todo.exe list

# Output in JSON format
.\todo.exe list --format json

# Pipe to other commandlets
.\todo.exe list --format json | jq
```

### Global Storage Example
```bash
# Work with global todos
.\todo.exe --global add "Review project proposal"
.\todo.exe --global add "Update documentation"
.\todo.exe --global list

# Toggle completion
.\todo.exe --global toggle 1
.\todo.exe --global list
```

## Development
### Prerequisites
- Go 1.24.5 or later

### Running Tests
To run tests, use:
```bash
go test ./...

# Or use the PowerShell build script
.\makefile.ps1 test
```

### Build Targets
The project includes both a Makefile and PowerShell script for cross-platform development:

```bash
# Using PowerShell script (Windows)
.\makefile.ps1 build      # Build the application
.\makefile.ps1 test       # Run tests
.\makefile.ps1 coverage   # Generate coverage report
.\makefile.ps1 check      # Run all quality checks

# Using Makefile (Unix-like systems)
make build
make test
make coverage
```

## Storage Locations

- **Local**: `.todos.json` in the current working directory
- **Global**: `~/.todo/todos.json` in the user's home directory

The global storage directory (`~/.todo/`) is automatically created when first used and can be used for future configuration files and extensions.