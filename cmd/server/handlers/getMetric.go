package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) getMetric(w http.ResponseWriter, r *http.Request) {
	var (
		metric models.Metric
		err    error
	)

	if err = json.NewDecoder(r.Body).Decode(&metric); err != nil {
		fmt.Println(r.URL.Path)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.service.GetMetric(&metric); err != nil {
		re, ok := err.(*models.RequestError)
		if ok {
			http.Error(w, re.Err.Error(), re.StatusCode)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
