package main

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
)

func TestMetricSender_collectMetricsGopsutil(t *testing.T) {
	type args struct {
		poll      string
		chMetrics chan []models.Metric
	}
	tests := []struct {
		name    string
		args    args
		metrics string
		wantErr bool
	}{
		{
			name:    "Positive test",
			args:    args{poll: "5s", chMetrics: make(chan []models.Metric)},
			metrics: "TotalMemory, FreeMemory, CPUutilization1",
			wantErr: false,
		},
		{
			name:    "Negative test",
			args:    args{poll: "5s", chMetrics: make(chan []models.Metric)},
			metrics: "Total1Memory",
			wantErr: true,
		},
	}
	ms := &MetricSender{
		cfg:    models.Config{},
		client: &http.Client{},
		ctx:    context.Background(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var g *errgroup.Group

			tickerPoll, _ := getTicker(tt.args.poll)

			ctx, done := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
			// ctx, done := context.WithCancel(context.Background())
			g, ms.ctx = errgroup.WithContext(ctx)

			g.Go(func() error {
				err := ms.collectMetricsGopsutil(tickerPoll, tt.args.chMetrics)

				return err
			})

			select {
			case <-ms.ctx.Done():
				tickerPoll.Stop()

			case metrics := <-tt.args.chMetrics:
				for _, item := range metrics {
					assert.Equal(t, strings.Contains(tt.metrics, item.ID), !tt.wantErr)
				}
			}

			done()

			err := g.Wait()
			if err != nil {
				assert.Equal(t, err.Error(), "context canceled")
			}
		})
	}
}
