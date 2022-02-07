package storage

import "github.com/ilnurmamatkazin/go-devops/cmd/server/models"

type Metric interface {
	ReadGauge(name string) (float64, error)
	SetGauge(metric models.MetricGauge) error
	ReadCounter(name string) (int, error)
	SetCounter(metric models.MetricCounter) error
	Info() string
}
