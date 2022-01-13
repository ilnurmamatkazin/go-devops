package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) parseGaugeMetric(w http.ResponseWriter, r *http.Request) {
	valueMetric := chi.URLParam(r, "valueMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	f, err := strconv.ParseFloat(valueMetric, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric := models.MetricGauge{Name: nameMetric, Value: f}
	_ = h.repository.SetGauge(metric)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
