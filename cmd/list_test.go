package cmd

import (
	"context"
	"fmt"
	"io"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/fsouza/fake-gcs-server/fakestorage"
)

func setupGCS(t *testing.T, bucketName string) (*fakestorage.Server, *storage.Client) {
	t.Helper() // Marks this func as a helper for clearer error logs

	server := fakestorage.NewServer([]fakestorage.Object{})

	server.CreateBucketWithOpts(fakestorage.CreateBucketOpts{
		Name:                  "test-data-encrypted",
		VersioningEnabled:     false,
		DefaultEventBasedHold: false,
	})

	// t.Cleanup runs automatically at the end of the test that calls this helper
	t.Cleanup(func() {
		server.Stop()
	})

	return server, server.Client()
}

func TestRead(t *testing.T) {
	ctx := context.Background()

	bucketName := "test-data-encrypted"
	objectName := "file.json"
	data := `{"hello": "world"}` // Use standard double quotes for valid JSON

	_, client := setupGCS(t, bucketName)

	// --- STEP 1: UPLOAD THE DATA ---
	// You must write the file to the fake bucket before you can read it.
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.WriteString(wc, data); err != nil {
		t.Fatalf("failed to write data: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	// --- STEP 2: READ THE DATA ---
	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		// Use t.Fatalf instead of log.Fatalf in tests so the test fails gracefully
		t.Fatalf("failed to create reader: %v", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("failed to read content: %v", err)
	}

	fmt.Printf("Successfully read from fake GCS: %s\n", string(content))

	// Optional: Add an assertion to make it a real test
	if string(content) != data {
		t.Errorf("got %s, want %s", string(content), data)
	}
}

func TestList(t *testing.T) {
	ctx := context.Background()

	bucketName := "test-data-encrypted"
	objectName := "file.json"
	data := `{"hello": "world"}` // Use standard double quotes for valid JSON

	_, client := setupGCS(t, bucketName)

	// You must write the file to the fake bucket before you can read it.
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.WriteString(wc, data); err != nil {
		t.Fatalf("failed to write data: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	// Read the data
	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		// Use t.Fatalf instead of log.Fatalf in tests so the test fails gracefully
		t.Fatalf("failed to create reader: %v", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("failed to read content: %v", err)
	}

	fmt.Printf("Successfully read from fake GCS: %s\n", string(content))

	// Optional: Add an assertion to make it a real test
	if string(content) != data {
		t.Errorf("got %s, want %s", string(content), data)
	}
}
