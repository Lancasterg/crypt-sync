/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
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
	Short: "Encrypt a file using the age master key",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ageHome := os.Getenv("AGE_HOME")
		keyName := "master.txt"

		if ageHome == "" {
			log.Fatalf("AGE_HOME environment variable not set")
		}

		inputFile := args[0]
		outputFile := args[1]
		keyPath := ageHome + "/" + keyName

		fmt.Println("Reading key from:", keyPath)

		content, err := os.ReadFile(keyPath)

		if err != nil {
			fmt.Println("File not found:", keyPath)
			os.Exit(1)
		}

		keyContent := string(content)
		re := regexp.MustCompile(`# public key: (age[a-z0-9]+)`)
		matches := re.FindStringSubmatch(keyContent)

		inputFileRead, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatalf("Failed to read input file: %w", err)
		}

		if len(matches) > 1 {

			publicKey := matches[1]
			fmt.Println("Extracted Public Key: ", publicKey)
			encryptedBytes, err := EncryptInMemory(publicKey, inputFileRead)
			log.Println(encryptedBytes)

			if err != nil {
				log.Fatalf("%w", err)
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
}
