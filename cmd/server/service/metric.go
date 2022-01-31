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
	if s.cfg.Key != "" && metric.Hash != nil {
		hash, err := hex.DecodeString(*metric.Hash)
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

	return
}

func (s *Service) SetArrayMetrics(metrics []models.Metric) (err error) {
	for _, metric := range metrics {
		// var (
		// 	i int64   = -100
		// 	f float64 = -100
		// )

		// if metric.Delta != nil {
		// 	i = *metric.Delta
		// }

		// if metric.Value != nil {
		// 	f = *metric.Value
		// }

		// fmt.Println("&&&&increment11 SetArrayMetrics&&&", metric, i, f)

		if err = checkHash(s.cfg.Key, metric); err != nil {
			// fmt.Println("&&&&increment11 SetArrayMetrics checkHash&&&", err.Error())

			return
		}
	}

	err = s.repository.SetArrayMetrics(metrics)

	return
}

func (s *Service) GetMetric(metric *models.Metric) (err error) {
	if err = s.GetOldMetric(metric); err != nil {
		return
	}

	// var (
	// 	i int64   = -100
	// 	f float64 = -100
	// )

	// if metric.Delta != nil {
	// 	i = *metric.Delta
	// }

	// if metric.Value != nil {
	// 	f = *metric.Value
	// }

	// fmt.Println("&&&&increment6 GetMetric&&&", metric, i, f)

	hash := utils.SetEncodeHesh(metric.ID, metric.MetricType, s.cfg.Key, metric.Delta, metric.Value)
	// fmt.Println("&&&&increment6 GetMetric&&&", metric, i, f, hash)

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
			fmt.Println("&&&&&&&", key, err.Error())
			return &models.RequestError{
				StatusCode: http.StatusConflict,
				Err:        errors.New(err.Error()),
			}
		}

		sign := utils.SetHesh(metric.ID, metric.MetricType, key, metric.Delta, metric.Value)

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

			fmt.Println("&&&&444444&&&", metric.ID, metric.MetricType, i, f, metric.Hash, key, sign, hash)

			return &models.RequestError{
				StatusCode: http.StatusBadGateway,
				Err:        errors.New("подпись неверна"),
			}
		}
	}

	return
}
