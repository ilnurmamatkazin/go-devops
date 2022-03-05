package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
	service_mocks "github.com/ilnurmamatkazin/go-devops/cmd/server/service/mock_service"

	// "github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(r *service_mocks.MockMetric, metric *models.Metric)

	type want struct {
		code        int
		response    string
		contentType string
	}

	var metric models.Metric

	tests := []struct {
		name      string
		inputBody string
		// inputMetric          models.Metric
		mockBehavior mockBehavior
		want         want
	}{
		{
			name:      "Ok",
			inputBody: `{"id": "Alloc", "type": "gauge", "value": 123.5}`,
			// inputMetric: models.Metric{
			// 	ID:         "Alloc",
			// 	MetricType: "gauge",
			// 	Value:      &val,
			// },
			mockBehavior: func(r *service_mocks.MockMetric, metric *models.Metric) {
				r.EXPECT().GetMetric(metric).Return(nil)
			},
			want: want{
				code:        200,
				response:    `{"id": "Alloc", "type": "gauge", "value": 123.5}`,
				contentType: "application/json",
			},
		},
		// {
		// 	name:                 "Wrong Input",
		// 	inputBody:            `{"username": "username"}`,
		// 	inputMetric:          models.Metric{},
		// 	mockBehavior:         func(r *service_mocks.MockMetric, metric *models.Metric) {},
		// 	expectedStatusCode:   400,
		// 	expectedResponseBody: `{"message":"invalid input body"}`,
		// },
		// {
		// 	name:      "Service Error",
		// 	inputBody: `{"id": "Alloc", "type": "gauge", "value": 123.5}`,
		// 	inputMetric: models.Metric{
		// 		ID:         "Alloc",
		// 		MetricType: "gauge",
		// 		Value:      &val,
		// 	},
		// 	mockBehavior: func(r *service_mocks.MockMetric, metric *models.Metric) {
		// 		r.EXPECT().GetMetric(metric).Return(errors.New("something went wrong"))
		// 	},
		// 	expectedStatusCode:   500,
		// 	expectedResponseBody: "something went wrong\n",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = json.Unmarshal([]byte(tt.inputBody), &metric)

			c := gomock.NewController(t)
			defer c.Finish()

			repo := service_mocks.NewMockMetric(c)
			tt.mockBehavior(repo, &metric)

			services := &service.Service{Metric: repo}
			handler := Handler{services}

			// Init Endpoint
			r := chi.NewRouter()
			r.Route("/", func(r chi.Router) {
				r.Post("/value/", handler.getMetric)
			})

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/value/", bytes.NewBufferString(tt.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, res.StatusCode, tt.want.code)

			// получаем и проверяем тело запроса
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			// тело ответа
			assert.JSONEq(t, string(resBody), tt.want.response)

			// заголовок ответа
			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)
		})
	}
}
