package service

import (
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/memory"
)

type Metric interface {
	SetMetric(metric models.Metric) error
	GetMetric(metric *models.Metric) error
	GetInfo() string
}

type Service struct {
	repository *memory.MemoryRepository
}

func NewService() *Service {
	return &Service{
		repository: memory.NewMemoryRepository(),
	}
}
