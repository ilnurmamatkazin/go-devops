//go:build ignore
// +build ignore

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
	"github.com/stretchr/testify/mock"
	"golang.org/x/sync/errgroup"
)

type FakeRequestSend struct {
	mock.Mock
}

func (mock *FakeRequestSend) Send(ctx context.Context, data interface{}, layout string) error {
	// args := mock.Called(ctx, data, layout)
	return nil //args.Error(0)
}

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

			cfg := parseConfig()

			sr := &RequestSend{
				cfg:    cfg,
				client: createClient(),
			}

			ctx := context.Background()

			var metric models.Metric
			_ = json.Unmarshal([]byte(tt.args.data), &metric)

			err := sr.Send(ctx, metric, "http://%s"+tt.args.layout)
			assert.NoError(t, err)
		})
	}
}

func BenchmarkSendMetrics(b *testing.B) {
	cfg := parseConfig()

	ms := MetricSend{
		cfg: cfg,
		sender: &RequestSend{
			cfg:    cfg,
			client: createClient(),
		},
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

func TestMetricSender_sendMetrics(t *testing.T) {
	type args struct {
		chMetrics   models.Metric
		chMetricsGU models.Metric
	}

	tests := []struct {
		name   string
		layout string
		args   args
		err    error
	}{
		{
			name:   "Positive",
			layout: "http://%s/update",
			args: args{
				chMetrics:   models.Metric{ID: "Alloc", MetricType: "gauge", Value: func(val float64) *float64 { return &val }(123.4)},
				chMetricsGU: models.Metric{ID: "Alloc", MetricType: "gauge", Value: func(val float64) *float64 { return &val }(123.4)},
			},
			err: nil,
		},
	}
	var (
		g        *errgroup.Group
		ctxGroup context.Context
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			chMetrics := make(chan []models.Metric)
			chMetricsGopsutil := make(chan []models.Metric)

			ctx, done := context.WithTimeout(context.Background(), time.Duration(4)*time.Second)
			g, ctxGroup = errgroup.WithContext(ctx)

			sender := &FakeRequestSend{}
			// sender.On("Send", ctxGroup, tt.args.chMetrics, tt.layout).Return(tt.err)
			sender.On("Send", ctxGroup, mock.Anything, mock.Anything).Return(tt.err)

			ms := &MetricSend{
				cfg:    models.Config{},
				sender: sender,
			}

			tickerReport := time.NewTicker(time.Duration(2) * time.Second)

			g.Go(func() error {
				metrics := []models.Metric{
					{ID: "Alloc", MetricType: "gauge", Value: func(val float64) *float64 { return &val }(123.4)},
				}

				chMetrics <- metrics
				chMetricsGopsutil <- metrics

				<-ctxGroup.Done()
				return ctxGroup.Err()

			})

			g.Go(func() error {
				err := ms.sendMetrics(ctxGroup, tickerReport, chMetrics, chMetricsGopsutil)
				done()
				return err
			})

			err := g.Wait()
			if err != nil {
				t.Logf("received error: %v", err)

				if err.Error() == "context deadline exceeded" {
					err = nil
				}
			}

			assert.NoError(t, err)
		})
	}
}
