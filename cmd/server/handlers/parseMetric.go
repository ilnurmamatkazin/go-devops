package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) parseMetric(w http.ResponseWriter, r *http.Request) {
	var (
		metric models.Metric
		err    error
	)
	fmt.Println("parseMetric", r.URL.Path)

	if err = json.NewDecoder(r.Body).Decode(&metric); err != nil {
		fmt.Println("!!!!!!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(metric)

	if err = h.service.SetMetric(metric); err != nil {
		fmt.Println("#####", metric)

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
}
