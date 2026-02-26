package cmd

import (
	"fmt"
	"io"
	"os"

	"filippo.io/age"
)

func decryptFile(secretKey string, inputFile string) ([]byte, error) {
	identity, err := age.ParseX25519Identity(secretKey)
	if err != nil {
		return nil, err
	}

	in, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}

	r, err := age.Decrypt(in, identity)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(r)
}
func encryptFile(publicKey string, inputFile string, outputFile string) error {
	recipient, err := age.ParseX25519Recipient(publicKey)
	if err != nil {
		return fmt.Errorf("failed to parse recipient: %w", err)
	}

	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	// We wrap the closure in a bit of logic to catch errors on close
	defer out.Close()

	w, err := age.Encrypt(out, recipient)
	if err != nil {
		return err
	}

	// We MUST close 'w' to flush the final age block/MAC.
	// We use a defer with a named error or a manual close to be safe.
	closeErr := func() error {
		_, copyErr := io.Copy(w, in)
		if copyErr != nil {
			w.Close() // Close anyway to clean up, but return the copy error
			return copyErr
		}
		return w.Close()
	}()

	if closeErr != nil {
		return fmt.Errorf("encryption failed: %w", closeErr)
	}

	return nil
}
