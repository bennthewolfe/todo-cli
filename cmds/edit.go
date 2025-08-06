package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
)

// NewEditCommand creates a new edit command for urfave/cli
func NewEditCommand() *cli.Command {
	return &cli.Command{
		Name:      "edit",
		Usage:     "Edit a todo item by ID",
		Aliases:   []string{"e"},
		ArgsUsage: "<id> <new_task>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 2 {
				return fmt.Errorf("ID and new task description are required\nUsage: todo-cli edit <id> <new_task>")
			}

			id, err := strconv.Atoi(c.Args().First())
			if err != nil {
				return fmt.Errorf("invalid ID: %s must be a number", c.Args().First())
			}

			if id <= 0 {
				return fmt.Errorf("ID must be greater than 0")
			}

			// Join all arguments after the ID as the new task
			newTask := strings.Join(c.Args().Slice()[1:], " ")

			// Initialize todo list and storage
			todoList, storage, err := initializeTodoList()
			if err != nil {
				return err
			}

			// Update the item
			if err := todoList.Update(id-1, newTask); err != nil { // Convert to 0-based index
				return err
			}

			// Save the updated todo list
			if err := storage.Save(*todoList); err != nil {
				return fmt.Errorf("error saving todos: %w", err)
			}

			fmt.Printf("Updated todo item %d: %s\n", id, newTask)
			return nil
		},
	}
}

// Legacy command struct for backward compatibility
type EditCommand struct{}

func init() {
	RegisterCommand(&EditCommand{})
}

func (c *EditCommand) Name() string {
	return "edit"
}

func (c *EditCommand) Description() string {
	return "Edit a todo item by ID"
}

func (c *EditCommand) Usage() string {
	return "todo-cli edit <id> <new_task>"
}

func (c *EditCommand) Execute(args []string, todoList TodoListInterface) error {
	if len(args) < 2 {
		return fmt.Errorf("ID and new task description are required\nUSAGE: %s", c.Usage())
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid ID: %s must be a number", args[0])
	}

	if id <= 0 {
		return fmt.Errorf("ID must be greater than 0")
	}

	newTask := args[1]
	return todoList.Update(id-1, newTask) // Convert to 0-based index
}
