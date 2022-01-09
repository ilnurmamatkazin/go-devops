package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func getCounter(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")

	mutexCounter.Lock()
	value := storageCounter[nameMetric]
	mutexCounter.Unlock()

	if value == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(strconv.Itoa(value)))

}
