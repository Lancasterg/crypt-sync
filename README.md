# go-crypt-sync
v0.0.1
# go-crypt-sync 🔒

## MVP
A CLI tool for securely syncing encrypted files to Google Cloud Storage.

## Why it exists

In an era where data privacy is paramount, trusting cloud providers with sensitive plaintext data is a risk. `go-crypt-sync` was built on the principle of "Zero Trust".

It ensures that your files—specifically sensitive configuration or credential files—are encrypted **before** they leave your machine. The cloud provider only ever sees encrypted binary blobs. Decryption occurs locally on your trusted machine, in memory, ensuring no unencrypted data is ever written to disk during the sync process unless explicitly requested.

## Installation

To install the binary directly:

```bash
go install github.com/lancaster0180/go-crypt-sync@latest
