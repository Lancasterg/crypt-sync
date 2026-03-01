package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
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

type JsonEntry struct {
	Service  string     `json:"service"`
	Login    Login      `json:"login"`
	Recovery []Recovery `json:"recovery"`
}

func DownloadFromGCSBucket(bucketName string, objectName string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
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

	fmt.Printf("Downloaded %s from bucket %s\n", objectName, bucketName)
	// Process data as needed
	return data, nil
}

func UploadToGCSBucket(bucketName string, objectName string, data []byte) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
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

func NewJsonEntry(service string, username string, password string, recovery []Recovery) JsonEntry {
	return JsonEntry{
		Service:  service,
		Login:    Login{Username: username, Password: *secretstring.New(password)},
		Recovery: recovery,
	}
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

func EncryptInMemory(publicKey string, data []byte) ([]byte, error) {
	// 1. Parse the recipient
	recipient, err := age.ParseX25519Recipient(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipient: %w", err)
	}

	// 2. Prepare a buffer to hold the encrypted output
	out := &bytes.Buffer{}

	// 3. Set up the age writer
	// This wraps our 'out' buffer
	w, err := age.Encrypt(out, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to create age encryptor: %w", err)
	}

	// 4. Write the cleartext data to the age writer
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write data: %w", err)
	}

	// 5. CRITICAL: Close the age writer to finalize the MAC and flush bytes
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close age writer: %w", err)
	}

	// 6. Return the raw bytes from the buffer
	return out.Bytes(), nil
}
