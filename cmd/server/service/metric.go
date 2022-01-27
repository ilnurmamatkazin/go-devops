package service

import (
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

func (s *Service) SetOldMetric(metric models.Metric) {
	s.repository.SetOldMetric(metric)
}

func (s *Service) GetOldMetric(metric *models.Metric) (err error) {
	err = s.repository.ReadMetric(metric)

	// switch metric.MetricType {
	// case "gauge":
	// 	metric.Value = &f
	// case "counter":
	// 	i := int64(f)
	// 	metric.Delta = &i
	// default:
	// 	err = &models.RequestError{
	// 		StatusCode: http.StatusNotImplemented,
	// 		Err:        errors.New(http.StatusText(http.StatusNotImplemented)),
	// 	}
	// }

	return
}

func (s *Service) SetMetric(metric models.Metric) (err error) {
	if s.cfg.Key != "" {
		hash, err := hex.DecodeString(metric.Hash)
		if err != nil {
			fmt.Println("&&&&&&&", s.cfg.Key, err.Error())
			return &models.RequestError{
				StatusCode: http.StatusConflict,
				Err:        errors.New(err.Error()),
			}
		}

		sign := utils.SetHesh(metric.ID, metric.MetricType, s.cfg.Key, metric.Delta, metric.Value)

		if !hmac.Equal(sign, hash) {
			var (
				i int64
				f float64
			)

			if metric.Value != nil {
				f = *metric.Value
			} else {
				f = 0
			}

			if metric.Delta != nil {
				i = *metric.Delta
			} else {
				i = 0
			}

			fmt.Println("&&&&444444&&&", metric.ID, metric.MetricType, i, f, metric.Hash, s.cfg.Key, sign, hash)

			return &models.RequestError{
				StatusCode: http.StatusBadGateway,
				Err:        errors.New("подпись неверна"),
			}
		}
	}

	err = s.repository.SetMetric(metric)

	// switch metric.MetricType {
	// case "gauge":
	// 	metricGauge := models.MetricGauge{Name: metric.ID, Value: *metric.Value}
	// 	err = s.repository.SetGauge(metricGauge)
	// case "counter":
	// 	metricCounter := models.MetricCounter{Name: metric.ID, Value: *metric.Delta}
	// 	err = s.repository.SetCounter(metricCounter)
	// default:
	// 	err = &models.RequestError{
	// 		StatusCode: http.StatusNotImplemented,
	// 		Err:        errors.New(http.StatusText(http.StatusNotImplemented)),
	// 	}
	// }

	return
}

func (s *Service) GetMetric(metric *models.Metric) (err error) {
	if err = s.GetOldMetric(metric); err != nil {
		return
	}

	metric.Hash = utils.SetEncodeHesh(metric.ID, metric.MetricType, s.cfg.Key, metric.Delta, metric.Value)

	return
}

func (s *Service) GetInfo() string {
	return s.repository.Info()
}
