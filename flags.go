package main

import (
	"flag"
	"fmt"
)

func init() {
	flag.Usage = func() {
		fmt.Println("Todo CLI Application:")
		fmt.Println("A simple command-line interface for managing todo items.")
		fmt.Println("This is based on the tutorial at https://codingwithpatrik.dev/posts/how-to-build-a-cli-todo-app-in-go/")
		fmt.Println("\nUsage: todo-cli [OPTIONS]")
		fmt.Println("\nOptions:")
		fmt.Println("  --add, -a <task>       Add a new todo item")
		fmt.Println("  --delete, -d <index>   Delete a todo item by index")
		fmt.Println("  --edit, -e             Edit a todo item (requires --id and --task)")
		fmt.Println("  --id, -i <id>          Specify the ID of the todo item to act on")
		fmt.Println("  --task, -t <task>      Specify the new task for the todo item")
		fmt.Println("  --toggle, -c <index>   Toggle completion status of a todo item by index")
		fmt.Println("  --list, -l             List all todo items")
		fmt.Println("  --format, -f <format>      Specify output format (table, json, pretty, none)")
		fmt.Println("  --debug                Enable debug mode")
		fmt.Println("  --help                 Show this help message")
		fmt.Println("\nExamples:")
		fmt.Println("  todo-cli --add \"Buy groceries\"")
		fmt.Println("  todo-cli --delete 2")
		fmt.Println("  todo-cli --edit --id 1 --task \"Read a book\"")
		fmt.Println("  todo-cli --list --format json | jq '[.[] | select(.completed == false)]'")
	}
}

type CmdFlag struct {
	Add    string `flag:"add" help:"Add a new todo item"`
	Delete int    `flag:"delete" help:"Delete a todo item by index"`
	Edit   bool   `flag:"edit" help:"Edit a todo item by index and add a new task"`
	ID     int    `flag:"id" help:"ID of the todo item to act on"`
	Task   string `flag:"task" help:"Task for the todo item"`
	Toggle int    `flag:"toggle" help:"Toggle completion status of a todo item by index"`
	List   bool   `flag:"list" help:"List todo items"`
	Format string `flag:"format" help:"Output format for listing todo items (table, json, pretty, none)"`
	Debug  bool   `flag:"debug" help:"Enable debug mode"`
}

func NewCmdFlag() *CmdFlag {
	cf := &CmdFlag{}

	flag.StringVar(&cf.Add, "add", "", "Add a new todo item")
	flag.StringVar(&cf.Add, "a", "", "Alias for --add")

	flag.IntVar(&cf.Delete, "delete", -1, "Delete a todo item by index")
	flag.IntVar(&cf.Delete, "del", -1, "Alias for --delete")
	flag.IntVar(&cf.Delete, "d", -1, "Alias for --delete")

	flag.BoolVar(&cf.Edit, "edit", false, "Edit a todo item by index and add a new task")
	flag.BoolVar(&cf.Edit, "e", false, "Alias for --edit")

	flag.IntVar(&cf.ID, "id", -1, "ID of the todo item to act on")
	flag.IntVar(&cf.ID, "i", -1, "Alias for --id")

	flag.StringVar(&cf.Task, "task", "", "Task for the todo item")
	flag.StringVar(&cf.Task, "t", "", "Alias for --task")

	flag.IntVar(&cf.Toggle, "toggle", -1, "Toggle completion status of a todo item by index")
	flag.IntVar(&cf.Toggle, "c", -1, "Alias for --toggle (complete)")

	flag.BoolVar(&cf.List, "list", false, "List todo items")
	flag.BoolVar(&cf.List, "l", false, "Alias for --list")

	flag.StringVar(&cf.Format, "format", "table", "Output format for listing todo items (table, json, pretty, none)")
	flag.StringVar(&cf.Format, "f", "table", "Alias for --format")

	flag.BoolVar(&cf.Debug, "debug", false, "Enable debug mode")

	flag.Parse()

	return cf
}

// ParseFlags parses global flags and returns the remaining arguments and the value of the --list flag.
func ParseFlags() (bool, []string) {
	listFlag := flag.Bool("list", false, "Display the todo list after executing the command")
	flag.Parse()
	return *listFlag, flag.Args()
}
