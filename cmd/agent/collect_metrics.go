package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
)

func (ms *MetricSender) collectMetrics(tickerPoll *time.Ticker, chMetrics chan []models.Metric) (err error) {
	var (
		rtm       runtime.MemStats
		pollCount int64
	)

	metrics := make([]models.Metric, 30, 30)

	fmt.Println("&&&&&&&", len(metrics), cap(metrics))

	for {
		select {
		case <-ms.ctx.Done():
			return ms.ctx.Err()

		case <-tickerPoll.C:
			runtime.ReadMemStats(&rtm)

			// metrics = metrics[:0]

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
			metrics = append(metrics, ms.createMetric("gauge", "Alloc", float64(rtm.Alloc), 0))
			metrics = append(metrics, ms.createMetric("gauge", "RandomValue", rand.Float64(), 0))

			// metrics[0] = *ms.createMetric("gauge", "Alloc", float64(rtm.Alloc), 0)
			// metrics[1] = *ms.createMetric("gauge", "BuckHashSys", float64(rtm.BuckHashSys), 0)
			// metrics[2] = *ms.createMetric("gauge", "Frees", float64(rtm.Frees), 0)
			// metrics[3] = *ms.createMetric("gauge", "GCCPUFraction", rtm.GCCPUFraction, 0)
			// metrics[4] = *ms.createMetric("gauge", "GCSys", float64(rtm.GCSys), 0)
			// metrics[5] = *ms.createMetric("gauge", "HeapAlloc", float64(rtm.HeapAlloc), 0)
			// metrics[6] = *ms.createMetric("gauge", "HeapIdle", float64(rtm.HeapIdle), 0)
			// metrics[7] = *ms.createMetric("gauge", "HeapInuse", float64(rtm.HeapInuse), 0)
			// metrics[8] = *ms.createMetric("gauge", "HeapObjects", float64(rtm.HeapObjects), 0)
			// metrics[9] = *ms.createMetric("gauge", "HeapReleased", float64(rtm.HeapReleased), 0)
			// metrics[10] = *ms.createMetric("gauge", "HeapSys", float64(rtm.HeapSys), 0)
			// metrics[11] = *ms.createMetric("gauge", "LastGC", float64(rtm.LastGC), 0)
			// metrics[12] = *ms.createMetric("gauge", "Lookups", float64(rtm.Lookups), 0)
			// metrics[13] = *ms.createMetric("gauge", "MCacheInuse", float64(rtm.MCacheInuse), 0)
			// metrics[14] = *ms.createMetric("gauge", "MCacheSys", float64(rtm.MCacheSys), 0)
			// metrics[15] = *ms.createMetric("gauge", "MSpanInuse", float64(rtm.MSpanInuse), 0)
			// metrics[16] = *ms.createMetric("gauge", "MSpanSys", float64(rtm.MSpanSys), 0)
			// metrics[17] = *ms.createMetric("gauge", "Mallocs", float64(rtm.Mallocs), 0)
			// metrics[18] = *ms.createMetric("gauge", "NextGC", float64(rtm.NextGC), 0)
			// metrics[19] = *ms.createMetric("gauge", "NumForcedGC", float64(rtm.NumForcedGC), 0)
			// metrics[20] = *ms.createMetric("gauge", "NumGC", float64(rtm.NumGC), 0)
			// metrics[21] = *ms.createMetric("gauge", "OtherSys", float64(rtm.OtherSys), 0)
			// metrics[22] = *ms.createMetric("gauge", "PauseTotalNs", float64(rtm.PauseTotalNs), 0)
			// metrics[23] = *ms.createMetric("gauge", "TotalAlloc", float64(rtm.TotalAlloc), 0)
			// metrics[24] = *ms.createMetric("gauge", "StackInuse", float64(rtm.StackInuse), 0)
			// metrics[25] = *ms.createMetric("gauge", "StackSys", float64(rtm.StackSys), 0)
			// metrics[26] = *ms.createMetric("gauge", "Sys", float64(rtm.Sys), 0)

			// pollCount++
			// metrics[27] = *ms.createMetric("counter", "PollCount", 0, pollCount)
			// metrics[28] = *ms.createMetric("gauge", "Alloc", float64(rtm.Alloc), 0)
			// metrics[29] = *ms.createMetric("gauge", "RandomValue", rand.Float64(), 0)

			chMetrics <- metrics
		}
	}
}
