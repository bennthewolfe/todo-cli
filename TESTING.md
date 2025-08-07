# Testing Documentation for Todo CLI

This document describes the comprehensive test suite for the Todo CLI project.

## Test Structure

The test suite includes:

### 1. Unit Tests
- **`storage_test.go`** - Tests for the storage functionality (Save/Load operations)
- **`todo_test.go`** - Tests for TodoList operations (Add, Delete, Update, Toggle)
- **`cmds/commands_test.go`** - Tests for command creation and structure

### 2. Integration Tests
- **`integration_test.go`** - End-to-end CLI testing with actual binary execution

### 3. Benchmark Tests
- **`main_test.go`** - Performance benchmarks for core operations

## Running Tests

### Cross-Platform Build Scripts

This project provides both **Unix/Linux Makefile** and **Windows PowerShell** build scripts with complete feature parity.

#### On Unix/Linux/macOS (using Makefile)

```bash
# Run all tests
make test

# Generate coverage report
make coverage

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run benchmark tests
make bench

# Run tests with race detection
make test-race

# Format, vet, lint, and test
make check

# Default workflow (test then build)
make all

# Show all available targets
make help
```

#### On Windows (using PowerShell script)

```powershell
# Run all tests
.\makefile.ps1 test

# Generate coverage report
.\makefile.ps1 coverage

# Run unit tests only
.\makefile.ps1 test-unit

# Run integration tests only
.\makefile.ps1 test-integration

# Run benchmark tests
.\makefile.ps1 bench

# Run tests with race detection
.\makefile.ps1 test-race

# Format, vet, lint, and test
.\makefile.ps1 check

# Default workflow (test then build)
.\makefile.ps1 all

# Show all available targets
.\makefile.ps1 help
```

### Basic Test Commands (Direct Go)

If you prefer using Go commands directly:

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out -v ./...

# View coverage report in browser
go tool cover -html=coverage.out

# Run only unit tests (fast)
go test -v -short ./...

# Run only integration tests
go test -v -run TestCLI ./...

# Run benchmark tests
go test -bench=. -run=^$ .

# Run benchmark tests with memory allocation stats
go test -bench=Benchmark -benchmem -run=^$

# Run tests with race condition detection
go test -v -race ./...
```

## Available Build Targets

Both build systems support the following targets with identical functionality:

| Target | Description |
|--------|-------------|
| `all` | Run default workflow (test then build) |
| `build` | Build the application |
| `clean` | Clean build files |
| `test` | Run all tests |
| `test-unit` | Run unit tests only |
| `test-integration` | Run integration tests only |
| `coverage` | Run tests with coverage report |
| `bench` | Run benchmark tests only |
| `bench-verbose` | Run benchmark tests with all tests |
| `test-race` | Run tests with race condition detection |
| `lint` | Lint the code (requires golangci-lint) |
| `fmt` / `format` | Format the code |
| `vet` | Vet the code |
| `check` | Run all quality checks (fmt, vet, lint, test) |
| `deps` | Download dependencies |
| `install` | Install the application |
| `help` | Show help with all available targets |

## Test Coverage

Current test coverage:
- **Main package**: ~44% coverage
- **Commands package**: ~13% coverage

### What's Tested

✅ **Storage Operations**
- Save/Load functionality
- JSON marshaling/unmarshaling
- File creation and error handling
- Invalid JSON handling

✅ **TodoList Operations**
- Adding tasks
- Deleting tasks by index
- Updating task content
- Toggling completion status
- Index validation
- Timestamp management

✅ **Command Structure**
- Command creation and configuration
- Command aliases and usage text
- Flag definitions

✅ **CLI Integration**
- Error handling and exit codes
- Command-line argument parsing
- End-to-end workflows
- Help system functionality

✅ **Performance Benchmarks**
- Add operation performance
- Delete operation performance
- Storage save/load performance

### Test Examples

#### Adding a Todo Item
```go
func TestTodoList_Add(t *testing.T) {
    todoList := &TodoList{}
    err := todoList.Add("Test task")
    // Verifies task is added with correct ID, timestamps, etc.
}
```

#### CLI Integration Test
```go
func TestCLIWorkflow(t *testing.T) {
    // Tests complete workflow: add -> list -> edit -> toggle -> delete
    // Builds actual binary and executes commands
}
```

#### Storage Test
```go
func TestStorage_Save(t *testing.T) {
    // Tests saving TodoList to JSON file
    // Verifies file creation and content
}
```

## Test Utilities

### Test Helpers
- `setupTestEnvironment()` - Creates isolated test directories
- `testArgs` - Mock implementation of CLI Args interface
- Temporary file/directory management

### Test Isolation
- Each integration test runs in its own temporary directory
- No test data pollution between tests
- Automatic cleanup after test completion

## Continuous Integration

The test suite is designed to be CI-friendly:
- Fast unit tests (< 1 second)
- Isolated integration tests
- Deterministic results
- Clear error messages
- Exit codes for automation

## Performance Benchmarks

Current benchmark results (as of August 6, 2025):
```
goos: windows
goarch: amd64
pkg: github.com/bennthewolfe/todo-cli
cpu: 11th Gen Intel(R) Core(TM) i7-1165G7 @ 2.80GHz
BenchmarkTodoList_Add-8       2219103       516.7 ns/op     580 B/op       4 allocs/op
BenchmarkTodoList_Delete-8   1000000000         0.0004443 ns/op       0 B/op       0 allocs/op
BenchmarkStorage_Save-8          3231    319042 ns/op   42001 B/op       6 allocs/op
BenchmarkStorage_Load-8          4354    284935 ns/op   49128 B/op     321 allocs/op
```

### Benchmark Analysis
- **Add Operation**: ~517 ns per operation with 4 memory allocations (580 bytes)
- **Delete Operation**: Extremely fast at ~0.0004 ns per operation with zero allocations
- **Storage Save**: ~319 μs per operation with 6 allocations (42KB)
- **Storage Load**: ~285 μs per operation with 321 allocations (49KB)

## Adding New Tests

### For New Commands
1. Add command structure tests in `cmds/commands_test.go`
2. Add integration tests in `integration_test.go`
3. Update test documentation

### For New Features
1. Add unit tests for the functionality
2. Add integration tests if CLI-facing
3. Add benchmarks for performance-critical code

### Test Naming Convention
- `Test<FunctionName>` for unit tests
- `TestCLI<Feature>` for integration tests
- `Benchmark<Operation>` for benchmarks

## Troubleshooting Tests

### Common Issues
1. **Test isolation**: Make sure tests don't depend on shared state
2. **File permissions**: Tests create temporary files - ensure write permissions
3. **Path issues**: Tests use absolute paths and temporary directories
4. **Timing issues**: Some timestamp tests may be sensitive to system speed

### Debug Mode
```powershell
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v -run TestSpecificFunction ./...

# Enable race detection
go test -v -race ./...
```
