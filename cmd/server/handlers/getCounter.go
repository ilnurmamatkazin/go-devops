package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func getCounter(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")

	mutexCounter.Lock()
	value := storageCounter[nameMetric]
	mutexCounter.Unlock()

	if value == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	valueText := []string{}

	for _, item := range value {
		valueText = append(valueText, strconv.Itoa(item))
	}

	w.Write([]byte(strings.Join(valueText, ", ")))

}
