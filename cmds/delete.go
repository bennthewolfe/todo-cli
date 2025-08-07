package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"
)

// NewDeleteCommand creates a new delete command for urfave/cli
func NewDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete a todo item by ID",
		Aliases:   []string{"del", "rm"},
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

			// Delete the item
			if err := todoList.Delete(id - 1); err != nil { // Convert to 0-based index
				return cli.Exit(fmt.Sprintf("failed to delete task: %v", err), 1)
			}

			// Save the updated todo list
			if err := storage.Save(*todoList); err != nil {
				return cli.Exit(fmt.Sprintf("error saving todos: %v", err), 2)
			}

			fmt.Printf("Deleted todo item with ID: %d\n", id)
			return nil
		},
	}
}

// Legacy command struct for backward compatibility
type DeleteCommand struct{}

func init() {
	RegisterCommand(&DeleteCommand{})
}

func (c *DeleteCommand) Name() string {
	return "delete"
}

func (c *DeleteCommand) Description() string {
	return "Delete a todo item by ID"
}

func (c *DeleteCommand) Usage() string {
	return "todo-cli delete <id>"
}

func (c *DeleteCommand) Execute(args []string, todoList TodoListInterface) error {
	if len(args) != 1 {
		return fmt.Errorf("exactly one ID is required\nUSAGE: %s", c.Usage())
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid ID: %s must be a number", args[0])
	}

	if id <= 0 {
		return fmt.Errorf("ID must be greater than 0")
	}

	return todoList.Delete(id - 1) // Convert to 0-based index
}
