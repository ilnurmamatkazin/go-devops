package service

import (
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage"
)

type Metric interface {
	SetMetric(metric models.Metric) error
	GetMetric(metric *models.Metric) error
	GetInfo() string
}

type Service struct {
	repository *storage.Storage
	cfg        models.Config
}

func New(cfg models.Config, repository *storage.Storage) *Service {
	return &Service{
		repository: repository,
		cfg:        cfg,
	}
}
