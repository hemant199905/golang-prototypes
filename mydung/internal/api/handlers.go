package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mydung/internal/models"
)

// Helper function (private rakha kyunki sirf yahi use hoga)
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		fmt.Printf("[Error] JSON encode fail hua: %v\n", err)
	}
}

func MakeJobHandler(db models.JobStore) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Query parameter se queue ka naam nikalna
        queueName := r.URL.Query().Get("queue")

        // 2. Production-grade practice: Agar name khali hai toh default set karein
        if queueName == "" {
            queueName = "default" 
        }

        jobID := fmt.Sprintf("job_%d", time.Now().Unix())
        newJob := models.Job{ID: jobID, Status: "PENDING"}

        // Context pass karte hue save aur enqueue karna
        db.SaveJob(r.Context(), &newJob)
        err := db.Enqueue(r.Context(), queueName, &newJob)
        
        if err != nil {
            respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Queue full!"})
            return
        }

        respondJSON(w, http.StatusAccepted, map[string]string{
            "job_id": jobID,
            "queue":  queueName,
        })
    }
}

// // S capital rakha
// func StatusHandler(db models.JobStore) http.HandlerFunc {	return func(w http.ResponseWriter, r *http.Request) {
// 		// URL se "id" parameter nikalna
// 		jobID := r.URL.Query().Get("id")

// 	if jobID == "" {
// 		respondJSON(w, http.StatusBadRequest, map[string]string{
// 			"error": "Job ID dena zaroori hai!",
// 		})
// 		return
// 	}

// 	// Hamare Map (Database) mein check karna
// 	job, exists := db.GetJob(jobID)
// 	if !exists {
// 		respondJSON(w, http.StatusNotFound, map[string]string{
// 			"error": "Wrong ID! Kripya sahi Job ID dein.",
// 		})
// 		return
// 	}

// 	// Agar mil gaya, toh uska pura data (aur status) user ko bhej do
// 	respondJSON(w, http.StatusOK, job)
// }}