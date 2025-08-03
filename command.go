package main

import (
	"flag"
	"fmt"
	"os"
)

type CmdFlag struct {
	Add    string `flag:"add" help:"Add a new todo item"`
	Delete int    `flag:"delete" help:"Delete a todo item by index"`
	Edit   bool   `flag:"edit" help:"Edit a todo item by index and add a new title"`
	ID     int    `flag:"id" help:"ID of the todo item to act on"`
	Title  string `flag:"title" help:"Title for the todo item"`
	Toggle int    `flag:"toggle" help:"Toggle completion status of a todo item by index"`
	List   bool   `flag:"list" help:"List todo items"`
	Format string `flag:"format" help:"Output format for listing todo items (table, json, pretty, none)"`
	Debug  bool   `flag:"debug" help:"Enable debug mode"`
}

func NewCmdFlag() *CmdFlag {
	cf := &CmdFlag{}

	flag.StringVar(&cf.Add, "add", "", "Add a new todo item")
	flag.IntVar(&cf.Delete, "delete", -1, "Delete a todo item by index")
	flag.BoolVar(&cf.Edit, "edit", false, "Edit a todo item by index and add a new title")
	flag.IntVar(&cf.ID, "id", -1, "ID of the todo item to act on")
	flag.StringVar(&cf.Title, "title", "", "Title for the todo item")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Toggle completion status of a todo item by index")
	flag.BoolVar(&cf.List, "list", false, "List todo items")
	flag.StringVar(&cf.Format, "format", "table", "Output format for listing todo items (table, json, pretty, none)")
	flag.BoolVar(&cf.Debug, "debug", false, "Enable debug mode")

	flag.Parse()

	return cf
}

func (cf *CmdFlag) Validate() error {
	if cf.Delete < -1 {
		return fmt.Errorf("invalid delete index: %d", cf.Delete)
	}
	if cf.Edit == true && cf.ID < -1 && cf.Title == "" {
		return fmt.Errorf("invalid edit arguments: %d, %s", cf.ID, cf.Title)
	}
	if cf.Toggle < -1 {
		return fmt.Errorf("invalid toggle index: %d", cf.Toggle)
	}
	if cf.List && cf.Format != "" {
		allowedFormats := []string{"table", "json", "pretty", "none"}
		valid := false
		for _, format := range allowedFormats {
			if cf.Format == format {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid list format: %s", cf.List)
		}
	}

	return nil
}

func (cf *CmdFlag) Execute(todoList *TodoList) error {
	if err := cf.Validate(); err != nil {
		return err
	}

	if len(os.Args) <= 1 {
		cf.List = true
	}

	if cf.Debug {
		fmt.Println("Executing command with flags:")
		fmt.Printf("Add: %s\n", cf.Add)
		fmt.Printf("Delete: %d\n", cf.Delete)
		fmt.Printf("Edit: %d\n", cf.Edit)
		fmt.Printf("Toggle: %d\n", cf.Toggle)
		fmt.Printf("List: %s\n", cf.List)
		fmt.Printf("Debug: %t\n", cf.Debug)
		fmt.Println("\n")
	}

	switch {
	case cf.Add != "":
		todoList.add(cf.Add)
	case cf.Delete >= 0:
		if err := todoList.delete(cf.Delete - 1); err != nil {
			return err
		}
	case cf.Edit:
		if err := todoList.update(cf.ID-1, cf.Title); err != nil {
			return err
		}
	case cf.Toggle >= 0:
		if err := todoList.toggle(cf.Toggle); err != nil {
			return err
		}
	}

	if cf.List {
		todoList.view(cf.Format)
	}

	return nil
}
