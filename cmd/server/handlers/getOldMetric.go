package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) getOldMetric(w http.ResponseWriter, r *http.Request) {
	metric := getMetricFromRequest(r)

	if checkMetricType(metric.MetricType) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}

	if err := h.service.GetOldMetric(&metric); err != nil {
		re, ok := err.(*models.RequestError)
		if ok {
			http.Error(w, re.Err.Error(), re.StatusCode)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	var (
		i int64
		f float64
	)

	if metric.Value != nil {
		f = *metric.Value
	} else {
		f = 0
	}

	if metric.Delta != nil {
		i = *metric.Delta
	} else {
		i = 0
	}

	fmt.Println("^^^^SetOldMetric^^^^^", metric.ID, i, f)

	var (
		httpStatus int
		strValue   string
	)

	switch metric.MetricType {
	case "counter":
		httpStatus = http.StatusOK
		strValue = strconv.Itoa(int(*metric.Delta))
	case "gauge":
		httpStatus = http.StatusOK
		strValue = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	default:
		httpStatus = http.StatusNotImplemented

	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(httpStatus)

	if httpStatus == http.StatusOK {
		w.Write([]byte(strValue))
	}

}
