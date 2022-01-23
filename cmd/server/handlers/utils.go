package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func getMetricFromRequest(r *http.Request) (metric models.Metric) {
	metric.ID = chi.URLParam(r, "nameMetric")
	metric.MetricType = chi.URLParam(r, "typeMetric")

	return
}

func checkMetricType(metricType string) (err error) {
	if (metricType != "counter") || (metricType != "gauge") {
		err = &models.RequestError{
			StatusCode: http.StatusNotImplemented,
			Err:        errors.New(http.StatusText(http.StatusNotImplemented)),
		}
	}

	return
}
