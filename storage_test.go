package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStorage_Save(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "todo_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test_todos.json")
	storage := NewStorage[TodoList](testFile)

	// Test data
	testTodos := TodoList{
		{ID: 1, Task: "Test task 1", Completed: false, CreatedAt: "2023-01-01T00:00:00Z", UpdatedAt: "2023-01-01T00:00:00Z"},
		{ID: 2, Task: "Test task 2", Completed: true, CreatedAt: "2023-01-01T00:00:00Z", UpdatedAt: "2023-01-01T00:00:00Z", CompletedAt: "2023-01-01T01:00:00Z"},
	}

	// Test saving
	err = storage.Save(testTodos)
	if err != nil {
		t.Errorf("Save() error = %v, want nil", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Save() did not create file")
	}
}

func TestStorage_Load(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "todo_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test_todos.json")
	storage := NewStorage[TodoList](testFile)

	// Test loading from non-existent file (should create empty file)
	todos, err := storage.Load()
	if err != nil {
		t.Errorf("Load() error = %v, want nil", err)
	}

	if len(todos) != 0 {
		t.Errorf("Load() from empty file returned %d items, want 0", len(todos))
	}

	// Test saving and loading data
	testTodos := TodoList{
		{ID: 1, Task: "Test task 1", Completed: false, CreatedAt: "2023-01-01T00:00:00Z", UpdatedAt: "2023-01-01T00:00:00Z"},
		{ID: 2, Task: "Test task 2", Completed: true, CreatedAt: "2023-01-01T00:00:00Z", UpdatedAt: "2023-01-01T00:00:00Z", CompletedAt: "2023-01-01T01:00:00Z"},
	}

	err = storage.Save(testTodos)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loadedTodos, err := storage.Load()
	if err != nil {
		t.Errorf("Load() error = %v, want nil", err)
	}

	if len(loadedTodos) != len(testTodos) {
		t.Errorf("Load() returned %d items, want %d", len(loadedTodos), len(testTodos))
	}

	// Verify content
	for i, todo := range loadedTodos {
		if todo.ID != testTodos[i].ID {
			t.Errorf("Load() ID = %d, want %d", todo.ID, testTodos[i].ID)
		}
		if todo.Task != testTodos[i].Task {
			t.Errorf("Load() Task = %s, want %s", todo.Task, testTodos[i].Task)
		}
		if todo.Completed != testTodos[i].Completed {
			t.Errorf("Load() Completed = %t, want %t", todo.Completed, testTodos[i].Completed)
		}
	}
}

func TestStorage_LoadInvalidJSON(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "todo_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "invalid.json")

	// Create file with invalid JSON
	err = os.WriteFile(testFile, []byte(`{"incomplete": "json"`), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	storage := NewStorage[TodoList](testFile)

	// Test loading invalid JSON
	_, err = storage.Load()
	if err == nil {
		t.Errorf("Load() with invalid JSON should return error")
	}
}
