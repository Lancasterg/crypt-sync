package cmd

import (
	"os"
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJSONEntry(t *testing.T) {
	t.Run("SuccessfulCreation", func(t *testing.T) {
		recovery := []Recovery{{Question: "q", Answer: "a"}}
		entry, err := NewJSONEntry("SomeService", "my-username", "mypassword", recovery, nil)
		require.NoError(t, err)

		assert.Equal(t, "SomeService", entry.Service)
		assert.Equal(t, "my-username", entry.Login.Username)
		assert.Equal(t, "mypassword", entry.Login.Password)
		assert.Equal(t, recovery, entry.Recovery)
		assert.NotZero(t, entry.MetaData.CreatedAt)
		assert.NotZero(t, entry.MetaData.UpdatedAt)
		assert.Equal(t, 1, entry.MetaData.VersionNumber)
		assert.NotZero(t, entry.MetaData.FingerPrint)
	})

	t.Run("FingerprintMatching", func(t *testing.T) {
		jsonEntry1, err := NewJSONEntry("SomeService", "my-username", "mypassword", nil, nil)
		require.NoError(t, err)

		jsonEntry2, err := NewJSONEntry("SomeService", "my-username", "mypassword", nil, nil)
		require.NoError(t, err)

		// Timestamps will differ, so we compare the fingerprints which should be identical
		// for identical sensitive data.
		assert.Equal(t, jsonEntry1.MetaData.FingerPrint, jsonEntry2.MetaData.FingerPrint)
	})

	t.Run("UpdateEntrySuccess", func(t *testing.T) {
		// Create initial entry
		initialEntry, err := NewJSONEntry("Service", "User", "Pass", nil, nil)
		require.NoError(t, err)

		// Simulate an update
		updatedEntry, err := NewJSONEntry("Service", "User", "Pass", nil, &initialEntry.MetaData)
		require.NoError(t, err)

		assert.Equal(t, initialEntry.MetaData.CreatedAt, updatedEntry.MetaData.CreatedAt)
		assert.True(t, updatedEntry.MetaData.UpdatedAt.After(initialEntry.MetaData.UpdatedAt))
		assert.Equal(t, initialEntry.MetaData.VersionNumber+1, updatedEntry.MetaData.VersionNumber)
	})

	t.Run("UpdateEntryFingerprintMismatch", func(t *testing.T) {
		entry, err := NewJSONEntry("Service", "User", "Pass", nil, nil)
		require.NoError(t, err)

		_, err = NewJSONEntry("Service", "User", "DifferentPass", nil, &entry.MetaData)
		expectedErr := "fingerprints do not match - file may have been maliciously edited."
		assert.EqualError(t, err, expectedErr)
	})
}

func TestGetDefaultKeyPath(t *testing.T) {
	t.Run("FromFlagValue", func(t *testing.T) {
		flagVal := "/custom/path/key.txt"
		value, err := GetDefaultKeyPath(flagVal)
		require.NoError(t, err)
		assert.Equal(t, flagVal, value)
	})

	t.Run("FromEnvVar", func(t *testing.T) {
		t.Setenv("AGE_HOME", "/path/to/.config/age")
		value, err := GetDefaultKeyPath("")
		require.NoError(t, err)
		assert.Equal(t, "/path/to/.config/age/master.txt", value)
	})

	t.Run("NoFlagNoEnvVar", func(t *testing.T) {
		// Unset AGE_HOME if it's set in the environment
		originalAgeHome, isSet := os.LookupEnv("AGE_HOME")
		if isSet {
			os.Unsetenv("AGE_HOME")
		}
		defer func() {
			if isSet {
				os.Setenv("AGE_HOME", originalAgeHome)
			}
		}()

		_, err := GetDefaultKeyPath("")
		assert.Error(t, err)
		assert.EqualError(t, err, "AGE_HOME environment variable not set")
	})
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	// Keys for unit tests are ephemeral (Generated on the fly)
	dummyKey, err := age.GenerateX25519Identity()
	require.NoError(t, err)
	dummyPublicKey := dummyKey.Recipient().String()
	dummySecretKey := dummyKey.String()

	bytesToEncrypt := []byte("Hello, world! This is a roundtrip test.")

	// Encrypt
	encryptedBytes, err := EncryptInMemory(dummyPublicKey, bytesToEncrypt)
	require.NoError(t, err)
	require.NotEmpty(t, encryptedBytes)

	// Decrypt
	unencryptedBytes, err := decryptData(dummySecretKey, encryptedBytes)
	require.NoError(t, err)

	assert.Equal(t, string(bytesToEncrypt), string(unencryptedBytes))
}
