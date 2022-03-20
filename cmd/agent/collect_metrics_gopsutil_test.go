package main

// import (
// 	"context"
// 	"log"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"golang.org/x/sync/errgroup"

// 	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
// )

// // TestMetricSender_collectMetricsGopsutil функция, тестирующая отправку gopsutil метрик.
// func TestMetricSender_collectMetricsGopsutil(t *testing.T) {
// 	type args struct {
// 		poll      string
// 		chMetrics chan []models.Metric
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		metrics string
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Positive test",
// 			args:    args{poll: "5s", chMetrics: make(chan []models.Metric)},
// 			metrics: "TotalMemory, FreeMemory, CPUutilization1",
// 			wantErr: false,
// 		},
// 		{
// 			name:    "Negative test",
// 			args:    args{poll: "5s", chMetrics: make(chan []models.Metric)},
// 			metrics: "Total1Memory",
// 			wantErr: true,
// 		},
// 	}
// 	ms := &MetricSend{}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var (
// 				g        *errgroup.Group
// 				ctxGroup context.Context
// 			)

// 			tickerPoll, _ := getTicker(tt.args.poll)

// 			ctx, done := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
// 			// ctx, done := context.WithCancel(context.Background())
// 			g, ctxGroup = errgroup.WithContext(ctx)

// 			g.Go(func() error {
// 				err := ms.collectMetricsGopsutil(ctxGroup, tickerPoll, tt.args.chMetrics)

// 				return err
// 			})

// 			select {
// 			case <-ctxGroup.Done():
// 				tickerPoll.Stop()

// 			case metrics := <-tt.args.chMetrics:
// 				for _, item := range metrics {
// 					assert.Equal(t, strings.Contains(tt.metrics, item.ID), !tt.wantErr)
// 				}
// 			}

// 			done()

// 			err := g.Wait()
// 			if err != nil {
// 				assert.Equal(t, err.Error(), "context canceled")
// 			}
// 		})
// 	}
// }

// func BenchmarkCollectMetricsGopsutil(b *testing.B) {
// 	ms := &MetricSend{}

// 	tickerPoll := time.NewTicker(time.Duration(2) * time.Second)
// 	chMetrics := make(chan []models.Metric)

// 	for i := 0; i < b.N; i++ {
// 		var (
// 			g        *errgroup.Group
// 			ctxGroup context.Context
// 		)

// 		ctx, done := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
// 		g, ctxGroup = errgroup.WithContext(ctx)

// 		g.Go(func() error {
// 			err := ms.collectMetricsGopsutil(ctxGroup, tickerPoll, chMetrics)
// 			log.Println(err.Error())
// 			return err
// 		})

// 		select {
// 		case <-ctxGroup.Done():
// 			tickerPoll.Stop()

// 		case <-chMetrics:
// 		}

// 		done()

// 		err := g.Wait()
// 		if err != nil {
// 			log.Println(err.Error())
// 		}
// 	}
// }
