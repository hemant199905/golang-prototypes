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

// M capital rakha
func MakeJobHandler(jobQueue chan<- models.Job, db models.JobStore) http.HandlerFunc {	return func(w http.ResponseWriter, r *http.Request) {
		jobID := fmt.Sprintf("job_%d", time.Now().Unix())
		newJob := models.Job{ID: jobID, Name: "on_demand_report", Status: "PENDING"}

		// 🌟 Doosre package se function call kiya!
		db.SaveJob(&newJob)

		jobQueue <- newJob
		respondJSON(w, http.StatusAccepted, map[string]string{"job_id": jobID})
	}
}

// S capital rakha
func StatusHandler(db models.JobStore) http.HandlerFunc {	return func(w http.ResponseWriter, r *http.Request) {
		// URL se "id" parameter nikalna
		jobID := r.URL.Query().Get("id")

	if jobID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Job ID dena zaroori hai!",
		})
		return
	}

	// Hamare Map (Database) mein check karna
	job, exists := db.GetJob(jobID)
	if !exists {
		respondJSON(w, http.StatusNotFound, map[string]string{
			"error": "Wrong ID! Kripya sahi Job ID dein.",
		})
		return
	}

	// Agar mil gaya, toh uska pura data (aur status) user ko bhej do
	respondJSON(w, http.StatusOK, job)
}}