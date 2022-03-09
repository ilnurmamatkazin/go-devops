package storage

import (
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/pg"
)

// Storage структура, описывающая слой работы с базой данных.
type Storage struct {
	Metric
}

// NewStorage конструктор, создающий структуру слоя работы с базой данных.
func NewStorage(cfg *models.Config, db *pg.Repository) *Storage {
	return &Storage{
		Metric: NewStorageMetric(cfg, db),
	}
}
