package handlers

import (
	"net/http"
)

// getInfo функция для получения html страницы со значением всех метрик, собранных системой.
func (h *Handler) getInfo(w http.ResponseWriter, r *http.Request) {
	html := h.service.GetInfo()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(html))
}
