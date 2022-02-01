package service

import (
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

func (s *Service) SetOldMetric(metric models.Metric) {
	s.repository.SetOldMetric(metric)
}

func (s *Service) GetOldMetric(metric *models.Metric) (err error) {
	return s.repository.ReadMetric(metric)
}

func (s *Service) SetMetric(metric models.Metric) (err error) {
	if s.cfg.Key != "" && metric.Hash != nil {
		hash, err := hex.DecodeString(*metric.Hash)
		if err != nil {
			return &models.RequestError{
				StatusCode: http.StatusConflict,
				Err:        errors.New(err.Error()),
			}
		}

		sign := utils.SetHesh(metric.ID, metric.MetricType, s.cfg.Key, metric.Delta, metric.Value)

		if !hmac.Equal(sign, hash) {
			return &models.RequestError{
				StatusCode: http.StatusBadGateway,
				Err:        errors.New("подпись неверна"),
			}
		}
	}

	err = s.repository.SetMetric(metric)

	return
}

func (s *Service) SetArrayMetrics(metrics []models.Metric) (err error) {
	for _, metric := range metrics {
		if err = checkHash(s.cfg.Key, metric); err != nil {
			return
		}
	}

	return s.repository.SetArrayMetrics(metrics)
}

func (s *Service) GetMetric(metric *models.Metric) (err error) {
	if err = s.GetOldMetric(metric); err != nil {
		return
	}

	hash := utils.SetEncodeHesh(metric.ID, metric.MetricType, s.cfg.Key, metric.Delta, metric.Value)
	metric.Hash = &hash

	return
}

func (s *Service) GetInfo() string {
	return s.repository.Info()
}

func checkHash(key string, metric models.Metric) (err error) {
	if key != "" && metric.Hash != nil {
		hash, err := hex.DecodeString(*metric.Hash)
		if err != nil {
			return &models.RequestError{
				StatusCode: http.StatusConflict,
				Err:        errors.New(err.Error()),
			}
		}

		sign := utils.SetHesh(metric.ID, metric.MetricType, key, metric.Delta, metric.Value)

		if !hmac.Equal(sign, hash) {
			return &models.RequestError{
				StatusCode: http.StatusBadGateway,
				Err:        errors.New("подпись неверна"),
			}
		}
	}

	return
}
