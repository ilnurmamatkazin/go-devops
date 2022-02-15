package main

import (
	// "context"
	// "fmt"
	// "net/http"
	// "strings"
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	// "github.com/stretchr/testify/assert"
	// "golang.org/x/sync/errgroup"
)

func TestMetricSender_collectMetrics(t *testing.T) {
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
			name: "Positive test",
			args: args{poll: "2s", chMetrics: make(chan []models.Metric)},
			metrics: `Alloc, BuckHashSys, Frees, GCCPUFraction, GCSys, HeapAlloc, HeapIdle,
			HeapInuse, HeapObjects, HeapReleased, HeapSys, LastGC, Lookups, MCacheInuse,
			MCacheSys, MSpanInuse, MSpanSys, Mallocs, NextGC, NumForcedGC, NumGC, OtherSys,
			PauseTotalNs, TotalAlloc, StackInuse, StackSys, Sys, PollCount, RandomValue`,
			wantErr: false,
		},
		{
			name:    "Negative test",
			args:    args{poll: "2s", chMetrics: make(chan []models.Metric)},
			metrics: "Total1Memory",
			wantErr: true,
		},
	}

	// ms := &MetricSender{
	// 	cfg:    models.Config{},
	// 	client: &http.Client{},
	// 	ctx:    context.Background(),
	// }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// var g *errgroup.Group
			// var done context.CancelFunc

			// ctx, done := context.WithCancel(context.Background())
			// g, ms.ctx = errgroup.WithContext(ctx)

			// g.Go(func() error {
			// 	return ms.collectMetrics(tt.args.poll, tt.args.chMetrics)
			// })

			// metrics := <-tt.args.chMetrics

			// for _, item := range metrics {
			// 	fmt.Println(tt.metrics, item.ID, strings.Contains(tt.metrics, item.ID), !tt.wantErr)
			// 	assert.Equal(t, strings.Contains(tt.metrics, item.ID), !tt.wantErr)
			// }

			// done()
			// err := g.Wait()
			// assert.Equal(t, err.Error(), "context canceled")
		})
	}
}
