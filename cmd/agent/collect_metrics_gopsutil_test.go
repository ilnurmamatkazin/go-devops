package main

import (
	"context"
	"net/http"
	"strings"
	"testing"

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
			args:    args{poll: "2s", chMetrics: make(chan []models.Metric)},
			metrics: "TotalMemory, FreeMemory, CPUutilization1",
			wantErr: false,
		},
		{
			name:    "Negative test",
			args:    args{poll: "2s", chMetrics: make(chan []models.Metric)},
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
			var done context.CancelFunc

			ctx, done := context.WithCancel(context.Background())
			g, ms.ctx = errgroup.WithContext(ctx)

			g.Go(func() error {
				return ms.collectMetricsGopsutil(tt.args.poll, tt.args.chMetrics)
			})

			metrics := <-tt.args.chMetrics

			for _, item := range metrics {
				assert.Equal(t, strings.Contains(tt.metrics, item.ID), !tt.wantErr)
			}

			done()
			err := g.Wait()
			assert.Equal(t, err.Error(), "context canceled")

		})
	}
}
