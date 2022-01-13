package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) parseCounterMetric(w http.ResponseWriter, r *http.Request) {
	valueMetric := chi.URLParam(r, "valueMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	i, err := strconv.ParseInt(valueMetric, 10, 64) //strconv.Atoi(valueMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metric := models.MetricCounter{Name: nameMetric, Value: i}
	_ = h.repository.SetCounter(metric)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}
