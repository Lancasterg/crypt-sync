/*
Copyright © 2026 GEORGE LANCASTER <lancaster0180@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "go-crypt-sync",
	Version: "0.0.3",
	Short:   "A simple CLI tool for encrypting files before storing them in a GCS bucket, and then decrypting them once they are ready to be viewed again.",
	Long: `This tool provides a secure way to manage sensitive files in Google Cloud Storage (GCS). 
	It is built on the principle of "Zero Trust" regarding cloud providers. By encrypting files 
	locally using the 'age' encryption format before they ever reach the network, you ensure 
	that the cloud provider only ever hosts opaque, encrypted blobs. Even in the event of a 
	provider-side data breach or unauthorized access, your data remains protected as the 
	decryption keys never leave your local machine.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
