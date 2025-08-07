package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

// NewListCommand creates a new list command for urfave/cli
func NewListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List all todo items",
		Aliases: []string{"l", "ls"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Output format (table, json, pretty, none)",
				Value:   "table",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			format := c.String("format")

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
				return cli.Exit(fmt.Sprintf("invalid format: %s. Allowed formats: %s", format, strings.Join(allowedFormats, ", ")), 1)
			}

			// Get the appropriate storage path based on global flag
			storagePath, err := GetStoragePath(c.Bool("global"))
			if err != nil {
				return cli.Exit(fmt.Sprintf("error getting storage path: %v", err), 2)
			}

			// Initialize todo list and storage
			todoList, _, err := initializeTodoListWithPath(storagePath)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to initialize todo list: %v", err), 2)
			}

			todoList.View(format)
			return nil
		},
	}
}

// Legacy command struct for backward compatibility
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
