package main

import (
	"context"
	"fmt"
	"os"

	commands "github.com/bennthewolfe/todo-cli/cmds"
	"github.com/bennthewolfe/todo-cli/config"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:    "Todo CLI",
		Usage:   "A simple command-line interface for managing todo items",
		Version: config.Version,
		Description: "Todo CLI is a command-line application for managing a to-do list. " +
			"It allows users to add, view, and manage tasks efficiently. " +
			"This project is inspired by the tutorial from https://codingwithpatrik.dev/posts/how-to-build-a-cli-todo-app-in-go/.",
		UsageText: "todo [global options] command [command options] [arguments...]",
		// Global flags
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug mode",
			},
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "Use global todo storage in user's home directory (~/.todo/todos.json)",
			},
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List all todo items (overrides other commands)",
			},
			&cli.BoolFlag{
				Name:    "archive",
				Aliases: []string{"a"},
				Usage:   "Work with archive files instead of main todo list (only list and delete commands supported)",
			},
		},

		// Default action when no command is specified
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Bool("debug") {
				fmt.Println("DEBUG: Debug mode enabled")
				fmt.Printf("DEBUG: Args: %v\n", c.Args().Slice())
				fmt.Printf("DEBUG: Global flag: %v\n", c.Bool("global"))
				fmt.Printf("DEBUG: List flag: %v\n", c.Bool("list"))
			}

			// If --list flag is set, show the list regardless of other arguments
			if c.Bool("list") {
				if c.Bool("debug") {
					fmt.Println("DEBUG: --list flag detected, showing todo list")
				}
			}

			// Get the appropriate storage path
			storagePath, err := commands.GetStoragePath(c.Bool("global"))
			if err != nil {
				return cli.Exit(fmt.Sprintf("error getting storage path: %v", err), 2)
			}

			// Default to list command with table format (or when --list flag is used)
			// Initialize todo list directly
			todoList := &commands.TodoList{}
			storage := commands.NewStorage[commands.TodoList](storagePath)

			loadedList, err := storage.Load()
			if err != nil {
				return cli.Exit(fmt.Sprintf("error loading todos: %v", err), 2)
			}

			*todoList = loadedList
			todoList.View("table")
			return nil
		},

		Commands: []*cli.Command{
			commands.NewAddCommand(),
			commands.NewArchiveCommand(),
			commands.NewCleanupCommand(),
			commands.NewDeleteCommand(),
			commands.NewEditCommand(),
			commands.NewListCommand(),
			commands.NewToggleCommand(),
			commands.NewVersionCommand(),
			// Removed NewHelpCommand() - using urfave/cli built-in help instead
		},
	}

	// Append examples to global help
	cli.RootCommandHelpTemplate = fmt.Sprintf(`%s
EXAMPLES:
	todo add "Buy groceries"
	todo delete 2
	todo edit 1 "Read a book"
	todo toggle 1
	todo list --format json
	todo list --format json | jq '[.[] | select(.completed == false)]'
	`, cli.RootCommandHelpTemplate)

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
