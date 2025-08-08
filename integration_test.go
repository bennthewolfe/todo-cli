package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
		globalStoragePath := filepath.Join(mockHomeDir, ".todo", "todos.json")
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

// TestCLIArchive tests the archive command functionality
func TestCLIArchive(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")
	if runtime.GOOS != "windows" {
		buildPath = filepath.Join(t.TempDir(), "todo")
	}

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "todo_archive_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Test archive without items
	t.Run("archive_invalid_id", func(t *testing.T) {
		cmd := exec.Command(buildPath, "archive", "1")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected archive to fail with no items, but it succeeded")
		}

		expectedMsg := "invalid ID"
		if !strings.Contains(string(output), expectedMsg) {
			t.Errorf("Expected %q in error output, got: %s", expectedMsg, output)
		}
	})

	// Add some test items
	t.Run("setup_for_archive", func(t *testing.T) {
		// Add test items
		testItems := []string{"Item to archive", "Another item", "Third item"}
		for _, item := range testItems {
			cmd := exec.Command(buildPath, "add", item)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to add test item: %v\nOutput: %s", err, output)
			}
		}
	})

	// Test successful archive
	t.Run("archive_item", func(t *testing.T) {
		cmd := exec.Command(buildPath, "archive", "2")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to archive item: %v\nOutput: %s", err, output)
		}

		expectedMsg := "Archived todo item: Another item"
		if !strings.Contains(string(output), expectedMsg) {
			t.Errorf("Expected %q in output, got: %s", expectedMsg, output)
		}

		// Verify archive file was created
		if _, err := os.Stat(".todos.archive.json"); os.IsNotExist(err) {
			t.Errorf("Archive file was not created")
		}

		// Verify item was removed from main list
		cmd = exec.Command(buildPath, "list", "--format", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list todos after archive: %v\nOutput: %s", err, output)
		}

		if strings.Contains(string(output), "Another item") {
			t.Errorf("Item still appears in main list after archiving")
		}

		// Verify only 2 items remain
		outputStr := string(output)
		if strings.Count(outputStr, `"task":`) != 2 {
			t.Errorf("Expected 2 items in main list after archive, got output: %s", outputStr)
		}
	})

	// Test archive without arguments
	t.Run("archive_without_args", func(t *testing.T) {
		cmd := exec.Command(buildPath, "archive")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected archive to fail without arguments, but it succeeded")
		}

		expectedMsg := "exactly one ID is required"
		if !strings.Contains(string(output), expectedMsg) {
			t.Errorf("Expected %q in error output, got: %s", expectedMsg, output)
		}
	})

	// Test archive with invalid ID format
	t.Run("archive_invalid_format", func(t *testing.T) {
		cmd := exec.Command(buildPath, "archive", "abc")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected archive to fail with invalid ID format, but it succeeded")
		}

		expectedMsg := "invalid ID: abc must be a number"
		if !strings.Contains(string(output), expectedMsg) {
			t.Errorf("Expected %q in error output, got: %s", expectedMsg, output)
		}
	})
}

// TestCLIGlobalArchive tests archive command with global flag
func TestCLIGlobalArchive(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")
	if runtime.GOOS != "windows" {
		buildPath = filepath.Join(t.TempDir(), "todo")
	}

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	// Create temporary directory to use as mock home
	mockHomeDir, err := os.MkdirTemp("", "mock_home")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(mockHomeDir)

	// Save original environment variables
	oldHome := os.Getenv("HOME")
	oldUserProfile := os.Getenv("USERPROFILE")
	defer func() {
		os.Setenv("HOME", oldHome)
		os.Setenv("USERPROFILE", oldUserProfile)
	}()

	// Set mock home directory
	os.Setenv("HOME", mockHomeDir)
	os.Setenv("USERPROFILE", mockHomeDir)

	// Add a global todo for testing
	t.Run("setup_global_item", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--global", "add", "Global item to archive")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to add global todo: %v\nOutput: %s", err, output)
		}
	})

	// Test global archive
	t.Run("archive_global_item", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--global", "archive", "1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to archive global item: %v\nOutput: %s", err, output)
		}

		expectedMsg := "Archived todo item: Global item to archive"
		if !strings.Contains(string(output), expectedMsg) {
			t.Errorf("Expected %q in output, got: %s", expectedMsg, output)
		}

		// Verify global archive file was created
		globalArchivePath := filepath.Join(mockHomeDir, ".todo", "todos.archive.json")
		if _, err := os.Stat(globalArchivePath); os.IsNotExist(err) {
			t.Errorf("Global archive file was not created at %s", globalArchivePath)
		}

		// Verify item was removed from global list
		cmd = exec.Command(buildPath, "--global", "list")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list global todos after archive: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if strings.Contains(outputStr, "Global item to archive") {
			t.Errorf("Item still appears in global list after archiving")
		}
	})
}

// TestCLIListFlag tests the --list flag functionality
func TestCLIListFlag(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	// Create temp directory for test data
	tempDir, err := os.MkdirTemp("", "todo_list_flag_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Clean up any existing todos
	os.Remove(".todos.json")

	t.Run("list_flag_with_add", func(t *testing.T) {
		cmd := exec.Command(buildPath, "add", "Test item for list flag", "--list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to add task with --list flag: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		// Should contain both the add confirmation and the list output
		if !strings.Contains(outputStr, "Added task: Test item for list flag") {
			t.Errorf("Expected add confirmation in output, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Test item for list flag") {
			t.Errorf("Expected task to appear in list output, got: %s", outputStr)
		}
		// Should contain table headers
		if !strings.Contains(outputStr, "ID") || !strings.Contains(outputStr, "Task") {
			t.Errorf("Expected table headers in list output, got: %s", outputStr)
		}
	})

	t.Run("list_flag_before_command", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--list", "toggle", "1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to toggle with --list flag: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		// Should contain both the toggle confirmation and the list output
		if !strings.Contains(outputStr, "Toggled completion status for todo item with ID: 1") {
			t.Errorf("Expected toggle confirmation in output, got: %s", outputStr)
		}
		// Should contain table output after the action
		if !strings.Contains(outputStr, "Test item for list flag") {
			t.Errorf("Expected task to appear in list output after toggle, got: %s", outputStr)
		}
	})

	t.Run("list_flag_with_global", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--global", "add", "Global item with list flag", "--list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to add global task with --list flag: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		// Should contain both the add confirmation and the global list output
		if !strings.Contains(outputStr, "Added task: Global item with list flag") {
			t.Errorf("Expected add confirmation in output, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Global item with list flag") {
			t.Errorf("Expected global task to appear in list output, got: %s", outputStr)
		}
	})

	t.Run("list_flag_with_global_before_command", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--global", "--list", "edit", "1", "Updated global item")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to edit global task with --list flag: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		// Should contain both the edit confirmation and the list output
		if !strings.Contains(outputStr, "Updated todo item 1: Updated global item") {
			t.Errorf("Expected edit confirmation in output, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Updated global item") {
			t.Errorf("Expected updated task to appear in list output, got: %s", outputStr)
		}
	})

	t.Run("list_flag_with_version", func(t *testing.T) {
		cmd := exec.Command(buildPath, "version", "--list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run version with --list flag: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		// Should contain both the version info and the list output
		if !strings.Contains(outputStr, "TODO CLI Version:") {
			t.Errorf("Expected version info in output, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Test item for list flag") {
			t.Errorf("Expected task to appear in list output after version, got: %s", outputStr)
		}
	})

	t.Run("help_contains_list_flag", func(t *testing.T) {
		cmd := exec.Command(buildPath, "help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to get help: %v\nOutput: %s", err, output)
		}

		expectedFlags := []string{
			"--list",
			"-l",
			"List all todo items (overrides other commands)",
		}

		outputStr := string(output)
		for _, flag := range expectedFlags {
			if !strings.Contains(outputStr, flag) {
				t.Errorf("Expected %q in help output, got: %s", flag, outputStr)
			}
		}
	})
}

// TestCLICleanup tests the cleanup command functionality
func TestCLICleanup(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	// Create temp directory for test data
	tempDir, err := os.MkdirTemp("", "todo_cleanup_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Clean up any existing todos
	os.Remove(".todos.json")

	t.Run("cleanup_no_completed_items", func(t *testing.T) {
		// Add some incomplete items
		cmd := exec.Command(buildPath, "add", "Incomplete task 1")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		cmd = exec.Command(buildPath, "add", "Incomplete task 2")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		// Run cleanup with --force (to avoid interactive prompt)
		cmd = exec.Command(buildPath, "cleanup", "--force")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run cleanup: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "No completed items found to archive") {
			t.Errorf("Expected 'No completed items found' message, got: %s", outputStr)
		}
	})

	t.Run("cleanup_with_completed_items_force", func(t *testing.T) {
		// Add and complete some items
		cmd := exec.Command(buildPath, "add", "Completed task 1")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		cmd = exec.Command(buildPath, "add", "Completed task 2")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		cmd = exec.Command(buildPath, "add", "Incomplete task")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		// Complete first two tasks
		cmd = exec.Command(buildPath, "toggle", "3") // Complete first added task
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to toggle task: %v", err)
		}

		cmd = exec.Command(buildPath, "toggle", "4") // Complete second added task
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to toggle task: %v", err)
		}

		// Run cleanup with --force
		cmd = exec.Command(buildPath, "cleanup", "--force")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run cleanup: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Successfully archived 2 completed item(s)") {
			t.Errorf("Expected success message for 2 items, got: %s", outputStr)
		}

		// Verify archive file was created
		if _, err := os.Stat(".todos.archive.json"); os.IsNotExist(err) {
			t.Errorf("Archive file was not created")
		}

		// Verify remaining todos only contains incomplete items
		cmd = exec.Command(buildPath, "list", "--format", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list remaining todos: %v\nOutput: %s", err, output)
		}

		outputStr = string(output)
		if strings.Contains(outputStr, "Completed task 1") || strings.Contains(outputStr, "Completed task 2") {
			t.Errorf("Completed tasks should have been archived, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Incomplete task") {
			t.Errorf("Incomplete task should remain, got: %s", outputStr)
		}
	})

	t.Run("cleanup_with_list_flag", func(t *testing.T) {
		// Add and complete a task
		cmd := exec.Command(buildPath, "add", "Task for list test")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		cmd = exec.Command(buildPath, "toggle", "2") // Complete the newly added task
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to toggle task: %v", err)
		}

		// Run cleanup with --force and --list
		cmd = exec.Command(buildPath, "cleanup", "--force", "--list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run cleanup with --list: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		// Should contain both cleanup confirmation and list output
		if !strings.Contains(outputStr, "Successfully archived 1 completed item(s)") {
			t.Errorf("Expected cleanup confirmation, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "ID") || !strings.Contains(outputStr, "Task") {
			t.Errorf("Expected list table headers, got: %s", outputStr)
		}
	})

	t.Run("cleanup_without_args", func(t *testing.T) {
		cmd := exec.Command(buildPath, "cleanup")
		output, err := cmd.CombinedOutput()
		// This should work (no args required), but will be interactive
		// Since we can't test interactive mode easily, we test that it doesn't error on the command structure
		if err != nil && !strings.Contains(string(output), "Found") {
			t.Fatalf("Cleanup command should accept no arguments: %v\nOutput: %s", err, output)
		}
	})

	t.Run("cleanup_with_delete_flag", func(t *testing.T) {
		// Clean up any existing todos first
		os.Remove(".todos.json")
		os.Remove(".todos.archive.json")

		// Add and complete some items for delete testing
		cmd := exec.Command(buildPath, "add", "Delete test task 1")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		cmd = exec.Command(buildPath, "add", "Delete test task 2")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		cmd = exec.Command(buildPath, "add", "Keep this task")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		// Complete first two tasks (they should be IDs 1 and 2)
		cmd = exec.Command(buildPath, "toggle", "1")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to toggle task 1: %v", err)
		}

		cmd = exec.Command(buildPath, "toggle", "2")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to toggle task 2: %v", err)
		}

		// Run cleanup with --delete and --force
		cmd = exec.Command(buildPath, "cleanup", "--delete", "--force")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run cleanup --delete: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Successfully deleted 2 completed item(s)") {
			t.Errorf("Expected delete success message for 2 items, got: %s", outputStr)
		}

		// Verify remaining todos only contains incomplete items
		cmd = exec.Command(buildPath, "list", "--format", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list remaining todos: %v\nOutput: %s", err, output)
		}

		outputStr = string(output)
		if strings.Contains(outputStr, "Delete test task 1") || strings.Contains(outputStr, "Delete test task 2") {
			t.Errorf("Deleted tasks should not appear in list, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Keep this task") {
			t.Errorf("Incomplete task should remain, got: %s", outputStr)
		}

		// Verify that no archive file was created (since we deleted, not archived)
		if _, err := os.Stat(".todos.archive.json"); err == nil {
			// Archive file exists, but check if our deleted items are NOT in it
			archiveContent, err := os.ReadFile(".todos.archive.json")
			if err == nil {
				archiveStr := string(archiveContent)
				if strings.Contains(archiveStr, "Delete test task 1") || strings.Contains(archiveStr, "Delete test task 2") {
					t.Errorf("Deleted tasks should not be in archive file, got: %s", archiveStr)
				}
			}
		}
	})

	t.Run("help_contains_cleanup_command", func(t *testing.T) {
		cmd := exec.Command(buildPath, "help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to get help: %v\nOutput: %s", err, output)
		}

		expectedHelp := []string{
			"cleanup",
			"clean",
			"Archive or delete all completed todo items",
		}

		outputStr := string(output)
		for _, help := range expectedHelp {
			if !strings.Contains(outputStr, help) {
				t.Errorf("Expected %q in help output, got: %s", help, outputStr)
			}
		}
	})
}

// TestCLIGlobalCleanup tests cleanup command with global flag
func TestCLIGlobalCleanup(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	// Create temp directory and mock home directory
	tempDir, err := os.MkdirTemp("", "todo_global_cleanup_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create mock home directory
	mockHomeDir := filepath.Join(tempDir, "home")
	if err := os.MkdirAll(mockHomeDir, 0755); err != nil {
		t.Fatalf("Failed to create mock home directory: %v", err)
	}

	// Set HOME environment variable
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", mockHomeDir)
	defer os.Setenv("HOME", oldHome)

	// Also set USERPROFILE for Windows
	oldUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("USERPROFILE", mockHomeDir)
	defer os.Setenv("USERPROFILE", oldUserProfile)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	t.Run("setup_global_items", func(t *testing.T) {
		// Add some global items
		cmd := exec.Command(buildPath, "--global", "add", "Global completed item")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add global item: %v", err)
		}

		cmd = exec.Command(buildPath, "--global", "add", "Global incomplete item")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add global item: %v", err)
		}

		// Complete first item
		cmd = exec.Command(buildPath, "--global", "toggle", "1")
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to toggle global item: %v", err)
		}
	})

	t.Run("cleanup_global_items", func(t *testing.T) {
		cmd := exec.Command(buildPath, "--global", "cleanup", "--force")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to cleanup global items: %v\nOutput: %s", err, output)
		}

		expectedMsg := "Successfully archived 1 completed item(s)"
		if !strings.Contains(string(output), expectedMsg) {
			t.Errorf("Expected %q in output, got: %s", expectedMsg, output)
		}

		// Verify global archive file was created
		globalArchivePath := filepath.Join(mockHomeDir, ".todo", "todos.archive.json")
		if _, err := os.Stat(globalArchivePath); os.IsNotExist(err) {
			t.Errorf("Global archive file was not created at %s", globalArchivePath)
		}

		// Verify completed item was removed from global list
		cmd = exec.Command(buildPath, "--global", "list", "--format", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list global todos after cleanup: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if strings.Contains(outputStr, "Global completed item") {
			t.Errorf("Completed item still appears in global list after cleanup")
		}
		if !strings.Contains(outputStr, "Global incomplete item") {
			t.Errorf("Incomplete item should remain in global list")
		}
	})
}

// TestCLIArchiveFlag tests the --archive global flag functionality
func TestCLIArchiveFlag(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	// Create temp directory for test data
	tempDir, err := os.MkdirTemp("", "todo_archive_flag_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Clean up any existing files
	os.Remove(".todos.json")
	os.Remove(".todos.archive.json")

	t.Run("archive_flag_with_list_command", func(t *testing.T) {
		// Create some archive data
		archiveData := `[
			{"internal_id":"test1","task":"Archived task 1","completed":true,"created_at":"2025-08-08T00:00:00Z","updated_at":"2025-08-08T00:00:00Z","completed_at":"2025-08-08T00:00:00Z"},
			{"internal_id":"test2","task":"Archived task 2","completed":true,"created_at":"2025-08-08T00:00:00Z","updated_at":"2025-08-08T00:00:00Z","completed_at":"2025-08-08T00:00:00Z"}
		]`
		err := os.WriteFile(".todos.archive.json", []byte(archiveData), 0644)
		if err != nil {
			t.Fatalf("Failed to create archive file: %v", err)
		}

		// Test listing archive with --archive flag
		cmd := exec.Command(buildPath, "--archive", "list", "--format", "json")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run --archive list: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Archived task 1") || !strings.Contains(outputStr, "Archived task 2") {
			t.Errorf("Expected archived tasks in output, got: %s", outputStr)
		}
	})

	t.Run("archive_flag_with_delete_command", func(t *testing.T) {
		// Test deleting from archive
		cmd := exec.Command(buildPath, "--archive", "delete", "1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run --archive delete: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Deleted todo item with ID: 1") {
			t.Errorf("Expected delete confirmation, got: %s", output)
		}

		// Verify item was deleted
		cmd = exec.Command(buildPath, "--archive", "list", "--format", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list archive after delete: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if strings.Contains(outputStr, "Archived task 1") {
			t.Errorf("Deleted task should not appear in archive, got: %s", outputStr)
		}
		if !strings.Contains(outputStr, "Archived task 2") {
			t.Errorf("Remaining task should still be in archive, got: %s", outputStr)
		}
	})

	t.Run("archive_flag_blocks_unsupported_commands", func(t *testing.T) {
		// Test that add command is blocked
		cmd := exec.Command(buildPath, "--archive", "add", "Should fail")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected error when using --archive with add command")
		}

		if !strings.Contains(string(output), "--archive flag is only supported with 'list' and 'delete' commands") {
			t.Errorf("Expected validation error message, got: %s", output)
		}

		// Test that toggle command is blocked
		cmd = exec.Command(buildPath, "--archive", "toggle", "1")
		output, err = cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected error when using --archive with toggle command")
		}

		if !strings.Contains(string(output), "--archive flag is only supported with 'list' and 'delete' commands") {
			t.Errorf("Expected validation error message, got: %s", output)
		}

		// Test that edit command is blocked
		cmd = exec.Command(buildPath, "--archive", "edit", "1", "Should fail")
		output, err = cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected error when using --archive with edit command")
		}

		if !strings.Contains(string(output), "--archive flag is only supported with 'list' and 'delete' commands") {
			t.Errorf("Expected validation error message, got: %s", output)
		}

		// Test that archive command is blocked
		cmd = exec.Command(buildPath, "--archive", "archive", "1")
		output, err = cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected error when using --archive with archive command")
		}

		if !strings.Contains(string(output), "--archive flag is only supported with 'list' and 'delete' commands") {
			t.Errorf("Expected validation error message, got: %s", output)
		}

		// Test that cleanup command is blocked
		cmd = exec.Command(buildPath, "--archive", "cleanup")
		output, err = cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected error when using --archive with cleanup command")
		}

		if !strings.Contains(string(output), "--archive flag is only supported with 'list' and 'delete' commands") {
			t.Errorf("Expected validation error message, got: %s", output)
		}
	})

	t.Run("archive_flag_help_shows_flag", func(t *testing.T) {
		cmd := exec.Command(buildPath, "help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to get help: %v\nOutput: %s", err, output)
		}

		expectedFlags := []string{
			"--archive, -a",
			"Work with archive files instead of main todo list",
		}

		outputStr := string(output)
		for _, flag := range expectedFlags {
			if !strings.Contains(outputStr, flag) {
				t.Errorf("Expected %q in help output, got: %s", flag, outputStr)
			}
		}
	})

	t.Run("archive_flag_with_global_flag", func(t *testing.T) {
		// Create temp home directory
		mockHomeDir := filepath.Join(tempDir, "home")
		if err := os.MkdirAll(filepath.Join(mockHomeDir, ".todo"), 0755); err != nil {
			t.Fatalf("Failed to create mock home directory: %v", err)
		}

		// Set HOME environment variable
		oldHome := os.Getenv("HOME")
		oldUserProfile := os.Getenv("USERPROFILE")
		defer func() {
			os.Setenv("HOME", oldHome)
			os.Setenv("USERPROFILE", oldUserProfile)
		}()
		os.Setenv("HOME", mockHomeDir)
		os.Setenv("USERPROFILE", mockHomeDir)

		// Create global archive data
		globalArchiveData := `[
			{"internal_id":"global1","task":"Global archived task","completed":true,"created_at":"2025-08-08T00:00:00Z","updated_at":"2025-08-08T00:00:00Z","completed_at":"2025-08-08T00:00:00Z"}
		]`
		globalArchivePath := filepath.Join(mockHomeDir, ".todo", "todos.archive.json")
		err := os.WriteFile(globalArchivePath, []byte(globalArchiveData), 0644)
		if err != nil {
			t.Fatalf("Failed to create global archive file: %v", err)
		}

		// Test listing global archive with --global --archive flags
		cmd := exec.Command(buildPath, "--global", "--archive", "list", "--format", "json")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run --global --archive list: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Global archived task") {
			t.Errorf("Expected global archived task in output, got: %s", outputStr)
		}

		// Test deleting from global archive
		cmd = exec.Command(buildPath, "--global", "--archive", "delete", "1")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run --global --archive delete: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Deleted todo item with ID: 1") {
			t.Errorf("Expected delete confirmation, got: %s", output)
		}

		// Verify item was deleted from global archive
		cmd = exec.Command(buildPath, "--global", "--archive", "list", "--format", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list global archive after delete: %v\nOutput: %s", err, output)
		}

		outputStr = string(output)
		if strings.Contains(outputStr, "Global archived task") {
			t.Errorf("Deleted task should not appear in global archive, got: %s", outputStr)
		}
	})

	t.Run("archive_flag_with_empty_archive", func(t *testing.T) {
		// Remove archive file
		os.Remove(".todos.archive.json")

		// Test listing empty archive
		cmd := exec.Command(buildPath, "--archive", "list")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to run --archive list with empty archive: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "No todos found") {
			t.Errorf("Expected 'No todos found' message for empty archive, got: %s", outputStr)
		}

		// Test deleting from empty archive
		cmd = exec.Command(buildPath, "--archive", "delete", "1")
		output, err = cmd.CombinedOutput()
		if err == nil {
			t.Errorf("Expected error when deleting from empty archive")
		}

		if !strings.Contains(string(output), "failed to delete task: invalid index: 0") {
			t.Errorf("Expected invalid index error message, got: %s", output)
		}
	})
}

// TestCLIFilterFlag tests the --filter flag functionality for list command
func TestCLIFilterFlag(t *testing.T) {
	// Build the CLI for testing
	buildPath := filepath.Join(t.TempDir(), "todo.exe")

	cmd := exec.Command("go", "build", "-o", buildPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build CLI: %v", err)
	}

	tempDir, err := os.MkdirTemp("", "todo_filter_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	t.Run("filter_mixed_completion_status", func(t *testing.T) {
		// Clean up any existing todos file
		os.Remove(".todos.json")

		// Add multiple tasks
		tasks := []string{"Task 1", "Task 2", "Task 3", "Task 4"}
		for _, task := range tasks {
			cmd := exec.Command(buildPath, "add", task)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to add task '%s': %v\nOutput: %s", task, err, output)
			}
		}

		// Toggle some tasks to completed
		for _, id := range []string{"1", "3"} {
			cmd := exec.Command(buildPath, "toggle", id)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to toggle task %s: %v\nOutput: %s", id, err, output)
			}
		}

		// List all tasks (without filter)
		cmd := exec.Command(buildPath, "list", "--format", "json")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list all tasks: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		allTaskCount := strings.Count(outputStr, `"task":`)
		if allTaskCount != 4 {
			t.Errorf("Expected 4 tasks in unfiltered list, got %d: %s", allTaskCount, outputStr)
		}

		// List with filter (should only show incomplete tasks)
		cmd = exec.Command(buildPath, "list", "--filter", "--format", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list filtered tasks: %v\nOutput: %s", err, output)
		}

		filteredOutputStr := string(output)
		filteredTaskCount := strings.Count(filteredOutputStr, `"task":`)
		if filteredTaskCount != 2 {
			t.Errorf("Expected 2 tasks in filtered list, got %d: %s", filteredTaskCount, filteredOutputStr)
		}

		// Verify the correct tasks remain (Task 2 and Task 4 should be incomplete)
		if !strings.Contains(filteredOutputStr, "Task 2") {
			t.Errorf("Filtered list should contain 'Task 2': %s", filteredOutputStr)
		}
		if !strings.Contains(filteredOutputStr, "Task 4") {
			t.Errorf("Filtered list should contain 'Task 4': %s", filteredOutputStr)
		}

		// Verify completed tasks are not shown
		if strings.Contains(filteredOutputStr, "Task 1") {
			t.Errorf("Filtered list should not contain completed 'Task 1': %s", filteredOutputStr)
		}
		if strings.Contains(filteredOutputStr, "Task 3") {
			t.Errorf("Filtered list should not contain completed 'Task 3': %s", filteredOutputStr)
		}
	})

	t.Run("filter_all_completed", func(t *testing.T) {
		// Clean up any existing todos file
		os.Remove(".todos.json")

		// Add tasks and complete all of them
		tasks := []string{"Completed Task 1", "Completed Task 2"}
		for _, task := range tasks {
			cmd := exec.Command(buildPath, "add", task)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to add task '%s': %v\nOutput: %s", task, err, output)
			}
		}

		// Complete all tasks
		for _, id := range []string{"1", "2"} {
			cmd := exec.Command(buildPath, "toggle", id)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to toggle task %s: %v\nOutput: %s", id, err, output)
			}
		}

		// List with filter (should show no tasks)
		cmd := exec.Command(buildPath, "list", "--filter", "--format", "table")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list filtered tasks: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "No todos found") {
			t.Errorf("Expected 'No todos found' message when all tasks are completed, got: %s", outputStr)
		}
	})

	t.Run("filter_none_completed", func(t *testing.T) {
		// Clean up any existing todos file
		os.Remove(".todos.json")

		// Add tasks but don't complete any
		tasks := []string{"Incomplete Task 1", "Incomplete Task 2"}
		for _, task := range tasks {
			cmd := exec.Command(buildPath, "add", task)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Failed to add task '%s': %v\nOutput: %s", task, err, output)
			}
		}

		// List with filter (should show all tasks)
		cmd := exec.Command(buildPath, "list", "--filter", "--format", "json")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list filtered tasks: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		taskCount := strings.Count(outputStr, `"task":`)
		if taskCount != 2 {
			t.Errorf("Expected 2 tasks in filtered list when none completed, got %d: %s", taskCount, outputStr)
		}

		// Verify both tasks are shown
		if !strings.Contains(outputStr, "Incomplete Task 1") {
			t.Errorf("Filtered list should contain 'Incomplete Task 1': %s", outputStr)
		}
		if !strings.Contains(outputStr, "Incomplete Task 2") {
			t.Errorf("Filtered list should contain 'Incomplete Task 2': %s", outputStr)
		}
	})

	t.Run("filter_with_different_formats", func(t *testing.T) {
		// Clean up any existing todos file
		os.Remove(".todos.json")

		// Add mixed tasks
		cmd := exec.Command(buildPath, "add", "Format Test Complete")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to add task: %v\nOutput: %s", err, output)
		}

		cmd = exec.Command(buildPath, "add", "Format Test Incomplete")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to add task: %v\nOutput: %s", err, output)
		}

		// Complete first task
		cmd = exec.Command(buildPath, "toggle", "1")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to toggle task: %v\nOutput: %s", err, output)
		}

		// Test filter with table format
		cmd = exec.Command(buildPath, "list", "--filter", "--format", "table")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list filtered tasks (table): %v\nOutput: %s", err, output)
		}

		tableOutput := string(output)
		if !strings.Contains(tableOutput, "Format Test Incomplete") {
			t.Errorf("Table filtered list should contain incomplete task: %s", tableOutput)
		}
		if strings.Contains(tableOutput, "Format Test Complete") {
			t.Errorf("Table filtered list should not contain completed task: %s", tableOutput)
		}

		// Test filter with pretty format
		cmd = exec.Command(buildPath, "list", "--filter", "--format", "pretty")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list filtered tasks (pretty): %v\nOutput: %s", err, output)
		}

		prettyOutput := string(output)
		if !strings.Contains(prettyOutput, "Format Test Incomplete") {
			t.Errorf("Pretty filtered list should contain incomplete task: %s", prettyOutput)
		}
		if strings.Contains(prettyOutput, "Format Test Complete") {
			t.Errorf("Pretty filtered list should not contain completed task: %s", prettyOutput)
		}

		// Test filter with json format
		cmd = exec.Command(buildPath, "list", "--filter", "--format", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list filtered tasks (json): %v\nOutput: %s", err, output)
		}

		jsonOutput := string(output)
		if !strings.Contains(jsonOutput, "Format Test Incomplete") {
			t.Errorf("JSON filtered list should contain incomplete task: %s", jsonOutput)
		}
		if strings.Contains(jsonOutput, "Format Test Complete") {
			t.Errorf("JSON filtered list should not contain completed task: %s", jsonOutput)
		}
	})

	t.Run("filter_with_global_flag", func(t *testing.T) {
		// Clean up any existing files
		homeDir, _ := os.UserHomeDir()
		globalTodosPath := filepath.Join(homeDir, ".todo", "todos.json")
		os.Remove(globalTodosPath)

		// Add global tasks
		cmd := exec.Command(buildPath, "--global", "add", "Global Complete Task")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to add global task: %v\nOutput: %s", err, output)
		}

		cmd = exec.Command(buildPath, "--global", "add", "Global Incomplete Task")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to add global task: %v\nOutput: %s", err, output)
		}

		// Complete first global task
		cmd = exec.Command(buildPath, "--global", "toggle", "1")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to toggle global task: %v\nOutput: %s", err, output)
		}

		// List filtered global tasks
		cmd = exec.Command(buildPath, "--global", "list", "--filter", "--format", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to list filtered global tasks: %v\nOutput: %s", err, output)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Global Incomplete Task") {
			t.Errorf("Global filtered list should contain incomplete task: %s", outputStr)
		}
		if strings.Contains(outputStr, "Global Complete Task") {
			t.Errorf("Global filtered list should not contain completed task: %s", outputStr)
		}
	})
}
