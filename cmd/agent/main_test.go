package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
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

	cfg := models.Config{
		Address:        Address,
		ReportInterval: ReportInterval,
		PollInterval:   PollInterval,
	}

	metricSender := MetricSender{
		cfg:    parseConfig(),
		client: createClient(),
		ctx:    context.Background(),
	}

	if err := env.Parse(&cfg); err != nil {
		log.Println("Ошибка чтения конфигурации")
		os.Exit(2)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := metricSender.sendMetric(tt.args.typeMetric, tt.args.nameMetric, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("sendMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
