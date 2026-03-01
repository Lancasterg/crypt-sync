package cmd

import (
	"testing"
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
