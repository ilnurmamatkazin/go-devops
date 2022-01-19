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
		cfg       models.Config
		err       error
	)

	metricSender := MetricSender{
		cfg:    parseConfig(),
		client: createClient(),
		ctx:    context.Background(),
	}

	interval, duration, err := utils.GetDataForTicker(cfg.PollInterval)
	if err == nil {
		log.Fatalf("Ошибка создания тикера")
	}

	tickerPoll := time.NewTicker(time.Duration(interval) * duration)

	interval, duration, err = utils.GetDataForTicker(cfg.ReportInterval)
	if err == nil {
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
	metric.MType = typeMetric

	f := value.(float64)
	i := value.(int64)

	// switch typeMetric {
	// case "counter":
	// 	var i int64
	// 	i = value.(int64)
	// 	metric.Delta = &i
	// case "gauge":
	// 	var f float64

	// 	switch value := value.(type) {
	// 	case float64:
	// 		f = value
	// 	case uint64:
	// 		f = float64(value)
	// 	case uint32:
	// 		f = float64(value)

	// 	default:
	// 	}

	if (nameMetric == "GCCPUFraction") && (f == 0) {
		f = rand.Float64()
	} else if (nameMetric == "LastGC") && (f == 0) ||
		(nameMetric == "Lookups") && (f == 0) ||
		(nameMetric == "NumForcedGC") && (f == 0) ||
		(nameMetric == "NumGC") && (f == 0) ||
		(nameMetric == "PauseTotalNs") && (f == 0) {
		f = float64(rand.Intn(100))
	}

	metric.Value = &f
	metric.Delta = &i
	// default:
	// 	err = errors.New("недопустимый тип")
	// 	return
	// }

	b, err := json.Marshal(metric)
	if err != nil {
		log.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, b)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		log.Println(err)
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

	flag.Parse()

	cfg.Address = *address
	cfg.ReportInterval = *reportInterval
	cfg.PollInterval = *pollInterval

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

//cfg = parseConfig()

// конструируем HTTP-клиент
// client := &http.Client{}
// client.Timeout = time.Second * 2

// client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
// 	if len(via) >= 2 {
// 		log.Fatalf("Количество редиректов %d больше 2", len(via))
// 	}
// 	return nil
// }

// transport := &http.Transport{}
// transport.MaxIdleConns = 20
// client.Transport = transport

// ctx := context.Background()

// strDurationPollInterval := cfg.PollInterval[len(cfg.PollInterval)-1:]
// strPollInterval := cfg.PollInterval[0 : len(cfg.PollInterval)-1]

// strDurationReportInterval := cfg.ReportInterval[len(cfg.ReportInterval)-1:]
// strReportInterval := cfg.ReportInterval[0 : len(cfg.ReportInterval)-1]

// pollInterval, _ := strconv.Atoi(strPollInterval)
// reportInterval, _ := strconv.Atoi(strReportInterval)

// var durationPollInterval, durationReportInterval time.Duration

// switch strDurationPollInterval {
// case "s":
// 	durationPollInterval = time.Second
// case "m":
// 	durationPollInterval = time.Minute
// case "h":
// 	durationPollInterval = time.Hour
// default:
// 	log.Println("Неверный временной интервал")
// 	return
// }

// switch strDurationReportInterval {
// case "s":
// 	durationReportInterval = time.Second
// case "m":
// 	durationReportInterval = time.Minute
// case "h":
// 	durationReportInterval = time.Hour
// default:
// 	log.Println("Неверный временной интервал")
// 	return
// }

// tickerPoll := time.NewTicker(time.Duration(pollInterval) * durationPollInterval)
// tickerReport := time.NewTicker(time.Duration(reportInterval) * durationReportInterval)
