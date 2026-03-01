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
	Use:   "encrypt [input_file] [output_file]",
	Args:  cobra.ExactArgs(2),
	Short: "Encrypt a file using the age master key and upload the encrypted file to a GCS bucket.",
	Long: `This key is typically found at /$HOME/.config/age"
	You can export the AGE_HOME environment variable to point to this location.
	By default, the tool looks for the master.txt file in the same directory.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		flagKey, _ := cmd.Flags().GetString("encryption-key")
		keyPath, err := GetDefaultKeyPath(flagKey)
		if err != nil {
			return err
		}
		inputFile := args[0]
		outputFile := args[1]

		log.Println("Reading key from:", keyPath)

		content, err := os.ReadFile(keyPath)

		if err != nil {
			return fmt.Errorf("file not found: %s", keyPath)
		}

		keyContent := string(content)
		re := regexp.MustCompile(`# public key: (age[a-z0-9]+)`)
		matches := re.FindStringSubmatch(keyContent)

		inputFileRead, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %v", err)
		}

		if len(matches) > 1 {
			publicKey := matches[1]

			encryptedBytes, err := EncryptInMemory(publicKey, inputFileRead)
			if err != nil {
				return err
			}

			// TODO: Make bucket name configurable via flag
			log.Printf("Uploading file to %s/%s", "encrypted-files-home", outputFile)
			err = UploadObject(cmd.Context(), "encrypted-files-home", outputFile, encryptedBytes)
			return err
		} else {
			return fmt.Errorf("public key not found in key file")
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringP("encryption-key", "k", "", "Specify an encryption key (optional)")
}
