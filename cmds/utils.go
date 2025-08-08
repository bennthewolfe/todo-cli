package commands

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/aquasecurity/table"
	"github.com/liamg/tml"
	"github.com/urfave/cli/v3"
)

// Storage represents the storage interface for TodoList
type Storage[T any] struct {
	filename string
}

// NewStorage creates a new storage instance
func NewStorage[T any](filename string) *Storage[T] {
	return &Storage[T]{filename: filename}
}

// Save saves data to the storage file
func (s *Storage[T]) Save(data T) error {
	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data to JSON: %w", err)
	}
	return os.WriteFile(s.filename, fileData, 0644)
}

// Load loads data from the storage file
func (s *Storage[T]) Load() (T, error) {
	var data T

	// Check if the file exists, if not create an empty one
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		emptyFile, err := os.Create(s.filename)
		if err != nil {
			return data, fmt.Errorf("error creating file: %w", err)
		}
		emptyFile.Close()
	}

	// Read the file content
	fileData, err := os.ReadFile(s.filename)
	if err != nil {
		return data, fmt.Errorf("error reading file: %w", err)
	}

	// Check if the file is empty
	if len(fileData) == 0 {
		return data, nil
	}

	if err := json.Unmarshal(fileData, &data); err != nil {
		return data, fmt.Errorf("error unmarshaling JSON data: %w", err)
	}

	return data, nil
}

// generateShortGUID generates a short GUID-like identifier
func generateShortGUID() string {
	bytes := make([]byte, 6) // 12 character hex string
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Todo represents a single todo item
type Todo struct {
	InternalID  string `json:"internal_id"` // Hidden GUID for internal tracking
	Task        string `json:"task"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	CompletedAt string `json:"completed_at,omitempty"`
}

// TodoList type for the commands package
type TodoList []Todo

// Interface implementation for TodoListInterface
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

// Private methods
func (todoList *TodoList) validateIndex(index int) error {
	if index < 0 || index >= len(*todoList) {
		return fmt.Errorf("invalid index: %d", index)
	}
	return nil
}

func (todoList *TodoList) add(task string) error {
	todo := Todo{
		InternalID:  generateShortGUID(),
		Task:        task,
		Completed:   false,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
		CompletedAt: "",
	}

	// Add a new Todo item to the list
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

	// If list is empty, output null to match expected behavior
	if len(t) == 0 {
		fmt.Println("null")
		return
	}

	// Create display version with index-based IDs
	type DisplayTodo struct {
		ID          int    `json:"id"`
		Task        string `json:"task"`
		Completed   bool   `json:"completed"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		CompletedAt string `json:"completed_at,omitempty"`
	}

	displayTodos := make([]DisplayTodo, len(t))
	for index, todo := range t {
		displayTodos[index] = DisplayTodo{
			ID:          index + 1, // Use array index + 1 as display ID
			Task:        todo.Task,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
			CompletedAt: todo.CompletedAt,
		}
	}

	var jsonOutput []byte
	var err error

	if style == "pretty" {
		jsonOutput, err = json.MarshalIndent(displayTodos, "", "  ")
	} else {
		jsonOutput, err = json.Marshal(displayTodos)
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

	// Dynamically generate headers, but skip InternalID and add ID at the beginning
	var headers []string
	headers = append(headers, "ID") // Add ID as first column

	for i := 0; i < todoType.NumField(); i++ {
		field := todoType.Field(i)
		// Skip the InternalID field since it's for internal use only
		if field.Name != "InternalID" {
			headers = append(headers, field.Name)
		}
	}

	t := table.New(os.Stdout)

	// Table options
	t.SetRowLines(false)
	// t.SetDividers(table.MarkdownDividers)

	t.SetHeaders(headers...)

	for index, todo := range *todoList {
		// Use index + 1 as the display ID
		displayID := index + 1

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

		// Add row with ID first, then other fields (excluding InternalID)
		t.AddRow(
			fmt.Sprintf("%d", displayID), // ID column
			todo.Task,                    // Task column
			completedEmoji,               // Completed column
			createdAtStr,                 // CreatedAt column
			updatedAtStr,                 // UpdatedAt column
			tml.Sprintf("<green>%s</green>", completedAtStr), // CompletedAt column
		)
	}

	t.Render()
}

// GetStoragePath returns the appropriate storage path based on the global flag
func GetStoragePath(isGlobal bool) (string, error) {
	if !isGlobal {
		return ".todos.json", nil
	}

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user home directory: %w", err)
	}

	// Create ~/.todo directory if it doesn't exist
	todoDir := filepath.Join(homeDir, ".todo")
	if err := os.MkdirAll(todoDir, 0755); err != nil {
		return "", fmt.Errorf("unable to create todo directory: %w", err)
	}

	return filepath.Join(todoDir, "todos.json"), nil
}

// GetEffectiveStoragePath returns the appropriate storage path based on both global and archive flags
func GetEffectiveStoragePath(isGlobal, isArchive bool) (string, error) {
	if isArchive {
		return GetArchivePath(isGlobal)
	}
	return GetStoragePath(isGlobal)
}

// IsCommandAllowedWithArchive checks if a command is allowed when using the --archive flag
func IsCommandAllowedWithArchive(commandName string) bool {
	allowedCommands := map[string]bool{
		"list":   true,
		"delete": true,
	}
	return allowedCommands[commandName]
}

// ValidateArchiveFlagUsage validates that --archive flag is only used with supported commands
func ValidateArchiveFlagUsage(c *cli.Command, commandName string) error {
	if c.Bool("archive") && !IsCommandAllowedWithArchive(commandName) {
		return fmt.Errorf("--archive flag is only supported with 'list' and 'delete' commands, not '%s'", commandName)
	}
	return nil
}

// GetArchivePath returns the appropriate archive storage path based on the global flag
func GetArchivePath(isGlobal bool) (string, error) {
	if !isGlobal {
		return ".todos.archive.json", nil
	}

	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user home directory: %w", err)
	}

	// Create ~/.todo directory if it doesn't exist
	todoDir := filepath.Join(homeDir, ".todo")
	if err := os.MkdirAll(todoDir, 0755); err != nil {
		return "", fmt.Errorf("unable to create todo directory: %w", err)
	}

	return filepath.Join(todoDir, "todos.archive.json"), nil
}

// CheckAndExecuteListFlag checks if the --list flag is set and executes list command if so
// Returns true if list should be executed after the main command, false otherwise
func CheckAndExecuteListFlag(c *cli.Command) bool {
	return c.Bool("list")
}

// ExecuteListCommand executes the list command with table format
func ExecuteListCommand(c *cli.Command) error {
	if c.Bool("debug") {
		fmt.Println("DEBUG: Executing list command after main action")
	}

	// Get the appropriate storage path based on global and archive flags
	storagePath, err := GetEffectiveStoragePath(c.Bool("global"), c.Bool("archive"))
	if err != nil {
		return fmt.Errorf("error getting storage path: %w", err)
	}

	// Initialize todo list and storage
	todoList, _, err := initializeTodoListWithPath(storagePath)
	if err != nil {
		return fmt.Errorf("failed to initialize todo list: %w", err)
	}

	// Display todos with table format (default)
	fmt.Println() // Add a blank line before list output
	todoList.View("table")
	return nil
}

// initializeTodoList initializes the todo list and storage with custom path
func initializeTodoListWithPath(storagePath string) (*TodoList, *Storage[TodoList], error) {
	todoList := &TodoList{}
	storage := NewStorage[TodoList](storagePath)

	loadedList, err := storage.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("error loading todos: %w", err)
	}

	*todoList = loadedList
	return todoList, storage, nil
}

// initializeTodoList initializes the todo list and storage (legacy function for compatibility)
func initializeTodoList() (*TodoList, *Storage[TodoList], error) {
	return initializeTodoListWithPath(".todos.json")
}
