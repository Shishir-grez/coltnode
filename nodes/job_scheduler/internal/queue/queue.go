package queue

import (
	"sync"

	"github.com/Shishir_grez/coltnode/nodes/job_scheduler/pkg/models"
)

// JobQueue represents a thread-safe FIFO queue for jobs
type JobQueue struct {
	jobs []*models.Job
	mu   sync.Mutex
}

// NewJobQueue creates a new empty job queue
func NewJobQueue() *JobQueue {
	return &JobQueue{
		jobs: make([]*models.Job, 0),
	}
}

// Enqueue adds a job to the queue
func (q *JobQueue) Enqueue(job *models.Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	q.jobs = append(q.jobs, job)
}

// Dequeue removes and returns the next job from the queue
// Returns nil if queue is empty
func (q *JobQueue) Dequeue() *models.Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	if len(q.jobs) == 0 {
		return nil
	}
	
	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return job
}

// Peek returns the next job without removing it
// Returns nil if queue is empty
func (q *JobQueue) Peek() *models.Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	if len(q.jobs) == 0 {
		return nil
	}
	
	return q.jobs[0]
}

// Size returns the number of jobs in the queue
func (q *JobQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	
	return len(q.jobs)
}