package utils

import (
	"testing"
	"time"
)

func TestGetDataForTicker(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name         string
		args         args
		wantInterval int
		wantDuration time.Duration
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInterval, gotDuration, err := GetDataForTicker(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataForTicker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInterval != tt.wantInterval {
				t.Errorf("GetDataForTicker() gotInterval = %v, want %v", gotInterval, tt.wantInterval)
			}
			if gotDuration != tt.wantDuration {
				t.Errorf("GetDataForTicker() gotDuration = %v, want %v", gotDuration, tt.wantDuration)
			}
		})
	}
}
