package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) parseCounterMetric(w http.ResponseWriter, r *http.Request) {
	valueMetric := chi.URLParam(r, "valueMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	i, err := strconv.Atoi(valueMetric)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s value is not integer", nameMetric), http.StatusBadRequest)
		return
	}

	metric := models.MetricCounter{Name: nameMetric, Value: i}
	_ = h.repository.CreateCounter(metric)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}
