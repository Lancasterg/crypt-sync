# go-crypt-sync v0.0.3 🔒

A CLI tool for securely syncing encrypted files to Google Cloud Storage.

## Why it exists

In an era where data privacy is paramount, trusting cloud providers with sensitive plaintext data is a risk. `go-crypt-sync` was built on the principle of "Zero Trust".

It ensures that your files—specifically sensitive configuration or credential files—are encrypted **before** they leave your machine. The cloud provider only ever sees encrypted binary blobs. Decryption occurs locally on your trusted machine, in memory, ensuring no unencrypted data is ever written to disk during the sync process unless explicitly requested.

The latest binary can be downloaded [here](https://github.com/Lancasterg/crypt-sync/blob/main/bin/go-crypt-sync)

## Flags
### Download

| Command | Flag | Short | Description | Required |
| :--- | :--- | :--- | :--- | :--- |
| `download` | `bucket-name` | `Positional` | The name of the bucket to download from | Yes |
| `download` | `file-name` | `Positional` | The name of the file to download | Yes |
| `download` | `--output` | `-o` | Specify an output file path. If no, outout will be writtent to stdout | No |
| `download` | `--decrypt` | `-d` | Decrypt the file after downloading | No |

### Encrypt
| Command | Flag | Short | Description | Required |
| :--- | :--- | :--- | :--- | :--- |
| `encrypt` | `input-file` | `Positional` | The name of the file to encrypt | Yes |
| `encrypt` | `output-file` | `Positional` | The name of the file to upload once encrypted | Yes |

### List
| Command | Flag | Short | Description | Required |
| :--- | :--- | :--- | :--- | :--- |
| `list` | `bucket-name` | `--bucket` | `The bucket to view the contents of` | Yes |

### Rm
| Command | Flag | Short | Description | Required |
| :--- | :--- | :--- | :--- | :--- |
| `rm` | `bucket-name` | `--bucket` | `The bucket to view the contents of` | Yes |

## Installation (for devs)

Ensure you have Actually Good Encryption (AGE) installed on your machine.

Arch (💯) users
```bash
$ sudo pacman -S age
```

Mac (🧪)
```bash
$ brew install age
```

Ubuntu / Debian (🕹️)
```bash
$ sudo apt install age
```

``` bash
# Clone the repo
$ git clone git@github.com:Lancasterg/crypt-sync.git

# Set your AGE_HOME dir
$ export AGE_HOME="$HOME/.config/age"

# Encrypt and upload your first file
$ go run main.go encrypt [local_file.json] [uploaded_file.enc]

# Download the file you just uploaded
$ go run main.go download [bucket-name] [file_name] [--output dev_tools/test123.json] [--decrypt true]

```
