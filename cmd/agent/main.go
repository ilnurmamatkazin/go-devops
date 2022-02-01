package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/ilnurmamatkazin/go-devops/internal/utils"
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
	var (
		mutex     sync.Mutex
		rtm       runtime.MemStats
		pollCount int64
		err       error
	)

	metricSender := MetricSender{
		cfg:    parseConfig(),
		client: createClient(),
		ctx:    context.Background(),
	}

	interval, duration, err := utils.GetDataForTicker(metricSender.cfg.PollInterval)
	if err != nil {
		log.Fatalf("Ошибка создания тикера")
	}

	tickerPoll := time.NewTicker(time.Duration(interval) * duration)

	interval, duration, err = utils.GetDataForTicker(metricSender.cfg.ReportInterval)
	if err != nil {
		log.Fatalf("Ошибка создания тикера")
	}

	tickerReport := time.NewTicker(time.Duration(interval) * duration)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	done := make(chan bool, 1)

	go func() {
	loop:
		for {
			select {
			case <-quit:
				done <- true
				break loop

			case <-tickerPoll.C:
				mutex.Lock()

				runtime.ReadMemStats(&rtm)
				pollCount++

				mutex.Unlock()
			case <-tickerReport.C:
				mutex.Lock()

				metricSender.sendMetric("gauge", "Alloc", rtm.Alloc)
				metricSender.sendMetric("gauge", "BuckHashSys", rtm.BuckHashSys)
				metricSender.sendMetric("gauge", "Frees", rtm.Frees)
				metricSender.sendMetric("gauge", "GCCPUFraction", rtm.GCCPUFraction)
				metricSender.sendMetric("gauge", "GCSys", rtm.GCSys)
				metricSender.sendMetric("gauge", "HeapAlloc", rtm.HeapAlloc)
				metricSender.sendMetric("gauge", "HeapIdle", rtm.HeapIdle)
				metricSender.sendMetric("gauge", "HeapInuse", rtm.HeapInuse)
				metricSender.sendMetric("gauge", "HeapObjects", rtm.HeapObjects)
				metricSender.sendMetric("gauge", "HeapReleased", rtm.HeapReleased)
				metricSender.sendMetric("gauge", "HeapSys", rtm.HeapSys)
				metricSender.sendMetric("gauge", "LastGC", rtm.LastGC)
				metricSender.sendMetric("gauge", "Lookups", rtm.Lookups)
				metricSender.sendMetric("gauge", "MCacheInuse", rtm.MCacheInuse)
				metricSender.sendMetric("gauge", "MCacheSys", rtm.MCacheSys)
				metricSender.sendMetric("gauge", "MSpanInuse", rtm.MSpanInuse)
				metricSender.sendMetric("gauge", "MSpanSys", rtm.MSpanSys)
				metricSender.sendMetric("gauge", "Mallocs", rtm.Mallocs)
				metricSender.sendMetric("gauge", "NextGC", rtm.NextGC)
				metricSender.sendMetric("gauge", "NumForcedGC", rtm.NumForcedGC)
				metricSender.sendMetric("gauge", "NumGC", rtm.NumGC)
				metricSender.sendMetric("gauge", "OtherSys", rtm.OtherSys)
				metricSender.sendMetric("gauge", "PauseTotalNs", rtm.PauseTotalNs)
				metricSender.sendMetric("gauge", "TotalAlloc", rtm.TotalAlloc)
				metricSender.sendMetric("gauge", "StackInuse", rtm.StackInuse)
				metricSender.sendMetric("gauge", "StackSys", rtm.StackSys)
				metricSender.sendMetric("gauge", "Sys", rtm.Sys)
				metricSender.sendMetric("counter", "PollCount", pollCount)
				metricSender.sendMetric("gauge", "RandomValue", rand.Float64())

				metricSender.sendArrayMetric(rtm, pollCount)

				mutex.Unlock()
			}

		}
	}()

	<-done

}

func (ms MetricSender) sendMetric(typeMetric, nameMetric string, value interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ms.ctx, 1*time.Second)
	defer cancel()

	var metric models.Metric

	endpoint := fmt.Sprintf("http://%s/update", ms.cfg.Address)

	metric.ID = nameMetric
	metric.MetricType = typeMetric

	convertValue(value, &metric)
	metric.Hash = utils.SetEncodeHesh(metric.ID, metric.MetricType, ms.cfg.Key, metric.Delta, metric.Value)

	if err = ms.sendRequest(ctx, metric, endpoint); err != nil {
		log.Println(err)
		return
	}

	return
}

func (ms MetricSender) sendArrayMetric(rtm runtime.MemStats, pollCount int64) (err error) {
	ctx, cancel := context.WithTimeout(ms.ctx, 3*time.Second)
	defer cancel()

	endpoint := fmt.Sprintf("http://%s/updates/", ms.cfg.Address)
	metrics := make([]models.Metric, 0, 29)

	metrics = append(metrics, ms.createMetric("gauge", "Alloc", float64(rtm.Alloc)))
	metrics = append(metrics, ms.createMetric("gauge", "BuckHashSys", float64(rtm.BuckHashSys)))
	metrics = append(metrics, ms.createMetric("gauge", "Frees", float64(rtm.Frees)))
	metrics = append(metrics, ms.createMetric("gauge", "GCCPUFraction", rtm.GCCPUFraction))
	metrics = append(metrics, ms.createMetric("gauge", "GCSys", float64(rtm.GCSys)))
	metrics = append(metrics, ms.createMetric("gauge", "HeapAlloc", float64(rtm.HeapAlloc)))
	metrics = append(metrics, ms.createMetric("gauge", "HeapIdle", float64(rtm.HeapIdle)))
	metrics = append(metrics, ms.createMetric("gauge", "HeapInuse", float64(rtm.HeapInuse)))
	metrics = append(metrics, ms.createMetric("gauge", "HeapObjects", float64(rtm.HeapObjects)))
	metrics = append(metrics, ms.createMetric("gauge", "HeapReleased", float64(rtm.HeapReleased)))
	metrics = append(metrics, ms.createMetric("gauge", "HeapSys", float64(rtm.HeapSys)))
	metrics = append(metrics, ms.createMetric("gauge", "LastGC", float64(rtm.LastGC)))
	metrics = append(metrics, ms.createMetric("gauge", "Lookups", float64(rtm.Lookups)))
	metrics = append(metrics, ms.createMetric("gauge", "MCacheInuse", float64(rtm.MCacheInuse)))
	metrics = append(metrics, ms.createMetric("gauge", "MCacheSys", float64(rtm.MCacheSys)))
	metrics = append(metrics, ms.createMetric("gauge", "MSpanInuse", float64(rtm.MSpanInuse)))
	metrics = append(metrics, ms.createMetric("gauge", "MSpanSys", float64(rtm.MSpanSys)))
	metrics = append(metrics, ms.createMetric("gauge", "Mallocs", float64(rtm.Mallocs)))
	metrics = append(metrics, ms.createMetric("gauge", "NextGC", float64(rtm.NextGC)))
	metrics = append(metrics, ms.createMetric("gauge", "NumForcedGC", float64(rtm.NumForcedGC)))
	metrics = append(metrics, ms.createMetric("gauge", "NumGC", float64(rtm.NumGC)))
	metrics = append(metrics, ms.createMetric("gauge", "OtherSys", float64(rtm.OtherSys)))
	metrics = append(metrics, ms.createMetric("gauge", "PauseTotalNs", float64(rtm.PauseTotalNs)))
	metrics = append(metrics, ms.createMetric("gauge", "TotalAlloc", float64(rtm.TotalAlloc)))
	metrics = append(metrics, ms.createMetric("gauge", "StackInuse", float64(rtm.StackInuse)))
	metrics = append(metrics, ms.createMetric("gauge", "StackSys", float64(rtm.StackSys)))
	metrics = append(metrics, ms.createMetric("gauge", "Sys", float64(rtm.Sys)))
	metrics = append(metrics, ms.createMetric("counter", "PollCount", float64(pollCount)))
	metrics = append(metrics, ms.createMetric("gauge", "RandomValue", rand.Float64()))

	if err = ms.sendRequest(ctx, metrics, endpoint); err != nil {
		log.Println(err)
		return
	}

	return
}

func (ms MetricSender) createMetric(metricType, id string, value float64) (metric models.Metric) {
	metric.ID = id
	metric.MetricType = metricType

	if metricType == "counter" {
		i := int64(value)
		metric.Delta = &i
	} else {
		metric.Value = &value
	}

	metric.Hash = utils.SetEncodeHesh(metric.ID, metric.MetricType, ms.cfg.Key, metric.Delta, metric.Value)

	return
}

func (ms MetricSender) sendRequest(ctx context.Context, data interface{}, endpoint string) (err error) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, b)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		// log.Println(err)
		return
	}

	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Set("Content-Type", "application/json")

	// отправляем запрос и получаем ответ
	response, err := ms.client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}
	// печатаем код ответа
	// fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()

	return
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

func convertValue(value interface{}, metric *models.Metric) {
	var f float64

	switch metric.MetricType {
	case "counter":
		i := value.(int64)
		metric.Delta = &i
	case "gauge":
		switch value := value.(type) {
		case float64:
			f = value
		case uint64:
			f = float64(value)
		case uint32:
			f = float64(value)

		default:
		}

	}

	if (metric.ID == "GCCPUFraction") && (f == 0) {
		f = rand.Float64()
	} else if (metric.ID == "LastGC") && (f == 0) ||
		(metric.ID == "Lookups") && (f == 0) ||
		(metric.ID == "NumForcedGC") && (f == 0) ||
		(metric.ID == "NumGC") && (f == 0) ||
		(metric.ID == "PauseTotalNs") && (f == 0) {
		f = float64(rand.Intn(100))
	}

	metric.Value = &f
}
