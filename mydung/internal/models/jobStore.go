package models

// Yeh hamara Contract hai
type JobStore interface {
	SaveJob(job *Job)
	GetJob(jobID string) (*Job, bool)
}
