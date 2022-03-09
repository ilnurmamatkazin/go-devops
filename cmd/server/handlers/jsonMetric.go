package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

// GetMetric функция получения значения метрики по имени и типу метрики.
// Имя и тип метрики получают из теле http запроса.
func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	var (
		metric models.Metric
		err    error
	)
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.Service.GetMetric(&metric); err != nil {
		re, ok := err.(*models.RequestError)
		if ok {
			http.Error(w, re.Err.Error(), re.StatusCode)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	sendOkJSONData(w, metric)
}

// parseMetric функция сохранения метрики в системе.
// Имя, тип и значение метрики получают из теле http запроса.
func (h *Handler) parseMetric(w http.ResponseWriter, r *http.Request) {
	var (
		metric models.Metric
		err    error
	)
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.Service.SetMetric(metric); err != nil {
		sendError(w, err)
		return
	}

	sendOkJSONData(w, metric)
}

// parseMetrics функция группового сохранения метрик в системе.
// Массив метрик получают из теле http запроса.
func (h *Handler) parseMetrics(w http.ResponseWriter, r *http.Request) {
	var (
		metrics []models.Metric
		err     error
		status  models.Status
	)
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.Service.SetArrayMetrics(metrics); err != nil {
		sendError(w, err)
		return
	}

	status.Status = http.StatusText(http.StatusOK)

	sendOkJSONData(w, status)
}

// sendOkJSONData вспомогательная функция, формирующая ответ 200.
func sendOkJSONData(w http.ResponseWriter, object interface{}) {
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(object); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// sendError вспомогательная функция, формирующая ответы с ошибками.
func sendError(w http.ResponseWriter, err error) {
	re, ok := err.(*models.RequestError)
	if ok {
		http.Error(w, re.Err.Error(), re.StatusCode)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
