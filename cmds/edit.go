package commands

import (
	"fmt"
	"strconv"
)

// EditCommand handles editing todo items
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
