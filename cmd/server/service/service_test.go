package service

import (
	"crypto/hmac"
	"encoding/hex"
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage"
	"github.com/stretchr/testify/assert"
)

func TestService_SetMetric(t *testing.T) {
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
			metric:  metric{ID: "TotalAlloc", MetricType: "gauge", Value: 175368},
			wantErr: true,
		},
		{
			name:    "pisitive counter test",
			metric:  metric{ID: "PollCount", MetricType: "counter", Delta: 5, Hash: "b54c435bba4ef9334d8b0ca2938c912a75660e5a152a0fb68a4177dfccdaf9e9"},
			wantErr: false,
		},
		{
			name:    "negative counter test",
			metric:  metric{ID: "PollCount", MetricType: "counter", Delta: 5},
			wantErr: true,
		},
	}

	cfg := models.Config{Key: "qwerty"}
	r := storage.NewStorage(&cfg, nil)
	s := *NewService(&cfg, r)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			newMetric := models.Metric{
				ID:         tt.metric.ID,
				MetricType: tt.metric.MetricType,
				Hash:       &tt.metric.Hash,
			}

			if newMetric.MetricType == "counter" {
				newMetric.Delta = &tt.metric.Delta
			} else {
				newMetric.Value = &tt.metric.Value
			}

			err := s.SetMetric(newMetric)

			assert.Equal(t, (err != nil), tt.wantErr)
		})
	}
}

func TestService_GetMetric(t *testing.T) {
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
			name:    "pisitive counter test",
			metric:  metric{ID: "PollCount", MetricType: "counter", Delta: 5, Hash: "b54c435bba4ef9334d8b0ca2938c912a75660e5a152a0fb68a4177dfccdaf9e9"},
			wantErr: false,
		},
	}

	cfg := models.Config{Key: "qwerty"}
	r := storage.NewStorage(&cfg, nil)
	s := *NewService(&cfg, r)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			newMetric := models.Metric{
				ID:         tt.metric.ID,
				MetricType: tt.metric.MetricType,
				Hash:       &tt.metric.Hash,
			}

			if newMetric.MetricType == "counter" {
				newMetric.Delta = &tt.metric.Delta
			} else {
				newMetric.Value = &tt.metric.Value
			}

			err := s.SetMetric(newMetric)
			assert.NoError(t, err)
			err = s.GetMetric(&newMetric)
			assert.NoError(t, err)

			newHash, err := hex.DecodeString(*newMetric.Hash)
			assert.NoError(t, err)

			hash, err := hex.DecodeString(tt.metric.Hash)
			assert.NoError(t, err)

			assert.Equal(t, !hmac.Equal(newHash, hash), tt.wantErr)
		})
	}
}
