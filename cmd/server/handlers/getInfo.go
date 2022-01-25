package handlers

import (
	"net/http"
)

func (h *Handler) getInfo(w http.ResponseWriter, r *http.Request) {
	html := h.service.GetInfo()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(html))
}
