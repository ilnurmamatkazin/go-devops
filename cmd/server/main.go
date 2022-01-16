package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/handlers"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/memory"
)

const (
	ADDRESS        = "localhost:8080"
	STORE_INTERVAL = 300
	STORE_FILE     = "./tmp/devops-metrics-db.json"
	RESTORE        = true
)

func main() {
	cfg := models.Config{
		Address:       ADDRESS,
		StoreInterval: STORE_INTERVAL,
		StoreFile:     STORE_FILE,
		Restore:       RESTORE,
	}

	if err := env.Parse(&cfg); err != nil {
		os.Exit(2)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	m := memory.NewMemoryRepository(cfg)
	s := service.NewService(m)
	h := handlers.New(s)
	r := h.NewRouter()

	go http.ListenAndServe(":"+strings.Split(cfg.Address, ":")[1], r)

	// fmt.Println("Server started...")

	<-quit

	m.SaveToFile()

	// fmt.Println("Server shutdown")
}
