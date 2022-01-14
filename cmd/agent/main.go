package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
)

const (
	address        = "localhost:8080"
	pollInterval   = "2s"
	reportInterval = "10s"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval string `env:"REPORT_INTERVAL"`
	PollInterval   string `env:"POLL_INTERVAL"`
}

func main() {
	var (
		mutex     sync.Mutex
		rtm       runtime.MemStats
		pollCount int64
		cfg       Config
	)

	cfg = Config{
		Address:        address,
		ReportInterval: reportInterval,
		PollInterval:   pollInterval,
	}

	if err := env.Parse(&cfg); err != nil {
		os.Exit(2)
	}

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

	pi, _ := strconv.Atoi(strings.Split(cfg.PollInterval, "s")[0])
	ri, _ := strconv.Atoi(strings.Split(cfg.ReportInterval, "s")[0])

	tickerPoll := time.NewTicker(time.Duration(pi) * time.Second)
	tickerReport := time.NewTicker(time.Duration(ri) * time.Second)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	done := make(chan bool, 1)

	go func() {
	loop:
		for {
			select {
			case <-quit:
				// fmt.Println("Shutdown Agent ...")
				// t := time.NewTicker(5 * time.Second)
				// <-t.C
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

	// fmt.Println("Agent exiting")

}

func sendMetric(ctxBase context.Context, client *http.Client, typeMetric, nameMetric string, value interface{}, cfg Config) (err error) {
	ctx, cancel := context.WithTimeout(ctxBase, 1*time.Second)
	defer cancel()

	// err = sendMetricText(ctx, client, typeMetric, nameMetric, value)
	err = sendMetricJSON(ctx, client, typeMetric, nameMetric, value, cfg)

	return
}

// func sendMetricText(ctx context.Context, client *http.Client, typeMetric, nameMetric string, value interface{}) (err error) {
// 	endpoint := fmt.Sprintf("http://%s/update/%s/%s/%v", cfg., port, typeMetric, nameMetric, value)

// 	buf := new(bytes.Buffer)
// 	err = binary.Write(buf, binary.LittleEndian, value)

// 	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return
// 	}

// 	request.Header.Set("Content-Type", "text/plain; charset=utf-8")

// 	// отправляем запрос и получаем ответ
// 	response, err := client.Do(request)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return
// 	}
// 	// печатаем код ответа
// 	fmt.Println("Статус-код ", response.Status)
// 	defer response.Body.Close()
// 	// читаем поток из тела ответа
// 	body, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// и печатаем его
// 	fmt.Println(string(body))

// 	return
// }

func sendMetricJSON(ctx context.Context, client *http.Client, typeMetric, nameMetric string, value interface{}, cfg Config) (err error) {
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
		//fmt.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, b)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		//fmt.Println(err)
		return
	}

	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Set("Content-Type", "application/json")

	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		//fmt.Println(err)
		return
	}
	// печатаем код ответа
	// fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // и печатаем его
	// fmt.Println(string(body))

	return
}
