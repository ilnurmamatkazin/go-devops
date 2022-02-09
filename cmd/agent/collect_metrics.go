package main

import (
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

func collectMetrics(poll string, quit chan os.Signal, chMetrics chan []models.Metric) {
	var (
		rtm       runtime.MemStats
		value     float64
		pollCount int64
	)

	interval, duration, err := utils.GetDataForTicker(poll)
	if err != nil {
		log.Fatalf("Ошибка создания тикера")
		return
	}

	metrics := make([]models.Metric, 0, 29)

	tickerPoll := time.NewTicker(time.Duration(interval) * duration)

	for {
		select {
		case <-quit:
			return

		case <-tickerPoll.C:
			value = float64(rtm.Alloc)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "Alloc", Value: &value})

			value = float64(rtm.BuckHashSys)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "BuckHashSys", Value: &value})

			value = float64(rtm.Frees)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "Frees", Value: &value})

			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "GCCPUFraction", Value: &rtm.GCCPUFraction})

			value = float64(rtm.GCSys)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "GCSys", Value: &value})

			value = float64(rtm.HeapAlloc)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "HeapAlloc", Value: &value})

			value = float64(rtm.HeapIdle)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "HeapIdle", Value: &value})

			value = float64(rtm.HeapInuse)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "HeapInuse", Value: &value})

			value = float64(rtm.HeapObjects)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "HeapObjects", Value: &value})

			value = float64(rtm.HeapReleased)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "HeapReleased", Value: &value})

			value = float64(rtm.HeapSys)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "HeapSys", Value: &value})

			value = float64(rtm.LastGC)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "LastGC", Value: &value})

			value = float64(rtm.Lookups)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "Lookups", Value: &value})

			value = float64(rtm.MCacheInuse)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "MCacheInuse", Value: &value})

			value = float64(rtm.MCacheSys)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "MCacheSys", Value: &value})

			value = float64(rtm.MSpanInuse)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "MSpanInuse", Value: &value})

			value = float64(rtm.MSpanSys)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "MSpanSys", Value: &value})

			value = float64(rtm.Mallocs)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "Mallocs", Value: &value})

			value = float64(rtm.NextGC)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "NextGC", Value: &value})

			value = float64(rtm.NumForcedGC)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "NumForcedGC", Value: &value})

			value = float64(rtm.NumGC)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "NumGC", Value: &value})

			value = float64(rtm.OtherSys)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "OtherSys", Value: &value})

			value = float64(rtm.PauseTotalNs)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "PauseTotalNs", Value: &value})

			value = float64(rtm.TotalAlloc)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "TotalAlloc", Value: &value})

			value = float64(rtm.StackInuse)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "StackInuse", Value: &value})

			value = float64(rtm.StackSys)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "StackSys", Value: &value})

			value = float64(rtm.Sys)
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "Sys", Value: &value})

			pollCount++
			metrics = append(metrics, models.Metric{MetricType: "counter", ID: "PollCount", Delta: &pollCount})

			value = rand.Float64()
			metrics = append(metrics, models.Metric{MetricType: "gauge", ID: "RandomValue", Value: &value})

		}
	}
}
