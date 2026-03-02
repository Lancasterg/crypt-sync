# go-crypt-sync v0.0.4 🔒

A CLI tool for securely syncing encrypted files to Google Cloud Storage.

## Why it exists

In an era where data privacy is paramount, trusting cloud providers with sensitive plaintext data is a risk. `go-crypt-sync` was built on the principle of "Zero Trust".

It ensures that your files—specifically sensitive configuration or credential files—are encrypted **before** they leave your machine. The cloud provider only ever sees encrypted binary blobs. Decryption occurs locally on your trusted machine, in memory, ensuring no unencrypted data is ever written to disk during the sync process unless explicitly requested.

The latest binary can be downloaded [here](https://github.com/Lancasterg/crypt-sync/blob/main/bin/go-crypt-sync)

## A word on keys 🔑

This tool uses [age](https://github.com/FiloSottile/age) for encryption. The security of your data depends entirely on your private key (identity file). 

1.  **Do not share your private key.**
2.  **Do not commit your private key to version control.**
3.  **Back up your private key securely.**

If you lose your private key, your data is irretrievable. If a malicious actor gains access to your private key, they can decrypt everything. To put it bluntly: if someone gets your private key, you are fucked. I would recommend putting a copy of your key file in Google Secret Manager (or a password manager like [proton pass](https://proton.me/pass/security), or maybe multiple **trusted** places). If you dropped your laptop tomorrow after encrypting half your files and uploading them to a GCS bucket (without storing the key elsewhere) you would really kick yourself, so please don't do it.

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
## TODO 
### Padding
Add padding to the encrypted json files 
Encryption hides the content, but it usually preserves the approximate length. If an attacker knows you are storing a "Status" JSON, they can guess the value based on the byte count:

`{"status": "active"} → 20 bytes`

`{"status": "terminated_for_cause"} → 34 bytes`
      
If the encrypted file is 128 bytes, the attacker knows it’s the "active" status. If it's 142 bytes, they know it's the longer string. In the context of a time-critical mission like a drone strike, this could reveal if a drone is "Idle," "Armed," or "Mission_Failed" just by looking at the metadata of the encrypted packet. JSON is "verbose" and highly structured. Because it uses keys like "velocity_x": and "battery_level":, the length of the values often falls into predictable patterns. Unlike binary formats, JSON's heavy use of whitespace and repetitive keys makes file-size analysis very effective for attackers.

Padding allows the user to round the file size up to a fixed increment (e.g., the nearest 1KB or 4KB block). Thus in the eyes of the defender, they cannot determine the status of the drone as the interceptd packets are all the same size.

`Original: 20 bytes → Padded: 1,024 bytes`

`Original: 34 bytes → Padded: 1,024 bytes`
```go
// Pseudocode for some kind of implementation
// Calculate how much padding is needed to hit a 1KB boundary
paddingSize := 1024 - (len(jsonData) % 1024)
paddedData := append(jsonData, make([]byte, paddingSize)...)
// Then encrypt 'paddedData'
```

### Fix all tests and reach code coverage of 90%
Speaks for itself really.

### Authenticated Encryption with Additional Data (AEAD)
Standard encryption hides data, but it doesn't always prevent an attacker from "flipping bits" in the ciphertext to change the message (e.g., changing a land: false to land: true).

The Idea: Use ChaCha20-Poly1305 (which age uses internally). It attaches a "MAC" (Message Authentication Code). If even one bit of the encrypted JSON is changed during transit, the decryption will fail entirely.

The Pro Move: Use the AAD (Additional Authenticated Data) field to bind the encrypted JSON to a specific Drone ID or Mission ID. This prevents "Replay Attacks" where an attacker captures a valid "Land" command and sends it back to the drone later.


### "Chaffing and Winnowing" (Traffic Camouflage)
Even with padding, an attacker can see when you are sending data. If they see a burst of 1KB packets every time the drone nears a target, they know you're taking photos.

The Idea: Send "Chaff"—fake, encrypted JSON files—at perfectly regular intervals (e.g., exactly every 500ms), regardless of whether you have real data to send.

How it works: 1.  The drone sends a 1KB packet every 0.5s.
2.  If there's real telemetry, it’s in the packet.
3.  If there's no news, the packet contains random noise encrypted to look like a real JSON.

Result: To an observer, the drone's "heartbeat" never changes, making it impossible to tell when it’s actually performing an action.


### Ephemeral Key Rotation (Perfect Forward Secrecy)
If you use the same age key for a month and the key is compromised, an attacker can decrypt all the files they recorded over that month.

The Idea: Generate a new key for every single flight session.

Implementation (drone): 
1.  Use a "Master Key" to exchange a "Session Key" at takeoff.
2.  Encrypt all in-flight JSONs with the Session Key.
3.  Delete the Session Key from the drone's RAM as soon as it lands (or if it detects it’s being tampered with).
4.  Now, even if the drone is captured, the recorded radio traffic from yesterday remains unreadable.

Implementation (Cloud Storage):
1. Use a master key to encrypt a file before sending to GCP
2. Use this key for a month
3. Once the key has expired, the private key can still be used to decrypt old messages, but the public key is now useless.
4. This gives a window of opportunity (once a month) to download and wipe the GCS bucket. A new private key can then be generated to re-encrypt the files.