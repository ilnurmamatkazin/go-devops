package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
	"golang.org/x/sync/errgroup"
)

const (
	Address        = "127.0.0.1:8080"
	PollInterval   = "2s"
	ReportInterval = "10s"
	Key            = ""
)

type MetricSender struct {
	cfg    models.Config
	client *http.Client
	ctx    context.Context
}

func main() {
	var g *errgroup.Group

	metricSender := MetricSender{
		cfg:    parseConfig(),
		client: createClient(),
		//ctx:    context.Background(),
	}

	ctx, done := context.WithCancel(context.Background())
	g, metricSender.ctx = errgroup.WithContext(ctx)

	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case <-signalChannel:
			done()
		case <-metricSender.ctx.Done():
			return metricSender.ctx.Err()
		}

		return nil
	})

	chMetrics := make(chan []models.Metric)
	chMetricsGopsutil := make(chan []models.Metric)

	g.Go(func() error {
		return metricSender.collectMetrics(metricSender.cfg.PollInterval, chMetrics)
	})

	g.Go(func() error {
		return metricSender.collectMetricsGopsutil(metricSender.cfg.PollInterval, chMetricsGopsutil)
	})

	g.Go(func() error {
		return metricSender.sendMetrics(metricSender.cfg.PollInterval, chMetrics, chMetricsGopsutil)
	})

	if err := g.Wait(); err != nil && err != context.Canceled {
		if errors.Is(err, context.Canceled) {
			log.Printf("received error: %v", err)
		}
	}

	// done := make(chan bool, 1)

	// go func() {
	// loop:
	// 	for {
	// 		select {
	// 		case <-quit:
	// 			done <- true
	// 			break loop

	// 		case <-chMetrics:
	// 			mutex.Lock()

	// 			runtime.ReadMemStats(&rtm)
	// 			pollCount++

	// 			mutex.Unlock()
	// 		case <-tickerReport.C:
	// 			mutex.Lock()

	// 			metricSender.sendMetric("gauge", "Alloc", rtm.Alloc)
	// 			metricSender.sendMetric("gauge", "BuckHashSys", rtm.BuckHashSys)
	// 			metricSender.sendMetric("gauge", "Frees", rtm.Frees)
	// 			metricSender.sendMetric("gauge", "GCCPUFraction", rtm.GCCPUFraction)
	// 			metricSender.sendMetric("gauge", "GCSys", rtm.GCSys)
	// 			metricSender.sendMetric("gauge", "HeapAlloc", rtm.HeapAlloc)
	// 			metricSender.sendMetric("gauge", "HeapIdle", rtm.HeapIdle)
	// 			metricSender.sendMetric("gauge", "HeapInuse", rtm.HeapInuse)
	// 			metricSender.sendMetric("gauge", "HeapObjects", rtm.HeapObjects)
	// 			metricSender.sendMetric("gauge", "HeapReleased", rtm.HeapReleased)
	// 			metricSender.sendMetric("gauge", "HeapSys", rtm.HeapSys)
	// 			metricSender.sendMetric("gauge", "LastGC", rtm.LastGC)
	// 			metricSender.sendMetric("gauge", "Lookups", rtm.Lookups)
	// 			metricSender.sendMetric("gauge", "MCacheInuse", rtm.MCacheInuse)
	// 			metricSender.sendMetric("gauge", "MCacheSys", rtm.MCacheSys)
	// 			metricSender.sendMetric("gauge", "MSpanInuse", rtm.MSpanInuse)
	// 			metricSender.sendMetric("gauge", "MSpanSys", rtm.MSpanSys)
	// 			metricSender.sendMetric("gauge", "Mallocs", rtm.Mallocs)
	// 			metricSender.sendMetric("gauge", "NextGC", rtm.NextGC)
	// 			metricSender.sendMetric("gauge", "NumForcedGC", rtm.NumForcedGC)
	// 			metricSender.sendMetric("gauge", "NumGC", rtm.NumGC)
	// 			metricSender.sendMetric("gauge", "OtherSys", rtm.OtherSys)
	// 			metricSender.sendMetric("gauge", "PauseTotalNs", rtm.PauseTotalNs)
	// 			metricSender.sendMetric("gauge", "TotalAlloc", rtm.TotalAlloc)
	// 			metricSender.sendMetric("gauge", "StackInuse", rtm.StackInuse)
	// 			metricSender.sendMetric("gauge", "StackSys", rtm.StackSys)
	// 			metricSender.sendMetric("gauge", "Sys", rtm.Sys)
	// 			metricSender.sendMetric("counter", "PollCount", pollCount)
	// 			metricSender.sendMetric("gauge", "RandomValue", rand.Float64())

	// 			metricSender.sendArrayMetric(rtm, pollCount)

	// 			mutex.Unlock()
	// 		}

	// 	}
	// }()

	// <-done
	// wg.Wait()

}

func parseConfig() (cfg models.Config) {
	address := flag.String("a", Address, "a address")
	reportInterval := flag.String("r", ReportInterval, "a report_interval")
	pollInterval := flag.String("p", PollInterval, "a poll_interval")
	key := flag.String("k", Key, "a secret key")

	flag.Parse()

	cfg.Address = *address
	cfg.ReportInterval = *reportInterval
	cfg.PollInterval = *pollInterval
	cfg.Key = *key

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("env.Parse error: %s", err.Error())
	}

	return
}

func createClient() *http.Client {
	// конструируем HTTP-клиент
	client := &http.Client{}
	client.Timeout = time.Second * 2

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 2 {
			log.Fatalf("Количество редиректов %d больше 2", len(via))
		}
		return nil
	}

	transport := &http.Transport{}
	transport.MaxIdleConns = 20
	client.Transport = transport

	return client
}

func (ms *MetricSender) createMetric(metricType, id string, value float64, delta int64) (metric models.Metric) {
	if metricType == "counter" {
		metric = models.Metric{MetricType: metricType, ID: id, Delta: &delta}
	} else {
		metric = models.Metric{MetricType: metricType, ID: id, Value: &value}
	}

	metric.Hash = utils.SetEncodeHash(metric.ID, metric.MetricType, ms.cfg.Key, metric.Delta, metric.Value)

	return
}
