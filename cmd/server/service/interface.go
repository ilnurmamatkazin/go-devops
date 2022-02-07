package service

import "github.com/ilnurmamatkazin/go-devops/cmd/server/models"

type Metric interface {
	SetMetric(metric models.Metric) error
	GetMetric(metric *models.Metric) error
	GetInfo() string
}
