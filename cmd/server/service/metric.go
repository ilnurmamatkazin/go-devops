package service

import (
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

// ServiceMetric структура, описывающая слой бизнес-логики по работе с метрикой.
type ServiceMetric struct {
	storage *storage.Storage
	cfg     *models.Config
}

// NewServiceMetric конструктор, создающий структуру слоя бизнес-логики по работе с метрикой.
func NewServiceMetric(cfg *models.Config, storage *storage.Storage) *ServiceMetric {
	return &ServiceMetric{
		cfg:     cfg,
		storage: storage,
	}
}

// SetOldMetric устаревшия функция сохранения метрики в системе.
func (s *ServiceMetric) SetOldMetric(metric models.Metric) {
	s.storage.SetOldMetric(metric)
}

// GetOldMetric устаревшия функция получения метрики.
func (s *ServiceMetric) GetOldMetric(metric *models.Metric) (err error) {
	return s.storage.ReadMetric(metric)
}

// SetMetric функция сохранения метрики в системе.
func (s *ServiceMetric) SetMetric(metric models.Metric) (err error) {
	if s.cfg.Key != "" && metric.Hash != nil {
		hash, err := hex.DecodeString(*metric.Hash)
		if err != nil {
			return &models.RequestError{
				StatusCode: http.StatusConflict,
				Err:        errors.New(err.Error()),
			}
		}

		sign := utils.SetHash(metric.ID, metric.MetricType, s.cfg.Key, metric.Delta, metric.Value)

		if !hmac.Equal(sign, hash) {
			return &models.RequestError{
				StatusCode: http.StatusBadGateway,
				Err:        errors.New("подпись неверна"),
			}
		}
	}

	err = s.storage.SetMetric(metric)

	return
}

// SetMetric функция сохранения массива метрик в системе.
func (s *ServiceMetric) SetArrayMetrics(metrics []models.Metric) (err error) {
	for _, metric := range metrics {
		if err = checkHash(s.cfg.Key, metric); err != nil {
			return
		}
	}

	return s.storage.SetArrayMetrics(metrics)
}

// GetMetric функция получения метрики.
func (s *ServiceMetric) GetMetric(metric *models.Metric) (err error) {
	if err = s.GetOldMetric(metric); err != nil {
		return
	}

	hash := utils.SetEncodeHash(metric.ID, metric.MetricType, s.cfg.Key, metric.Delta, metric.Value)
	metric.Hash = &hash

	return
}

// GetInfo функция получения html страницы со списком метрик.
func (s *ServiceMetric) GetInfo() string {
	return s.storage.Info()
}

// checkHash внутренняя функция проверки подписи метрики.
func checkHash(key string, metric models.Metric) (err error) {
	if key != "" && metric.Hash != nil {
		hash, err := hex.DecodeString(*metric.Hash)
		if err != nil {
			return &models.RequestError{
				StatusCode: http.StatusConflict,
				Err:        errors.New(err.Error()),
			}
		}

		sign := utils.SetHash(metric.ID, metric.MetricType, key, metric.Delta, metric.Value)

		if !hmac.Equal(sign, hash) {
			return &models.RequestError{
				StatusCode: http.StatusBadGateway,
				Err:        errors.New("подпись неверна"),
			}
		}
	}

	return
}
