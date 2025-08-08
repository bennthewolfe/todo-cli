package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

// NewCleanupCommand creates a new cleanup command for urfave/cli
func NewCleanupCommand() *cli.Command {
	return &cli.Command{
		Name:      "cleanup",
		Usage:     "Archive or delete all completed todo items",
		Aliases:   []string{"clean"},
		ArgsUsage: " ",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Skip confirmation prompt",
			},
			&cli.BoolFlag{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "Delete completed items instead of archiving them",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			// Determine the action based on flags
			isDelete := c.Bool("delete")

			// Get the appropriate storage paths based on global flag
			storagePath, err := GetStoragePath(c.Bool("global"))
			if err != nil {
				return cli.Exit(fmt.Sprintf("error getting storage path: %v", err), 2)
			}

			var archivePath string
			var archiveList *TodoList
			var archiveStorage *Storage[TodoList]

			// Only initialize archive if we're archiving (not deleting)
			if !isDelete {
				archivePath, err = GetArchivePath(c.Bool("global"))
				if err != nil {
					return cli.Exit(fmt.Sprintf("error getting archive path: %v", err), 2)
				}

				// Initialize archive list and storage
				archiveList, archiveStorage, err = initializeTodoListWithPath(archivePath)
				if err != nil {
					return cli.Exit(fmt.Sprintf("failed to initialize archive list: %v", err), 2)
				}
			}

			// Initialize todo list and storage
			todoList, storage, err := initializeTodoListWithPath(storagePath)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to initialize todo list: %v", err), 2)
			}

			// Find all completed items
			var completedItems []Todo
			var remainingItems TodoList

			for _, item := range *todoList {
				if item.Completed {
					completedItems = append(completedItems, item)
				} else {
					remainingItems = append(remainingItems, item)
				}
			}

			// Check if there are any completed items to process
			if len(completedItems) == 0 {
				if isDelete {
					fmt.Println("No completed items found to delete.")
				} else {
					fmt.Println("No completed items found to archive.")
				}
				return nil
			}

			// Show confirmation unless --force flag is used
			if !c.Bool("force") {
				if isDelete {
					fmt.Printf("Found %d completed item(s) to delete:\n", len(completedItems))
				} else {
					fmt.Printf("Found %d completed item(s) to archive:\n", len(completedItems))
				}
				for i, item := range completedItems {
					fmt.Printf("  %d. %s\n", i+1, item.Task)
				}

				if isDelete {
					fmt.Printf("\nAre you sure you want to delete these %d completed item(s)? (y/N): ", len(completedItems))
				} else {
					fmt.Printf("\nAre you sure you want to archive these %d completed item(s)? (y/N): ", len(completedItems))
				}

				reader := bufio.NewReader(os.Stdin)
				response, err := reader.ReadString('\n')
				if err != nil {
					return cli.Exit(fmt.Sprintf("error reading confirmation: %v", err), 2)
				}

				response = strings.TrimSpace(strings.ToLower(response))
				if response != "y" && response != "yes" {
					if isDelete {
						fmt.Println("Delete cancelled.")
					} else {
						fmt.Println("Cleanup cancelled.")
					}
					return nil
				}
			}

			if isDelete {
				// Delete mode: just remove completed items (don't archive)
				*todoList = remainingItems
			} else {
				// Archive mode: add completed items to archive, then remove from main list
				for _, item := range completedItems {
					if err := archiveList.Add(item.Task); err != nil {
						return cli.Exit(fmt.Sprintf("failed to add item to archive: %v", err), 1)
					}

					// Update the archived item to match the original (preserve timestamps and completion status)
					archiveIndex := len(*archiveList) - 1
					(*archiveList)[archiveIndex] = item
				}

				// Update the main todo list to only contain non-completed items
				*todoList = remainingItems

				// Save the archive
				if err := archiveStorage.Save(*archiveList); err != nil {
					return cli.Exit(fmt.Sprintf("error saving archive: %v", err), 2)
				}
			}

			// Save the main todo list (always needed)
			if err := storage.Save(*todoList); err != nil {
				return cli.Exit(fmt.Sprintf("error saving todos: %v", err), 2)
			}

			// Show success message
			if isDelete {
				fmt.Printf("Successfully deleted %d completed item(s).\n", len(completedItems))
			} else {
				fmt.Printf("Successfully archived %d completed item(s).\n", len(completedItems))
			}

			// Check if --list flag is set and execute list command after cleanup
			if CheckAndExecuteListFlag(c) {
				if err := ExecuteListCommand(c); err != nil {
					return cli.Exit(fmt.Sprintf("error executing list: %v", err), 2)
				}
			}

			return nil
		},
	}
}

// Legacy command struct for backward compatibility
type CleanupCommand struct{}

func init() {
	RegisterCommand(&CleanupCommand{})
}

func (c *CleanupCommand) Name() string {
	return "cleanup"
}

func (c *CleanupCommand) Description() string {
	return "Archive all completed todo items"
}

func (c *CleanupCommand) Usage() string {
	return "todo-cli cleanup [--force]"
}

func (c *CleanupCommand) Execute(args []string, todoList TodoListInterface) error {
	// Note: Legacy interface doesn't support cleanup functionality
	// This would need to be implemented if using the legacy system
	return fmt.Errorf("cleanup functionality not supported in legacy interface")
}
