/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [optional bucket-name]",
	Short: "List the contents of a GCS bucket",
	Long:  `List the contents of a GCS bucket. The default bucket is encrypted-files-home`,
	Run: func(cmd *cobra.Command, args []string) {

		bucketName, err := cmd.Flags().GetString("bucket-name")

		ctx := context.Background()

		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		// Set a timeout for the listing operation
		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()

		fmt.Printf("Listing objects in bucket: %s\n", bucketName)
		fmt.Println("------------------------------------------")

		// Get an iterator for the objects in the bucket
		// You can pass nil for the query to get everything
		it := client.Bucket(bucketName).Objects(ctx, nil)

		log.Printf("File name\t\t\t   | Size     \t| Created\n")

		for {
			// Advance the iterator
			attrs, err := it.Next()
			if err == iterator.Done {
				break // End of the list
			}
			if err != nil {
				log.Fatalf("Error iterating objects: %v", err)
			}

			// Print the object metadata
			log.Printf("%-30s | %d bytes\t| %s\n",
				attrs.Name,
				attrs.Size,
				attrs.Created.Format("2006-01-02 15:04"))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("bucket-name", "b", "encrypted-files-home", "Specify a bucket name (optional)")
}
