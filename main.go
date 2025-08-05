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

	args := os.Args[1:]

	// Parse flags
	commandFlags, additionalArgs := NewCmdFlag()

	if commandFlags.Debug {
		fmt.Printf("Args: %v\n", args)
		fmt.Printf("Additional Args: %v\n", additionalArgs)
		fmt.Printf("Flags: %+v\n", commandFlags)
	}

	// If no arguments provided, default to list command
	if len(args) == 0 {
		args = []string{"list"}
	}

	// Create command registry
	registry := commands.GetRegistry()

	// Execute the command
	if err := registry.Execute(args[0], args[1:], &todoList); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Save the updated todo list
	if err := storage.Save(todoList); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving todos: %v\n", err)
		os.Exit(1)
	}

	// Check the --list flag and display the todo list if set
	if commandFlags.List {
		todoList.View("table")
	}
}
