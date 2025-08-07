package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to create a temporary test environment
func setupTestEnvironment(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "todo_cmd_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Change to temp directory
	oldWd, _ := os.Getwd()
	os.Chdir(tempDir)

	cleanup := func() {
		os.Chdir(oldWd)
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestAddCommand_Creation(t *testing.T) {
	cmd := NewAddCommand()

	if cmd.Name != "add" {
		t.Errorf("NewAddCommand() Name = %s, want 'add'", cmd.Name)
	}

	if cmd.Usage != "Add a new todo item" {
		t.Errorf("NewAddCommand() Usage = %s, want 'Add a new todo item'", cmd.Usage)
	}

	if len(cmd.Aliases) == 0 || cmd.Aliases[0] != "a" {
		t.Errorf("NewAddCommand() should have alias 'a'")
	}
}

func TestDeleteCommand_Creation(t *testing.T) {
	cmd := NewDeleteCommand()

	if cmd.Name != "delete" {
		t.Errorf("NewDeleteCommand() Name = %s, want 'delete'", cmd.Name)
	}

	if cmd.Usage != "Delete a todo item by ID" {
		t.Errorf("NewDeleteCommand() Usage = %s, want 'Delete a todo item by ID'", cmd.Usage)
	}

	expectedAliases := []string{"del", "rm"}
	if len(cmd.Aliases) != 2 || cmd.Aliases[0] != "del" || cmd.Aliases[1] != "rm" {
		t.Errorf("NewDeleteCommand() Aliases = %v, want %v", cmd.Aliases, expectedAliases)
	}
}

func TestEditCommand_Creation(t *testing.T) {
	cmd := NewEditCommand()

	if cmd.Name != "edit" {
		t.Errorf("NewEditCommand() Name = %s, want 'edit'", cmd.Name)
	}

	if cmd.Usage != "Edit a todo item by ID" {
		t.Errorf("NewEditCommand() Usage = %s, want 'Edit a todo item by ID'", cmd.Usage)
	}

	if len(cmd.Aliases) == 0 || cmd.Aliases[0] != "e" {
		t.Errorf("NewEditCommand() should have alias 'e'")
	}
}

func TestToggleCommand_Creation(t *testing.T) {
	cmd := NewToggleCommand()

	if cmd.Name != "toggle" {
		t.Errorf("NewToggleCommand() Name = %s, want 'toggle'", cmd.Name)
	}

	if cmd.Usage != "Toggle completion status of a todo item by ID" {
		t.Errorf("NewToggleCommand() Usage incorrect")
	}

	expectedAliases := []string{"t", "complete"}
	if len(cmd.Aliases) != 2 || cmd.Aliases[0] != "t" || cmd.Aliases[1] != "complete" {
		t.Errorf("NewToggleCommand() Aliases = %v, want %v", cmd.Aliases, expectedAliases)
	}
}

func TestListCommand_Creation(t *testing.T) {
	cmd := NewListCommand()

	if cmd.Name != "list" {
		t.Errorf("NewListCommand() Name = %s, want 'list'", cmd.Name)
	}

	if cmd.Usage != "List all todo items" {
		t.Errorf("NewListCommand() Usage = %s, want 'List all todo items'", cmd.Usage)
	}

	expectedAliases := []string{"l", "ls"}
	if len(cmd.Aliases) != 2 || cmd.Aliases[0] != "l" || cmd.Aliases[1] != "ls" {
		t.Errorf("NewListCommand() Aliases = %v, want %v", cmd.Aliases, expectedAliases)
	}

	// Check if format flag is present
	if len(cmd.Flags) == 0 {
		t.Errorf("NewListCommand() should have format flag")
	}
}

func TestVersionCommand_Creation(t *testing.T) {
	cmd := NewVersionCommand()

	if cmd.Name != "version" {
		t.Errorf("NewVersionCommand() Name = %s, want 'version'", cmd.Name)
	}

	if cmd.Usage != "Display the version of the application" {
		t.Errorf("NewVersionCommand() Usage incorrect")
	}

	if len(cmd.Aliases) == 0 || cmd.Aliases[0] != "v" {
		t.Errorf("NewVersionCommand() should have alias 'v'")
	}
}

// Test the commands with actual functionality (requires proper CLI setup)
func TestCommandsWithTodoList(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test initializeTodoList function
	todoList, storage, err := initializeTodoList()
	if err != nil {
		t.Fatalf("initializeTodoList() error = %v", err)
	}

	if todoList == nil {
		t.Fatalf("initializeTodoList() returned nil todoList")
	}

	if storage == nil {
		t.Fatalf("initializeTodoList() returned nil storage")
	}

	// Test adding a task through the interface
	err = todoList.Add("Test task")
	if err != nil {
		t.Errorf("TodoList.Add() error = %v", err)
	}

	if len(*todoList) != 1 {
		t.Errorf("TodoList.Add() length = %d, want 1", len(*todoList))
	}

	// Test saving
	err = storage.Save(*todoList)
	if err != nil {
		t.Errorf("Storage.Save() error = %v", err)
	}

	// Test loading
	loadedList, err := storage.Load()
	if err != nil {
		t.Errorf("Storage.Load() error = %v", err)
	}

	if len(loadedList) != 1 {
		t.Errorf("Storage.Load() length = %d, want 1", len(loadedList))
	}

	if loadedList[0].Task != "Test task" {
		t.Errorf("Storage.Load() task = %s, want 'Test task'", loadedList[0].Task)
	}
}

// TestGetStoragePath tests the GetStoragePath function
func TestGetStoragePath(t *testing.T) {
	tests := []struct {
		name     string
		global   bool
		expected string
	}{
		{
			name:     "local storage",
			global:   false,
			expected: ".todos.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.global {
				// Test local storage path
				path, err := GetStoragePath(tt.global)
				if err != nil {
					t.Errorf("GetStoragePath() error = %v", err)
				}
				if path != tt.expected {
					t.Errorf("GetStoragePath() = %s, want %s", path, tt.expected)
				}
			}
		})
	}

	// Test global storage path with mock home directory
	t.Run("global storage", func(t *testing.T) {
		// Create a temporary directory to use as mock home
		tempDir, err := os.MkdirTemp("", "mock_home")
		if err != nil {
			t.Fatalf("Failed to create temp directory: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Save original environment variables
		oldHome := os.Getenv("HOME")
		oldUserProfile := os.Getenv("USERPROFILE")
		defer func() {
			os.Setenv("HOME", oldHome)
			os.Setenv("USERPROFILE", oldUserProfile)
		}()

		// Set mock home directory
		os.Setenv("HOME", tempDir)
		os.Setenv("USERPROFILE", tempDir)

		// Test global storage path
		path, err := GetStoragePath(true)
		if err != nil {
			t.Errorf("GetStoragePath() error = %v", err)
		}

		// Verify path contains expected components
		if !strings.Contains(path, ".todos") {
			t.Errorf("GetStoragePath() = %s, should contain '.todos'", path)
		}
		if !strings.Contains(path, "todos.json") {
			t.Errorf("GetStoragePath() = %s, should contain 'todos.json'", path)
		}

		// Verify .todos directory was created
		todosDirPath := filepath.Join(tempDir, ".todos")
		if _, err := os.Stat(todosDirPath); os.IsNotExist(err) {
			t.Errorf(".todos directory was not created at %s", todosDirPath)
		}
	})
}

// TestInitializeTodoListWithPath tests the initializeTodoListWithPath function
func TestInitializeTodoListWithPath(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a custom storage file path
	customPath := filepath.Join(tempDir, "custom.json")

	// Initialize with custom path
	todoList, storage, err := initializeTodoListWithPath(customPath)
	if err != nil {
		t.Errorf("initializeTodoListWithPath() error = %v", err)
	}

	if todoList == nil {
		t.Error("initializeTodoListWithPath() todoList is nil")
	}

	if storage == nil {
		t.Error("initializeTodoListWithPath() storage is nil")
	}

	// Add a task and verify it uses the custom path
	err = todoList.Add("Test custom path")
	if err != nil {
		t.Errorf("TodoList.Add() error = %v", err)
	}

	err = storage.Save(*todoList)
	if err != nil {
		t.Errorf("Storage.Save() error = %v", err)
	}

	// Verify file was created at custom path
	if _, err := os.Stat(customPath); os.IsNotExist(err) {
		t.Errorf("Custom storage file was not created at %s", customPath)
	}
}
