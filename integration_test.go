package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLIIntegration tests the CLI application end-to-end
func TestCLIIntegration(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	// Create temp directory for test data
	tempDir, err := os.MkdirTemp("", "todo_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
		expectedCode   int
	}{
		{
			name:           "add task",
			args:           []string{"add", "Test task"},
			expectedOutput: "Added task: Test task",
			expectError:    false,
		},
		{
			name:           "add without args",
			args:           []string{"add"},
			expectedOutput: "task description is required",
			expectError:    true,
			expectedCode:   1,
		},
		{
			name:           "list empty",
			args:           []string{"list", "--format", "json"},
			expectedOutput: "null",
			expectError:    false,
		},
		{
			name:           "list invalid format",
			args:           []string{"list", "--format", "invalid"},
			expectedOutput: "invalid format",
			expectError:    true,
			expectedCode:   1,
		},
		{
			name:           "delete without args",
			args:           []string{"delete"},
			expectedOutput: "exactly one ID is required",
			expectError:    true,
			expectedCode:   1,
		},
		{
			name:           "delete invalid ID",
			args:           []string{"delete", "abc"},
			expectedOutput: "invalid ID: abc must be a number",
			expectError:    true,
			expectedCode:   1,
		},
		{
			name:           "edit without args",
			args:           []string{"edit"},
			expectedOutput: "ID and new task description are required",
			expectError:    true,
			expectedCode:   1,
		},
		{
			name:           "toggle without args",
			args:           []string{"toggle"},
			expectedOutput: "exactly one ID is required",
			expectError:    true,
			expectedCode:   1,
		},
		{
			name:           "version command",
			args:           []string{"version"},
			expectedOutput: "TODO CLI Version:",
			expectError:    false,
		},
		{
			name:           "help command",
			args:           []string{"help"},
			expectedOutput: "Todo CLI - A simple command-line interface",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a separate temp directory for each test
			testTempDir, err := os.MkdirTemp("", "todo_integration_test_"+tt.name)
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(testTempDir)

			// Change to test directory
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd)
			os.Chdir(testTempDir)

			cmd := exec.Command(buildPath, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := strings.TrimSpace(string(output))

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				}
				if exitError, ok := err.(*exec.ExitError); ok {
					if tt.expectedCode != 0 && exitError.ExitCode() != tt.expectedCode {
						t.Errorf("Expected exit code %d, got %d", tt.expectedCode, exitError.ExitCode())
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			if !strings.Contains(outputStr, tt.expectedOutput) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, outputStr)
			}
		})
	}
}

// TestCLIWorkflow tests a complete workflow
func TestCLIWorkflow(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	// Create temp directory for test data
	tempDir, err := os.MkdirTemp("", "todo_workflow_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Test workflow: add -> list -> edit -> toggle -> delete

	// 1. Add a task
	cmd = exec.Command(buildPath, "add", "Buy groceries")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to add task: %v, output: %s", err, output)
	}
	if !strings.Contains(string(output), "Added task: Buy groceries") {
		t.Errorf("Add command output unexpected: %s", output)
	}

	// 2. Add another task
	cmd = exec.Command(buildPath, "add", "Walk the dog")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to add second task: %v", err)
	}
	if !strings.Contains(string(output), "Added task: Walk the dog") {
		t.Errorf("Second add command output unexpected: %s", output)
	}

	// 3. List tasks in JSON format
	cmd = exec.Command(buildPath, "list", "--format", "json")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}
	if !strings.Contains(string(output), "Buy groceries") || !strings.Contains(string(output), "Walk the dog") {
		t.Errorf("List command should show both tasks: %s", output)
	}

	// 4. Edit first task
	cmd = exec.Command(buildPath, "edit", "1", "Buy organic groceries")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to edit task: %v", err)
	}
	if !strings.Contains(string(output), "Updated todo item 1") {
		t.Errorf("Edit command output unexpected: %s", output)
	}

	// 5. Toggle first task to completed
	cmd = exec.Command(buildPath, "toggle", "1")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to toggle task: %v", err)
	}
	if !strings.Contains(string(output), "Toggled completion status for todo item with ID: 1") {
		t.Errorf("Toggle command output unexpected: %s", output)
	}

	// 6. Verify task is completed
	cmd = exec.Command(buildPath, "list", "--format", "json")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to list tasks after toggle: %v", err)
	}
	if !strings.Contains(string(output), "Buy organic groceries") {
		t.Errorf("Task should be updated: %s", output)
	}

	// 7. Delete second task
	cmd = exec.Command(buildPath, "delete", "2")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}
	if !strings.Contains(string(output), "Deleted todo item with ID: 2") {
		t.Errorf("Delete command output unexpected: %s", output)
	}

	// 8. Verify only one task remains
	cmd = exec.Command(buildPath, "list", "--format", "json")
	output, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to list tasks after delete: %v", err)
	}
	if strings.Contains(string(output), "Walk the dog") {
		t.Errorf("Deleted task should not appear: %s", output)
	}
	if !strings.Contains(string(output), "Buy organic groceries") {
		t.Errorf("Remaining task should still exist: %s", output)
	}
}

// TestCLIHelp tests help functionality
func TestCLIHelp(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
	}{
		{
			name: "general help",
			args: []string{"help"},
			expectedOutput: []string{
				"Todo CLI - A simple command-line interface",
				"COMMANDS:",
				"add",
				"delete",
				"edit",
				"list",
				"toggle",
				"version",
				"help",
			},
		},
		{
			name: "help flag",
			args: []string{"--help"},
			expectedOutput: []string{
				"Todo CLI",
				"USAGE:",
				"COMMANDS:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(buildPath, tt.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Help command failed: %v", err)
			}

			outputStr := string(output)
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Help output should contain %q, got: %s", expected, outputStr)
				}
			}
		})
	}
}

// TestCLIGlobalStorage tests the global storage functionality
func TestCLIGlobalStorage(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	// Create temp directory for test data
	tempDir, err := os.MkdirTemp("", "todo_global_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock home directory for testing
	mockHomeDir := filepath.Join(tempDir, "home")
	if err := os.MkdirAll(mockHomeDir, 0755); err != nil {
		t.Fatalf("Failed to create mock home directory: %v", err)
	}

	// Set HOME environment variable for the test
	oldHome := os.Getenv("HOME")
	oldUserProfile := os.Getenv("USERPROFILE")
	defer func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("USERPROFILE", oldUserProfile)
	}()
	os.Setenv("HOME", mockHomeDir)
	os.Setenv("USERPROFILE", mockHomeDir)

	// Change to temp directory for local storage
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Test adding todo to global storage
	t.Run("add_global_todo", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--global", "add", "Global test task")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to add global todo: %v\nOutput: %s", err, output)
		}

		expectedMsg := "Added task: Global test task"
		if !strings.Contains(string(output), expectedMsg) {
			t.Errorf("Expected %q in output, got: %s", expectedMsg, output)
		}

		// Verify global storage file was created
		globalStoragePath := filepath.Join(mockHomeDir, ".todos", "todos.json")
		if _, err := os.Stat(globalStoragePath); os.IsNotExist(err) {
			t.Errorf("Global storage file was not created at %s", globalStoragePath)
		}
	})

	// Test listing global todos
	t.Run("list_global_todos", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--global", "list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list global todos: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Global test task") {
			t.Errorf("Expected global todo in output, got: %s", outputStr)
		}
	})

	// Test adding local todo
	t.Run("add_local_todo", func(t *testing.T) {
		cmd := exec.Command(buildPath, "add", "Local test task")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to add local todo: %v\nOutput: %s", err, output)
		}

		expectedMsg := "Added task: Local test task"
		if !strings.Contains(string(output), expectedMsg) {
			t.Errorf("Expected %q in output, got: %s", expectedMsg, output)
		}
	})

	// Test that local and global storage are separate
	t.Run("storage_separation", func(t *testing.T) {
		// List local todos
		cmd := exec.Command(buildPath, "list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list local todos: %v\nOutput: %s", err, output)
		}

		localOutput := string(output)
		if !strings.Contains(localOutput, "Local test task") {
			t.Errorf("Expected local todo in local output, got: %s", localOutput)
		}
		if strings.Contains(localOutput, "Global test task") {
			t.Errorf("Global todo should not appear in local output, got: %s", localOutput)
		}

		// List global todos
		cmd = exec.Command(buildPath, "--global", "list")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list global todos: %v\nOutput: %s", err, output)
		}

		globalOutput := string(output)
		if !strings.Contains(globalOutput, "Global test task") {
			t.Errorf("Expected global todo in global output, got: %s", globalOutput)
		}
		if strings.Contains(globalOutput, "Local test task") {
			t.Errorf("Local todo should not appear in global output, got: %s", globalOutput)
		}
	})

	// Test default action with global flag
	t.Run("default_action_global", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--global")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run default action with global flag: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Global test task") {
			t.Errorf("Expected global todo in default action output, got: %s", outputStr)
		}
	})

	// Test global flag with other commands
	t.Run("global_toggle", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--global", "toggle", "1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to toggle global todo: %v\nOutput: %s", err, output)
		}

		// Verify the todo was toggled
		cmd = exec.Command(buildPath, "--global", "list", "--format", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list global todos after toggle: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, `"completed":true`) {
			t.Errorf("Expected todo to be marked as completed, got: %s", outputStr)
		}
	})

	// Test help contains global flag information
	t.Run("help_contains_global_flag", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to get help: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		expectedTexts := []string{
			"--global",
			"-g",
			"global todo storage",
			"home directory",
		}

		for _, expected := range expectedTexts {
			if !strings.Contains(outputStr, expected) {
				t.Errorf("Expected help to contain %q, got: %s", expected, outputStr)
			}
		}
	})
}
