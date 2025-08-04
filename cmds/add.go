package commands

import (
	"fmt"
	"strings"
)

// AddCommand handles adding new todo items
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
