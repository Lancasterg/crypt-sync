/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
*/
package cmd

import (
	"context"
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
	Run: func(cmd *cobra.Command, args []string) {
		bucketName, err := cmd.Flags().GetString("bucket-name")

		ctx := context.Background()

		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		fileName, err := cmd.Flags().GetString("file-name")
		if err != nil {
			log.Fatalf("Error getting file name: %v", err)
		}

		err = client.Bucket(bucketName).Object(fileName).Delete(ctx)
		if err != nil {
			log.Fatalf("Failed to delete object: %v", err)
		}

		log.Printf("Object %s deleted from bucket %s\n", fileName, bucketName)

	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().StringP("file-name", "f", "", "Specify a file name (required)")
	rmCmd.Flags().StringP("bucket-name", "b", "encrypted-files-home", "Specify a bucket name (optional)")
}
