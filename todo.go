package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/aquasecurity/table"
	"github.com/liamg/tml"
)

type Todo struct {
	ID          int    `json:"id"`
	Task        string `json:"task"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CompletedAt string `json:"completed_at,omitempty"`
}

// TodoListInterface defines the interface for TodoList operations that commands can use
type TodoListInterface interface {
	Add(task string) error
	Delete(index int) error
	Update(index int, task string) error
	Toggle(index int) error
	View(format string)
}

type TodoList []Todo

// Interface wrapper methods (capitalize first letter to make them public)
func (todoList *TodoList) Add(task string) error {
	return todoList.add(task)
}

func (todoList *TodoList) Delete(index int) error {
	return todoList.delete(index)
}

func (todoList *TodoList) Update(index int, task string) error {
	return todoList.update(index, task)
}

func (todoList *TodoList) Toggle(index int) error {
	return todoList.toggle(index)
}

func (todoList *TodoList) View(format string) {
	todoList.view(format)
}

func (todoList *TodoList) validateIndex(index int) error {
	if index < 0 || index >= len(*todoList) {
		err := errors.New("invalid index")
		return fmt.Errorf(err.Error(), "%d", index)
	}
	return nil
}

func (todoList *TodoList) add(task string) error {
	todo := Todo{
		Task:        task,
		Completed:   false,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
		CompletedAt: "",
	}

	// Add a new Todo item to the list
	todo.ID = len(*todoList) + 1
	todo.CreatedAt = time.Now().Format(time.RFC3339)
	*todoList = append(*todoList, todo)

	return nil
}

func (todoList *TodoList) delete(index int) error {
	t := *todoList

	// Validate the index before attempting to delete
	if err := todoList.validateIndex(index); err != nil {
		return err
	}

	// Remove the Todo item at the specified index
	t = append(t[:index], t[index+1:]...)
	*todoList = t

	return nil
}

func (todoList *TodoList) toggle(index int) error {
	t := *todoList
	index-- // Adjust for 0-based index

	// Validate the index before attempting to toggle
	if err := t.validateIndex(index); err != nil {
		return err
	}

	// Toggle the completion status of the Todo item
	t[index].Completed = !t[index].Completed
	t[index].UpdatedAt = time.Now().Format(time.RFC3339)

	if t[index].Completed {
		t[index].CompletedAt = time.Now().Format(time.RFC3339)
	} else {
		t[index].CompletedAt = ""
	}

	return nil
}

func (todoList *TodoList) update(index int, task string) error {
	t := *todoList

	// Validate the index before attempting to update
	if err := t.validateIndex(index); err != nil {
		return err
	}

	// Update the task of the Todo item
	t[index].Task = task
	t[index].UpdatedAt = time.Now().Format(time.RFC3339)

	return nil
}

func (todoList *TodoList) view(format string) {
	t := *todoList

	switch format {
	case "json":
		t.viewJSON("raw")
		return
	case "pretty":
		t.viewJSON("pretty")
		return
	case "table":
		t.viewTable()
		return
	case "none":
		return
	default:
		t.viewJSON("raw")
		return
	}
}

func (todoList *TodoList) viewJSON(style string) {
	t := *todoList

	var jsonOutput []byte
	var err error

	if style == "pretty" {
		jsonOutput, err = json.MarshalIndent(t, "", "  ")
	} else {
		jsonOutput, err = json.Marshal(t)
	}

	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println(string(jsonOutput))
}

func (todoList *TodoList) viewTable() {
	if len(*todoList) == 0 {
		fmt.Println("No todos found.")
		return
	}

	todoType := reflect.TypeOf(Todo{})

	timeFormat := "2006-01-02"

	var headers []string
	for i := 0; i < todoType.NumField(); i++ {
		headers = append(headers, todoType.Field(i).Name)
	}

	t := table.New(os.Stdout)

	// Table options
	t.SetRowLines(false)
	// t.SetDividers(table.MarkdownDividers)

	t.SetHeaders(headers...)

	for _, todo := range *todoList {
		// Handle all time fields consistently
		var createdAtStr, updatedAtStr, completedAtStr string

		// CreatedAt
		if createdAt, err := time.Parse(time.RFC3339, todo.CreatedAt); err == nil {
			createdAtStr = createdAt.Format(timeFormat)
		} else {
			createdAtStr = "Invalid"
		}

		// UpdatedAt
		if updatedAt, err := time.Parse(time.RFC3339, todo.UpdatedAt); err == nil {
			updatedAtStr = updatedAt.Format(timeFormat)
		} else {
			updatedAtStr = "Invalid"
		}

		// CompletedAt
		if todo.CompletedAt != "" {
			if completedAt, err := time.Parse(time.RFC3339, todo.CompletedAt); err == nil {
				completedAtStr = completedAt.Format(timeFormat)
			} else {
				completedAtStr = "Invalid"
			}
		} else {
			completedAtStr = ""
		}

		// Completion to emoji
		var completedEmoji string
		if todo.Completed {
			completedEmoji = "✅"
		} else {
			completedEmoji = "❌"
		}

		t.AddRow(
			fmt.Sprintf("%d", todo.ID),
			todo.Task,
			completedEmoji,
			createdAtStr,
			updatedAtStr,
			tml.Sprintf("<green>%s</green>", completedAtStr),
		)
	}

	t.Render()
}
