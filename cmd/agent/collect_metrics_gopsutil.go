package main

import (
	"log"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
)

//https://github.com/shirou/gopsutil/blob/master/cpu/cpu.go
//https://www.socketloop.com/tutorials/golang-get-hardware-information-such-as-disk-memory-and-cpu-usage
// https://developpaper.com/get-system-performance-data-with-go-language-gopsutil-library/
// https://stackoverflow.com/questions/61201928/how-to-get-total-ram-from-golang-code-on-mac

func (ms *MetricSender) collectMetricsGopsutil(poll string, chMetrics chan []models.Metric) (err error) {
	var (
		rtm gopsutil.MemStats
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
			// runtime.ReadMemStats(&rtm)

			metrics := make([]models.Metric, 0, 29)

			metrics = append(metrics, ms.createMetric("gauge", "Alloc", float64(rtm.Alloc), 0))
			metrics = append(metrics, ms.createMetric("gauge", "BuckHashSys", float64(rtm.BuckHashSys), 0))
			metrics = append(metrics, ms.createMetric("gauge", "Frees", float64(rtm.Frees), 0))
			metrics = append(metrics, ms.createMetric("gauge", "GCCPUFraction", rtm.GCCPUFraction, 0))

			chMetrics <- metrics
		}
	}
}
