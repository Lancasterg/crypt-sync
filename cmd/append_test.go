package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendCmd(t *testing.T) {
	// Because `recoveryPairs` is a global variable, its state can persist between
	// test runs if not explicitly reset.
	// We run sub-tests to manage state and test different scenarios.

	t.Run("AppendToNewFile", func(t *testing.T) {
		recoveryPairs = []string{} // Reset global state
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test_passwords.json")

		// Initialize file with empty array
		initialData := []JSONEntry{}
		initBytes, _ := json.Marshal(initialData)
		err := os.WriteFile(tmpFile, initBytes, 0644)
		require.NoError(t, err)

		// Configure and execute the command
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)
		rootCmd.SetArgs([]string{
			"append",
			"Netflix",
			"user@example.com",
			"password123",
			"--file-name", tmpFile,
			"--recovery", "Pet:Dog",
			"--recovery", "City:London",
		})

		err = rootCmd.Execute()
		require.NoError(t, err, "Command failed: %v", err)

		// Verify results
		content, err := os.ReadFile(tmpFile)
		require.NoError(t, err)

		var entries []JSONEntry
		err = json.Unmarshal(content, &entries)
		require.NoError(t, err)

		require.Len(t, entries, 1)
		entry := entries[0]
		assert.Equal(t, "Netflix", entry.Service)
		assert.Equal(t, "user@example.com", entry.Login.Username)
		require.Len(t, entry.Recovery, 2)
		assert.Equal(t, "Pet", entry.Recovery[0].Question)
		assert.Equal(t, "Dog", entry.Recovery[0].Answer)
		assert.Equal(t, "City", entry.Recovery[1].Question)
		assert.Equal(t, "London", entry.Recovery[1].Answer)
	})

	t.Run("AppendToExistingFile", func(t *testing.T) {
		recoveryPairs = []string{} // Reset global state
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "test_passwords_existing.json")

		// Initialize file with one entry
		existingEntry, _ := NewJSONEntry("Google", "test@google.com", "oldpass", nil, nil)
		initialData := []JSONEntry{*existingEntry}
		initBytes, _ := json.Marshal(initialData)
		err := os.WriteFile(tmpFile, initBytes, 0644)
		require.NoError(t, err)

		// Configure and execute the command to append another entry
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)
		rootCmd.SetArgs([]string{
			"append",
			"Amazon",
			"test@amazon.com",
			"newpass",
			"--file-name", tmpFile,
		})

		err = rootCmd.Execute()
		require.NoError(t, err, "Command failed: %v", err)

		// Verify results
		content, err := os.ReadFile(tmpFile)
		require.NoError(t, err)

		var entries []JSONEntry
		err = json.Unmarshal(content, &entries)
		require.NoError(t, err)

		require.Len(t, entries, 2)
		assert.Equal(t, "Google", entries[0].Service)
		assert.Equal(t, "Amazon", entries[1].Service)
		assert.Equal(t, "test@amazon.com", entries[1].Login.Username)
	})
}
