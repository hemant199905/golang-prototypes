package models

import "context"

type JobStore interface {
	// Persistence (Status check ke liye)
	SaveJob(ctx context.Context, job *Job) error
	GetJob(ctx context.Context, jobID string) (*Job, bool)

	// Queue (Job processing ke liye)
	Enqueue(ctx context.Context, queueName string, job *Job) error
	Dequeue(ctx context.Context, queueName string) (*Job, error)
}