package storage

import (
	"errors"
	"sync"

	"github.com/Shishir_grez/coltnode/nodes/job_scheduler/pkg/models"
)

// MemoryStorage implements in-memory storage for jobs and workers
type MemoryStorage struct {
	jobs    map[string]*models.Job
	workers map[string]*models.Worker
	mu      sync.RWMutex
}

// NewMemoryStorage creates a new memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		jobs:    make(map[string]*models.Job),
		workers: make(map[string]*models.Worker),
	}
}

// SaveJob stores a job in memory
func (s *MemoryStorage) SaveJob(job *models.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.jobs[job.ID] = job
	return nil
}

// GetJob retrieves a job by ID
func (s *MemoryStorage) GetJob(id string) (*models.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	job, exists := s.jobs[id]
	if !exists {
		return nil, errors.New("job not found")
	}
	return job, nil
}

// UpdateJob updates an existing job
func (s *MemoryStorage) UpdateJob(job *models.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	_, exists := s.jobs[job.ID]
	if !exists {
		return errors.New("job not found")
	}
	
	s.jobs[job.ID] = job
	return nil
}

// SaveWorker stores a worker in memory
func (s *MemoryStorage) SaveWorker(worker *models.Worker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.workers[worker.ID] = worker
	return nil
}

// GetWorker retrieves a worker by ID
func (s *MemoryStorage) GetWorker(id string) (*models.Worker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	worker, exists := s.workers[id]
	if !exists {
		return nil, errors.New("worker not found")
	}
	return worker, nil
}

// UpdateWorker updates an existing worker
func (s *MemoryStorage) UpdateWorker(worker *models.Worker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	_, exists := s.workers[worker.ID]
	if !exists {
		return errors.New("worker not found")
	}
	
	s.workers[worker.ID] = worker
	return nil
}

// GetAvailableWorkers returns all workers with "active" status
func (s *MemoryStorage) GetAvailableWorkers() ([]*models.Worker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	available := make([]*models.Worker, 0)
	for _, worker := range s.workers {
		if worker.Status == "active" {
			available = append(available, worker)
		}
	}
	return available, nil
}

// GetAllJobs returns all jobs in the storage
func (s *MemoryStorage) GetAllJobs() ([]*models.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	jobs := make([]*models.Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	return jobs, nil
}