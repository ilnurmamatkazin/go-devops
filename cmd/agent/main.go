package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
)

const (
	url            = "127.0.0.1"
	port           = 8080
	pollInterval   = 2
	reportInterval = 10
)

func main() {
	var (
		mutex     sync.Mutex
		rtm       runtime.MemStats
		pollCount int64
	)

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

	tickerPoll := time.NewTicker(pollInterval * time.Second)
	tickerReport := time.NewTicker(reportInterval * time.Second)

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

				sendMetric(ctx, client, "gauge", "Alloc", rtm.Alloc)
				sendMetric(ctx, client, "gauge", "BuckHashSys", rtm.BuckHashSys)
				sendMetric(ctx, client, "gauge", "Frees", rtm.Frees)
				sendMetric(ctx, client, "gauge", "GCCPUFraction", rtm.GCCPUFraction)
				sendMetric(ctx, client, "gauge", "GCSys", rtm.GCSys)
				sendMetric(ctx, client, "gauge", "HeapAlloc", rtm.HeapAlloc)
				sendMetric(ctx, client, "gauge", "HeapIdle", rtm.HeapIdle)
				sendMetric(ctx, client, "gauge", "HeapInuse", rtm.HeapInuse)
				sendMetric(ctx, client, "gauge", "HeapObjects", rtm.HeapObjects)
				sendMetric(ctx, client, "gauge", "HeapReleased", rtm.HeapReleased)
				sendMetric(ctx, client, "gauge", "HeapSys", rtm.HeapSys)
				sendMetric(ctx, client, "gauge", "LastGC", rtm.LastGC)
				sendMetric(ctx, client, "gauge", "Lookups", rtm.Lookups)
				sendMetric(ctx, client, "gauge", "MCacheInuse", rtm.MCacheInuse)
				sendMetric(ctx, client, "gauge", "MCacheSys", rtm.MCacheSys)
				sendMetric(ctx, client, "gauge", "MSpanInuse", rtm.MSpanInuse)
				sendMetric(ctx, client, "gauge", "MSpanSys", rtm.MSpanSys)
				sendMetric(ctx, client, "gauge", "Mallocs", rtm.Mallocs)
				sendMetric(ctx, client, "gauge", "NextGC", rtm.NextGC)
				sendMetric(ctx, client, "gauge", "NumForcedGC", rtm.NumForcedGC)
				sendMetric(ctx, client, "gauge", "NumGC", rtm.NumGC)
				sendMetric(ctx, client, "gauge", "OtherSys", rtm.OtherSys)
				sendMetric(ctx, client, "gauge", "PauseTotalNs", rtm.PauseTotalNs)
				sendMetric(ctx, client, "gauge", "TotalAlloc", rtm.TotalAlloc)
				sendMetric(ctx, client, "gauge", "StackInuse", rtm.StackInuse)
				sendMetric(ctx, client, "gauge", "StackSys", rtm.StackSys)
				sendMetric(ctx, client, "gauge", "Sys", rtm.Sys)
				sendMetric(ctx, client, "counter", "PollCount", pollCount)
				sendMetric(ctx, client, "gauge", "RandomValue", rand.Float64())

				mutex.Unlock()
			}

		}
	}()

	<-done

	// fmt.Println("Agent exiting")

}

func sendMetric(ctxBase context.Context, client *http.Client, typeMetric, nameMetric string, value interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctxBase, 1*time.Second)
	defer cancel()

	// err = sendMetricText(ctx, client, typeMetric, nameMetric, value)
	err = sendMetricJSON(ctx, client, typeMetric, nameMetric, value)

	return
}

func sendMetricText(ctx context.Context, client *http.Client, typeMetric, nameMetric string, value interface{}) (err error) {
	endpoint := fmt.Sprintf("http://%s:%d/update/%s/%s/%v", url, port, typeMetric, nameMetric, value)

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, value)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		//fmt.Println(err)
		return
	}

	request.Header.Set("Content-Type", "text/plain; charset=utf-8")

	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		//fmt.Println(err)
		return
	}
	// печатаем код ответа
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	// читаем поток из тела ответа
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// и печатаем его
	fmt.Println(string(body))

	return
}

func sendMetricJSON(ctx context.Context, client *http.Client, typeMetric, nameMetric string, value interface{}) (err error) {
	var metric models.Metric

	endpoint := fmt.Sprintf("http://%s:%d/update", url, port)

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
