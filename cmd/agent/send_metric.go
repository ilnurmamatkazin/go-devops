package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

func (ms *MetricSender) sendMetrics(report string, chMetrics chan []models.Metric, chMetricsGopsutil chan []models.Metric) error {
	var (
		metrics         []models.Metric
		metricsGopsutil []models.Metric
	)

	interval, duration, err := utils.GetDataForTicker(report)
	if err != nil {
		log.Fatalf("Ошибка создания тикера")
		return err
	}

	tickerReport := time.NewTicker(time.Duration(interval) * duration)

	for {
		select {
		case <-ms.ctx.Done():
			return ms.ctx.Err()

		case metrics = <-chMetrics:
		case metricsGopsutil = <-chMetricsGopsutil:

		case <-tickerReport.C:
			for _, metric := range metrics {
				if err = ms.sendRequest(metric, "http://%s/update"); err != nil {
					return err
				}
			}

			for _, metric := range metricsGopsutil {
				if err = ms.sendRequest(metric, "http://%s/update"); err != nil {
					return err
				}
			}

			if err = ms.sendRequest(metrics, "http://%s/updates/"); err != nil {
				return err
			}

			if err = ms.sendRequest(metricsGopsutil, "http://%s/updates/"); err != nil {
				return err
			}

		}

	}
}

// func (ms *MetricSender) sendMetric(gctx context.Context, metric models.Metric) (err error) {
// 	ctx, cancel := context.WithTimeout(gctx, 1*time.Second)
// 	defer cancel()

// 	endpoint := fmt.Sprintf("http://%s/update", ms.cfg.Address)

// 	metric.Hash = utils.SetEncodeHash(metric.ID, metric.MetricType, ms.cfg.Key, metric.Delta, metric.Value)

// 	if err = ms.sendRequest(ctx, metric, endpoint); err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	return
// }

// func (ms MetricSender) sendArrayMetrics(rtm runtime.MemStats, pollCount int64) (err error) {
// 	ctx, cancel := context.WithTimeout(ms.ctx, 3*time.Second)
// 	defer cancel()

// 	endpoint := fmt.Sprintf("http://%s/updates/", ms.cfg.Address)
// 	metrics := make([]models.Metric, 0, 29)

// 	if err = ms.sendRequest(ctx, metrics, endpoint); err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	return
// }

// func (ms MetricSender) createMetric(metricType, id string, value float64) (metric models.Metric) {
// 	metric.ID = id
// 	metric.MetricType = metricType

// 	if metricType == "counter" {
// 		i := int64(value)
// 		metric.Delta = &i
// 	} else {
// 		metric.Value = &value
// 	}

// 	metric.Hash = utils.SetEncodeHash(metric.ID, metric.MetricType, ms.cfg.Key, metric.Delta, metric.Value)

// 	return
// }

func (ms MetricSender) sendRequest(data interface{}, layout string) (err error) {
	ctx, cancel := context.WithTimeout(ms.ctx, 5*time.Second)
	defer cancel()

	endpoint := fmt.Sprintf(layout, ms.cfg.Address)

	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, b)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		log.Println(err)
		return
	}

	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Set("Content-Type", "application/json")

	// отправляем запрос и получаем ответ
	response, err := ms.client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	// печатаем код ответа
	// fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()

	return
}

// func convertValue(value interface{}, metric *models.Metric) {
// 	var f float64

// 	switch metric.MetricType {
// 	case "counter":
// 		i := value.(int64)
// 		metric.Delta = &i
// 	case "gauge":
// 		switch value := value.(type) {
// 		case float64:
// 			f = value
// 		case uint64:
// 			f = float64(value)
// 		case uint32:
// 			f = float64(value)

// 		default:
// 		}

// 	}

// 	if (metric.ID == "GCCPUFraction") && (f == 0) {
// 		f = rand.Float64()
// 	} else if (metric.ID == "LastGC") && (f == 0) ||
// 		(metric.ID == "Lookups") && (f == 0) ||
// 		(metric.ID == "NumForcedGC") && (f == 0) ||
// 		(metric.ID == "NumGC") && (f == 0) ||
// 		(metric.ID == "PauseTotalNs") && (f == 0) {
// 		f = float64(rand.Intn(100))
// 	}

// 	metric.Value = &f
// }
