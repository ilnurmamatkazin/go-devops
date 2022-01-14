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

func New() *Handler {
	return &Handler{
		service: service.NewService(),
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
	// r.Use(middleware.AllowContentType("application/json"))

	// r.Route("/", func(r chi.Router) {
	// 	r.Get("/", h.getInfo)
	// })

	r.Route("/update", func(r chi.Router) {
		r.Post("/counter/{nameMetric}/{valueMetric}", h.parseCounterMetric)
		r.Post("/gauge/{nameMetric}/{valueMetric}", h.parseGaugeMetric)
		r.Post("/{unknown}/{nameMetric}/{valueMetric}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotImplemented)
		})

		r.With(middleware.AllowContentType("application/json")).Post("/", h.parseMetric)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", h.getInfo)
		r.Get("/value/gauge/{nameMetric}", h.getGauge)
		r.Get("/value/counter/{nameMetric}", h.getCounter)

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
