package storage

import (
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/stretchr/testify/assert"
)

func TestStorageMetrick_ReadMetric(t *testing.T) {
	type fields struct {
		metrics map[string]models.Metric
		*sync.RWMutex
	}

	type args struct {
		metric        string
		metricStorage string
	}

	var metric models.Metric

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Posiitve get counter",
			fields: fields{
				metrics: make(map[string]models.Metric, 1),
			},
			args: args{
				metric:        `{"id": "PollCount", "type": "counter", "delta": 12345}`,
				metricStorage: `{"id": "PollCount", "type": "counter", "delta": 12345}`,
			},
			wantErr: false,
		},
		{
			name: "Negative get counter",
			fields: fields{
				metrics: make(map[string]models.Metric, 1),
			},
			args: args{
				metric:        `{"id": "PollCount111111", "type": "counter", "delta": 12345}`,
				metricStorage: `{"id": "PollCount", "type": "counter", "delta": 12345}`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StorageMetrick{
				metrics: tt.fields.metrics,
			}

			_ = json.Unmarshal([]byte(tt.args.metricStorage), &metric)

			s.Lock()
			s.metrics[metric.ID] = metric
			s.Unlock()

			_ = json.Unmarshal([]byte(tt.args.metric), &metric)

			err := s.ReadMetric(&metric)

			assert.Equal(t, (err != nil), tt.wantErr)

		})
	}
}

func TestStorageMetrick_Info(t *testing.T) {
	type fields struct {
		metrics map[string]models.Metric
		*sync.RWMutex
	}

	var metric models.Metric

	tests := []struct {
		name       string
		metricName string
		fields     fields
		metrics    []string
		wantErr    bool
	}{
		{
			name:       "Positive",
			metricName: "PollCount",
			fields: fields{
				metrics: make(map[string]models.Metric, 2),
			},
			metrics: []string{
				`{"id": "PollCount", "type": "counter", "delta": 12345}`,
				`{"id": "Alloc", "type": "gauge", "value": 1234.5}`,
			},
			wantErr: false,
		},
		{
			name:       "Negative",
			metricName: "PollCount",
			fields: fields{
				metrics: make(map[string]models.Metric),
			},
			metrics: make([]string, 0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StorageMetrick{
				metrics: tt.fields.metrics,
			}

			s.Lock()
			for _, item := range tt.metrics {
				_ = json.Unmarshal([]byte(item), &metric)

				s.metrics[metric.ID] = metric
			}
			s.Unlock()

			html := s.Info()

			assert.Equal(t, strings.Contains(html, tt.metricName), !tt.wantErr)

		})
	}
}
