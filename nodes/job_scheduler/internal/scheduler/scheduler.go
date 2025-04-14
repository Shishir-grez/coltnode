package scheduler

import (
	"sync"
	"time"

	"github.com/Shishir_grez/coltnode/nodes/job_scheduler/internal/queue"
	"github.com/Shishir_grez/coltnode/nodes/job_scheduler/pkg/models"
)

// Scheduler manages job assignments to available workers
type Scheduler struct {
	jobQueue    *queue.JobQueue
	workers     []*models.Worker
	workerIndex int // For round-robin assignment
	mu          sync.Mutex
	storage     Storage // Interface for persistence
}

// Storage defines the interface for job and worker persistence
type Storage interface {
	SaveJob(*models.Job) error
	GetJob(id string) (*models.Job, error)
	UpdateJob(*models.Job) error
	SaveWorker(*models.Worker) error
	GetWorker(id string) (*models.Worker, error)
	UpdateWorker(*models.Worker) error
	GetAvailableWorkers() ([]*models.Worker, error)
}

// NewScheduler creates a new scheduler with the given queue and storage
func NewScheduler(jobQueue *queue.JobQueue, storage Storage) *Scheduler {
	return &Scheduler{
		jobQueue:    jobQueue,
		workers:     make([]*models.Worker, 0),
		workerIndex: 0,
		storage:     storage,
	}
}

// RegisterWorker adds a new worker to the scheduler
func (s *Scheduler) RegisterWorker(worker *models.Worker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.workers = append(s.workers, worker)
	return s.storage.SaveWorker(worker)
}

// ScheduleJob assigns a job to a worker using round-robin
func (s *Scheduler) ScheduleJob(job *models.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	availableWorkers, err := s.storage.GetAvailableWorkers()
	if err != nil {
		return err
	}
	
	if len(availableWorkers) == 0 {
		// No workers available, keep job in queue
		return nil
	}
	
	// Simple round-robin worker selection
	worker := availableWorkers[s.workerIndex]
	s.workerIndex = (s.workerIndex + 1) % len(availableWorkers)
	
	// Update job status
	job.Status = "running"
	err = s.storage.UpdateJob(job)
	if err != nil {
		return err
	}
	
	// In a real system, you'd send the job to the worker here
	// For now, we'll just simulate it
	go s.simulateJobExecution(job, worker)
	
	return nil
}

// simulateJobExecution simulates a job running on a worker
func (s *Scheduler) simulateJobExecution(job *models.Job, worker *models.Worker) {
	// Simulate job execution time
	time.Sleep(5 * time.Second)
	
	// Update job status to completed
	job.Status = "completed"
	s.storage.UpdateJob(job)
}

// Start begins the scheduling process
func (s *Scheduler) Start() {
	go func() {
		for {
			// Check for jobs in the queue
			job := s.jobQueue.Dequeue()
			if job != nil {
				s.ScheduleJob(job)
			}
			
			// Don't busy-wait
			time.Sleep(1 * time.Second)
		}
	}()
}