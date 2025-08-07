# PowerShell build script for todo-cli project

param(
    [Parameter(Position=0)]
    [string]$Target = "help"
)

function All {
    Write-Host "Running default workflow: test then build..." -ForegroundColor Green
    Test
    Build
}

function Build {
    Write-Host "Building application..." -ForegroundColor Green
    go build -o todo-cli.exe -v .
}

function Clean {
    Write-Host "Cleaning build files..." -ForegroundColor Green
    go clean
    Remove-Item -Path "todo-cli.exe" -ErrorAction SilentlyContinue
    Remove-Item -Path "coverage.out" -ErrorAction SilentlyContinue
    Remove-Item -Path "coverage.html" -ErrorAction SilentlyContinue
}

function Test {
    Write-Host "Running all tests..." -ForegroundColor Green
    go test -v ./...
}

function Coverage {
    Write-Host "Running tests with coverage..." -ForegroundColor Green
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    Write-Host "Coverage report generated: coverage.html" -ForegroundColor Yellow
}

function Test-Unit {
    Write-Host "Running unit tests..." -ForegroundColor Green
    go test -v -short ./...
}

function Test-Integration {
    Write-Host "Running integration tests..." -ForegroundColor Green
    go test -v -run TestCLI ./...
}

function Bench {
    Write-Host "Running benchmark tests..." -ForegroundColor Green
    go test -bench=. -benchmem ./...
}

function Test-Race {
    Write-Host "Running tests with race detection..." -ForegroundColor Green
    go test -v -race ./...
}

function Format {
    Write-Host "Formatting code..." -ForegroundColor Green
    go fmt ./...
}

function Fmt {
    Write-Host "Formatting code..." -ForegroundColor Green
    go fmt ./...
}

function Vet {
    Write-Host "Vetting code..." -ForegroundColor Green
    go vet ./...
}

function Lint {
    Write-Host "Linting code..." -ForegroundColor Green
    if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
        golangci-lint run
    } else {
        Write-Warning "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    }
}

function Dependencies {
    Write-Host "Downloading dependencies..." -ForegroundColor Green
    go mod download
    go mod tidy
}

function Check {
    Write-Host "Running all quality checks..." -ForegroundColor Green
    Fmt
    Vet
    Lint
    Test
}

function Install {
    Write-Host "Installing application..." -ForegroundColor Green
    go install .
}

function Show-Help {
    Write-Host "Available targets:" -ForegroundColor Cyan
    Write-Host "  all           - Run default workflow (test then build)"
    Write-Host "  build         - Build the application"
    Write-Host "  clean         - Clean build files"
    Write-Host "  test          - Run all tests"
    Write-Host "  test-unit     - Run unit tests only"
    Write-Host "  test-integration - Run integration tests only"
    Write-Host "  coverage      - Run tests with coverage report"
    Write-Host "  bench         - Run benchmark tests"
    Write-Host "  test-race     - Run tests with race condition detection"
    Write-Host "  lint          - Lint the code"
    Write-Host "  fmt           - Format the code"
    Write-Host "  format        - Format the code (alias for fmt)"
    Write-Host "  vet           - Vet the code"
    Write-Host "  check         - Run all quality checks"
    Write-Host "  deps          - Download dependencies"
    Write-Host "  install       - Install the application"
    Write-Host "  help          - Show this help"
    Write-Host ""
    Write-Host "Usage: .\makefile.ps1 [target]" -ForegroundColor Yellow
    Write-Host "Example: .\makefile.ps1 coverage" -ForegroundColor Yellow
}

# Main execution
switch ($Target.ToLower()) {
    "all" { All }
    "build" { Build }
    "clean" { Clean }
    "test" { Test }
    "coverage" { Coverage }
    "test-unit" { Test-Unit }
    "test-integration" { Test-Integration }
    "bench" { Bench }
    "test-race" { Test-Race }
    "format" { Format }
    "fmt" { Fmt }
    "vet" { Vet }
    "lint" { Lint }
    "deps" { Dependencies }
    "check" { Check }
    "install" { Install }
    "help" { Show-Help }
    default { Show-Help }
}
