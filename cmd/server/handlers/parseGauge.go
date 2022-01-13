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

	var metric models.Metric

	metric.ID = nameMetric
	metric.MType = "gauge"
	metric.Value = &f

	if err = h.service.SetMetric(metric); err != nil {
		re, ok := err.(*models.RequestError)
		if ok {
			http.Error(w, re.Err.Error(), re.StatusCode)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
