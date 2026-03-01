/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
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

func DecryptFile(keyName string, fileContent string) ([]byte, error) {
	ageHome := os.Getenv("AGE_HOME")

	if ageHome == "" {
		fmt.Println("AGE_HOME environment variable not set")
		os.Exit(1)
	}

	keyPath := ageHome + "/" + keyName

	fmt.Println("Reading key from:", keyPath)

	content, err := os.ReadFile(keyPath)

	if err != nil {
		fmt.Println("File not found:", keyPath)
		os.Exit(1)
	}

	keyContent := string(content)
	re := regexp.MustCompile(`(AGE-SECRET-KEY-[A-Z0-9]+)`)

	matches := re.FindStringSubmatch(keyContent)

	if len(matches) > 1 {

		privateKey := matches[1]
		bytes, err := decryptData(privateKey, []byte(fileContent))

		if err != nil {
			formattedErr := fmt.Errorf("decryption failed: %w", err)
			fmt.Println(formattedErr)
			return nil, err

		} else {
			return bytes, nil
		}
	} else {
		fmt.Println("Private key not found in key file")
		os.Exit(1)
	}
	return nil, err

}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use: "download [bucket-name] [file-name]",
	// Args:  cobra.ExactArgs(2),
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download called")
		bucketName := args[0]
		fileName := args[1]

		// Download from GCS bucket
		file, err := DownloadFromGCSBucket(bucketName, fileName)
		if err != nil {
			log.Fatalf("Failed to download from GCS: %v", err)
		}

		// Decrypt (returns []byte)
		decryptedRaw, err := DecryptFile("master.txt", string(file))
		if err != nil {
			log.Fatalf("Decryption failed: %v", err)
		}

		// Skip Base64 Decoding entirely.
		// DecryptFile already gave us the "cleartext" bytes.
		var entries []JsonEntry
		err = json.Unmarshal(decryptedRaw, &entries)
		if err != nil {
			// If it fails here, print the string to see what you actually got
			log.Printf("Actual content: %s\n", string(decryptedRaw))
			log.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		// Success!
		if len(entries) > 0 {
			log.Printf("Service: %s\n", entries[0].Service)
			log.Printf("Password: %v\n", entries[0].Login.Password)
		}

		outputFile, err := cmd.Flags().GetString("output")
		if err != nil || outputFile == "" {
			log.Printf("No output flag specified, writing to stdout\n")
			fmt.Println(string(decryptedRaw))
		} else {
			err = os.WriteFile(outputFile, decryptedRaw, 0644)
			if err != nil {
				log.Fatalf("Failed to write to file: %v", err)
			}
			log.Printf("Wrote to file: %s\n", outputFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("output", "o", "", "Specify an output file path (optional)")

}
