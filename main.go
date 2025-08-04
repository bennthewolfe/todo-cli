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

	// Parse flags
	listFlag, args := ParseFlags()

	// If no arguments provided, default to list command
	if len(args) == 0 {
		args = []string{"list"}
	}

	// Extract command name and arguments
	commandName := args[0]
	commandArgs := args[1:]

	// Create command registry
	registry := commands.GetRegistry()

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

	// Check the --list flag and display the todo list if set
	if listFlag {
		todoList.View("table")
	}
}
