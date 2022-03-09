package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
	"github.com/stretchr/testify/assert"
)

func TestMetricSender_sendRequest(t *testing.T) {
	type args struct {
		data   string
		layout string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				data:   `{"id": "Alloc", "type": "gauge", "value": 123.5}`,
				layout: "/update",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				var metric models.Metric

				assert.Equal(t, req.URL.String(), tt.args.layout)

				err := json.NewDecoder(req.Body).Decode(&metric)
				assert.NoError(t, err)

				strMetric, _ := json.Marshal(metric)
				assert.JSONEq(t, string(strMetric), tt.args.data)
			}))
			defer server.Close()

			ms := MetricSender{
				cfg:    models.Config{Address: server.URL},
				client: server.Client(),
				ctx:    context.Background(),
			}

			var metric models.Metric
			_ = json.Unmarshal([]byte(tt.args.data), &metric)

			err := ms.sendRequest(metric, "%s"+tt.args.layout)
			assert.NoError(t, err)
		})
	}
}
