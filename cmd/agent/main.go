// Сервис сбора системных метрик и отправки их на сервер.
package main

import (
	"context"
	"encoding/json"
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
	"github.com/ilnurmamatkazin/go-devops/cmd/agent/transport/grpc"
	"github.com/ilnurmamatkazin/go-devops/internal/model"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
	"golang.org/x/sync/errgroup"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

const (
	Address = "127.0.0.1:8080" // адрес принимающего сервера
	// PollInterval   = "20000000n"      // период сбора  метрик
	// ReportInterval = "100000000n"     // период отправки метрик
	PollInterval   = "2s"                 // период сбора  метрик
	ReportInterval = "10s"                // период отправки метрик
	Key            = ""                   // ключ для формирования подписи
	PublicKey      = "../keys/public.pem" // открытый ключ для шифрования
	Config         = "./config.json"      // имя json файла с конфигурацией
	NeedGRPC       = true                 // флаг, указывающий протокол передачи данных
	AddressGRPC    = "localhost:8000"     // адрес grpc сервера
)

type RequestSender interface {
	Send(ctx context.Context, data interface{}, layout string) error
}

type RequestSend struct {
	cfg        models.Config    // поле с конфигурационными данными
	client     *http.Client     // поле с созданным http клиентом, для отправки данных на сервер
	grpcClient *grpc.GRPCClient // поле с созданным grpc клиентом, для отправки данных на сервер
}

// MetricSend вспомогательная структура, для проброса вспомогательных структур
type MetricSend struct {
	cfg    models.Config // поле с конфигурационными данными
	sender RequestSender
}

func main() {
	build := model.NewBuild(buildVersion, buildDate, buildCommit)
	build.Print()

	go http.ListenAndServe(":6060", nil)

	var (
		g        *errgroup.Group
		ctxGroup context.Context
	)

	cfg := parseConfig()

	grpcClient, err := grpc.NewGRPCClient(cfg)
	if err != nil {
		log.Println(err)
	}

	metricSend := MetricSend{
		cfg: cfg,
		sender: &RequestSend{
			cfg:        cfg,
			client:     createClient(),
			grpcClient: grpcClient,
		},
	}

	chMetrics := make(chan []models.Metric)
	chMetricsGopsutil := make(chan []models.Metric)

	ctx, done := context.WithCancel(context.Background())
	g, ctxGroup = errgroup.WithContext(ctx)

	tickerPoll, err := getTicker(metricSend.cfg.PollInterval)
	if err != nil {
		log.Fatalf("Ошибка создания тикера")
		return
	}

	tickerReport, err := getTicker(metricSend.cfg.ReportInterval)
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
		case <-ctxGroup.Done():
			tickerPoll.Stop()
			tickerReport.Stop()

			for i := range chMetrics {
				log.Println(i)
			}

			for i := range chMetricsGopsutil {
				log.Println(i)
			}

			return ctxGroup.Err()
		}

		return nil
	})

	g.Go(func() error {
		err := metricSend.collectMetrics(ctxGroup, tickerPoll, chMetrics)
		close(chMetrics)

		return err
	})

	g.Go(func() error {
		err := metricSend.collectMetricsGopsutil(ctxGroup, tickerPoll, chMetricsGopsutil)
		close(chMetricsGopsutil)

		return err
	})

	g.Go(func() error {
		err := metricSend.sendMetrics(ctxGroup, tickerReport, chMetrics, chMetricsGopsutil)

		return err
	})

	if err := g.Wait(); err != nil {
		log.Printf("received error: %v", err)

	}

	grpcClient.Close()

}

// parseConfig парсит флаги командной строки и получает данные из env переменных.
// ENV переменные имеют приоритет перед флагами.
func parseConfig() (cfg models.Config) {
	config := flag.String("c", Config, "a json config")

	if *config != "" {
		data, err := os.ReadFile(*config)

		if err != nil {
			log.Printf("os.ReadFile error: %s", err.Error())
		} else {
			err = json.Unmarshal(data, &cfg)
			if err != nil {
				log.Printf("json.Unmarshal error: %s", err.Error())
			}
		}
	}

	address := flag.String("a", Address, "a address")
	reportInterval := flag.String("r", ReportInterval, "a report_interval")
	pollInterval := flag.String("p", PollInterval, "a poll_interval")
	key := flag.String("k", Key, "a secret key")
	publicKey := flag.String("crypto-key", PublicKey, "a crypto key")

	flag.Parse()

	cfg.Address = *address
	cfg.ReportInterval = *reportInterval
	cfg.PollInterval = *pollInterval
	cfg.Key = *key
	cfg.PublicKey = *publicKey

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
func (ms *MetricSend) createMetric(metricType, id string, value float64, delta int64) (metric models.Metric) {
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

	return time.NewTicker(time.Duration(interval) * duration), nil
}
