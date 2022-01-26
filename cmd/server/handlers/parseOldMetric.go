package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) parseOldMetric(w http.ResponseWriter, r *http.Request) {
	metric := getMetricFromRequest(r)

	fmt.Println(metric)
	fmt.Println((metric.MetricType != "counter") && (metric.MetricType != "gauge"))

	if checkMetricType(metric.MetricType) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}

	if err := setMetricValue(&metric, chi.URLParam(r, "valueMetric")); err != nil {
		http.Error(w, err.(*models.RequestError).Err.Error(), err.(*models.RequestError).StatusCode)
		return
	}

	h.service.SetOldMetric(metric)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}

func setMetricValue(metric *models.Metric, value string) (err error) {
	switch metric.MetricType {
	case "counter":
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return &models.RequestError{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New(http.StatusText(http.StatusBadRequest)),
			}
		}

		metric.Delta = &i
	case "gauge":
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return &models.RequestError{
				StatusCode: http.StatusBadRequest,
				Err:        errors.New(http.StatusText(http.StatusBadRequest)),
			}
		}

		metric.Value = &f
	default:
		err = &models.RequestError{
			StatusCode: http.StatusNotImplemented,
			Err:        errors.New(http.StatusText(http.StatusNotImplemented)),
		}

	}

	return
}
