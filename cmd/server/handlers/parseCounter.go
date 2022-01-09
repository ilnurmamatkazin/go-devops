package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseCounterMetric(w http.ResponseWriter, r *http.Request) {
	if len(storageCounter) == 0 {
		storageCounter = make(map[string][]int)
	}

	valueMetric := chi.URLParam(r, "valueMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	var counter []int

	fmt.Println("$$$$", r.URL.Path, valueMetric, nameMetric)

	i, err := strconv.Atoi(valueMetric)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s value is not integer", nameMetric), http.StatusBadRequest)
		return
	}

	mutexCounter.Lock()
	counter = storageCounter[nameMetric]

	if counter == nil {
		counter = make([]int, 0)
	}

	counter = append(counter, i)

	storageCounter[nameMetric] = counter

	mutexCounter.Unlock()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}
