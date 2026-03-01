/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use: "encrypt [input_file] [output_file]",
	// Args:  cobra.ExactArgs(2),
	Short: "Encrypt a file using the age master key and upload the encrypted file to a GCS bucket.",
	Long: `This key is typically found at /$HOME/.config/age"
	You can export the AGE_HOME environment variable to point to this location.
	By default, the tool looks for the master.txt file in the same directory.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		ageHome := os.Getenv("AGE_HOME")
		keyPath, err := cmd.Flags().GetString("encryption-key")
		if keyPath != "" {
			// Do nothing as key has been found
		} else if ageHome != "" {
			keyPath = ageHome + "/master.txt"
		} else {
			log.Fatalf("Either --encryption-key or AGE_HOME environment variable not set")
		}

		if ageHome == "" {
			log.Fatalf("AGE_HOME environment variable not set")
		}

		inputFile := args[0]
		outputFile := args[1]

		fmt.Println("Reading key from:", keyPath)

		content, err := os.ReadFile(keyPath)

		if err != nil {
			log.Fatalf("File not found: %s", keyPath)
		}

		keyContent := string(content)
		re := regexp.MustCompile(`# public key: (age[a-z0-9]+)`)
		matches := re.FindStringSubmatch(keyContent)

		inputFileRead, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatalf("Failed to read input file: %v", err)
		}

		if len(matches) > 1 {

			publicKey := matches[1]
			log.Println("Extracted Public Key: ", publicKey)
			encryptedBytes, err := EncryptInMemory(publicKey, inputFileRead)
			log.Println(encryptedBytes)

			if err != nil {
				log.Fatalf("%v", err)
			}

			// Write to GCS
			err = UploadToGCSBucket("encrypted-files-home", outputFile, encryptedBytes)

		} else {
			log.Fatalf("Public key not found in key file")
		}

	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringP("encryption-key", "k", "", "Specify an encryption key (optional)")
}
