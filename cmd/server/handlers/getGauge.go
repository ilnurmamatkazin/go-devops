package handlers

import (
	// "bytes"
	// "encoding/binary"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func getGauge(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")

	mutexGauge.Lock()
	value := storageGauge[nameMetric]
	mutexGauge.Unlock()

	if value == 0 {
		http.NotFound(w, r)
		return
	}

	fmt.Println(value)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))

}
