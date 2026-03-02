package cmd

import (
	"testing"

	"cloud.google.com/go/storage"
	"github.com/fsouza/fake-gcs-server/fakestorage"
)

func SetupGCSHelper(t *testing.T, bucketName string) (*fakestorage.Server, *storage.Client) {
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
