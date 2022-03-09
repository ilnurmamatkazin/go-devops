package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

// getMetricFromRequest функция получения имени и типа метрики из строки запроса.
func getMetricFromRequest(r *http.Request) (metric models.Metric) {
	metric.ID = chi.URLParam(r, "nameMetric")
	metric.MetricType = chi.URLParam(r, "typeMetric")

	return
}

// checkMetricType функция проверки типа метрики.
func checkMetricType(metricType string) bool {
	return (metricType != "counter") && (metricType != "gauge")
}
