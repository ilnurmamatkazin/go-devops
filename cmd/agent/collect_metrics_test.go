package main

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

// TestMetricSender_collectMetrics функция, тестирующая отправку runtime метрик.
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

	ms := &MetricSend{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				g        *errgroup.Group
				ctxGroup context.Context
			)

			tickerPoll, _ := getTicker(tt.args.poll)

			ctx, done := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
			// ctx, done := context.WithCancel(context.Background())
			g, ctxGroup = errgroup.WithContext(ctx)

			g.Go(func() error {
				err := ms.collectMetrics(ctxGroup, tickerPoll, tt.args.chMetrics)

				return err
			})

			select {
			case <-ctxGroup.Done():
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

func BenchmarkCollectMetrics(b *testing.B) {
	ms := &MetricSend{}

	tickerPoll := time.NewTicker(time.Duration(2) * time.Second)
	chMetrics := make(chan []models.Metric)

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

		select {
		case <-ctxGroup.Done():
			tickerPoll.Stop()

		case <-chMetrics:
		}

		done()

		err := g.Wait()
		if err != nil {
			log.Println(err.Error())
		}
	}
}
