package main

import (
	"fmt"
	"os"
	"testing"
)

// TestMain sets up and tears down test environment
func TestMain(m *testing.M) {
	// Setup: You can add any global test setup here
	fmt.Println("Setting up tests...")

	// Run all tests
	code := m.Run()

	// Teardown: You can add any global test cleanup here
	fmt.Println("Cleaning up after tests...")

	os.Exit(code)
}

// Benchmark tests for performance
func BenchmarkTodoList_Add(b *testing.B) {
	todoList := &TodoList{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		todoList.Add(fmt.Sprintf("Task %d", i))
	}
}

func BenchmarkTodoList_Delete(b *testing.B) {
	// Setup: Create a list with many items
	todoList := &TodoList{}
	for i := 0; i < 1000; i++ {
		todoList.Add(fmt.Sprintf("Task %d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N && len(*todoList) > 0; i++ {
		todoList.Delete(0)
	}
}

func BenchmarkStorage_Save(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "bench_test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data
	todoList := &TodoList{}
	for i := 0; i < 100; i++ {
		todoList.Add(fmt.Sprintf("Task %d", i))
	}

	storage := NewStorage[TodoList](tempDir + "/bench_todos.json")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.Save(*todoList)
	}
}

func BenchmarkStorage_Load(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "bench_test")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data and save it
	todoList := &TodoList{}
	for i := 0; i < 100; i++ {
		todoList.Add(fmt.Sprintf("Task %d", i))
	}

	storage := NewStorage[TodoList](tempDir + "/bench_todos.json")
	storage.Save(*todoList)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.Load()
	}
}
