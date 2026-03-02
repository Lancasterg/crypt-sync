package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLIEndToEnd is an integration test for the entire CLI workflow.
func TestCLIEndToEnd(t *testing.T) {
	// --- Setup ---
	// 1. Fake GCS
	bucketName := "integration-test-bucket"
	_, client := SetupGCSHelper(t, bucketName)

	// 2. Temp files and keys
	tmpDir := t.TempDir()
	keyFile := filepath.Join(tmpDir, "master.txt")
	localFile := filepath.Join(tmpDir, "secrets.json")
	downloadedFile := filepath.Join(tmpDir, "downloaded_secrets.json")

	originalContent := `[{"some":"secret-data"}]`
	remoteObjectName := "secrets.enc"

	// 3. Create local file to upload
	err := os.WriteFile(localFile, []byte(originalContent), 0644)
	require.NoError(t, err)

	// 4. Generate age key and write to key file
	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)
	publicKey := identity.Recipient().String()
	keyContent := "# created: 2024-01-01T00:00:00Z\n"
	keyContent += "# public key: " + publicKey + "\n"
	keyContent += identity.String()
	err = os.WriteFile(keyFile, []byte(keyContent), 0644)
	require.NoError(t, err)

	// --- Test Workflow ---

	// Step 1: Encrypt and upload the file
	t.Run("EncryptAndUpload", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		encryptCmd.ResetFlags()
		encryptCmd.Flags().StringP("encryption-key", "k", "", "Specify an encryption key (optional)")
		encryptCmd.Flags().StringP("bucket-name", "b", "encrypted-files-home", "Specify a bucket name (optional)")

		rootCmd.SetArgs([]string{"encrypt", localFile, remoteObjectName, "--encryption-key", keyFile, "--bucket-name", bucketName})

		err := encryptCmd.Execute()
		require.NoError(t, err, "encrypt command failed. Output: %s", buf.String())
		assert.Contains(t, buf.String(), "Uploading file to")
	})

	// Step 2: List the bucket to verify upload
	t.Run("ListBucket", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		listCmd.ResetFlags()
		listCmd.Flags().StringP("bucket-name", "b", "encrypted-files-home", "Specify a bucket name (optional)")

		rootCmd.SetArgs([]string{"list", "--bucket-name", bucketName})

		err := listCmd.Execute()
		require.NoError(t, err, "list command failed. Output: %s", buf.String())
		output := buf.String()
		assert.Contains(t, output, remoteObjectName)
	})

	// Step 3: Download and decrypt the file
	t.Run("DownloadAndDecrypt", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		downloadCmd.ResetFlags()
		downloadCmd.Flags().StringP("output", "o", "", "Specify an output file path (optional)")
		downloadCmd.Flags().StringP("decryption-key", "k", "", "Specify a decryption key (optional)")
		downloadCmd.Flags().BoolVarP(&shouldDecrypt, "decrypt", "d", false, "Decrypt the file after downloading")

		rootCmd.SetArgs([]string{"download", bucketName, remoteObjectName, "--decryption-key", keyFile, "--decrypt=true", "--output", downloadedFile})

		err := downloadCmd.Execute()
		require.NoError(t, err, "download command failed. Output: %s", buf.String())
		assert.Contains(t, buf.String(), "File saved to: "+downloadedFile)

		// Verify content of downloaded file
		downloadedContent, err := os.ReadFile(downloadedFile)
		require.NoError(t, err)
		assert.Equal(t, originalContent, string(downloadedContent))
	})

	// Step 4: Remove the file from the bucket
	t.Run("RemoveFile", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		rmCmd.ResetFlags()
		rmCmd.Flags().StringP("file-name", "f", "", "Specify a file name (required)")
		rmCmd.Flags().StringP("bucket-name", "b", "encrypted-files-home", "Specify a bucket name (optional)")

		rootCmd.SetArgs([]string{"rm", "--bucket-name", bucketName, "--file-name", remoteObjectName})

		err := rmCmd.Execute()
		require.NoError(t, err, "rm command failed. Output: %s", buf.String())
		assert.Contains(t, buf.String(), "deleted from bucket")

		// Verify object is gone from GCS
		_, err = client.Bucket(bucketName).Object(remoteObjectName).NewReader(t.Context())
		assert.Error(t, err, "object should not exist after deletion")
	})

	// Step 5: List again to confirm removal
	t.Run("ListBucketAfterRemove", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd.SetOut(&buf)
		rootCmd.SetErr(&buf)

		listCmd.ResetFlags()
		listCmd.Flags().StringP("bucket-name", "b", "encrypted-files-home", "Specify a bucket name (optional)")
		rootCmd.SetArgs([]string{"list", "--bucket-name", bucketName})

		err := listCmd.Execute()
		require.NoError(t, err, "list command failed. Output: %s", buf.String())
		output := buf.String()
		assert.NotContains(t, output, remoteObjectName)
	})
}
