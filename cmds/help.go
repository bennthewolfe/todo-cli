package commands

// HelpCommand handles displaying help information
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

	// Get version info from the registry's version command if available
	version := "1.1.0"
	releaseDate := "2025-08-04"

	// If no specific command is requested, show general help
	if len(args) == 0 {
		registry.ShowHelp("", version, releaseDate)
		return nil
	}

	// Show help for specific command
	commandName := args[0]
	registry.ShowHelp(commandName, version, releaseDate)
	return nil
}
