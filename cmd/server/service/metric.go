package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (s *Service) SetOldMetric(metric models.Metric) {
	s.repository.SetOldMetric(metric)
}

func (s *Service) GetOldMetric(metric *models.Metric) (err error) {
	f, err := s.repository.ReadMetric(metric.ID)

	switch metric.MetricType {
	case "gauge":
		metric.Value = &f
	case "counter":
		i := int64(f)
		metric.Delta = &i
	default:
		err = &models.RequestError{
			StatusCode: http.StatusNotImplemented,
			Err:        errors.New(http.StatusText(http.StatusNotImplemented)),
		}
	}

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

		h := hmac.New(sha256.New, []byte(s.cfg.Key))
		h.Write(hash)
		sign := h.Sum(nil)

		if !hmac.Equal(sign, hash) {
			fmt.Println("&&&&444444&&&", s.cfg.Key)

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

	setHesh(metric, s.cfg.Key)

	return
}

func (s *Service) GetInfo() string {
	return s.repository.Info()
}

func setHesh(metric *models.Metric, key string) {
	if key == "" {
		return
	}

	var hash []byte

	if metric.MetricType == "gauge" {
		hash = []byte(fmt.Sprintf("%s:gauge:%f", metric.ID, *metric.Value))
	} else {
		hash = []byte(fmt.Sprintf("%s:counter:%d", metric.ID, *metric.Delta))
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write(hash)

	metric.Hash = hex.EncodeToString(h.Sum(nil))
}
