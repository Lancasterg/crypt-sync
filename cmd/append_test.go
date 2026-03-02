package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAppendCmd(t *testing.T) {
	// Reset the state of the global flag variable before each run.
	// Because `recoveryPairs` is a global variable, its state can persist between
	// test runs if not explicitly reset.
	recoveryPairs = []string{}

	// Setup a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_passwords.json")

	// Initialize file with empty array
	initialData := []JSONEntry{}
	initBytes, _ := json.Marshal(initialData)
	if err := os.WriteFile(tmpFile, initBytes, 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Configure the command
	// We use SetArgs to simulate CLI arguments
	rootCmd.SetArgs([]string{
		"append",
		"Netflix",
		"user@example.com",
		"password123",
		"--file-name", tmpFile,
		"--recovery", "Pet:Dog",
	})

	// Capture stdout/stderr to keep test output clean
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Command failed: %v", err)
	}

	// Verify results
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var entries []JSONEntry
	err = json.Unmarshal(content, &entries)

	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	} else {
		if entries[0].Service != "Netflix" {
			t.Errorf("Expected service 'Netflix', got '%s'", entries[0].Service)
		}
		if len(entries[0].Recovery) != 1 {
			t.Errorf("Expected 1 recovery question")
		}
	}

}
