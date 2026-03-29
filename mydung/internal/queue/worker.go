package queue

import (
	"context"
	"fmt"
	"sync"

	"mydung/internal/models"
)

func Worker(ctx context.Context,
			 workerID int,
			  db models.JobStore,
			   wg *sync.WaitGroup,
			    queueName string) {
    defer wg.Done()
    for {
        job, err := db.Dequeue(ctx, queueName) // Channel ki jagah Dequeue 📤
        if err != nil { return } // Shutdown signal
        
        fmt.Printf("[Worker %d] Processing: %s\n", workerID, job.ID)
        job.Status = "COMPLETED"
        db.SaveJob(ctx, job)
    }
}