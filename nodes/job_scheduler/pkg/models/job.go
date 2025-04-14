package models

import (
	"time"
	"github.com/google/uuid"
)

// Job represents a task to be executed by a worker
type Job struct {
	ID         string    // Unique identifier for the job
	Name       string    // Human-readable name for the job
	Command    string    // Command to be executed
	Args       []string  // Arguments for the command
	Status     string    // Current status: pending, running, completed, failed
	SubmitTime time.Time // Time when the job was submitted
}

func generateUniqueID() string {
	return uuid.New().String()
}

// NewJob creates a new Job with default values
func NewJob(name, command string, args []string) *Job {
	return &Job{
		ID:         generateUniqueID(), // You'll need to implement this
		Name:       name,
		Command:    command,
		Args:       args,
		Status:     "pending",
		SubmitTime: time.Now(),
	}
}