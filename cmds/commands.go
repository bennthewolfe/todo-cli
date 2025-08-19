package commands

import (
	"fmt"
)

var globalRegistry *CommandRegistry

// TodoListInterface defines the interface for TodoList operations that commands can use
type TodoListInterface interface {
	Add(task string) error
	Delete(index int) error
	Update(index int, task string) error
	Toggle(index int) error
	View(format string)
	FilterIncomplete()
}

// Command interface that all commands must implement
type Command interface {
	Name() string
	Description() string
	Usage() string
	Execute(args []string, todoList TodoListInterface) error
}

// CommandRegistry holds all available commands
type CommandRegistry struct {
	commands map[string]Command
}

func RegisterCommand(cmd Command) {
	if globalRegistry == nil {
		globalRegistry = &CommandRegistry{commands: make(map[string]Command)}
	}
	globalRegistry.Register(cmd)
}

func GetRegistry() *CommandRegistry {
	return globalRegistry
}

// Register adds a command to the registry
func (cr *CommandRegistry) Register(cmd Command) {
	cr.commands[cmd.Name()] = cmd
}

// Execute runs the specified command with the given arguments
func (cr *CommandRegistry) Execute(commandName string, args []string, todoList TodoListInterface) error {
	cmd, exists := cr.commands[commandName]
	if !exists {
		return fmt.Errorf("unknown command: %s", commandName)
	}

	return cmd.Execute(args, todoList)
}

// GetCommand returns a command by name
func (cr *CommandRegistry) GetCommand(name string) (Command, bool) {
	cmd, exists := cr.commands[name]
	return cmd, exists
}

// ListCommands returns all available commands
func (cr *CommandRegistry) ListCommands() []Command {
	var commands []Command
	for _, cmd := range cr.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// ShowHelp displays help for all commands or a specific command
func (cr *CommandRegistry) ShowHelp(commandName string, version, releaseDate string) {
	if commandName == "" {
		fmt.Println("Todo is a simple command-line interface for managing todo items.")
		fmt.Println("\nATTRIBUTION:")
		fmt.Println("  This project is inspired by the tutorial from https://codingwithpatrik.dev/posts/how-to-build-a-cli-todo-app-in-go/.")
		fmt.Println("\nVERSION:")
		fmt.Println("  todo version", version, "--", "(", releaseDate, ")")
		fmt.Println("\nUSAGE:")
		fmt.Println("  todo <COMMAND> [ARGUMENTS] [--FLAGS]")
		fmt.Println("\nCOMMANDS:")

		for _, cmd := range cr.ListCommands() {
			fmt.Printf("  %-10s %s\n", cmd.Name(), cmd.Description())
		}

		fmt.Println("\nUse 'todo help <COMMAND>' for more information about a command.")
		fmt.Println("\nEXAMPLES:")
		fmt.Println("  todo add \"Buy groceries\"")
		fmt.Println("  todo delete 2")
		fmt.Println("  todo edit 1 \"Read a book\"")
		fmt.Println("  todo toggle 1")
		fmt.Println("  todo list --format json")
		fmt.Println("  todo list --format json | jq '[.[] | select(.completed == false)]'")
		fmt.Println("\nGLOBAL OPTIONS:")
		fmt.Println("  --help, -h      Show help information")
		fmt.Println("  --version, -v   Show version information")
		return
	}

	cmd, exists := cr.GetCommand(commandName)
	if !exists {
		fmt.Printf("Unknown command: %s\n", commandName)
		return
	}

	fmt.Printf("Command: %s\n", cmd.Name())
	fmt.Printf("Description: %s\n", cmd.Description())
	fmt.Printf("Usage: %s\n", cmd.Usage())
}
