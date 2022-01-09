package handlers

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	mutexCounter, mutexGauge sync.Mutex
	storageCounter           map[string]int
	storageGauge             map[string]float64
)

type Router struct {
}

func NewRouter() *chi.Mux {
	// определяем роутер chi
	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/", getInfo)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/counter/{nameMetric}/{valueMetric}", ParseCounterMetric)
		r.Post("/gauge/{nameMetric}/{valueMetric}", ParseGaugeMetric)
		r.Post("/{unknown}/{nameMetric}/{valueMetric}", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotImplemented)
		})
	})

	// r.Route("/update/gauge", func(r chi.Router) {
	// 	r.Post("/{nameMetric}/{valueMetric}", ParseGaugeMetric)
	// })

	r.Route("/value/gauge", func(r chi.Router) {
		r.Get("/{nameMetric}", getGauge)
	})

	r.Route("/value/counter", func(r chi.Router) {
		r.Get("/{nameMetric}", getCounter)
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
