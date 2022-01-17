package main

import (
	"fmt"
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
	ADDRESS       = "localhost:8080"
	STOREINTERVAL = 300
	STOREFILE     = "devops-metrics-db.json"
	RESTORE       = true
)

func main() {
	cfg := models.Config{
		Address:       ADDRESS,
		StoreInterval: STOREINTERVAL,
		StoreFile:     STOREFILE,
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

	fmt.Println("Address ", cfg.Address, ":"+strings.Split(cfg.Address, ":")[1])

	go http.ListenAndServe(":"+strings.Split(cfg.Address, ":")[1], r)

	// fmt.Println("Server started...")

	<-quit

	fmt.Println("@@@@@")

	// m.SaveToFile()

	// fmt.Println("Server shutdown")
}
