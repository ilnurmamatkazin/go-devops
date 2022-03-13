package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

var (
	w http.ResponseWriter
	r *http.Request
	h *Handler
)

func ExampleHandler_GetMetric() {
	var metric models.Metric

	// Устанавливаем заголовок
	w.Header().Set("Content-Type", "application/json")

	// Передаем в теле запроса структуру вида:
	// {"id": "Alloc", "type": "gauge"}
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		// Обрабатываем ошибку
	}

	// Получаем значение метрики
	err = h.Service.GetMetric(&metric)
	if err != nil {
		// Обрабатываем ошибку
	}

	// Отправляем ответ
	sendOkJSONData(w, metric)

}

func ExampleHandler_ParseMetric() {
	var metric models.Metric

	// Устанавливаем заголовок
	w.Header().Set("Content-Type", "application/json")

	// Передаем в теле запроса структуру вида:
	// {"id": "Alloc", "type": "gauge", "value": 123.4}
	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		// Обрабатываем ошибку
	}

	//  Сохраняем значение метрики
	err = h.Service.SetMetric(metric)
	if err != nil {
		// Обрабатываем ошибку
	}

	// Отправляем ответ
	sendOkJSONData(w, metric)
}

func ExampleHandler_ParseMetrics() {
	var (
		metrics []models.Metric
		status  models.Status
	)

	// Устанавливаем заголовок
	w.Header().Set("Content-Type", "application/json")

	// Передаем в теле запроса массив структур вида:
	// [
	//		{"id": "Alloc", "type": "gauge", "value": 123.4},
	//		{"id": "PollCount", "type": "counter", "value": 1234}
	//	]
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		// Обрабатываем ошибку
	}

	//  Сохраняем значение списка метрик
	err = h.Service.SetArrayMetrics(metrics)
	if err != nil {
		// Обрабатываем ошибку
	}

	// Устанавливаем ответ
	status.Status = http.StatusText(http.StatusOK)

	// Отправляем ответ
	sendOkJSONData(w, status)
}

func ExampleHandler_GetInfo() {
	// Получаем данные
	html := h.Service.GetInfo()

	// Устанавливаем заголовки
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Отправляем ответ
	w.Write([]byte(html))
}

func ExampleHandler_GetOldMetric() {
	// Получаем данные: имя метрики и ее тип
	metric := getMetricFromRequest(r)

	// Проверяем тип данных
	check := checkMetricType(metric.MetricType)
	if !check {
		// Формируем ошибку и выходим из функции
		return
	}

	// Получаем данные из системы
	err := h.Service.GetOldMetric(&metric)
	if err != nil {
		// Обрабатываем ошибку
	}

	var (
		httpStatus int
		strValue   string
	)

	// Устанавливаем код ответа и значение ответа
	switch metric.MetricType {
	case "counter":
		httpStatus = http.StatusOK
		strValue = strconv.Itoa(int(*metric.Delta))
	case "gauge":
		httpStatus = http.StatusOK
		strValue = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	default:
		httpStatus = http.StatusNotImplemented

	}

	// Устанавливаем заголовки
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(httpStatus)

	// Отправляем ответ
	if httpStatus == http.StatusOK {
		w.Write([]byte(strValue))
	}

}

func ExampleHandler_ParseOldMetric() {
	// Получаем данные: имя метрики и ее тип
	metric := getMetricFromRequest(r)

	// Проверяем тип данных
	check := checkMetricType(metric.MetricType)
	if !check {
		// Формируем ошибку и выходим из функции
		return
	}

	// Получаем значение метрики
	err := setMetricValue(&metric, chi.URLParam(r, "valueMetric"))
	if err != nil {
		// Обрабатываем ошибку
	}

	// Сохраняем данные
	h.Service.SetOldMetric(metric)

	// // Устанавливаем заголовки и отаправляем ответ
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}

func ExampleHandler_Ping() {
	// Устанавливаем заголовки
	w.Header().Set("Content-Type", "application/json")

	// Проверяем соединение
	err := h.Service.Ping()
	if err != nil {
		// Обрабатываем ошибку
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
}
