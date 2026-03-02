package cmd

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDownloadCommand tests the 'download' cobra command.
func TestDownloadCommand(t *testing.T) {
	// --- Setup ---
	ctx := context.Background()
	bucketName := "test-download-bucket"
	_, client := SetupGCSHelper(t, bucketName)

	tmpDir := t.TempDir()
	keyFile := filepath.Join(tmpDir, "key.txt")
	objectName := "my-download-secret.enc"
	originalContent := "this is a downloadable secret"

	// 1. Generate key
	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)
	publicKey := identity.Recipient().String()
	keyContent := "# created: 2024-01-01T00:00:00Z\n"
	keyContent += "# public key: " + publicKey + "\n"
	keyContent += identity.String()
	err = os.WriteFile(keyFile, []byte(keyContent), 0644)
	require.NoError(t, err)

	// 2. Encrypt data and upload it to fake GCS
	encryptedData, err := EncryptInMemory(publicKey, []byte(originalContent))
	require.NoError(t, err)

	wc := client.Bucket("test-data-encrypted").Object(objectName).NewWriter(ctx)
	_, err = wc.Write(encryptedData)
	require.NoError(t, err)
	err = wc.Close()
	require.NoError(t, err)

	// --- Test Cases ---
	t.Run("DownloadOnly", func(t *testing.T) {
		// Execute download command without decryption
		var outBuf bytes.Buffer
		rootCmd.SetOut(&outBuf)
		rootCmd.SetErr(&outBuf)

		downloadCmd.ResetFlags()
		downloadCmd.Flags().StringP("output", "o", "", "Specify an output file path (optional)")
		downloadCmd.Flags().StringP("decryption-key", "k", "", "Specify a decryption key (optional)")
		downloadCmd.Flags().BoolVarP(&shouldDecrypt, "decrypt", "d", false, "Decrypt the file after downloading")

		rootCmd.SetArgs([]string{
			"download",
			"test-data-encrypted",
			objectName,
		})

		err := downloadCmd.Execute()
		require.NoError(t, err, "download command failed. Output: %s", outBuf.String())

		// Assert output is the encrypted data
		assert.Contains(t, outBuf.String(), string(encryptedData), "output should be the raw encrypted data")
	})

	t.Run("DownloadAndDecryptToStdout", func(t *testing.T) {
		// Execute download command with decryption
		var outBuf bytes.Buffer
		rootCmd.SetOut(&outBuf)
		rootCmd.SetErr(&outBuf)

		downloadCmd.ResetFlags()
		downloadCmd.Flags().StringP("output", "o", "", "Specify an output file path (optional)")
		downloadCmd.Flags().StringP("decryption-key", "k", "", "Specify a decryption key (optional)")
		downloadCmd.Flags().BoolVarP(&shouldDecrypt, "decrypt", "d", false, "Decrypt the file after downloading")

		rootCmd.SetArgs([]string{
			"download",
			"test-data-encrypted",
			objectName,
			"--decryption-key", keyFile,
			"--decrypt=true", // Also set via arg
		})

		err := downloadCmd.Execute()
		require.NoError(t, err, "download --decrypt command failed. Output: %s", outBuf.String())

		// Assert output is the decrypted data
		assert.Contains(t, outBuf.String(), originalContent, "output should be the decrypted content")
	})

	t.Run("DownloadAndDecryptToFile", func(t *testing.T) {
		outputFile := filepath.Join(tmpDir, "decrypted.txt")

		// Execute download command with decryption to a file
		var outBuf bytes.Buffer
		rootCmd.SetOut(&outBuf)
		rootCmd.SetErr(&outBuf)

		downloadCmd.ResetFlags()
		downloadCmd.Flags().StringP("output", "o", "", "Specify an output file path (optional)")
		downloadCmd.Flags().StringP("decryption-key", "k", "", "Specify a decryption key (optional)")
		downloadCmd.Flags().BoolVarP(&shouldDecrypt, "decrypt", "d", false, "Decrypt the file after downloading")

		rootCmd.SetArgs([]string{"download", "test-data-encrypted", objectName, "--decryption-key", keyFile, "--decrypt=true", "--output", outputFile})

		err := downloadCmd.Execute()
		require.NoError(t, err, "download --decrypt --output command failed. Output: %s", outBuf.String())

		// Assert file content
		fileBytes, err := os.ReadFile(outputFile)
		require.NoError(t, err)
		assert.Equal(t, originalContent, string(fileBytes), "decrypted file content is wrong")
		assert.Contains(t, outBuf.String(), "File saved to: "+outputFile)
	})
}
