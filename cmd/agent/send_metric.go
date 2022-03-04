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
)

// sendMetrics функция, реализующая отправку метрик на сервер.
func (ms *MetricSender) sendMetrics(tickerReport *time.Ticker, chMetrics chan []models.Metric, chMetricsGopsutil chan []models.Metric) (err error) {
	var (
		metrics         []models.Metric
		metricsGopsutil []models.Metric
	)

	for {
		select {
		case <-ms.ctx.Done():
			return ms.ctx.Err()

		case metrics = <-chMetrics:
		case metricsGopsutil = <-chMetricsGopsutil:

		case <-tickerReport.C:
			for _, metric := range metrics {
				if err = ms.sendRequest(metric, "http://%s/update"); err != nil {
					return
				}
			}

			for _, metric := range metricsGopsutil {
				if err = ms.sendRequest(metric, "http://%s/update"); err != nil {
					return
				}
			}

			if len(metrics) > 0 {
				if err = ms.sendRequest(metrics, "http://%s/updates/"); err != nil {
					return
				}
			}

			if len(metricsGopsutil) > 0 {
				if err = ms.sendRequest(metricsGopsutil, "http://%s/updates/"); err != nil {
					return
				}
			}

		}

	}
}

// sendRequest функция, реализующая создание запроса для отправки метрик на сервер.
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

	request.Header.Set("Content-Type", "application/json")

	// отправляем запрос и получаем ответ
	response, ok := ms.client.Do(request)
	if ok != nil {
		log.Println(ok)
	} else {
		// печатаем код ответа
		fmt.Println("Статус-код ", response.Status)
		defer response.Body.Close()
	}

	return
}
