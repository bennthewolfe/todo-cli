package commands

import (
	"context"
	"fmt"

	"github.com/bennthewolfe/todo-cli/config"
	"github.com/urfave/cli/v3"
)

// NewVersionCommand creates a new version command for urfave/cli
func NewVersionCommand() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Usage:   "Display the version of the application",
		Aliases: []string{"v"},
		Action: func(ctx context.Context, c *cli.Command) error {
			fmt.Println("TODO CLI Version:", config.Version)
			fmt.Println("Release Date:", config.ReleaseDate)
			return nil
		},
	}
}

// Legacy command struct for backward compatibility
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
