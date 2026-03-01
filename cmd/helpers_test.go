package cmd

import (
	"testing"

	"filippo.io/age"
)

func TestCreateNewJsonObject(t *testing.T) {
	jsonEntry1, err := NewJSONEntry("SomeService", "my-username", "mypassword", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	jsonEntry2, err := NewJSONEntry("SomeService", "my-username", "mypassword", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Timestamps will differ, so we compare the fingerprints
	if jsonEntry1.MetaData.FingerPrint != jsonEntry2.MetaData.FingerPrint {
		t.Errorf("Fingerprints do not match")
	}
}

func TestNewJSONEntryFail(t *testing.T) {
	JSONEntry1, err := NewJSONEntry("SomeService", "username", "password", []Recovery{{Question: "question", Answer: "answer"}}, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	JSONEntry2, err := NewJSONEntry("SomeService", "username", "password", []Recovery{{Question: "question", Answer: "answer"}}, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if JSONEntry1.MetaData.FingerPrint != JSONEntry2.MetaData.FingerPrint {
		t.Fatalf("Fingerprints should match")
	}
	if JSONEntry1.Login.Username != JSONEntry2.Login.Username {
		t.Fatalf("Usernames should match")
	}
	if JSONEntry1.Login.Password != JSONEntry2.Login.Password {
		t.Fatalf("Passwords should match")
	}

}

func TestGetDefaultKeyPathFromFilename(t *testing.T) {
	flagVal := "hello"
	value, err := GetDefaultKeyPath(flagVal)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	} else if value != flagVal {
		t.Errorf("Expected %s, got %s", flagVal, value)
	} else if value != "hello" {
		t.Errorf("Expected %s, got %s", "hello", value)
	} else {
		t.Logf("Test passed")
	}

}

func TestGetDefaultKeyPathFromEnvVar(t *testing.T) {
	t.Setenv("AGE_HOME", "/path/to/.config/age")
	value, err := GetDefaultKeyPath("")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	} else if value != "/path/to/.config/age/master.txt" {
		t.Errorf("Expected %s, got %s", "/path/to/.config/age/master.txt", value)
	} else {
		t.Logf("Test passed")
	}

}

func TestNewJSONEntryFingerprintMismatch(t *testing.T) {
	entry, err := NewJSONEntry("Service", "User", "Pass", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = NewJSONEntry("Service", "User", "DifferentPass", nil, &entry.MetaData)
	expectedErr := "fingerprints do not match - file may have been maliciously edited."

	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %v", expectedErr, err)
	}
}

func TestEncryptInMemory(t *testing.T) {
	// Keys for unit tests are ephemeral (Generated on the fly)
	// This means that there is a low likelyhood of them being stolen
	// And even if they are, they aren't able to be used to identify anyone.
	dummyKey, err := age.GenerateX25519Identity()
	dummyPublicKey := dummyKey.Recipient().String()
	dummySecretKey := dummyKey.String()

	bytesToEncrypt := []byte("Hello, world!")
	encryptedBytes, err := EncryptInMemory(dummyPublicKey, bytesToEncrypt)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	unencryptedBytes, err := decryptData(dummySecretKey, encryptedBytes)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Logf("Unencrypted bytes: %v", string(unencryptedBytes))

}
