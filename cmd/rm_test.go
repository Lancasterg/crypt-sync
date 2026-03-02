/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
*/
package cmd

import (
	"context"
	"errors"
	"io"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRMCommand(t *testing.T) {
	// Setup
	ctx := context.Background()
	bucketName := "test-rm-bucket"
	objectName := "file-to-delete.txt"
	data := `{"delete": "me"}`

	_, client := SetupGCSHelper(t, bucketName)

	// 1. Upload a file to be deleted
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	_, err := io.WriteString(wc, data)
	require.NoError(t, err)
	err = wc.Close()
	require.NoError(t, err)

	// Verify file exists before deletion
	_, err = client.Bucket(bucketName).Object(objectName).Attrs(ctx)
	require.NoError(t, err, "file should exist before deletion")

	// 2. Execute the 'rm' command
	rmCmd.ResetFlags()
	rmCmd.Flags().StringP("file-name", "f", "", "Specify a file name (required)")
	rmCmd.Flags().StringP("bucket-name", "b", "encrypted-files-home", "Specify a bucket name (optional)")
	rootCmd.SetArgs([]string{"rm", "--bucket-name", bucketName, "--file-name", objectName})
	err = rmCmd.Execute()
	require.NoError(t, err, "rm command failed")

	// 3. Assert that the file is gone
	_, err = client.Bucket(bucketName).Object(objectName).Attrs(ctx)
	assert.Error(t, err, "expected an error when getting attributes of a deleted object")
	assert.True(t, errors.Is(err, storage.ErrObjectNotExist), "expected error to be ObjectNotExist")
}
