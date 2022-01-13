package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
)

func (h *Handler) getCounter(w http.ResponseWriter, r *http.Request) {
	nameMetric := chi.URLParam(r, "nameMetric")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	value, err := h.repository.ReadCounter(nameMetric)
	if err != nil {
		re, ok := err.(*models.RequestError)
		if ok {
			w.WriteHeader(re.StatusCode)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(int(value))))
}
