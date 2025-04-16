package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	interactiveCmd = &cobra.Command{
		Use:   "interactive",
		Short: "Start interactive mode",
		Long:  `Start an interactive shell for the job scheduler CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			startInteractiveMode()
		},
	}
)

func startInteractiveMode() {
	fmt.Println("Welcome to ColtNode CLI Interactive Mode")
	fmt.Println("Type 'help' for available commands or 'exit' to quit")
	fmt.Println("Server URL:", serverURL)
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("coltnode> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "exit" || input == "quit" {
			fmt.Println("Exiting interactive mode")
			break
		}

		handleInteractiveCommand(input)
	}
}

func handleInteractiveCommand(input string) {
	// Split the input into command and arguments
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "help":
		showInteractiveHelp()
	case "job":
		handleJobCommand(args)
	case "worker":
		handleWorkerCommand(args)
	case "server":
		if len(args) > 0 {
			serverURL = args[0]
			fmt.Println("Server URL set to:", serverURL)
		} else {
			fmt.Println("Current server URL:", serverURL)
		}
	default:
		fmt.Printf("Unknown command: %s\nType 'help' for available commands\n", command)
	}
}

func showInteractiveHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help                           Show this help message")
	fmt.Println("  exit, quit                     Exit interactive mode")
	fmt.Println("  server [url]                   Show or set server URL")
	fmt.Println("  job create --name NAME --command CMD [--arg ARG]...")
	fmt.Println("                                 Create a new job")
	fmt.Println("  job get --id ID                Get information about a job")
	fmt.Println("  job list                       List all jobs")
	fmt.Println("  worker register --name NAME [--cpu N] [--memory M]")
	fmt.Println("                                 Register a new worker")
	fmt.Println("  worker list                    List all workers")
}

func handleJobCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Missing job subcommand. Available: create, get, list")
		return
	}

	subcommand := args[0]
	subargs := args[1:]

	switch subcommand {
	case "create":
		// Parse arguments for job creation
		jobName = ""
		jobCommand = ""
		jobArgs = []string{}

		for i := 0; i < len(subargs); i++ {
			if subargs[i] == "--name" && i+1 < len(subargs) {
				jobName = subargs[i+1]
				i++
			} else if subargs[i] == "--command" && i+1 < len(subargs) {
				jobCommand = subargs[i+1]
				i++
			} else if subargs[i] == "--arg" && i+1 < len(subargs) {
				jobArgs = append(jobArgs, subargs[i+1])
				i++
			}
		}

		if jobName == "" || jobCommand == "" {
			fmt.Println("Missing required arguments. Usage: job create --name NAME --command CMD [--arg ARG]...")
			return
		}

		createJob()

	case "get":
		// Parse arguments for job info
		jobID = ""

		for i := 0; i < len(subargs); i++ {
			if subargs[i] == "--id" && i+1 < len(subargs) {
				jobID = subargs[i+1]
				i++
			}
		}

		if jobID == "" {
			fmt.Println("Missing required argument. Usage: job get --id ID")
			return
		}

		getJob()

	case "list":
		listJobs()

	default:
		fmt.Printf("Unknown job subcommand: %s\nAvailable: create, get, list\n", subcommand)
	}
}

func handleWorkerCommand(args []string) {
	if len(args) == 0 {
		fmt.Println("Missing worker subcommand. Available: register, list")
		return
	}

	subcommand := args[0]
	subargs := args[1:]

	switch subcommand {
	case "register":
		// Parse arguments for worker registration
		workerName = ""
		workerCPU = 1
		workerMemory = 1024

		for i := 0; i < len(subargs); i++ {
			if subargs[i] == "--name" && i+1 < len(subargs) {
				workerName = subargs[i+1]
				i++
			} else if subargs[i] == "--cpu" && i+1 < len(subargs) {
				fmt.Sscanf(subargs[i+1], "%d", &workerCPU)
				i++
			} else if subargs[i] == "--memory" && i+1 < len(subargs) {
				fmt.Sscanf(subargs[i+1], "%d", &workerMemory)
				i++
			}
		}

		if workerName == "" {
			fmt.Println("Missing required argument. Usage: worker register --name NAME [--cpu N] [--memory M]")
			return
		}

		registerWorker()

	case "list":
		listWorkers()

	default:
		fmt.Printf("Unknown worker subcommand: %s\nAvailable: register, list\n", subcommand)
	}
}
