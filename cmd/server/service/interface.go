package service

import "github.com/ilnurmamatkazin/go-devops/cmd/server/models"

//go:generate mockgen -source=interface.go -destination=mock_service/mock.go

// Metric интерфейс, описывающий методы слоя бизнес-логики.
type Metric interface {
	SetMetric(metric models.Metric) error
	GetMetric(metric *models.Metric) error
	SetArrayMetrics(metrics []models.Metric) error
	GetInfo() string
	Ping() error
	SetOldMetric(metric models.Metric)
	GetOldMetric(metric *models.Metric) error
}
