package main

import (
	"fmt"
	"os"

	"github.com/Shishir_grez/coltnode/nodes/job_scheduler/cmd/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
