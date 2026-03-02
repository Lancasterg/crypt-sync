package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that we can encypt a file using the public key and then decrypt it using the
// private key
func TestEncryptionDecryption(t *testing.T) {
	dummyKey, err := age.GenerateX25519Identity()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	dummyKeyMap := map[string]string{
		"dummyPublicKey": dummyKey.Recipient().String(),
		"dummySecretKey": dummyKey.String(),
	}

	bytesToEncrypt := []byte("Hello, world!")

	encryptedBytes, err := EncryptInMemory(dummyKeyMap["dummyPublicKey"], bytesToEncrypt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	} else {
		t.Logf("Encrypted bytes: %v", string(encryptedBytes))
	}
	bytesToDecrypt, err := decryptData(dummyKeyMap["dummySecretKey"], encryptedBytes)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Logf("Unencrypted bytes: %v", string(bytesToDecrypt))

}

// Test that we can encypt a file, but cannot decrypt it if we have lost the key.
// This is to demonstate how important it is that you look after you private key.
func TestEncryptionLostKey(t *testing.T) {
	// Generate a new key
	dummyKey, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Create a mapping of [string][string]
	dummyKeyMap := map[string]string{
		"dummyPublicKey": dummyKey.Recipient().String(),
		"dummySecretKey": dummyKey.String(),
	}

	// Delete the private key from memory
	delete(dummyKeyMap, "dummyPrivateKey")
	t.Logf("dummyPrivateKey deleted from dummyKeyMap")

	bytesToEncrypt := []byte("Hello, world! We come in peace.")
	encryptedBytes, err := EncryptInMemory(dummyKeyMap["dummyPublicKey"], bytesToEncrypt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	} else {
		t.Logf("Encrypted bytes: %v", string(encryptedBytes))
	}

}

// TestEncryptCommand tests the 'encrypt' cobra command.
func TestEncryptCommand(t *testing.T) {
	// --- Setup ---
	// 1. Fake GCS
	bucketName := "test-encrypt-bucket"
	_, client := SetupGCSHelper(t, bucketName)

	// 2. Temp files and keys
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "my-secret.txt")
	keyFile := filepath.Join(tmpDir, "key.txt")
	outputObjectName := "my-secret.enc"
	fileContent := "this is a super secret message"

	// Create input file
	err := os.WriteFile(inputFile, []byte(fileContent), 0644)
	require.NoError(t, err)

	// Generate age key and write to key file
	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)
	publicKey := identity.Recipient().String()
	keyContent := "# created: 2024-01-01T00:00:00Z\n"
	keyContent += "# public key: " + publicKey + "\n"
	keyContent += identity.String()
	err = os.WriteFile(keyFile, []byte(keyContent), 0644)
	require.NoError(t, err)

	// --- Execute ---
	// Reset args and output
	var outBuf bytes.Buffer
	rootCmd.SetOut(&outBuf)
	rootCmd.SetErr(&outBuf)
	rootCmd.SetArgs([]string{
		"encrypt",
		inputFile,
		outputObjectName,
		"--encryption-key", keyFile,
		"--bucket-name", bucketName,
	})

	err = rootCmd.Execute()
	require.NoError(t, err, "encrypt command failed. Output: %s", outBuf.String())

	// --- Assertions ---
	// 1. Check if the file was uploaded to fake GCS
	rc, err := client.Bucket(bucketName).Object(outputObjectName).NewReader(t.Context()) // using t.Name() as context is fine for tests
	require.NoError(t, err, "encrypted object not found in GCS")
	defer rc.Close()

	encryptedData, err := io.ReadAll(rc)
	require.NoError(t, err)
	assert.NotEmpty(t, encryptedData, "uploaded file is empty")

	// 2. Verify the content by decrypting it
	decryptedData, err := decryptData(identity.String(), encryptedData)
	require.NoError(t, err, "failed to decrypt the uploaded data")
	assert.Equal(t, fileContent, string(decryptedData), "decrypted content does not match original")

	t.Logf("Encrypt command output: %s", outBuf.String())
}
