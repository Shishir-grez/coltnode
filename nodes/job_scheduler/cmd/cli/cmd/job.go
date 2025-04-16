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
	jobName    string
	jobCommand string
	jobArgs    []string
	jobID      string

	jobCmd = &cobra.Command{
		Use:   "job",
		Short: "Manage jobs in the scheduler",
		Long:  `Create, list, and get information about jobs in the scheduler.`,
	}

	createJobCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new job",
		Long:  `Create a new job with the specified name, command, and arguments.`,
		Run: func(cmd *cobra.Command, args []string) {
			createJob()
		},
	}

	getJobCmd = &cobra.Command{
		Use:   "get",
		Short: "Get information about a job",
		Long:  `Get detailed information about a specific job by ID.`,
		Run: func(cmd *cobra.Command, args []string) {
			getJob()
		},
	}

	listJobsCmd = &cobra.Command{
		Use:   "list",
		Short: "List all jobs",
		Long:  `List all jobs in the scheduler.`,
		Run: func(cmd *cobra.Command, args []string) {
			listJobs()
		},
	}
)

func init() {
	// Add subcommands to job command
	jobCmd.AddCommand(createJobCmd)
	jobCmd.AddCommand(getJobCmd)
	jobCmd.AddCommand(listJobsCmd)

	// Flags for create job command
	createJobCmd.Flags().StringVar(&jobName, "name", "", "Name of the job (required)")
	createJobCmd.Flags().StringVar(&jobCommand, "command", "", "Command to execute (required)")
	createJobCmd.Flags().StringArrayVar(&jobArgs, "arg", []string{}, "Arguments for the command (can be specified multiple times)")
	createJobCmd.MarkFlagRequired("name")
	createJobCmd.MarkFlagRequired("command")

	// Flags for get job command
	getJobCmd.Flags().StringVar(&jobID, "id", "", "ID of the job to get information about (required)")
	getJobCmd.MarkFlagRequired("id")
}

func createJob() {
	// Prepare request body
	requestBody, err := json.Marshal(map[string]interface{}{
		"name":    jobName,
		"command": jobCommand,
		"args":    jobArgs,
	})
	if err != nil {
		exitWithError("Failed to create request: %v", err)
	}

	// Make API request
	resp, err := http.Post(serverURL+"/jobs", "application/json", bytes.NewBuffer(requestBody))
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
		exitWithError("Failed to create job: %s", body)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		exitWithError("Failed to parse response: %v", err)
	}

	// Print job ID
	fmt.Printf("Job created successfully. ID: %s\n", response["job_id"])
}

func getJob() {
	// Make API request
	resp, err := http.Get(serverURL + "/jobs/" + jobID)
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
		exitWithError("Failed to get job: %s", body)
	}

	// Parse response
	var job map[string]interface{}
	if err := json.Unmarshal(body, &job); err != nil {
		exitWithError("Failed to parse response: %v", err)
	}

	// Pretty print job information
	prettyJSON, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		exitWithError("Failed to format response: %v", err)
	}

	fmt.Println(string(prettyJSON))
}

func listJobs() {
	// Make API request
	resp, err := http.Get(serverURL + "/jobs")
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
		exitWithError("Failed to list jobs: %s", body)
	}

	// Parse response
	var jobs []map[string]interface{}
	if err := json.Unmarshal(body, &jobs); err != nil {
		exitWithError("Failed to parse response: %v", err)
	}

	// Pretty print jobs
	prettyJSON, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		exitWithError("Failed to format response: %v", err)
	}

	fmt.Println(string(prettyJSON))
}
