package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Shishir_grez/coltnode/nodes/job_scheduler/internal/queue"
	"github.com/Shishir_grez/coltnode/nodes/job_scheduler/internal/scheduler"
	"github.com/Shishir_grez/coltnode/nodes/job_scheduler/internal/storage"
	"github.com/Shishir_grez/coltnode/nodes/job_scheduler/pkg/models"
)

func main() {
	// Initialize components
	jobQueue := queue.NewJobQueue()
	memoryStorage := storage.NewMemoryStorage()
	jobScheduler := scheduler.NewScheduler(jobQueue, memoryStorage)
	
	// Start the scheduler
	jobScheduler.Start()
	
	// Set up Gin router
	router := gin.Default()
	
	// API endpoints
	router.POST("/jobs", func(c *gin.Context) {
		var jobRequest struct {
			Name    string   `json:"name" binding:"required"`
			Command string   `json:"command" binding:"required"`
			Args    []string `json:"args"`
		}
		
		if err := c.ShouldBindJSON(&jobRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		job := models.NewJob(jobRequest.Name, jobRequest.Command, jobRequest.Args)
		job.ID = uuid.New().String()
		
		// Save the job
		if err := memoryStorage.SaveJob(job); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save job"})
			return
		}
		
		// Add to queue
		jobQueue.Enqueue(job)
		
		c.JSON(http.StatusCreated, gin.H{
			"job_id": job.ID,
			"status": job.Status,
		})
	})
	
	router.GET("/jobs/:id", func(c *gin.Context) {
		jobID := c.Param("id")
		
		job, err := memoryStorage.GetJob(jobID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		
		c.JSON(http.StatusOK, job)
	})
	
	router.POST("/workers", func(c *gin.Context) {
		var workerRequest struct {
			Name     string `json:"name" binding:"required"`
			CPUCores int    `json:"cpu_cores" binding:"required"`
			MemoryMB int    `json:"memory_mb" binding:"required"`
		}
		
		if err := c.ShouldBindJSON(&workerRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		worker := models.NewWorker(workerRequest.Name, workerRequest.CPUCores, workerRequest.MemoryMB)
		
		// Register the worker
		if err := jobScheduler.RegisterWorker(worker); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register worker"})
			return
		}
		
		c.JSON(http.StatusCreated, gin.H{
			"worker_id": worker.ID,
			"status": worker.Status,
		})
	})
	
	router.GET("/workers", func(c *gin.Context) {
		workers, err := memoryStorage.GetAvailableWorkers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get workers"})
			return
		}
		
		c.JSON(http.StatusOK, workers)
	})
	
	// Start the server
	router.Run(":8080")
}