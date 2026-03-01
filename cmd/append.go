/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var recoveryPairs []string

// appendCmd represents the append command
var appendCmd = &cobra.Command{
	Use:   "append [service] [username] [password]",
	Short: "Append a new json object to the list of objects in a file.",
	Long: `Append a new json object to the list of objects in a file. 
eg: go run main.go append SomeService my-username mypassword --file-name=dev_tools/my-passwords.json --recovery=hello:goodbye --recovery=smile:yay
go run main.go append [service_name] [username] [password] --recovery=question:answer
	`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {

		filePath, err := cmd.Flags().GetString("file-name")
		if err != nil {
			return err
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("file not found: %s", filePath)
		}

		var entries []JSONEntry
		if err := json.Unmarshal(content, &entries); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
		log.Printf("Successfully loaded %d entries from %s", len(entries), filePath)

		service := args[0]
		username := args[1]
		password := args[2]

		var recoveryItems []Recovery
		for _, pair := range recoveryPairs {
			parts := strings.Split(pair, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid recovery format '%s'. Use 'Question:Answer'", pair)
			}
			recoveryItems = append(recoveryItems, Recovery{
				Question: parts[0],
				Answer:   parts[1],
			})
		}

		jsonEntry, err := NewJSONEntry(service, username, password, recoveryItems, nil)
		if err != nil {
			return fmt.Errorf("failed to create entry: %w", err)
		}
		entries = append(entries, *jsonEntry)

		jsonBytes, err := json.MarshalIndent(entries, "", "\t")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}

		if err := os.WriteFile(filePath, jsonBytes, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(appendCmd)
	appendCmd.Flags().StringP("file-name", "f", "", "Specify a file name (required)")
	_ = appendCmd.MarkFlagRequired("file-name")
	appendCmd.Flags().StringSliceVarP(&recoveryPairs, "recovery", "r", []string{}, "Recovery pairs in 'Question:Answer' format")
}
