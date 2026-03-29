package storage

import (
    "context"
    "sync"
    "mydung/internal/models"
)

type MemoryStore struct {
    data   map[string]*models.Job
    queues map[string]chan *models.Job
    mutex  sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        data:   make(map[string]*models.Job),
        queues: make(map[string]chan *models.Job),
    }
}

// Helper for lazy initialization
func (m *MemoryStore) getOrCreateQueue(name string) chan *models.Job {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    if q, exists := m.queues[name]; exists {
        return q
    }
    q := make(chan *models.Job, 100)
    m.queues[name] = q
    return q
}

func (m *MemoryStore) SaveJob(ctx context.Context, job *models.Job) error {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    m.data[job.ID] = job
    return nil
}

func (m *MemoryStore) GetJob(ctx context.Context, jobID string) (*models.Job, bool) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    job, exists := m.data[jobID]
    return job, exists
}

func (m *MemoryStore) Enqueue(ctx context.Context, queueName string, job *models.Job) error {
    q := m.getOrCreateQueue(queueName)
    select {
    case q <- job: return nil
    case <-ctx.Done(): return ctx.Err()
    }
}

func (m *MemoryStore) Dequeue(ctx context.Context, queueName string) (*models.Job, error) {
    q := m.getOrCreateQueue(queueName) // Lazy create taaki worker "fate" nahi 🛡️
    select {
    case job := <-q: return job, nil
    case <-ctx.Done(): return nil, ctx.Err()
    }
}