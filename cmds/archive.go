package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v3"
)

// NewArchiveCommand creates a new archive command for urfave/cli
func NewArchiveCommand() *cli.Command {
	return &cli.Command{
		Name:      "archive",
		Usage:     "Archive a todo item by ID, or archive all items if no ID provided",
		Aliases:   []string{"ar"},
		ArgsUsage: "[id]",
		Description: "Archive todo items to move them out of the active list. " +
			"If no ID is provided, all items will be archived after confirmation. " +
			"Use --force to skip confirmation prompt (required in some terminals).",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Skip confirmation prompt when archiving all items",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			// Validate archive flag usage
			if err := ValidateArchiveFlagUsage(c, "archive"); err != nil {
				return cli.Exit(err.Error(), 1)
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

			// Handle different argument cases
			if c.Args().Len() == 0 {
				// Archive all items (check for empty list here)
				if len(*todoList) == 0 {
					fmt.Println("No todos found to archive.")
					return nil
				}
				return archiveAllItems(c, todoList, archiveList, storage, archiveStorage)
			} else if c.Args().Len() == 1 {
				// Archive specific item by ID (don't check for empty here, let ID validation handle it)
				return archiveItemByID(c, todoList, archiveList, storage, archiveStorage)
			} else {
				return cli.Exit("too many arguments: provide either no arguments to archive all, or one ID to archive a specific item", 1)
			}
		},
	}
}

// archiveAllItems handles archiving all todo items
func archiveAllItems(c *cli.Command, todoList *TodoList, archiveList *TodoList, storage *Storage[TodoList], archiveStorage *Storage[TodoList]) error {
	itemCount := len(*todoList)

	// Show confirmation unless --force flag is used
	if !c.Bool("force") {
		fmt.Printf("Found %d item(s) to archive:\n", itemCount)
		for i, item := range *todoList {
			statusStr := ""
			if item.Completed {
				statusStr = " ✅"
			}
			fmt.Printf("  %d. %s%s\n", i+1, item.Task, statusStr)
		}

		fmt.Print("\nAre you sure you want to archive all ", itemCount, " item(s)? (y/N): ")
		
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil {
			fmt.Printf("\nNote: Interactive confirmation not supported in this terminal.\n")
			fmt.Printf("Use --force to archive without confirmation: todo archive --force\n")
			return cli.Exit("", 1)
		}
		
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Archive cancelled.")
			return nil
		}
	}

	// Archive all items
	for _, item := range *todoList {
		// Add to archive
		if err := archiveList.Add(item.Task); err != nil {
			return cli.Exit(fmt.Sprintf("failed to add item to archive: %v", err), 1)
		}

		// Update the archived item to match the original (preserve timestamps and completion status)
		archiveIndex := len(*archiveList) - 1
		(*archiveList)[archiveIndex] = item
	}

	// Clear the main todo list
	*todoList = TodoList{}

	// Save both lists
	if err := storage.Save(*todoList); err != nil {
		return cli.Exit(fmt.Sprintf("error saving todos: %v", err), 2)
	}

	if err := archiveStorage.Save(*archiveList); err != nil {
		return cli.Exit(fmt.Sprintf("error saving archive: %v", err), 2)
	}

	fmt.Printf("Successfully archived %d item(s).\n", itemCount)

	// Check if --list flag is set and execute list command after archive
	if CheckAndExecuteListFlag(c) {
		if err := ExecuteListCommand(c); err != nil {
			return cli.Exit(fmt.Sprintf("error executing list: %v", err), 2)
		}
	}

	return nil
}

// archiveItemByID handles archiving a specific item by ID
func archiveItemByID(c *cli.Command, todoList *TodoList, archiveList *TodoList, storage *Storage[TodoList], archiveStorage *Storage[TodoList]) error {
	id, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return cli.Exit(fmt.Sprintf("invalid ID: %s must be a number", c.Args().First()), 1)
	}

	if id <= 0 {
		return cli.Exit("ID must be greater than 0", 1)
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

	// Check if --list flag is set and execute list command after archive
	if CheckAndExecuteListFlag(c) {
		if err := ExecuteListCommand(c); err != nil {
			return cli.Exit(fmt.Sprintf("error executing list: %v", err), 2)
		}
	}

	return nil
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
