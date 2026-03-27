package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"mydung/internal/models"
)

type MemoryStore struct {
	data  map[string]*models.Job
	// Queue ke liye hum ek map of channels use karenge
	queues map[string]chan *models.Job
	mutex  sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data:   make(map[string]*models.Job),
		queues: make(map[string]chan *models.Job),
		mutex:  sync.RWMutex{},
	}
}

// SaveJob implementation
func (m *MemoryStore) SaveJob(ctx context.Context, job *models.Job) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.data[job.ID] = job
	return nil
}

// GetJob implementation
func (m *MemoryStore) GetJob(ctx context.Context, jobID string) (*models.Job, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	job, exists := m.data[jobID]
	return job, exists
}

// Enqueue: Job ko channel mein dalta hai
func (m *MemoryStore) Enqueue(ctx context.Context, queueName string, job *models.Job) error {
	m.mutex.Lock()
	q, exists := m.queues[queueName]
	if !exists {
		// Agar queue nahi hai, toh 100 size ka buffer banate hain
		q = make(chan *models.Job, 100)
		m.queues[queueName] = q
	}
	m.mutex.Unlock()

	select {
	case q <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("queue %s is full", queueName)
	}
}

// Dequeue: Channel se job nikalta hai (Blocking like BRPop)
func (m *MemoryStore) Dequeue(ctx context.Context, queueName string) (*models.Job, error) {
	m.mutex.RLock()
	q, exists := m.queues[queueName]
	m.mutex.RUnlock()

	if !exists {
		return nil, errors.New("queue not found")
	}

	select {
	case job := <-q:
		return job, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}