package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	workerName   string
	workerCPU    int
	workerMemory int
	workerID     string

	workerCmd = &cobra.Command{
		Use:   "worker",
		Short: "Manage workers in the scheduler",
		Long:  `Register, list, and get information about workers in the scheduler.`,
	}

	registerWorkerCmd = &cobra.Command{
		Use:   "register",
		Short: "Register a new worker",
		Long:  `Register a new worker with the specified name, CPU cores, and memory.`,
		Run: func(cmd *cobra.Command, args []string) {
			registerWorker()
		},
	}

	listWorkersCmd = &cobra.Command{
		Use:   "list",
		Short: "List all workers",
		Long:  `List all registered workers in the scheduler.`,
		Run: func(cmd *cobra.Command, args []string) {
			listWorkers()
		},
	}
)

func init() {
	// Add subcommands to worker command
	workerCmd.AddCommand(registerWorkerCmd)
	workerCmd.AddCommand(listWorkersCmd)

	// Flags for register worker command
	registerWorkerCmd.Flags().StringVar(&workerName, "name", "", "Name of the worker (required)")
	registerWorkerCmd.Flags().IntVar(&workerCPU, "cpu", 1, "Number of CPU cores")
	registerWorkerCmd.Flags().IntVar(&workerMemory, "memory", 1024, "Available memory in MB")
	registerWorkerCmd.MarkFlagRequired("name")
}

func registerWorker() {
	// Prepare request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"name":      workerName,
		"cpu_cores": workerCPU,
		"memory_mb": workerMemory,
	})
	if err != nil {
		exitWithError("Failed to create request: %v", err)
	}

	// Make API request
	resp, err := http.Post(serverURL+"/workers", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		exitWithError("Failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		exitWithError("Failed to read response: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		exitWithError("Failed to register worker: %s", body)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		exitWithError("Failed to parse response: %v", err)
	}

	// Print worker ID
	fmt.Printf("Worker registered successfully. ID: %s\n", response["worker_id"])
}

func listWorkers() {
	// Make API request
	resp, err := http.Get(serverURL + "/workers")
	if err != nil {
		exitWithError("Failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		exitWithError("Failed to read response: %v", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		exitWithError("Failed to list workers: %s", body)
	}

	// Parse response
	var workers []map[string]interface{}
	if err := json.Unmarshal(body, &workers); err != nil {
		exitWithError("Failed to parse response: %v", err)
	}

	// Pretty print workers
	prettyJSON, err := json.MarshalIndent(workers, "", "  ")
	if err != nil {
		exitWithError("Failed to format response: %v", err)
	}

	fmt.Println(string(prettyJSON))
}
