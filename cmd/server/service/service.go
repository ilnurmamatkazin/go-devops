package service

import (
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage"
)

// Service структура, описывающая слой бизнес-логики.
type Service struct {
	Metric
}

// NewService конструктор, создающий структуру слоя бизнес-логики.
func NewService(cfg *models.Config, repository *storage.Storage) *Service {
	return &Service{
		Metric: NewServiceMetric(cfg, repository),
	}
}
