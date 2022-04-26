package main

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"

// 	"github.com/shirou/gopsutil/v3/cpu"
// 	"github.com/shirou/gopsutil/v3/mem"
// )

// // collectMetricsGopsutil функция, реализуящая сбор gopsutil метрик.
// func (ms *MetricSend) collectMetricsGopsutil(ctx context.Context, tickerPoll *time.Ticker, chMetrics chan []models.Metric) (err error) {
// 	var (
// 		v              *mem.VirtualMemoryStat
// 		percentage     []float64
// 		cpuUtilization float64
// 	)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()

// 		case <-tickerPoll.C:
// 			metrics := make([]models.Metric, 0, 3)

// 			if v, err = mem.VirtualMemory(); err != nil {
// 				log.Print("Ошибка доступа к метрикам mem.VirtualMemoryStat")
// 				return
// 			}

// 			if _, err = cpu.Info(); err != nil {
// 				log.Print("Ошибка доступа к метрикам cpu.InfoStat")
// 				return
// 			}

// 			if percentage, err = cpu.Percent(0, true); err != nil {
// 				log.Print("Ошибка доступа к метрикам cpu.InfoStat")
// 				return
// 			}

// 			for _, percent := range percentage {
// 				cpuUtilization = cpuUtilization + percent
// 			}

// 			metrics = append(metrics, ms.createMetric("gauge", "TotalMemory", float64(v.Total), 0))
// 			metrics = append(metrics, ms.createMetric("gauge", "FreeMemory", float64(v.Free), 0))
// 			metrics = append(metrics, ms.createMetric("gauge", "CPUutilization1", cpuUtilization, 0))

// 			chMetrics <- metrics
// 		}
// 	}
// }
