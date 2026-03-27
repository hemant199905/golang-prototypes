package storage

import (
	"fmt"
	"mydung/internal/models"
)

// NewJobStore ek factory hai jo interface return karti hai
func NewJobStore(cfg models.Config) (models.JobStore, error) {
	switch cfg.Storage.Type {
	case "memory":
		return NewMemoryStore(), nil // Memory implementation
	case "redis":
		return NewRedisStore(cfg.Redis) // Redis implementation
	default:
		return nil, fmt.Errorf("invalid storage type: %s", cfg.Storage.Type)
	}
}
