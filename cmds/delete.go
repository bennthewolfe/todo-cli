package commands

import (
	"fmt"
	"strconv"
)

// DeleteCommand handles deleting todo items
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
