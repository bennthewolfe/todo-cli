Todo is a simple command-line interface for managing todo items.

ATTRIBUTION:
  This project is inspired by the tutorial from https://codingwithpatrik.dev/posts/how-to-build-a-cli-todo-app-in-go/.

VERSION:
  todo-cli version 2.0.0 -- ( 2025-08-06 )

USAGE:
  todo-cli <COMMAND> [ARGUMENTS] [--FLAGS]

COMMANDS:
  add        Add a new todo item
  delete     Delete a todo item by ID
  edit       Edit a todo item by ID
  list       List all todo items
  toggle     Toggle completion status of a todo item by ID
  version    Display the version of the application
  help       Show help information for commands

EXAMPLES:
  todo-cli add "Buy groceries"
  todo-cli delete 2
  todo-cli edit 1 "Read a book"
  todo-cli toggle 1
  todo-cli list --format json
  todo-cli list --format json | jq '[.[] | select(.completed == false)]'

GLOBAL OPTIONS:
  --help, -h      Show help information
  --debug         Enable debug mode
PS D:\Projects\todo-cli> go run . --help
NAME:
   Todo CLI - A simple command-line interface for managing todo items

USAGE:
   Todo CLI [global options] [command [command options]]

VERSION:
   2.0.0

DESCRIPTION:
   Todo CLI is a command-line application for managing a to-do list. It allows users to add, view, and manage tasks efficiently. This project is inspired by the tutorial from https://codingwithpatrik.dev/posts/how-to-build-a-cli-todo-app-in-go/.

COMMANDS:
   add, a               Add a new todo item
   delete, del, rm      Delete a todo item by ID
   edit, e              Edit a todo item by ID
   list, l, ls          List all todo items
   toggle, t, complete  Toggle completion status of a todo item by ID
   version, v           Display the version of the application
   help, h              Show help information for commands

GLOBAL OPTIONS:
   --debug        Enable debug mode (default: false)
   --help, -h     show help
   --version, -v  print the version
PS D:\Projects\todo-cli> go run . help
NAME:
   Todo CLI - A simple command-line interface for managing todo items

USAGE:
   Todo CLI [global options] [command [command options]]

VERSION:
   2.0.0

DESCRIPTION:
   Todo CLI is a command-line application for managing a to-do list. It allows users to add, view, and manage tasks efficiently. This project is inspired by the tutorial from https://codingwithpatrik.dev/posts/how-to-build-a-cli-todo-app-in-go/.

COMMANDS:
   add, a               Add a new todo item
   delete, del, rm      Delete a todo item by ID
   edit, e              Edit a todo item by ID
   list, l, ls          List all todo items
   toggle, t, complete  Toggle completion status of a todo item by ID
   version, v           Display the version of the application
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        Enable debug mode (default: false)
   --help, -h     show help
   --version, -v  print the version