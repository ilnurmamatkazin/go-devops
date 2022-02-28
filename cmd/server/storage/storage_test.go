package storage

import (
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/stretchr/testify/assert"
)

func TestSetOldMetric(t *testing.T) {
	type metric struct {
		ID         string
		MetricType string
		Hash       string
		Delta      int64
		Value      float64
	}
	tests := []struct {
		name    string
		metric  metric
		wantErr bool
	}{
		{
			name:    "pisitive gauge test",
			metric:  metric{ID: "TotalAlloc", MetricType: "gauge", Value: 175368, Hash: "ecd17becc5c0cb83489bbd3387aa351ea590ce64b44ac4b06645f058f36d20c2"},
			wantErr: false,
		},
		{
			name:    "negative gauge test",
			metric:  metric{ID: "TotalAlloc1", MetricType: "gauge", Value: 175368, Hash: "ecd17becc5c0cb83489bbd3387aa351ea590ce64b44ac4b06645f058f36d20c2"},
			wantErr: true,
		},
		{
			name:    "pisitive counter test",
			metric:  metric{ID: "PollCount", MetricType: "counter", Delta: 5, Hash: "b54c435bba4ef9334d8b0ca2938c912a75660e5a152a0fb68a4177dfccdaf9e9"},
			wantErr: false,
		},
		{
			name:    "negative counter test",
			metric:  metric{ID: "PollCount1", MetricType: "counter", Delta: 5, Hash: "b54c435bba4ef9334d8b0ca2938c912a75660e5a152a0fb68a4177dfccdaf9e9"},
			wantErr: true,
		},
	}

	cfg := models.Config{}
	r := NewStorage(&cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			newMetric := models.Metric{
				ID:         tt.metric.ID,
				MetricType: tt.metric.MetricType,
			}

			if newMetric.MetricType == "counter" {
				newMetric.Delta = &tt.metric.Delta
			} else {
				newMetric.Value = &tt.metric.Value
			}

			if !tt.wantErr {
				r.SetOldMetric(newMetric)
			}

			metric := models.Metric{
				ID:         tt.metric.ID,
				MetricType: tt.metric.MetricType,
			}

			err := r.ReadMetric(&metric)

			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, metric, newMetric)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
