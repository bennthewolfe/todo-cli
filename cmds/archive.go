package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"
)

// NewArchiveCommand creates a new archive command for urfave/cli
func NewArchiveCommand() *cli.Command {
	return &cli.Command{
		Name:      "archive",
		Usage:     "Archive a todo item by ID (moves to archive file)",
		Aliases:   []string{"ar"},
		ArgsUsage: "<id>",
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() != 1 {
				return cli.Exit("exactly one ID is required", 1)
			}

			id, err := strconv.Atoi(c.Args().First())
			if err != nil {
				return cli.Exit(fmt.Sprintf("invalid ID: %s must be a number", c.Args().First()), 1)
			}

			if id <= 0 {
				return cli.Exit("ID must be greater than 0", 1)
			}

			// Get the appropriate storage paths based on global flag
			storagePath, err := GetStoragePath(c.Bool("global"))
			if err != nil {
				return cli.Exit(fmt.Sprintf("error getting storage path: %v", err), 2)
			}

			archivePath, err := GetArchivePath(c.Bool("global"))
			if err != nil {
				return cli.Exit(fmt.Sprintf("error getting archive path: %v", err), 2)
			}

			// Initialize todo list and storage
			todoList, storage, err := initializeTodoListWithPath(storagePath)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to initialize todo list: %v", err), 2)
			}

			// Initialize archive list and storage
			archiveList, archiveStorage, err := initializeTodoListWithPath(archivePath)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to initialize archive list: %v", err), 2)
			}

			// Validate the ID exists
			if id-1 < 0 || id-1 >= len(*todoList) {
				return cli.Exit(fmt.Sprintf("invalid ID: %d (valid range: 1-%d)", id, len(*todoList)), 1)
			}

			// Get the item to archive
			todoItem := (*todoList)[id-1]

			// Add to archive
			if err := archiveList.Add(todoItem.Task); err != nil {
				return cli.Exit(fmt.Sprintf("failed to add item to archive: %v", err), 1)
			}

			// Update the archived item to match the original (preserve timestamps and completion status)
			archiveIndex := len(*archiveList) - 1
			(*archiveList)[archiveIndex] = todoItem

			// Remove from main list
			if err := todoList.Delete(id - 1); err != nil {
				return cli.Exit(fmt.Sprintf("failed to remove item from todo list: %v", err), 1)
			}

			// Save both lists
			if err := storage.Save(*todoList); err != nil {
				return cli.Exit(fmt.Sprintf("error saving todos: %v", err), 2)
			}

			if err := archiveStorage.Save(*archiveList); err != nil {
				return cli.Exit(fmt.Sprintf("error saving archive: %v", err), 2)
			}

			fmt.Printf("Archived todo item: %s\n", todoItem.Task)
			return nil
		},
	}
}

// Legacy command struct for backward compatibility
type ArchiveCommand struct{}

func init() {
	RegisterCommand(&ArchiveCommand{})
}

func (c *ArchiveCommand) Name() string {
	return "archive"
}

func (c *ArchiveCommand) Description() string {
	return "Archive a todo item by ID"
}

func (c *ArchiveCommand) Usage() string {
	return "todo-cli archive <id>"
}

func (c *ArchiveCommand) Execute(args []string, todoList TodoListInterface) error {
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

	// Note: Legacy interface doesn't support archive functionality
	// This would need to be implemented if using the legacy system
	return fmt.Errorf("archive functionality not supported in legacy interface")
}
