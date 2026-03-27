package queue

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"mydung/internal/models"
)

// wg ko as a pointer pass kiya
func Worker(workerID int, jobs <-chan models.Job, wg *sync.WaitGroup, db models.JobStore) {	// 🌟 MAGIC WORD: Yeh line automatically function ke END mein chalegi
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("[Worker %d] Processing Job: %s...\n", workerID, job.ID)
		time.Sleep(2 * time.Second) 

		mapJob, exists := db.GetJob(job.ID)
		if !exists { continue }

		// Simulate Error
		if rand.Float32() < 0.30 {
			mapJob.Status = "FAILED"
			mapJob.FailureReason = "Database timeout error 🚨"
			db.SaveJob(mapJob)
			fmt.Printf("[Worker %d] Job %s FAILED! ❌\n", workerID, job.ID)
			continue 
		}

		mapJob.Status = "COMPLETED"
		mapJob.FailureReason = ""
		db.SaveJob(mapJob)
		fmt.Printf("[Worker %d] Job %s COMPLETED! ✅\n", workerID, job.ID)
	}

	// Loop khatam hone ke baad yeh print hoga
	fmt.Printf("🚶‍♂️ [Worker %d] Channel closed, shift over. Going home!\n", workerID)
}