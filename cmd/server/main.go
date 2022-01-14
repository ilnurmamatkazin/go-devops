package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/handlers"
)

const (
	ADDRESS        = "localhost:8080"
	STORE_INTERVAL = 300
	STORE_FILE     = "./tmp/devops-metrics-db.json"
	RESTORE        = true
)

type Config struct {
	Address       string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	StoreFile     string `env:"STORE_FILE"`
	Restore       bool   `env:"RESTORE"`
}

func main() {
	cfg := Config{
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

	h := handlers.New()
	r := h.NewRouter()

	go http.ListenAndServe(":"+strings.Split(cfg.Address, ":")[1], r)

	// fmt.Println("Server started...")

	<-quit

	// fmt.Println("Server shutdown")
}
