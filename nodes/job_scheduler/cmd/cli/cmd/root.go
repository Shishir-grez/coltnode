package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	serverURL string
	rootCmd   = &cobra.Command{
		Use:   "coltnode",
		Short: "ColtNode CLI - A command-line interface for the job scheduler",
		Long: `ColtNode CLI is a comprehensive command-line tool for interacting with the job scheduler.
It supports both interactive and command modes for managing jobs and workers.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "Server URL for the job scheduler API")

	// Add commands
	rootCmd.AddCommand(jobCmd)
	rootCmd.AddCommand(workerCmd)
	rootCmd.AddCommand(interactiveCmd)
}

// exitWithError prints an error message and exits with code 1
func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
