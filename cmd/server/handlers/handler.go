package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) NewRouter() *chi.Mux {
	// определяем роутер chi
	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{typeMetric}/{nameMetric}/{valueMetric}", h.parseOldMetric)
		// r.Post("/gauge/{nameMetric}/{valueMetric}", h.parseGaugeMetric)
		// r.Post("/{unknown}/{nameMetric}/{valueMetric}", func(w http.ResponseWriter, r *http.Request) {
		// 	w.WriteHeader(http.StatusNotImplemented)
		// })

		r.With(middleware.AllowContentType("application/json")).Post("/", h.parseMetric)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", h.getInfo)
		r.Get("/value/{typeMetric}/{nameMetric}", h.getOldMetric)
		// r.Get("/value/counter/{nameMetric}", h.getOldMetric)

		r.With(middleware.AllowContentType("application/json")).Post("/value/", h.getMetric)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})

	return r

}
