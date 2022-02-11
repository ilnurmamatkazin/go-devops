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
	}

	chMetrics := make(chan []models.Metric)
	chMetricsGopsutil := make(chan []models.Metric)

	ctx, done := context.WithCancel(context.Background())
	g, metricSender.ctx = errgroup.WithContext(ctx)

	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case <-signalChannel:
			done()
		case <-metricSender.ctx.Done():
			<-chMetricsGopsutil
			<-chMetrics
			return metricSender.ctx.Err()
		}

		return nil
	})

	g.Go(func() error {
		return metricSender.collectMetrics(metricSender.cfg.PollInterval, chMetrics)
	})

	g.Go(func() error {
		return metricSender.collectMetricsGopsutil(metricSender.cfg.PollInterval, chMetricsGopsutil)
	})

	g.Go(func() error {
		return metricSender.sendMetrics(metricSender.cfg.ReportInterval, chMetrics, chMetricsGopsutil)
	})

	if err := g.Wait(); err != nil && err != context.Canceled {
		if errors.Is(err, context.Canceled) {
			log.Printf("received error: %v", err)
		}
	}
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
