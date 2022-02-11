package main

import (
	"log"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

//https://github.com/shirou/gopsutil/blob/master/cpu/cpu.go
//https://www.socketloop.com/tutorials/golang-get-hardware-information-such-as-disk-memory-and-cpu-usage
// https://developpaper.com/get-system-performance-data-with-go-language-gopsutil-library/
// https://stackoverflow.com/questions/61201928/how-to-get-total-ram-from-golang-code-on-mac

func (ms *MetricSender) collectMetricsGopsutil(poll string, chMetrics chan []models.Metric) (err error) {
	var (
		v *mem.VirtualMemoryStat
		// c              []cpu.InfoStat
		percentage     []float64
		cpuUtilization float64
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
			metrics := make([]models.Metric, 0, 3)

			if v, err = mem.VirtualMemory(); err != nil {
				log.Print("Ошибка доступа к метрикам mem.VirtualMemoryStat")
				return
			}

			if _, err = cpu.Info(); err != nil {
				log.Print("Ошибка доступа к метрикам cpu.InfoStat")
				return
			}

			if percentage, err = cpu.Percent(0, true); err != nil {
				log.Print("Ошибка доступа к метрикам cpu.InfoStat")
				return
			}

			for _, percent := range percentage {
				cpuUtilization = cpuUtilization + percent
			}

			metrics = append(metrics, ms.createMetric("gauge", "TotalMemory", float64(v.Total), 0))
			metrics = append(metrics, ms.createMetric("gauge", "FreeMemory", float64(v.Free), 0))
			metrics = append(metrics, ms.createMetric("gauge", "CPUutilization1", cpuUtilization, 0))

			chMetrics <- metrics
		}
	}
}