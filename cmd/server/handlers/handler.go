package handlers

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
)

// Handler структура, формирующая слой работы с функциями, принимающими апи запросы.
type Handler struct {
	service *service.Service
}

// NewHandler конструктор для структуры Handler.
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// NewRouter функция, которая создает роутер для приема апи запросов.
func (h *Handler) NewRouter() *chi.Mux {
	// определяем роутер chi
	r := chi.NewRouter()

	// зададим встроенные middleware, чтобы улучшить стабильность приложения
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewareGzip)

	r.Route("/update", func(r chi.Router) {
		r.Post("/{typeMetric}/{nameMetric}/{valueMetric}", h.parseOldMetric)
		r.With(middleware.AllowContentType("application/json")).Post("/", h.parseMetric)
	})

	r.Route("/ping", func(r chi.Router) {
		r.With(middleware.AllowContentType("application/json")).Get("/", h.ping)
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", h.getInfo)
		r.Get("/value/{typeMetric}/{nameMetric}", h.getOldMetric)
		r.With(middleware.AllowContentType("application/json")).Post("/value/", h.getMetric)
		r.With(middleware.AllowContentType("application/json")).Post("/update/", h.parseMetric)
		r.With(middleware.AllowContentType("application/json")).Post("/updates/", h.parseMetrics)
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

// gzipWriter структура для работы с запросами использующих сжатие.
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write функция отправки сжатого контента.
func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// middlewareGzip миделваер, обрабатывающий запросы с жатым контентом.
func middlewareGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
