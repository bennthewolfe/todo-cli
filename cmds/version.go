package commands

import (
	"fmt"

	"github.com/bennthewolfe/todo-cli/config"
)

// VersionCommand handles displaying the version number
type VersionCommand struct{}

func init() {
	RegisterCommand(&VersionCommand{})
}

func (c *VersionCommand) Name() string {
	return "version"
}

func (c *VersionCommand) Description() string {
	return "Display the version of the application"
}

func (c *VersionCommand) Usage() string {
	return "todo-cli version"
}

func (c *VersionCommand) Execute(args []string, todoList TodoListInterface) error {
	// Use the constants from the config package
	fmt.Println("TODO CLI Version:", config.Version)
	fmt.Println("Release Date:", config.ReleaseDate)
	return nil
}
