package storage

import (
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/pg"
)

type Storage struct {
	Metric
}

func NewStorage(cfg *models.Config, db *pg.Repository) *Storage {
	return &Storage{
		Metric: NewStorageMetric(cfg, db),
	}
}
