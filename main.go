package main

import (
	"fmt"
	"os"

	commands "github.com/bennthewolfe/todo-cli/cmds"
)

func main() {
	// Initialize todo list and storage
	todoList := TodoList{}

	storage := NewStorage[TodoList](".todos.json")
	loadedList, err := storage.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading todos: %v\n", err)
		os.Exit(1)
	}
	todoList = loadedList

	// Create command registry
	registry := commands.GetRegistry()

	// Parse command line arguments
	args := os.Args[1:]

	// If no arguments provided, default to list command
	if len(args) == 0 {
		args = []string{"list"}
	}

	// Handle --help and -h flags by converting them to help command
	if args[0] == "--help" || args[0] == "-h" {
		if len(args) > 1 {
			args = []string{"help", args[1]}
		} else {
			args = []string{"help"}
		}
	}

	// Handle --version and -v flags by converting them to version command
	if args[0] == "--version" || args[0] == "-v" {
		args = []string{"version"}
	}

	// Extract command name and arguments
	commandName := args[0]
	commandArgs := args[1:]

	// Execute the command
	if err := registry.Execute(commandName, commandArgs, &todoList); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Save the updated todo list
	if err := storage.Save(todoList); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving todos: %v\n", err)
		os.Exit(1)
	}
}
