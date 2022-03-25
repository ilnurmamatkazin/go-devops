//go:build ignore

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

func TestSetHash(t *testing.T) {
	type args struct {
		id         string
		metricType string
		key        string
		delta      int64
		value      float64
	}
	tests := []struct {
		name      string
		args      args
		want      []byte
		assertion assert.BoolAssertionFunc
	}{
		{
			name: "Positive gauge",
			args: args{
				id:         "Alloc",
				metricType: "gauge",
				key:        "key",
				value:      1234.5,
			},
			want:      []byte{132, 5, 131, 97, 157, 164, 4, 41, 153, 124, 96, 254, 135, 7, 244, 100, 32, 69, 174, 105, 74, 127, 27, 153, 229, 66, 164, 47, 252, 251, 233, 111},
			assertion: assert.True,
		},
		{
			name: "Positive counter",
			args: args{
				id:         "Alloc",
				metricType: "gauge",
				key:        "key",
				delta:      12345,
			},
			want:      []byte{197, 82, 183, 178, 226, 230, 247, 126, 107, 101, 95, 92, 65, 160, 77, 145, 235, 13, 77, 93, 4, 247, 247, 216, 109, 215, 41, 85, 46, 93, 54, 216},
			assertion: assert.True,
		},
		{
			name: "Negative gauge",
			args: args{
				id:         "Alloc",
				metricType: "gauge",
				key:        "key",
				value:      1234.5,
			},
			want:      []byte{197, 82, 183, 178, 226, 230, 247, 126, 107, 101, 95, 92, 65, 160, 77, 145, 235, 13, 77, 93, 4, 247, 247, 216, 109, 215, 41, 85, 46, 93, 54, 216},
			assertion: assert.False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetHash(tt.args.id, tt.args.metricType, tt.args.key, &tt.args.delta, &tt.args.value)

			tt.assertion(t, string(got) == string(tt.want))
		})
	}
}
