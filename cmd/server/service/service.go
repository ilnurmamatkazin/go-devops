package service

import (
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage"
)

type Service struct {
	repository *storage.Storage1
	cfg        models.Config
}

func New(cfg models.Config, repository *storage.Storage1) *Service {
	return &Service{
		repository: repository,
		cfg:        cfg,
	}
}
