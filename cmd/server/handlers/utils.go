package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func getMetricFromRequest(r *http.Request) (metric models.Metric) {
	metric.ID = chi.URLParam(r, "nameMetric")
	metric.MetricType = chi.URLParam(r, "typeMetric")

	return
}

func checkMetricType(metricType string) bool {
	return (metricType != "counter") && (metricType != "gauge")
}
