package storage

import (
	"sync"
	"mydung/internal/models"
)

// 1. Hamara naya Database Struct
type MemoryStore struct {
	data  map[string]*models.Job
	mutex sync.RWMutex
}

// 2. Database ko start karne ka function
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data:  make(map[string]*models.Job),
		mutex: sync.RWMutex{},
	}
}

// 3. SaveJob ab MemoryStore ka "Method" ban gaya hai
func (m *MemoryStore) SaveJob(job *models.Job) {
	m.mutex.Lock()
	m.data[job.ID] = job
	m.mutex.Unlock()
}

// 4. GetJob bhi ab MemoryStore ka "Method" ban gaya hai
func (m *MemoryStore) GetJob(jobID string) (*models.Job, bool) {
	m.mutex.RLock()
	job, exists := m.data[jobID]
	m.mutex.RUnlock()
	return job, exists
}