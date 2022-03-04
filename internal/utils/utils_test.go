package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetDataForTicker(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		wantInterval int
		wantDuration time.Duration
		wantErr      bool
	}{
		{name: "PositiveSecond", value: "2s", wantInterval: 2, wantDuration: time.Second, wantErr: false},
		{name: "PositiveMinute", value: "2m", wantInterval: 2, wantDuration: time.Minute, wantErr: false},
		{name: "PositiveHour", value: "2h", wantInterval: 2, wantDuration: time.Hour, wantErr: false},
		{name: "Negative", value: "2o", wantInterval: 2, wantDuration: time.Second, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInterval, gotDuration, err := GetDataForTicker(tt.value)

			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, gotInterval == tt.wantInterval, !tt.wantErr)
				assert.Equal(t, gotDuration == tt.wantDuration, !tt.wantErr)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
