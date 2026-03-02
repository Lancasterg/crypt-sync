/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a file from GCS",
	Long: `Delete a file from a GCS bucket. 
	We do not care if the file is encrypted or not, we simply want to delete it.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fileName, err := cmd.Flags().GetString("file-name")
		if err != nil || fileName == "" {
			return fmt.Errorf("file-name is required")
		}

		bucketName, err := cmd.Flags().GetString("bucket-name")
		if err != nil {
			return err
		}

		ctx := cmd.Context()

		client, err := storage.NewClient(ctx)
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer client.Close()

		err = client.Bucket(bucketName).Object(fileName).Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete object: %w", err)
		} else {
			log.Printf("Object %s deleted from bucket %s\n", fileName, bucketName)
			return nil

		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().StringP("file-name", "f", "", "Specify a file name (required)")
	rmCmd.Flags().StringP("bucket-name", "b", "encrypted-files-home", "Specify a bucket name (optional)")
}
