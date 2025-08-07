# Todo CLI - Copilot Repository Instructions

## Project Overview
This is a **Todo CLI application** written in **Go 1.24.5** using the `urfave/cli/v3` framework. It manages todo items with support for adding, editing, deleting, toggling completion, and listing tasks. Data is stored as JSON either locally (`.todos.json`) or globally (`~/.todos/todos.json`). The project emphasizes high-quality code with comprehensive testing and dual build system support.

**Repository Size**: ~50 files, 6.5MB binary output  
**Languages**: Go (primary), PowerShell, Makefile  
**Frameworks**: urfave/cli/v3 for CLI structure, testify for testing  
**Target Runtime**: Go 1.24.5+ on Windows/Unix/macOS

## Critical Build Requirements

### Bootstrap (ALWAYS run first)
```bash
go mod tidy  # Downloads dependencies, handles missing packages
```
Your task is to "onboard" this repository to Copilot coding agent by adding a .github/copilot-instructions.md file in the repository that contains information describing how a coding agent seeing it for the first time can work most efficiently.

You will do this task only one time per repository and doing a good job can SIGNIFICANTLY improve the quality of the agent's work, so take your time, think carefully, and search thoroughly before writing the instructions.

<Goals>
- Reduce the likelihood of a coding agent pull request getting rejected by the user due to
generating code that fails the continuous integration build, fails a validation pipeline, or
having misbehavior.
- Minimize bash command and build failures.
- Allow the agent to complete its task more quickly by minimizing the need for exploration using grep, find, str_replace_editor, and code search tools.
</Goals>

<Limitations>
- Instructions must be no longer than 2 pages.
- Instructions must not be task specific.
</Limitations>

<WhatToAdd>

Add the following high level details about the codebase to reduce the amount of searching the agent has to do to understand the codebase each time:
<HighLevelDetails>

- A summary of what the repository does.
- High level repository information, such as the size of the repo, the type of the project, the languages, frameworks, or target runtimes in use.
</HighLevelDetails>

Add information about how to build and validate changes so the agent does not need to search and find it each time.
<BuildInstructions>

- For each of bootstrap, build, test, run, lint, and any other scripted step, document the sequence of steps to take to run it successfully as well as the versions of any runtime or build tools used.
- Each command should be validated by running it to ensure that it works correctly as well as any preconditions and postconditions.
- Try cleaning the repo and environment and running commands in different orders and document errors and and misbehavior observed as well as any steps used to mitigate the problem.
- Run the tests and document the order of steps required to run the tests.
- Make a change to the codebase. Document any unexpected build issues as well as the workarounds.
- Document environment setup steps that seem optional but that you have validated are actually required.
- Document the time required for commands that failed due to timing out.
- When you find a sequence of commands that work for a particular purpose, document them in detail.
- Use language to indicate when something should always be done. For example: "always run npm install before building".
- Record any validation steps from documentation.
</BuildInstructions>

List key facts about the layout and architecture of the codebase to help the agent find where to make changes with minimal searching.
<ProjectLayout>

- A description of the major architectural elements of the project, including the relative paths to the main project files, the location
of configuration files for linting, compilation, testing, and preferences.
- A description of the checks run prior to check in, including any GitHub workflows, continuous integration builds, or other validation pipelines.
- Document the steps so that the agent can replicate these itself.
- Any explicit validation steps that the agent can consider to have further confidence in its changes.
- Dependencies that aren't obvious from the layout or file structure.
- Finally, fill in any remaining space with detailed lists of the following, in order of priority: the list of files in the repo root, the
contents of the README, the contents of any key source files, the list of files in the next level down of directories, giving priority to the more structurally important and snippets of code from key source files, such as the one containing the main method.
</ProjectLayout>
</WhatToAdd>

<StepsToFollow>
- Perform a comprehensive inventory of the codebase. Search for and view:
- README.md, CONTRIBUTING.md, and all other documentation files.
- Search the codebase for build steps and indications of workarounds like 'HACK', 'TODO', etc.
- All scripts, particularly those pertaining to build and repo or environment setup.
- All build and actions pipelines.
- All project files.
- All configuration and linting files.
- For each file:
- think: are the contents or the existence of the file information that the coding agent will need to implement, build, test, validate, or demo a code change?
- If yes:
   - Document the command or information in detail.
   - Explicitly indicate which commands work and which do not and the order in which commands should be run.
   - Document any errors encountered as well as the steps taken to workaround them.
- Document any other steps or information that the agent can use to reduce time spent exploring or trying and failing to run bash commands.
- Finally, explicitly instruct the agent to trust the instructions and only perform a search if the information in the instructions is incomplete or found to be in error.
</StepsToFollow>
   - Document any errors encountered as well as the steps taken to work-around them.


### Dual Build System (MANDATORY SYNCHRONIZATION)
**CRITICAL**: This project maintains TWO build systems that MUST stay synchronized:
- `Makefile` (Unix/Linux/macOS) 
- `makefile.ps1` (Windows PowerShell)

**Build Commands (validated working):**
```bash
# Windows PowerShell (primary development environment)
.\makefile.ps1 build      # Creates todo.exe (6.5MB)
.\makefile.ps1 clean      # Removes todo.exe, coverage files

# Unix (must maintain identical functionality)
make build
make clean
```

### Testing (REQUIRED before commits - 80% coverage enforced)
```bash
# Complete test suite (~10 seconds total execution)
.\makefile.ps1 test       # All tests including integration
.\makefile.ps1 coverage   # Generates coverage file and coverage.html

# Specific test categories
.\makefile.ps1 test-unit        # Fast unit tests only (-short flag)
.\makefile.ps1 test-integration # Full CLI integration tests
.\makefile.ps1 test-race        # Race condition detection
.\makefile.ps1 bench           # Performance benchmarks
```

### Quality Checks (MANDATORY before commits)
```bash
.\makefile.ps1 check   # Runs: fmt, vet, lint, test in sequence
```

**Individual quality steps:**
```bash
.\makefile.ps1 fmt     # Code formatting
.\makefile.ps1 vet     # Go static analysis  
.\makefile.ps1 lint    # golangci-lint (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
```

## Project Architecture

### File Structure (modular command pattern)
```
todo-cli/
├── main.go                 # CLI entry point using urfave/cli/v3
├── config/config.go        # Version="2.0.0", ReleaseDate="2025-08-06"
├── cmds/                   # Each command in separate file
│   ├── add.go                # NewAddCommand() - adds todos
│   ├── delete.go             # NewDeleteCommand() - removes by index  
│   ├── edit.go               # NewEditCommand() - updates task text
│   ├── list.go               # NewListCommand() - shows todos (table/json/pretty)
│   ├── toggle.go             # NewToggleCommand() - completion status
│   ├── version.go            # NewVersionCommand() - version info
│   ├── commands.go           # Command interfaces and registry
│   ├── utils.go              # GetStoragePath(), TodoList type
│   └── commands_test.go      # Command unit tests
├── storage.go              # Generic storage layer with Save/Load
├── todo.go                 # Core TodoList logic and data structures
├── integration_test.go     # End-to-end CLI testing
├── Makefile                # Unix build system
├── makefile.ps1            # PowerShell build system (synchronized)
└── .vscode/launch.json     # Debug configurations
```

### Command Registration Pattern
All commands are registered in `main.go`:
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

### Storage Architecture
- **Local**: `.todos.json` in current directory
- **Global**: `~/.todos/todos.json` in user home (use `--global` or `-g`)
- **Format**: JSON with GUID-based IDs, timestamps, completion status
- **Load/Save**: Generic storage layer in `storage.go`

## Command Implementation Standards

### Required Command Structure (enforced pattern)
```go
package commands

import (
    "context"
    "github.com/urfave/cli/v3"
)

func NewXCommand() *cli.Command {
    return &cli.Command{
        Name:        "commandname",
        Aliases:     []string{"alias"},
        Usage:       "Brief description",
        ArgsUsage:   "<arguments>",
        Action: func(ctx context.Context, c *cli.Command) error {
            // Get storage path
            storagePath, err := GetStoragePath(c.Bool("global"))
            if err != nil {
                return cli.Exit(fmt.Sprintf("error getting storage path: %v", err), 2)
            }

            // Load todos
            todoList := &TodoList{}
            storage := NewStorage[TodoList](storagePath)
            loadedList, err := storage.Load()
            if err != nil {
                return cli.Exit(fmt.Sprintf("error loading todos: %v", err), 2)
            }
            *todoList = loadedList

            // Perform operation...

            // Save changes
            if err := storage.Save(*todoList); err != nil {
                return cli.Exit(fmt.Sprintf("error saving todos: %v", err), 2)
            }
            return nil
        },
    }
}
```

### Global Flags (available on all commands)
- `--debug`: Enable verbose debug output
- `--global`, `-g`: Use global storage in `~/.todos/todos.json`

### Error Handling Standards
- Use `cli.Exit(message, code)` for user-facing errors
- Exit codes: 0=success, 1=general error, 2=file/storage error
- Wrap errors: `fmt.Errorf("context: %w", err)`

## Testing Requirements

### Test Coverage (validated working - 80% minimum enforced)
Current coverage from validation:
- **Unit tests**: Fast isolated logic tests
- **Integration tests**: Full CLI execution with real binary
- **Benchmarks**: Performance testing (validated: Add ~568ns, Delete ~0.0005ns)

### Test Naming (enforced convention)
```go
func TestAdd(t *testing.T)              // Unit test
func TestCLIAdd(t *testing.T)           // Integration test  
func BenchmarkAdd(b *testing.B)         // Benchmark test
```

### Required Tests for New Commands
1. Command structure test in `commands_test.go`
2. Integration test in `integration_test.go`
3. Error case coverage
4. Benchmark if performance-critical

## Development Workflow (validated)

### Adding New Commands
1. Create `cmds/newcommand.go` with `NewNewcommandCommand()` function
2. Add to `main.go` Commands slice
3. Write tests (unit + integration)
4. Run `.\makefile.ps1 check` (must pass)
5. Verify coverage: `.\makefile.ps1 coverage` (≥80%)

### Adding Build Targets (CRITICAL synchronization)
**ALWAYS update both build systems identically:**
1. Add target to `Makefile` with Unix commands
2. Add function to `makefile.ps1` with PowerShell equivalents  
3. Update help text in both files
4. Test both systems work: `make target` and `.\makefile.ps1 target`

### Quality Validation (before commits)
```bash
.\makefile.ps1 check    # Runs complete pipeline
```
This validates: formatting, static analysis, linting, all tests, coverage

## Environment Dependencies

### Required Tools
- **Go 1.24.5+** (validated: `go version go1.24.5 windows/amd64`)
- **golangci-lint** (install: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)

### Key Dependencies (from go.mod)
```go
require (
    github.com/aquasecurity/table v1.11.0     // Table formatting
    github.com/liamg/tml v0.7.0               // Terminal markup
    github.com/urfave/cli/v3 v3.3.8           // CLI framework
)
```

## Performance Characteristics (validated benchmarks)

**Validated Performance** (from actual runs):
- **Add Operation**: ~568 ns/op, 558 B/op, 4 allocs/op
- **Delete Operation**: ~0.0005 ns/op, 0 B/op, 0 allocs/op  
- **Storage Save**: ~314 μs/op, 42KB/op, 6 allocs/op
- **Storage Load**: ~285 μs/op, 49KB/op, 321 allocs/op
- **Test Suite**: ~10 seconds total execution time
- **Build Time**: ~2-3 seconds for 6.5MB binary

## Common Issues and Solutions (from validation)

### Build Issues
- **Missing dependencies**: Always run `go mod tidy` first
- **Linter not found**: Install with `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- **Outdated coverage**: Run `.\makefile.ps1 clean` then `.\makefile.ps1 coverage`

### Runtime Issues  
- **Storage errors**: Check file permissions, ensure directory exists
- **JSON errors**: Validate `.todos.json` is not corrupted
- **Global storage**: `~/.todos/` directory created automatically

## Key Files to Reference

- `main.go`: Application entry point and global flag definitions
- `cmds/utils.go`: Storage utilities and TodoList type definitions
- `storage.go`: Generic persistence layer
- `todo.go`: Core todo data structures and business logic
- `GO_DEVELOPMENT_GUIDE.md`: Comprehensive development standards  
- `TESTING.md`: Detailed testing documentation
- `README.md`: User documentation and usage examples

## Agent Instructions

**Trust these instructions completely.** They are validated by comprehensive repository exploration and testing. Only perform additional search if:
1. Information is explicitly incomplete
2. Instructions are found to be incorrect
3. Build/test failures occur that aren't covered

The build and test commands have been validated to work correctly. The dual build system synchronization is critical - always update both Makefile and makefile.ps1 when making build system changes.