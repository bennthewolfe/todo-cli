package commands

import (
	"fmt"
	"strconv"
)

// ToggleCommand handles toggling todo completion status
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
