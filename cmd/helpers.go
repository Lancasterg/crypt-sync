/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
*/

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"filippo.io/age"
	"github.com/chmller/secretstring"
)

type Login struct {
	Username string                    `json:"service"`
	Password secretstring.SecretString `json:"password"`
}

type Recovery struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type JSONEntry struct {
	Service  string     `json:"service"`
	Login    Login      `json:"login"`
	Recovery []Recovery `json:"recovery"`
}

// NewJSONEntry Creates a new JSONEntry Object
// Currently not used, but will be used in future to add to json documents
// prior to encrypting
func NewJSONEntry(service string, username string, password string, recovery []Recovery) JSONEntry {
	return JSONEntry{
		Service:  service,
		Login:    Login{Username: username, Password: *secretstring.New(password)},
		Recovery: recovery,
	}
}

// GetDefaultKeyPath resolves the key path from the flag or environment variable
func GetDefaultKeyPath(flagVal string) (string, error) {
	if flagVal != "" {
		return flagVal, nil
	}
	ageHome := os.Getenv("AGE_HOME")
	if ageHome == "" {
		return "", fmt.Errorf("AGE_HOME environment variable not set")
	}
	return filepath.Join(ageHome, "master.txt"), nil
}

// DownloadObject downloads an encrypted file from a GCS bucket.
// This function does not write to disk, it simply downloads the file and stores the result
// in memory as bytes.
func DownloadObject(ctx context.Context, bucketName string, objectName string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %w", objectName, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	return data, nil
}

func UploadObject(ctx context.Context, bucketName string, objectName string, data []byte) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// NewWriter does not return an error immediately
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	// 1. Write the data to the writer
	if _, err = wc.Write(data); err != nil {
		return fmt.Errorf("Object(%q).Write: %w", objectName, err)
	}

	// 2. Crucial: You MUST close the writer to flush the buffer
	// and finalize the upload to GCS.
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Object(%q).Close: %w", objectName, err)
	}

	return nil

}

func decryptData(secretKey string, data []byte) ([]byte, error) {
	identity, err := age.ParseX25519Identity(secretKey)
	if err != nil {
		return nil, err
	}

	r, err := age.Decrypt(bytes.NewReader(data), identity)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(r)
}

// Encrypt a file in memory
// The reason for doing this is that we want to leave no trace of the unencrypted file
// On the host machine.
func EncryptInMemory(publicKey string, data []byte) ([]byte, error) {
	// Parse the recipient
	recipient, err := age.ParseX25519Recipient(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipient: %w", err)
	}

	// Prepare a buffer to hold the encrypted output
	out := &bytes.Buffer{}

	// Set up the age writer
	// This wraps our 'out' buffer
	w, err := age.Encrypt(out, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to create age encryptor: %w", err)
	}

	// Write the cleartext data to the age writer
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write data: %w", err)
	}

	// CRITICAL: Close the age writer to finalize the MAC and flush bytes
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close age writer: %w", err)
	}

	// Return the raw bytes from the buffer
	return out.Bytes(), nil
}
