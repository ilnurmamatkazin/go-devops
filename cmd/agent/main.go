// Сервис сбора системных метрик и отправки их на сервер.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
	"golang.org/x/sync/errgroup"
)

const (
	Address        = "127.0.0.1:8080" // адрес принимающего сервера
	PollInterval   = "100s"           // период отправки метрик
	ReportInterval = "1000s"          // период сбора метрик
	Key            = ""               // ключ для формирования подписи
)

// MetricSender вспомогательная структура, для проброса вспомогательных структур
type MetricSender struct {
	cfg    models.Config   // поле с конфигурационными данными
	client *http.Client    // поле с созданным http клиентом, для отправки данных на сервер
	ctx    context.Context // поле с системны контекстом
}

func main() {
	go http.ListenAndServe(":6060", nil)

	var g *errgroup.Group

	metricSender := MetricSender{
		cfg:    parseConfig(),
		client: createClient(),
	}

	chMetrics := make(chan []models.Metric)
	chMetricsGopsutil := make(chan []models.Metric)

	ctx, done := context.WithCancel(context.Background())
	g, metricSender.ctx = errgroup.WithContext(ctx)

	tickerPoll, err := getTicker(metricSender.cfg.PollInterval)
	if err != nil {
		log.Fatalf("Ошибка создания тикера")
		return
	}

	tickerReport, err := getTicker(metricSender.cfg.ReportInterval)
	if err != nil {
		log.Fatalf("Ошибка создания тикера")
		return
	}

	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case <-signalChannel:
			done()
			for i := range chMetrics {
				log.Println(i)
			}

			for i := range chMetricsGopsutil {
				log.Println(i)
			}
		case <-metricSender.ctx.Done():
			tickerPoll.Stop()
			tickerReport.Stop()

			for i := range chMetrics {
				log.Println(i)
			}

			for i := range chMetricsGopsutil {
				log.Println(i)
			}

			return metricSender.ctx.Err()
		}

		return nil
	})

	g.Go(func() error {
		err := metricSender.collectMetrics(tickerPoll, chMetrics)
		close(chMetrics)

		return err
	})

	g.Go(func() error {
		err := metricSender.collectMetricsGopsutil(tickerPoll, chMetricsGopsutil)
		close(chMetricsGopsutil)

		return err
	})

	g.Go(func() error {
		err := metricSender.sendMetrics(tickerReport, chMetrics, chMetricsGopsutil)

		return err
	})

	if err := g.Wait(); err != nil {
		log.Printf("received error: %v", err)

	}

}

// parseConfig парсит флаги командной строки и получает данные из env переменных.
// ENV переменные имеют приоритет перед флагами.
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

// createClient конструируем HTTP-клиент.
func createClient() *http.Client {
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

// createMetric внутренняя функция со созданию метрики
func (ms *MetricSender) createMetric(metricType, id string, value float64, delta int64) (metric models.Metric) {
	if metricType == "counter" {
		metric = models.Metric{MetricType: metricType, ID: id, Delta: &delta}
	} else {
		metric = models.Metric{MetricType: metricType, ID: id, Value: &value}
	}

	metric.Hash = utils.SetEncodeHash(metric.ID, metric.MetricType, ms.cfg.Key, metric.Delta, metric.Value)

	return
}

// getTicker внутренняя функция по созданию тикера
func getTicker(strInterval string) (*time.Ticker, error) {
	interval, duration, err := utils.GetDataForTicker(strInterval)
	if err != nil {
		return nil, err
	}

	return time.NewTicker(time.Duration(interval) * duration / 1000), nil
}
