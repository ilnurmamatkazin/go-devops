package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseCounterMetric(w http.ResponseWriter, r *http.Request) {
	fmt.Println("$$$$$")
	if len(storageCounter) == 0 {
		storageCounter = make(map[string]int)
	}

	valueMetric := chi.URLParam(r, "valueMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	i, err := strconv.Atoi(valueMetric)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s value is not integer", nameMetric), http.StatusBadRequest)
		return
	}

	mutexCounter.Lock()
	storageCounter[nameMetric] = storageCounter[nameMetric] + i
	mutexCounter.Unlock()

	fmt.Println("@@@@@@", storageCounter)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}
