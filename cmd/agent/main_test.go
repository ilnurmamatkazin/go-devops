package main

import (
	"context"
	"net/http"
	"testing"
)

func Test_sendMetric(t *testing.T) {
	type args struct {
		ctxBase    context.Context
		client     *http.Client
		typeMetric string
		nameMetric string
		value      interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sendMetric(tt.args.ctxBase, tt.args.client, tt.args.typeMetric, tt.args.nameMetric, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("sendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
