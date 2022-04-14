package handlers

import (
	"net"
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

func checkTrustedSubnet(server, client string) bool {
	if server == "" {
		return true
	}

	_, ipnetServer, _ := net.ParseCIDR(server)
	ipClient := net.ParseIP(client)

	if !ipnetServer.Contains(ipClient) {
		return false
	}

	return true

}
