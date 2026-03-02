package cmd

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestListCommand tests the 'list' cobra command.
func TestListCommand(t *testing.T) {
	// Setup
	ctx := context.Background()
	bucketName := "test-list-bucket"
	_, client := SetupGCSHelper(t, bucketName)

	// Create some dummy files in the fake GCS bucket
	filesToUpload := map[string]string{
		"file1.txt":     "content1",
		"file2.txt":     "content2",
		"dir/file3.txt": "content3",
	}

	for name, data := range filesToUpload {
		wc := client.Bucket(bucketName).Object(name).NewWriter(ctx)
		_, err := io.WriteString(wc, data)
		require.NoError(t, err, "failed to write data for "+name)
		err = wc.Close()
		require.NoError(t, err, "failed to close writer for "+name)
	}

	// Execute the 'list' command
	// Capture stdout to check the output
	var outBuf bytes.Buffer
	rootCmd.SetOut(&outBuf)
	rootCmd.SetErr(&outBuf) // Capture stderr as well for debugging

	// Reset flags to avoid pollution from other tests
	listCmd.ResetFlags()
	listCmd.Flags().StringP("bucket-name", "b", "encrypted-files-home", "Specify a bucket name (optional)")

	rootCmd.SetArgs([]string{"list", "--bucket-name", bucketName})

	err := listCmd.Execute()
	require.NoError(t, err, "list command failed")

	// Assertions
	output := outBuf.String()
	t.Logf("List command output:\n%s", output)

	// Check that all uploaded files are listed
	for name := range filesToUpload {
		assert.Contains(t, output, name)
	}
}
