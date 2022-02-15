package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

func (ms *MetricSender) collectMetrics(poll string, chMetrics chan []models.Metric) (err error) {
	var (
		rtm       runtime.MemStats
		pollCount int64
	)

	interval, duration, err := utils.GetDataForTicker(poll)
	if err != nil {
		log.Fatalf("Ошибка создания тикера")
		return
	}

	tickerPoll := time.NewTicker(time.Duration(interval) * duration)

	for {
		select {
		case <-ms.ctx.Done():
			return ms.ctx.Err()

		case <-tickerPoll.C:
			runtime.ReadMemStats(&rtm)

			metrics := make([]models.Metric, 0, 29)

			metrics = append(metrics, ms.createMetric("gauge", "Alloc", float64(rtm.Alloc), 0))
			metrics = append(metrics, ms.createMetric("gauge", "BuckHashSys", float64(rtm.BuckHashSys), 0))
			metrics = append(metrics, ms.createMetric("gauge", "Frees", float64(rtm.Frees), 0))
			metrics = append(metrics, ms.createMetric("gauge", "GCCPUFraction", rtm.GCCPUFraction, 0))
			metrics = append(metrics, ms.createMetric("gauge", "GCSys", float64(rtm.GCSys), 0))
			metrics = append(metrics, ms.createMetric("gauge", "HeapAlloc", float64(rtm.HeapAlloc), 0))
			metrics = append(metrics, ms.createMetric("gauge", "HeapIdle", float64(rtm.HeapIdle), 0))
			metrics = append(metrics, ms.createMetric("gauge", "HeapInuse", float64(rtm.HeapInuse), 0))
			metrics = append(metrics, ms.createMetric("gauge", "HeapObjects", float64(rtm.HeapObjects), 0))
			metrics = append(metrics, ms.createMetric("gauge", "HeapReleased", float64(rtm.HeapReleased), 0))
			metrics = append(metrics, ms.createMetric("gauge", "HeapSys", float64(rtm.HeapSys), 0))
			metrics = append(metrics, ms.createMetric("gauge", "LastGC", float64(rtm.LastGC), 0))
			metrics = append(metrics, ms.createMetric("gauge", "Lookups", float64(rtm.Lookups), 0))
			metrics = append(metrics, ms.createMetric("gauge", "MCacheInuse", float64(rtm.MCacheInuse), 0))
			metrics = append(metrics, ms.createMetric("gauge", "MCacheSys", float64(rtm.MCacheSys), 0))
			metrics = append(metrics, ms.createMetric("gauge", "MSpanInuse", float64(rtm.MSpanInuse), 0))
			metrics = append(metrics, ms.createMetric("gauge", "MSpanSys", float64(rtm.MSpanSys), 0))
			metrics = append(metrics, ms.createMetric("gauge", "Mallocs", float64(rtm.Mallocs), 0))
			metrics = append(metrics, ms.createMetric("gauge", "NextGC", float64(rtm.NextGC), 0))
			metrics = append(metrics, ms.createMetric("gauge", "NumForcedGC", float64(rtm.NumForcedGC), 0))
			metrics = append(metrics, ms.createMetric("gauge", "NumGC", float64(rtm.NumGC), 0))
			metrics = append(metrics, ms.createMetric("gauge", "OtherSys", float64(rtm.OtherSys), 0))
			metrics = append(metrics, ms.createMetric("gauge", "PauseTotalNs", float64(rtm.PauseTotalNs), 0))
			metrics = append(metrics, ms.createMetric("gauge", "TotalAlloc", float64(rtm.TotalAlloc), 0))
			metrics = append(metrics, ms.createMetric("gauge", "StackInuse", float64(rtm.StackInuse), 0))
			metrics = append(metrics, ms.createMetric("gauge", "StackSys", float64(rtm.StackSys), 0))
			metrics = append(metrics, ms.createMetric("gauge", "Sys", float64(rtm.Sys), 0))

			pollCount++
			metrics = append(metrics, ms.createMetric("counter", "PollCount", 0, pollCount))
			metrics = append(metrics, ms.createMetric("gauge", "RandomValue", rand.Float64(), 0))

			chMetrics <- metrics
		}
	}
}
