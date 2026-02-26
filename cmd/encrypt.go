/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt [input_file] [output_file] [key_name]",
	Args:  cobra.ExactArgs(3),
	Short: "Encrypt a file using the age master key",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ageHome := os.Getenv("AGE_HOME")

		if ageHome == "" {
			fmt.Println("AGE_HOME environment variable not set")
			os.Exit(1)
		}

		inputFile := args[0]
		outputFile := args[1]
		keyName := args[2]

		keyPath := ageHome + "/" + keyName

		fmt.Println("Reading key from:", keyPath)

		content, err := os.ReadFile(keyPath)

		if err != nil {
			fmt.Println("File not found:", keyPath)
			os.Exit(1)
		}

		keyContent := string(content)
		re := regexp.MustCompile(`# public key: (age1[a-z0-9]+)`)

		matches := re.FindStringSubmatch(keyContent)

		if len(matches) > 1 {

			publicKey := matches[1]
			fmt.Println("Extracted Public Key: ", publicKey)
			err := encryptFile(publicKey, inputFile, outputFile)

			if err != nil {
				formattedErr := fmt.Errorf("...: %w", err)
				fmt.Println(formattedErr)
				// os.Exit(1)
			}
		} else {
			fmt.Println("Public key not found in key file")
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().String("output", "")

}
