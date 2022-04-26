package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	cr "github.com/ilnurmamatkazin/go-devops/cmd/agent/crypto"
	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
)

// sendMetrics функция, реализующая отправку метрик на сервер.
func (ms *MetricSend) sendMetrics(ctx context.Context, tickerReport *time.Ticker, chMetrics chan []models.Metric, chMetricsGopsutil chan []models.Metric) (err error) {
	var (
		metrics []models.Metric
		// metricsGopsutil []models.Metric
	)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case metrics = <-chMetrics:
		// case metricsGopsutil = <-chMetricsGopsutil:

		case <-tickerReport.C:
			if ms.cfg.NeedGRPC {
				if err = ms.sender.GRPCSendMetric(ctx, metrics); err != nil {
					return
				}
			} else {
				for _, metric := range metrics {
					ms.sender.OldSendMetric(ctx, metric)

					if err = ms.sender.Send(ctx, metric, "http://%s/update"); err != nil {
						return
					}
				}
			}

			// if ms.cfg.NeedGRPC {
			// 	if err = ms.sender.GRPCSendMetric(ctx, metricsGopsutil); err != nil {
			// 		return
			// 	}
			// } else {
			// 	for _, metric := range metricsGopsutil {
			// 		if err = ms.sender.Send(ctx, metric, "http://%s/update"); err != nil {
			// 			return
			// 		}
			// 	}
			// }

			if len(metrics) > 0 {
				if ms.cfg.NeedGRPC {
					if err = ms.sender.GRPCSendMetrics(ctx, metrics[:9]); err != nil {
						return
					}

					if err = ms.sender.GRPCSendMetrics(ctx, metrics[10:19]); err != nil {
						return
					}

					if err = ms.sender.GRPCSendMetrics(ctx, metrics[20:29]); err != nil {
						return
					}
				} else {
					if err = ms.sender.Send(ctx, metrics[:9], "http://%s/updates/"); err != nil {
						return
					}

					if err = ms.sender.Send(ctx, metrics[10:19], "http://%s/updates/"); err != nil {
						return
					}

					if err = ms.sender.Send(ctx, metrics[20:29], "http://%s/updates/"); err != nil {
						return
					}
				}
			}

			// if len(metricsGopsutil) > 0 {
			// 	if ms.cfg.NeedGRPC {
			// 		if err = ms.sender.GRPCSendMetrics(ctx, metricsGopsutil); err != nil {
			// 			return
			// 		}
			// 	} else {
			// 		if err = ms.sender.Send(ctx, metricsGopsutil, "http://%s/updates/"); err != nil {
			// 			return
			// 		}
			// 	}
			// }
		}

	}
}

// Send функция, реализующая создание запроса для отправки метрик на сервер.
func (ms *RequestSend) Send(ctx context.Context, data interface{}, layout string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var strIP string

	endpoint := fmt.Sprintf(layout, ms.cfg.Address)

	cryptoText, err := getCryptoText(ms.cfg.PublicKey, data)
	if err != nil {
		log.Println(err)
		return
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, []byte(cryptoText))

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		log.Println(err)
		return
	}

	ipv4 := getIPAdress()
	if ipv4 == nil {
		strIP = "неудалось получить ip адрес хоста"
		log.Println(strIP)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-REAL-IP", strIP)

	// отправляем запрос и получаем ответ
	response, ok := ms.client.Do(request)
	if ok != nil {
		log.Println(ok)
	} else {
		// печатаем код ответа
		// fmt.Println("Статус-код ", response.Status)
		defer response.Body.Close()
	}

	return
}

func (ms *RequestSend) GRPCSendMetric(ctx context.Context, metrics []models.Metric) error {
	cryptoMetrics := make([]string, 0, len(metrics))

	for _, metric := range metrics {
		cryptoText, err := getCryptoText(ms.cfg.PublicKey, metric)
		if err != nil {
			return err
		}

		cryptoMetrics = append(cryptoMetrics, string(cryptoText))
	}

	return ms.grpcClient.SendMetric(ctx, cryptoMetrics)
}

func (ms *RequestSend) GRPCSendMetrics(ctx context.Context, metrics []models.Metric) error {
	cryptoText, err := getCryptoText(ms.cfg.PublicKey, metrics)
	if err != nil {
		return err
	}

	return ms.grpcClient.SendMetrics(ctx, string(cryptoText))
}

func getIPAdress() net.IP {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		return addr.To4()
	}

	return nil
}

func getCryptoText(publicKey string, data interface{}) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return b, err
	}

	return cr.Encrypt(publicKey, b)
}

func (ms *RequestSend) OldSendMetric(ctxBase context.Context, metric models.Metric) (err error) {
	ctx, cancel := context.WithTimeout(ctxBase, 1*time.Second)
	defer cancel()

	var endpoint string

	if metric.MetricType == "gauge" {
		endpoint = fmt.Sprintf("http://127.0.0.1:8080/update/%s/%s/%v", metric.MetricType, metric.ID, metric.Value)
	} else {
		endpoint = fmt.Sprintf("http://127.0.0.1:8080/update/%s/%s/%v", metric.MetricType, metric.ID, metric.Delta)

	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Set("Content-Type", "text/plain; charset=UTF-8")

	// отправляем запрос и получаем ответ
	response, ok := ms.client.Do(request)
	if ok != nil {
		log.Println(ok)
	} else {
		// печатаем код ответа
		// fmt.Println("Статус-код ", response.Status)
		defer response.Body.Close()
	}

	return
}
