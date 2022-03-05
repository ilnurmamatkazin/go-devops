package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/agent/models"
)

func TestMetricSender_sendRequest(t *testing.T) {
	fmt.Println("^^^^^^^^^^^^^")
	type args struct {
		data   interface{}
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
				data:   `{"id": Alloc}`,
				layout: "http://%s/update",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// if err := ms.sendRequest(tt.args.data, tt.args.layout); (err != nil) != tt.wantErr {
			// 	t.Errorf("MetricSender.sendRequest() error = %v, wantErr %v", err, tt.wantErr)
			// }

			// expected := "dummy data"
			// Start a local HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				// equals(t, req.URL.String(), "/some/path")
				fmt.Println("@@@@", req.URL.String())
				// Send response to be tested
				rw.Write([]byte(`OK`))
			}))

			// Close the server when test finishes
			defer server.Close()

			ms := MetricSender{
				cfg:    models.Config{Address: server.URL},
				client: server.Client(),
				ctx:    context.Background(),
			}

			fmt.Println("####", ms)

			if err := ms.sendRequest(tt.args.data, tt.args.layout); (err != nil) != tt.wantErr {
				t.Errorf("MetricSender.sendRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
