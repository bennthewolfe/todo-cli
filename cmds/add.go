package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

// NewAddCommand creates a new add command for urfave/cli
func NewAddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a new todo item",
		Aliases:   []string{"a"},
		ArgsUsage: "<task>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				return fmt.Errorf("task description is required")
			}

			// Join all arguments as the task description
			task := strings.Join(c.Args().Slice(), " ")

			// Initialize todo list and storage
			todoList, storage, err := initializeTodoList()
			if err != nil {
				return err
			}

			// Add the task
			if err := todoList.Add(task); err != nil {
				return err
			}

			// Save the updated todo list
			if err := storage.Save(*todoList); err != nil {
				return fmt.Errorf("error saving todos: %w", err)
			}

			fmt.Printf("Added task: %s\n", task)
			return nil
		},
	}
}

// Legacy command struct for backward compatibility
type AddCommand struct{}

func init() {
	RegisterCommand(&AddCommand{})
}

func (c *AddCommand) Name() string {
	return "add"
}

func (c *AddCommand) Description() string {
	return "Add a new todo item"
}

func (c *AddCommand) Usage() string {
	return "todo-cli add <task>"
}

func (c *AddCommand) Execute(args []string, todoList TodoListInterface) error {
	if len(args) == 0 {
		return fmt.Errorf("task description is required\nUSAGE: %s", c.Usage())
	}

	task := strings.Join(args, " ")
	return todoList.Add(task)
}
