/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var shouldDecrypt bool

func decryptFileContent(keyPath string, fileContent []byte) ([]byte, error) {
	log.Println("Reading key from:", keyPath)

	content, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("key file not found: %s", keyPath)
	}

	keyContent := string(content)
	re := regexp.MustCompile(`(AGE-SECRET-KEY-[A-Z0-9]+)`)

	matches := re.FindStringSubmatch(keyContent)
	if len(matches) > 1 {
		privateKey := matches[1]
		bytes, err := decryptData(privateKey, fileContent)
		if err != nil {
			return nil, fmt.Errorf("decryption failed: %w", err)
		}
		return bytes, nil
	}
	return nil, fmt.Errorf("private key not found in key file")
}

var downloadCmd = &cobra.Command{
	Use:   "download [bucket-name] [file-name] [optional --decryption-key] [optional --output] [optional --decrypt (default: false)]",
	Short: "Download and optionally decrypt a GCS file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		bucketName := args[0]
		fileName := args[1]

		var finalData []byte

		flagKey, _ := cmd.Flags().GetString("decryption-key")
		keyPath, err := GetDefaultKeyPath(flagKey)
		if err != nil {
			return err
		}

		// Download from GCS bucket
		finalData, err = DownloadObject(cmd.Context(), bucketName, fileName)
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
		log.Printf("Downloaded %s from bucket %s\n", fileName, bucketName)

		// Optional Decrypt
		if shouldDecrypt {
			finalData, err = decryptFileContent(keyPath, finalData)
			if err != nil {
				return err
			}

			// Validate JSON structure
			var entries []JSONEntry
			if err := json.Unmarshal(finalData, &entries); err == nil && len(entries) > 0 {
				log.Printf("Successfully decrypted file and parsed JSON\n")
			}
		}

		// Output Logic is the same regardless if decrypted or not
		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			log.Printf("%s", string(finalData))
		} else {
			if err := os.WriteFile(outputFile, finalData, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			log.Printf("File saved to: %s\n", outputFile)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("output", "o", "", "Specify an output file path (optional)")
	downloadCmd.Flags().StringP("decryption-key", "k", "", "Specify a decryption key (optional)")
	downloadCmd.Flags().BoolVarP(&shouldDecrypt, "decrypt", "d", false, "Decrypt the file after downloading")
}
