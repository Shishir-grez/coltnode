package models

import (
	"time"
	"github.com/google/uuid"
)

// Resources represents the computing resources available on a worker
type Resources struct {
	CPUCores int // Number of CPU cores
	MemoryMB int // Available memory in MB
}

// Worker represents a node that can execute jobs
type Worker struct {
	ID           string    // Unique identifier for the worker
	Name         string    // Human-readable name for the worker
	Status       string    // Current status: active, offline, busy
	Resources    Resources // Available resources on this worker
	LastHeartbeat time.Time // Last time we heard from this worker
}

// NewWorker creates a new Worker with default values
func NewWorker(name string, cpuCores, memoryMB int) *Worker {
	return &Worker{
		ID:           uuid.New().String(),
		Name:         name,
		Status:       "active",
		Resources:    Resources{CPUCores: cpuCores, MemoryMB: memoryMB},
		LastHeartbeat: time.Now(),
	}
}