package handlers

import (
	"net/http"
)

// GetInfo функция для получения html страницы со значением всех метрик, собранных системой.
func (h *Handler) GetInfo(w http.ResponseWriter, r *http.Request) {
	html := h.Service.GetInfo()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(html))
}
