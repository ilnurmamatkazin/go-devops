package service

import (
	"errors"
	"net/http"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (s *Service) SetMetric(metric models.Metric) (err error) {
	switch metric.MType {
	case "gauge":
		metricGauge := models.MetricGauge{Name: metric.ID, Value: *metric.Value}
		err = s.repository.SetGauge(metricGauge)
	case "counter":
		metricCounter := models.MetricCounter{Name: metric.ID, Value: *metric.Delta}
		err = s.repository.SetCounter(metricCounter)
	default:
		err = &models.RequestError{
			StatusCode: http.StatusNotImplemented,
			Err:        errors.New(http.StatusText(http.StatusNotImplemented)),
		}
	}

	return
}

func (s *Service) GetMetric(metric *models.Metric) (err error) {
	switch metric.MType {
	case "gauge":
		var f float64
		f, err = s.repository.ReadGauge(metric.ID)
		metric.Value = &f
	case "counter":
		var i int64
		i, err = s.repository.ReadCounter(metric.ID)
		metric.Delta = &i
	default:
		err = &models.RequestError{
			StatusCode: http.StatusNotImplemented,
			Err:        errors.New(http.StatusText(http.StatusNotImplemented)),
		}
	}

	return
}

func (s *Service) GetInfo() string {
	return s.repository.Info()
}
