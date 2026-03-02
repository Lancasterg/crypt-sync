/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
*/
package cmd

import (
	"context"
	"io"
	"testing"
)

func TestRM(t *testing.T) {
	ctx := context.Background()

	bucketName := "test-data-encrypted"
	objectName := "file.json"
	data := `{"hello": "world"}` // Use standard double quotes for valid JSON

	_, client := setupGCS(t, "test-data-encrypted")

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err := io.WriteString(wc, data); err != nil {
		t.Fatalf("failed to write data: %v", err)
	}
	if err := wc.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}
	if err := client.Bucket(bucketName).Object(objectName).Delete(ctx); err != nil {

		if err != nil {
			t.Fatalf("failed to delete object: %v", err)
		} else {
			t.Logf("Object %s deleted from bucket %s\n", objectName, bucketName)
		}
	}

}
