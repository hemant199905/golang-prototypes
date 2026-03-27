package models
import "context"

// JobStore interface ko context-aware banate hain
type JobStore interface {
	SaveJob(ctx context.Context, job *Job) error        // Context add kiya aur error return type bhi
	GetJob(ctx context.Context, jobID string) (*Job, bool)
}