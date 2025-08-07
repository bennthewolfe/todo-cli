# Todo CLI

## Overview
Todo CLI is a command-line application for managing a to-do list. It allows users to add, view, and manage tasks efficiently. The application is written in Go and uses JSON for data serialization.

This based on the tutorial https://codingwithpatrik.dev/posts/how-to-build-a-cli-todo-app-in-go/

## Features
- Add tasks to the to-do list
- View all tasks
- Mark tasks as completed
- Delete tasks

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
Run the application using the following command:
```bash
go run .

# or if you have jq installed for a prettified result
# go run . | jq
```

## Installation
Add the application to your PATH, and then call it with:
```bash
todo
```

I created a symlink in a folder that is in my PATH with the following powershell:
```pwsh
New-Item -ItemType SymbolicLink -Path <LINK> -Target <ACTUAL SOURCE>
```

### Example
```pwsh
$ todo list --format json | jq
[
  {
    "id": 0,
    "title": "Buy groceries",
    "completed": false,
    "created_at": "2025-08-01T16:05:59-04:00",
    "updated_at": "2025-08-01T16:05:59-04:00"
  },
  {
    "id": 2,
    "title": "Buy groceries",
    "completed": false,
    "created_at": "2025-08-01T16:05:59-04:00",
    "updated_at": "2025-08-01T16:05:59-04:00"
  },
  {
    "id": 0,
    "title": "Walk the dog",
    "completed": false,
    "created_at": "2025-08-01T16:05:59-04:00",
    "updated_at": "2025-08-01T16:05:59-04:00"
  },
  {
    "id": 4,
    "title": "Walk the dog",
    "completed": false,
    "created_at": "2025-08-01T16:05:59-04:00",
    "updated_at": "2025-08-01T16:05:59-04:00"
  }
]
```

## Todos
[See .todos.json](./.todos.json)

## Development
### Prerequisites
- Go 1.24.5 or later

### Running Tests
To run tests, use:
```bash
go test ./...
```