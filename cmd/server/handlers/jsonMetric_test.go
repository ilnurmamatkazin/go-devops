package handlers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/golang/mock/gomock"
// 	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
// 	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
// 	service_mocks "github.com/ilnurmamatkazin/go-devops/cmd/server/service/mock_service"
// 	"github.com/stretchr/testify/assert"
// )

// func TestHandler_getMetric(t *testing.T) {
// 	type mockBehavior func(r *service_mocks.MockMetric, metric *models.Metric)

// 	type want struct {
// 		code        int
// 		response    string
// 		contentType string
// 	}

// 	var metric models.Metric

// 	tests := []struct {
// 		name         string
// 		inputBody    string
// 		mockBehavior mockBehavior
// 		want         want
// 	}{
// 		{
// 			name:      "Ok",
// 			inputBody: `{"id": "Alloc", "type": "gauge", "value": 123.5}`,
// 			mockBehavior: func(r *service_mocks.MockMetric, metric *models.Metric) {
// 				r.EXPECT().GetMetric(metric).Return(nil)
// 			},
// 			want: want{
// 				code:        200,
// 				response:    `{"id": "Alloc", "type": "gauge", "value": 123.5}`,
// 				contentType: "application/json",
// 			},
// 		},
// 		{
// 			name:      "Service Error",
// 			inputBody: `{"id": "Alloc", "type": "gauge", "value": 123.5}`,
// 			mockBehavior: func(r *service_mocks.MockMetric, metric *models.Metric) {
// 				r.EXPECT().GetMetric(metric).Return(fmt.Errorf(`{"message":"something went wrong"}`))
// 			},
// 			want: want{
// 				code:        500,
// 				response:    `{"message":"something went wrong"}`,
// 				contentType: "text/plain; charset=utf-8",
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_ = json.Unmarshal([]byte(tt.inputBody), &metric)

// 			c := gomock.NewController(t)
// 			defer c.Finish()

// 			repo := service_mocks.NewMockMetric(c)
// 			tt.mockBehavior(repo, &metric)

// 			services := &service.Service{Metric: repo}
// 			handler := Handler{services}

// 			r := chi.NewRouter()
// 			r.Route("/", func(r chi.Router) {
// 				r.Post("/value/", handler.GetMetric)
// 			})

// 			w := httptest.NewRecorder()
// 			req := httptest.NewRequest("POST", "/value/", bytes.NewBufferString(tt.inputBody))

// 			r.ServeHTTP(w, req)

// 			res := w.Result()

// 			// проверяем код ответа
// 			assert.Equal(t, res.StatusCode, tt.want.code)

// 			// получаем и проверяем тело запроса
// 			defer res.Body.Close()

// 			resBody, err := io.ReadAll(res.Body)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			// тело ответа
// 			assert.JSONEq(t, string(resBody), tt.want.response)

// 			// заголовок ответа
// 			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)
// 		})
// 	}
// }
