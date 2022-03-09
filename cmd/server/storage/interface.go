package storage

import "github.com/ilnurmamatkazin/go-devops/cmd/server/models"

//go:generate mockgen -source=interface.go -destination=mock_service/mock.go

type Metric interface {
	SetOldMetric(metric models.Metric)
	ReadMetric(metric *models.Metric) error
	SetMetric(metric models.Metric) error
	Info() string
	ConnectPG() error
	Save() error
	SetArrayMetrics(metric []models.Metric) error
	Ping() error
}
