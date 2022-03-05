package handlers

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
// 	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
// 	mocks "github.com/ilnurmamatkazin/go-devops/cmd/server/service/mock_service"
// )

// func TestHandler_getMetric(t *testing.T) {
// 	type mockBehavior func(r *mocks.MockMetric, metric *models.Metric)

// 	// type fields struct {
// 	// 	service *service.Service
// 	// }
// 	// type args struct {
// 	// 	w http.ResponseWriter
// 	// 	r *http.Request
// 	// }
// 	type want struct {
// 		code        int
// 		response    string
// 		contentType string
// 	}
// 	// tests := []struct {
// 	// 	name   string
// 	// 	want   want
// 	// 	fields fields
// 	// 	args   args
// 	// }{
// 	tests := []struct {
// 		name                 string
// 		inputBody            string
// 		id                   string
// 		metricType           string
// 		value                float64
// 		mockBehavior         mockBehavior
// 		want                 want
// 		expectedResponseBody string
// 	}{

// 		{
// 			name:       "Ok",
// 			inputBody:  `{"id": "Alloc", "metricType": "gauge", "value": 12345}`,
// 			id:         "Alloc",
// 			metricType: "gauge",
// 			value:      12345.6,
// 			mockBehavior: func(r *mocks.MockMetric, metric *models.Metric) {
// 				fmt.Println("#####", metric)
// 				r.EXPECT().GetMetric(metric).Return(nil)
// 			},
// 			want: want{
// 				code:        200,
// 				response:    `{"id": "Alloc", "metricType": "gauge", "value": 12345}`,
// 				contentType: "application/json",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Init Dependencies
// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			metric := &models.Metric{ID: tt.id, MetricType: tt.metricType, Value: &tt.value}
// 			repo := mocks.NewMockMetric(c)
// 			tt.mockBehavior(repo, metric)
// 			fmt.Println("@@@@@@@@")

// 			services := &service.Service{Metric: repo}
// 			handler := Handler{services}

// 			fmt.Println("@@@@11111@@@@")

// 			request := httptest.NewRequest(http.MethodGet, "/value/", bytes.NewBufferString(tt.inputBody))

// 			// создаём новый Recorder
// 			w := httptest.NewRecorder()
// 			// определяем хендлер
// 			h := http.HandlerFunc(handler.getMetric)
// 			// запускаем сервер
// 			h.ServeHTTP(w, request)
// 			res := w.Result()

// 			// проверяем код ответа
// 			if res.StatusCode != tt.want.code {
// 				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
// 			}

// 			// получаем и проверяем тело запроса
// 			defer res.Body.Close()
// 			resBody, err := io.ReadAll(res.Body)
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			if string(resBody) != tt.want.response {
// 				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
// 			}

// 			// заголовок ответа
// 			if res.Header.Get("Content-Type") != tt.want.contentType {
// 				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
// 			}

// 		})
// 	}
// }
