package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) getMetric(w http.ResponseWriter, r *http.Request) {
	var (
		metric models.Metric
		err    error
	)
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.service.GetMetric(&metric); err != nil {
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

	if err = h.service.SetMetric(metric); err != nil {
		sendError(w, err)
		return
	}

	sendOkJSONData(w, metric)
}

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

	if err = h.service.SetArrayMetrics(metrics); err != nil {
		sendError(w, err)
		return
	}

	status.Status = http.StatusText(http.StatusOK)

	sendOkJSONData(w, status)
}

func sendOkJSONData(w http.ResponseWriter, object interface{}) {
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(object); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func sendError(w http.ResponseWriter, err error) {
	re, ok := err.(*models.RequestError)
	if ok {
		http.Error(w, re.Err.Error(), re.StatusCode)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
