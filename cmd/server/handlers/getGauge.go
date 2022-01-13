package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) getGauge(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	var (
		metric models.Metric
		err    error
	)

	metric.ID = nameMetric
	metric.MType = "gauge"

	if err = h.service.GetMetric(&metric); err != nil {
		re, ok := err.(*models.RequestError)
		if ok {
			http.Error(w, re.Err.Error(), re.StatusCode)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatFloat(*metric.Value, 'f', -1, 64)))

}
