package commands

import (
	"fmt"
	"strings"
)

// ListCommand handles listing todo items
type ListCommand struct{}

func init() {
	RegisterCommand(&ListCommand{})
}

func (c *ListCommand) Name() string {
	return "list"
}

func (c *ListCommand) Description() string {
	return "List all todo items"
}

func (c *ListCommand) Usage() string {
	return "todo-cli list [--format <format>]\n  Formats: table, json, pretty"
}

func (c *ListCommand) Execute(args []string, todoList TodoListInterface) error {
	format := "table" // default format

	// Parse format flag if provided
	for i, arg := range args {
		if arg == "--format" || arg == "-f" {
			if i+1 < len(args) {
				format = args[i+1]
			} else {
				return fmt.Errorf("format value is required after --format flag")
			}
			break
		}
	}

	// Validate format
	allowedFormats := []string{"table", "json", "pretty", "none"}
	valid := false
	for _, allowedFormat := range allowedFormats {
		if format == allowedFormat {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid format: %s. Allowed formats: %s", format, strings.Join(allowedFormats, ", "))
	}

	todoList.View(format)
	return nil
}
