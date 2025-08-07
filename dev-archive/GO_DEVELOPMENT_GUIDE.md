# Go Development Guide - Todo CLI Project

## Table of Contents
1. [Project Overview](#project-overview)
2. [Command Structure Standards](#command-structure-standards)
3. [PowerShell & Makefile Synchronization](#powershell--makefile-synchronization)
4. [Testing Standards](#testing-standards)
5. [Code Organization](#code-organization)
6. [Development Workflow](#development-workflow)
7. [Quality Assurance](#quality-assurance)

## Project Overview

This is a Todo CLI application built with Go using the `urfave/cli/v3` framework. The project follows a modular command structure where each command is separated into its own file for maintainability and clarity.

### Key Technologies
- **Framework**: `urfave/cli/v3` for CLI structure
- **Testing**: Go's built-in testing package + `testify` for assertions
- **Build System**: Dual system with `Makefile` (Unix) and `makefile.ps1` (PowerShell)

## Command Structure Standards

### 1. File Organization
Each command MUST be in its own file within the `cmds/` directory:

```
cmds/
├── add.go          # Add command implementation
├── delete.go       # Delete command implementation
├── edit.go         # Edit command implementation
├── list.go         # List command implementation
├── toggle.go       # Toggle command implementation
├── version.go      # Version command implementation
├── commands.go     # Command registry and interfaces
├── commands_test.go # Tests for command functionality
└── utils.go        # Shared utilities for commands
```

### 2. Command File Structure
Each command file MUST follow this structure:

```go
package commands

import (
    "context"
    "github.com/urfave/cli/v3"
)

// NewXCommandName creates a new X command
func NewXCommandName() *cli.Command {
    return &cli.Command{
        Name:        "commandname",
        Aliases:     []string{"alias1", "alias2"},
        Usage:       "Brief description of what the command does",
        Description: "Detailed description with examples",
        Flags: []cli.Flag{
            // Command-specific flags
        },
        Action: func(ctx context.Context, c *cli.Command) error {
            // Command implementation
            return nil
        },
    }
}
```

### 3. Command Registration
All commands MUST be registered in `main.go` in the `Commands` slice:

```go
Commands: []*cli.Command{
    commands.NewAddCommand(),
    commands.NewDeleteCommand(),
    commands.NewEditCommand(),
    commands.NewListCommand(),
    commands.NewToggleCommand(),
    commands.NewVersionCommand(),
},
```

### 4. Naming Conventions
- **File names**: `{command}.go` (e.g., `add.go`, `delete.go`)
- **Function names**: `New{Command}Command()` (e.g., `NewAddCommand()`)
- **Command names**: Lowercase, descriptive (e.g., `"add"`, `"delete"`)
- **Aliases**: Single letter when possible (e.g., `"a"` for add, `"d"` for delete)

## PowerShell & Makefile Synchronization

**CRITICAL REQUIREMENT**: Both `Makefile` and `makefile.ps1` MUST be kept in perfect synchronization.

### Synchronization Rules

1. **Every target/function added to one MUST be added to the other**
2. **Target names MUST be identical** (case-sensitive)
3. **Functionality MUST be equivalent** between both systems
4. **Help output MUST match** exactly

### Adding New Build Targets

When adding a new build target, follow this process:

#### Step 1: Add to Makefile
```makefile
# New target description
new-target:
	@echo "Executing new target"
	# commands here
```

#### Step 2: Add to makefile.ps1
```powershell
function New-Target {
    Write-Host "Executing new target" -ForegroundColor Green
    # equivalent commands here
}

# Add to switch statement
"new-target" { New-Target }

# Add to help function
Write-Host "  new-target    - New target description"
```

### Current Synchronized Targets
- `all` / `all`
- `build` / `build`
- `clean` / `clean`
- `test` / `test`
- `coverage` / `coverage`
- `test-unit` / `test-unit`
- `test-integration` / `test-integration`
- `bench` / `bench`
- `bench-verbose` / `bench-verbose`
- `test-race` / `test-race`
- `lint` / `lint`
- `fmt` / `fmt` (also `format` alias in PowerShell)
- `vet` / `vet`
- `check` / `check`
- `deps` / `deps`
- `install` / `install`
- `help` / `help`

### PowerShell-Specific Considerations
- Use `Write-Host` with colors for better UX
- Handle missing tools gracefully with `Get-Command -ErrorAction SilentlyContinue`
- Use `Remove-Item -ErrorAction SilentlyContinue` for file cleanup
- Provide helpful warnings when tools are missing

## Testing Standards

### Test File Organization
- **Unit tests**: Co-located with source files (e.g., `commands_test.go`)
- **Integration tests**: In root directory (`integration_test.go`)
- **Component tests**: Specific test files for each component

### Test Naming Conventions
```go
// Unit tests
func TestAdd(t *testing.T) { ... }
func TestAddWithInvalidInput(t *testing.T) { ... }

// Integration tests  
func TestCLIAdd(t *testing.T) { ... }
func TestCLIIntegration(t *testing.T) { ... }

// Benchmark tests
func BenchmarkAdd(b *testing.B) { ... }
```

### Test Requirements

#### 1. Every New Command MUST Have Tests
When adding a new command:
- Unit tests for the command logic
- Integration tests for CLI behavior
- Error case testing
- Flag validation testing

#### 2. Test Categories
- **Unit tests** (`-short` flag): Fast, isolated tests
- **Integration tests**: Full CLI testing with real file system
- **Benchmark tests**: Performance testing

#### 3. Test Coverage Requirements
- Maintain **minimum 80% code coverage**
- Use `make coverage` or `.\makefile.ps1 coverage` to generate reports
- Review `coverage.html` for uncovered areas

#### 4. Test Structure
```go
func TestCommandName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid input", "test task", "Task added", false},
        {"empty input", "", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Running Tests
```bash
# All tests
make test
.\makefile.ps1 test

# Unit tests only
make test-unit
.\makefile.ps1 test-unit

# Integration tests only
make test-integration  
.\makefile.ps1 test-integration

# With coverage
make coverage
.\makefile.ps1 coverage

# With race detection
make test-race
.\makefile.ps1 test-race
```

## Code Organization

### Package Structure
```
todo-cli/
├── main.go              # CLI application entry point
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
├── Makefile            # Unix build system
├── makefile.ps1        # PowerShell build system
├── config/
│   └── config.go       # Application configuration
├── cmds/               # Command implementations
│   ├── commands.go     # Command interfaces and registry
│   ├── commands_test.go # Command tests
│   ├── utils.go        # Shared command utilities
│   ├── add.go          # Add command
│   ├── delete.go       # Delete command
│   ├── edit.go         # Edit command
│   ├── list.go         # List command
│   ├── toggle.go       # Toggle command
│   └── version.go      # Version command
├── storage.go          # Data persistence layer
├── storage_test.go     # Storage tests
├── todo.go             # Core todo logic
├── todo_test.go        # Todo logic tests
└── integration_test.go # End-to-end tests
```

### Import Standards
```go
// Standard library imports first
import (
    "context"
    "fmt"
    "os"
)

// Third-party imports second
import (
    "github.com/urfave/cli/v3"
    "github.com/stretchr/testify/assert"
)

// Local imports last
import (
    "github.com/bennthewolfe/todo-cli/config"
)
```

## Development Workflow

### 1. Feature Development Process
1. **Create feature branch**: `git checkout -b feature/command-name`
2. **Implement command**: Create new file in `cmds/`
3. **Add tests**: Unit and integration tests
4. **Update documentation**: Update this guide if needed
5. **Run quality checks**: `make check` or `.\makefile.ps1 check`
6. **Update build system**: If new targets needed, update both Makefile and makefile.ps1
7. **Commit and push**: Create pull request

### 2. Before Committing
Always run the complete quality check:
```bash
make check
# OR
.\makefile.ps1 check
```

This runs:
- Code formatting (`go fmt`)
- Code vetting (`go vet`) 
- Linting (`golangci-lint`)
- All tests

### 3. Adding New Dependencies
```bash
# Add new dependency
go get github.com/new/package

# Clean up dependencies
go mod tidy

# Commit both go.mod and go.sum
git add go.mod go.sum
git commit -m "Add new dependency: github.com/new/package"
```

## Quality Assurance

### 1. Code Standards
- **Formatting**: Use `gofmt` (run with `make fmt`)
- **Linting**: Use `golangci-lint` (run with `make lint`)
- **Vetting**: Use `go vet` (run with `make vet`)

### 2. Error Handling
```go
// Always handle errors explicitly
if err != nil {
    return cli.Exit(fmt.Sprintf("error message: %v", err), 1)
}

// Use cli.Exit for user-facing errors with appropriate exit codes
return cli.Exit("File not found", 2)
```

### 3. Logging and Debug Output
```go
// Use debug flag for verbose output
if c.Bool("debug") {
    fmt.Printf("DEBUG: Operation details\n")
}

// Use consistent error messages
return fmt.Errorf("failed to load todos: %w", err)
```

### 4. Git Hooks (Recommended)
Set up pre-commit hooks to ensure quality:

```bash
# .git/hooks/pre-commit (make executable)
#!/bin/sh
make check
```

### 5. Performance Considerations
- Run benchmarks for performance-critical code
- Use `make bench` to run benchmark tests
- Profile memory usage for large todo lists
- Monitor startup time for CLI responsiveness

## Common Patterns

### 1. Command Implementation Pattern
```go
func NewExampleCommand() *cli.Command {
    return &cli.Command{
        Name:    "example",
        Aliases: []string{"ex"},
        Usage:   "Example command usage",
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:     "input",
                Aliases:  []string{"i"},
                Usage:    "Input value",
                Required: true,
            },
        },
        Action: func(ctx context.Context, c *cli.Command) error {
            // Get storage path
            storagePath, err := GetStoragePath(c.Bool("global"))
            if err != nil {
                return cli.Exit(fmt.Sprintf("error getting storage path: %v", err), 2)
            }

            // Load todo list
            todoList := &TodoList{}
            storage := NewStorage[TodoList](storagePath)
            loadedList, err := storage.Load()
            if err != nil {
                return cli.Exit(fmt.Sprintf("error loading todos: %v", err), 2)
            }
            *todoList = loadedList

            // Perform operation
            input := c.String("input")
            if err := todoList.SomeOperation(input); err != nil {
                return cli.Exit(fmt.Sprintf("operation failed: %v", err), 1)
            }

            // Save changes
            if err := storage.Save(*todoList); err != nil {
                return cli.Exit(fmt.Sprintf("error saving todos: %v", err), 2)
            }

            return nil
        },
    }
}
```

### 2. Test Pattern
```go
func TestExampleCommand(t *testing.T) {
    // Setup
    tempDir := t.TempDir()
    todoFile := filepath.Join(tempDir, "todos.json")
    
    // Test cases
    tests := []struct {
        name    string
        args    []string
        want    string
        wantErr bool
    }{
        {
            name: "valid operation",
            args: []string{"--input", "test"},
            want: "success message",
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

---

## Summary

This guide ensures:
- **Consistent command structure** across all CLI commands
- **Perfect synchronization** between Makefile and makefile.ps1
- **Comprehensive testing** with every change
- **High code quality** through automated checks
- **Clear development workflow** for all contributors

Follow these standards religiously to maintain code quality and ensure smooth development experience across different platforms (Unix/Linux/macOS with Makefile, Windows with PowerShell).
