package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"
)

// NewToggleCommand creates a new toggle command for urfave/cli
func NewToggleCommand() *cli.Command {
	return &cli.Command{
		Name:      "toggle",
		Usage:     "Toggle completion status of a todo item by ID",
		Aliases:   []string{"t", "complete"},
		ArgsUsage: "<id>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 1 {
				return cli.Exit("exactly one ID is required", 1)
			}

			id, err := strconv.Atoi(c.Args().First())
			if err != nil {
				return cli.Exit(fmt.Sprintf("invalid ID: %s must be a number", c.Args().First()), 1)
			}

			if id <= 0 {
				return cli.Exit("ID must be greater than 0", 1)
			}

			// Get the appropriate storage path based on global flag
			storagePath, err := GetStoragePath(c.Bool("global"))
			if err != nil {
				return cli.Exit(fmt.Sprintf("error getting storage path: %v", err), 2)
			}

			// Initialize todo list and storage
			todoList, storage, err := initializeTodoListWithPath(storagePath)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to initialize todo list: %v", err), 2)
			}

			// Toggle the item
			if err := todoList.Toggle(id); err != nil { // toggle method expects 1-based index
				return cli.Exit(fmt.Sprintf("failed to toggle task: %v", err), 1)
			}

			// Save the updated todo list
			if err := storage.Save(*todoList); err != nil {
				return cli.Exit(fmt.Sprintf("error saving todos: %v", err), 2)
			}

			fmt.Printf("Toggled completion status for todo item with ID: %d\n", id)
			return nil
		},
	}
}

// Legacy command struct for backward compatibility
type ToggleCommand struct{}

func init() {
	RegisterCommand(&ToggleCommand{})
}

func (c *ToggleCommand) Name() string {
	return "toggle"
}

func (c *ToggleCommand) Description() string {
	return "Toggle completion status of a todo item by ID"
}

func (c *ToggleCommand) Usage() string {
	return "todo-cli toggle <id>"
}

func (c *ToggleCommand) Execute(args []string, todoList TodoListInterface) error {
	if len(args) != 1 {
		return fmt.Errorf("exactly one ID is required\nUsage: %s", c.Usage())
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid ID: %s must be a number", args[0])
	}

	if id <= 0 {
		return fmt.Errorf("ID must be greater than 0")
	}

	return todoList.Toggle(id) // toggle method expects 1-based index
}
