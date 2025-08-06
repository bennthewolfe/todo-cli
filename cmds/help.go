package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/bennthewolfe/todo-cli/config"
)

// NewHelpCommand creates a new help command for urfave/cli
func NewHelpCommand() *cli.Command {
	return &cli.Command{
		Name:      "help",
		Usage:     "Show help information for commands",
		Aliases:   []string{"h"},
		ArgsUsage: "[command]",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() == 0 {
				// Show general help
				showGeneralHelp()
				return nil
			}

			// Show help for a specific command
			commandName := c.Args().First()
			fmt.Printf("Help for command: %s\n", commandName)
			fmt.Printf("(Command-specific help would be shown here)\n")
			return nil
		},
	}
}

func showGeneralHelp() {
	version := config.Version
	releaseDate := config.ReleaseDate

	fmt.Println("Todo is a simple command-line interface for managing todo items.")
	fmt.Println("\nATTRIBUTION:")
	fmt.Println("  This project is inspired by the tutorial from https://codingwithpatrik.dev/posts/how-to-build-a-cli-todo-app-in-go/.")
	fmt.Println("\nVERSION:")
	fmt.Println("  todo-cli version", version, "--", "(", releaseDate, ")")
	fmt.Println("\nUSAGE:")
	fmt.Println("  todo-cli <COMMAND> [ARGUMENTS] [--FLAGS]")
	fmt.Println("\nCOMMANDS:")
	fmt.Println("  add        Add a new todo item")
	fmt.Println("  delete     Delete a todo item by ID")
	fmt.Println("  edit       Edit a todo item by ID")
	fmt.Println("  list       List all todo items")
	fmt.Println("  toggle     Toggle completion status of a todo item by ID")
	fmt.Println("  version    Display the version of the application")
	fmt.Println("  help       Show help information for commands")
	fmt.Println("\nEXAMPLES:")
	fmt.Println("  todo-cli add \"Buy groceries\"")
	fmt.Println("  todo-cli delete 2")
	fmt.Println("  todo-cli edit 1 \"Read a book\"")
	fmt.Println("  todo-cli toggle 1")
	fmt.Println("  todo-cli list --format json")
	fmt.Println("  todo-cli list --format json | jq '[.[] | select(.completed == false)]'")
	fmt.Println("\nGLOBAL OPTIONS:")
	fmt.Println("  --help, -h      Show help information")
	fmt.Println("  --debug         Enable debug mode")
}

// Legacy command struct for backward compatibility
type HelpCommand struct{}

func init() {
	RegisterCommand(&HelpCommand{})
}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Description() string {
	return "Show help information for commands"
}

func (c *HelpCommand) Usage() string {
	return "todo-cli help [COMMAND] or todo-cli [COMMAND] --help"
}

func (c *HelpCommand) Execute(args []string, todoList TodoListInterface) error {
	registry := GetRegistry()

	// Get version info from config
	version := config.Version
	releaseDate := config.ReleaseDate

	// If no specific command is requested, show general help
	if len(args) == 0 {
		fmt.Println("Todo is a simple command-line interface for managing todo items.")
		fmt.Println("\nATTRIBUTION:")
		fmt.Println("  This project is inspired by the tutorial from https://codingwithpatrik.dev/posts/how-to-build-a-cli-todo-app-in-go/.")
		fmt.Println("\nVERSION:")
		fmt.Println("  todo-cli version", version, "--", "(", releaseDate, ")")
		fmt.Println("\nUSAGE:")
		fmt.Println("  todo-cli <COMMAND> [ARGUMENTS] [--FLAGS]")
		fmt.Println("\nCOMMANDS:")

		for _, cmd := range registry.ListCommands() {
			fmt.Printf("  %-10s %s\n", cmd.Name(), cmd.Description())
		}

		fmt.Println("\nUse 'todo-cli help <COMMAND>' for more information about a command.")
		// todo: fmt.Println("\nTOPICS:")
		fmt.Println("\nEXAMPLES:")
		fmt.Println("  todo-cli add \"Buy groceries\"")
		fmt.Println("  todo-cli delete 2")
		fmt.Println("  todo-cli edit 1 \"Read a book\"")
		fmt.Println("  todo-cli toggle 1")
		fmt.Println("  todo-cli list --format json")
		fmt.Println("  todo-cli list --format json | jq '[.[] | select(.completed == false)]'")
		fmt.Println("\nGLOBAL OPTIONS:")
		fmt.Println("  --help, -h      Show help information")
		fmt.Println("  --version, -v   Show version information")
		return nil
	}

	// Show help for a specific command
	commandName := args[0]
	cmd, exists := registry.GetCommand(commandName)
	if !exists {
		fmt.Printf("Unknown command: %s\n", commandName)
		return nil
	}

	fmt.Printf("Command: %s\n", cmd.Name())
	fmt.Printf("Description: %s\n", cmd.Description())
	fmt.Printf("Usage: %s\n", cmd.Usage())
	return nil
}
