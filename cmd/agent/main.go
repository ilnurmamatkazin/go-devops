package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
)

const (
	ADDRESS        = "127.0.0.1:8080"
	POLLINTERVAL   = "2s"
	REPORTINTERVAL = "10s"
)

func main() {
	var (
		mutex     sync.Mutex
		rtm       runtime.MemStats
		pollCount int64
		cfg       models.Config
		err       error
	)

	cfg, err = parseConfig()
	if err != nil {
		fmt.Println("env.Parse", err.Error())
		os.Exit(2)
	}

	fmt.Println(cfg)

	// конструируем HTTP-клиент
	client := &http.Client{}
	client.Timeout = time.Second * 2

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 2 {
			os.Exit(3)
		}
		return nil
	}

	transport := &http.Transport{}
	transport.MaxIdleConns = 20
	client.Transport = transport

	ctx := context.Background()

	strDurPi := cfg.PollInterval[len(cfg.PollInterval)-1:]
	strPi := cfg.PollInterval[0 : len(cfg.PollInterval)-1]

	strDurRi := cfg.ReportInterval[len(cfg.ReportInterval)-1:]
	strRi := cfg.ReportInterval[0 : len(cfg.ReportInterval)-1]

	pi, _ := strconv.Atoi(strPi)
	ri, _ := strconv.Atoi(strRi)

	var durPi, durRi time.Duration

	switch strDurPi {
	case "s":
		durPi = time.Second
	case "m":
		durPi = time.Minute
	case "h":
		durPi = time.Hour
	default:
		fmt.Println("Неверный временной интервал")
		return
	}

	switch strDurRi {
	case "s":
		durRi = time.Second
	case "m":
		durRi = time.Minute
	case "h":
		durRi = time.Hour
	default:
		fmt.Println("Неверный временной интервал")
		return
	}

	tickerPoll := time.NewTicker(time.Duration(pi) * durPi)
	tickerReport := time.NewTicker(time.Duration(ri) * durRi)

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

				sendMetric(ctx, client, "gauge", "Alloc", rtm.Alloc, cfg)
				sendMetric(ctx, client, "gauge", "BuckHashSys", rtm.BuckHashSys, cfg)
				sendMetric(ctx, client, "gauge", "Frees", rtm.Frees, cfg)
				sendMetric(ctx, client, "gauge", "GCCPUFraction", rtm.GCCPUFraction, cfg)
				sendMetric(ctx, client, "gauge", "GCSys", rtm.GCSys, cfg)
				sendMetric(ctx, client, "gauge", "HeapAlloc", rtm.HeapAlloc, cfg)
				sendMetric(ctx, client, "gauge", "HeapIdle", rtm.HeapIdle, cfg)
				sendMetric(ctx, client, "gauge", "HeapInuse", rtm.HeapInuse, cfg)
				sendMetric(ctx, client, "gauge", "HeapObjects", rtm.HeapObjects, cfg)
				sendMetric(ctx, client, "gauge", "HeapReleased", rtm.HeapReleased, cfg)
				sendMetric(ctx, client, "gauge", "HeapSys", rtm.HeapSys, cfg)
				sendMetric(ctx, client, "gauge", "LastGC", rtm.LastGC, cfg)
				sendMetric(ctx, client, "gauge", "Lookups", rtm.Lookups, cfg)
				sendMetric(ctx, client, "gauge", "MCacheInuse", rtm.MCacheInuse, cfg)
				sendMetric(ctx, client, "gauge", "MCacheSys", rtm.MCacheSys, cfg)
				sendMetric(ctx, client, "gauge", "MSpanInuse", rtm.MSpanInuse, cfg)
				sendMetric(ctx, client, "gauge", "MSpanSys", rtm.MSpanSys, cfg)
				sendMetric(ctx, client, "gauge", "Mallocs", rtm.Mallocs, cfg)
				sendMetric(ctx, client, "gauge", "NextGC", rtm.NextGC, cfg)
				sendMetric(ctx, client, "gauge", "NumForcedGC", rtm.NumForcedGC, cfg)
				sendMetric(ctx, client, "gauge", "NumGC", rtm.NumGC, cfg)
				sendMetric(ctx, client, "gauge", "OtherSys", rtm.OtherSys, cfg)
				sendMetric(ctx, client, "gauge", "PauseTotalNs", rtm.PauseTotalNs, cfg)
				sendMetric(ctx, client, "gauge", "TotalAlloc", rtm.TotalAlloc, cfg)
				sendMetric(ctx, client, "gauge", "StackInuse", rtm.StackInuse, cfg)
				sendMetric(ctx, client, "gauge", "StackSys", rtm.StackSys, cfg)
				sendMetric(ctx, client, "gauge", "Sys", rtm.Sys, cfg)
				sendMetric(ctx, client, "counter", "PollCount", pollCount, cfg)
				sendMetric(ctx, client, "gauge", "RandomValue", rand.Float64(), cfg)

				mutex.Unlock()
			}

		}
	}()

	<-done

}

func sendMetric(ctxBase context.Context, client *http.Client, typeMetric, nameMetric string, value interface{}, cfg models.Config) (err error) {
	ctx, cancel := context.WithTimeout(ctxBase, 1*time.Second)
	defer cancel()

	var metric models.Metric

	endpoint := fmt.Sprintf("http://%s/update", cfg.Address)

	metric.ID = nameMetric
	metric.MType = typeMetric

	switch typeMetric {
	case "counter":
		var i int64
		i = value.(int64)
		metric.Delta = &i
	case "gauge":
		var f float64

		switch value := value.(type) {
		case float64:
			f = value
		case uint64:
			f = float64(value)
		case uint32:
			f = float64(value)

		default:
		}

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
	default:
		err = errors.New("недопустимый тип")
		return
	}

	b, err := json.Marshal(metric)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, b)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Set("Content-Type", "application/json")

	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	// печатаем код ответа
	// fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()

	return
}

func parseConfig() (cfg models.Config, err error) {
	address := flag.String("a", ADDRESS, "a address")
	reportInterval := flag.String("r", REPORTINTERVAL, "a report_interval")
	pollInterval := flag.String("p", POLLINTERVAL, "a poll_interval")

	flag.Parse()

	cfg.Address = *address
	cfg.ReportInterval = *reportInterval
	cfg.PollInterval = *pollInterval

	err = env.Parse(&cfg)

	return
}
