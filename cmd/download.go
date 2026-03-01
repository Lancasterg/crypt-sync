/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var shouldDecrypt bool

func DecryptFile(keyPath string, fileContent string) ([]byte, error) {

	log.Println("Reading key from:", keyPath)

	content, err := os.ReadFile(keyPath)

	if err != nil {
		log.Fatalf("File not found: %s", keyPath)
	}

	keyContent := string(content)
	re := regexp.MustCompile(`(AGE-SECRET-KEY-[A-Z0-9]+)`)

	matches := re.FindStringSubmatch(keyContent)

	if len(matches) > 1 {

		privateKey := matches[1]
		bytes, err := decryptData(privateKey, []byte(fileContent))

		if err != nil {
			log.Fatalf("decryption failed: %v", err)
		} else {
			return bytes, nil
		}
	} else {
		log.Fatalf("Private key not found in key file")
	}
	return nil, err

}

var downloadCmd = &cobra.Command{
	Use:   "download [bucket-name] [file-name] [optional --decryption-key] [optional --output] [optional --decrypt (default: false)]",
	Short: "Download and optionally decrypt a GCS file",
	Run: func(cmd *cobra.Command, args []string) {
		bucketName := args[0]
		fileName := args[1]

		var finalData []byte

		ageHome := os.Getenv("AGE_HOME")
		keyPath, err := cmd.Flags().GetString("decryption-key")
		if keyPath != "" {
			// Do nothing as key has been found
		} else if ageHome != "" {
			keyPath = ageHome + "/master.txt"
		} else {
			log.Fatalf("Either --decryption-key or AGE_HOME environment variable not set")
		}

		// Download from GCS bucket
		file, err := DownloadFromGCSBucket(bucketName, fileName)
		if err != nil {
			log.Fatalf("Download failed: %v", err)
		}
		finalData = file

		// Optional Decrypt
		if shouldDecrypt {
			finalData, err = DecryptFile(keyPath, string(file))
			if err != nil {
				log.Fatalf("Decryption failed: %v", err)
			}

			// Validate JSON structure
			var entries []JsonEntry
			if err := json.Unmarshal(finalData, &entries); err == nil && len(entries) > 0 {
				log.Printf("Successfully decrypted service: %s\n", entries[0].Service)
			}
		}

		// Output Logic is the same regardless if decrypted or not
		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			log.Printf("%s", string(finalData))
		} else {
			if err := os.WriteFile(outputFile, finalData, 0644); err != nil {
				log.Fatalf("Failed to write file: %v", err)
			}
			log.Printf("File saved to: %s\n", outputFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("output", "o", "", "Specify an output file path (optional)")
	downloadCmd.Flags().StringP("decryption-key", "k", "", "Specify a decryption key (optional)")
	downloadCmd.Flags().BoolVarP(&shouldDecrypt, "decrypt", "d", false, "Decrypt the file after downloading")
}
