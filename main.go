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

		// Global flags
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug mode",
			},
		},

		// Default action when no command is specified
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Bool("debug") {
				fmt.Println("DEBUG: Debug mode enabled")
				fmt.Printf("DEBUG: Args: %v\n", c.Args().Slice())
			}
			// Default to list command with table format
			// Initialize todo list directly
			todoList := &commands.TodoList{}
			storage := commands.NewStorage[commands.TodoList](".todos.json")

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
			commands.NewDeleteCommand(),
			commands.NewEditCommand(),
			commands.NewListCommand(),
			commands.NewToggleCommand(),
			commands.NewVersionCommand(),
			commands.NewHelpCommand(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
