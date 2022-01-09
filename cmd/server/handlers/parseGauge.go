package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseGaugeMetric(w http.ResponseWriter, r *http.Request) {
	if len(storageGauge) == 0 {
		storageGauge = make(map[string]float64)
	}

	valueMetric := chi.URLParam(r, "valueMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	f, err := strconv.ParseFloat(valueMetric, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s value is not float", nameMetric), http.StatusBadRequest)
		return
	}

	mutexGauge.Lock()
	storageGauge[nameMetric] = f
	mutexGauge.Unlock()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}
