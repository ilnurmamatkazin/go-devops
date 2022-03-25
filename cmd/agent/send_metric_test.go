//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestMetricSender_sendRequest(t *testing.T) {
	type args struct {
		data   string
		layout string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				data:   `{"id": "Alloc", "type": "gauge", "value": 123.5}`,
				layout: "/update",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				var metric models.Metric

				assert.Equal(t, req.URL.String(), tt.args.layout)

				err := json.NewDecoder(req.Body).Decode(&metric)
				assert.NoError(t, err)

				strMetric, _ := json.Marshal(metric)
				assert.JSONEq(t, string(strMetric), tt.args.data)
			}))
			defer server.Close()

			ms := MetricSender{
				cfg:    models.Config{Address: server.URL},
				client: server.Client(),
				// ctx:    context.Background(),
			}

			ctx := context.Background()

			var metric models.Metric
			_ = json.Unmarshal([]byte(tt.args.data), &metric)

			err := ms.sendRequest(ctx, metric, "%s"+tt.args.layout)
			assert.NoError(t, err)
		})
	}
}

func BenchmarkSendMetrics(b *testing.B) {
	ms := &MetricSender{
		cfg:    models.Config{Address: Address},
		client: &http.Client{},
		// ctx:    context.Background(),
	}

	tickerPoll := time.NewTicker(time.Duration(2) * time.Second)
	tickerReport := time.NewTicker(time.Duration(10) * time.Second)
	chMetrics := make(chan []models.Metric)
	chMetricsGU := make(chan []models.Metric)

	for i := 0; i < b.N; i++ {
		var (
			g        *errgroup.Group
			ctxGroup context.Context
		)

		ctx, done := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
		g, ctxGroup = errgroup.WithContext(ctx)

		g.Go(func() error {
			err := ms.collectMetrics(ctxGroup, tickerPoll, chMetrics)
			log.Println(err.Error())
			return err
		})

		g.Go(func() error {
			err := ms.collectMetricsGopsutil(ctxGroup, tickerPoll, chMetricsGU)
			log.Println(err.Error())
			return err
		})

		g.Go(func() error {
			err := ms.sendMetrics(ctxGroup, tickerReport, chMetrics, chMetricsGU)
			log.Println(err.Error())
			return err
		})

		<-ctxGroup.Done()

		tickerPoll.Stop()
		tickerReport.Stop()

		for i := range chMetrics {
			log.Println(i)
			close(chMetrics)
		}

		for i := range chMetricsGU {
			log.Println(i)
			close(chMetricsGU)

		}

		done()

		err := g.Wait()
		if err != nil {
			log.Println(err.Error())
		}
	}
}
