/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
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

func ListFiles(bucketName string, ctx context.Context) ([]string, error) {

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
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
	resultStore := []string{}
	for {
		// Advance the iterator
		attrs, err := it.Next()
		if err == iterator.Done {
			break // End of the list
		} else if err != nil {
			err := fmt.Errorf("error iterating objects: %w", err)
			return nil, err
		}

		// Print the object metadata
		metadataResultStorefmt := fmt.Sprintf("%-30s | %d bytes\t| %s\n",
			attrs.Name,
			attrs.Size,
			attrs.Created.Format("2006-01-02 15:04"),
		)

		resultStore = append(resultStore, metadataResultStorefmt)
	}
	return resultStore, nil
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [optional bucket-name]",
	Short: "List the contents of a GCS bucket",
	Long:  `List the contents of a GCS bucket. The default bucket is encrypted-files-home`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		bucketName, err := cmd.Flags().GetString("bucket-name")
		if err != nil {
			return err
		} else {
			_, err := ListFiles(bucketName, ctx)
			if err != nil {
				return err
			}

			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("bucket-name", "b", "encrypted-files-home", "Specify a bucket name (optional)")
}
