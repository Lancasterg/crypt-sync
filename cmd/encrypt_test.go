package cmd

import (
	"testing"

	"filippo.io/age"
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
