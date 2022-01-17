package main

import (
	"flag"
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
	STOREINTERVAL = "300s"
	STOREFILE     = "/tmp/devops-metrics-db.json"
	RESTORE       = true
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fmt.Println("env.Parse", err.Error())
		os.Exit(2)
	}

	fmt.Println(cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	m := memory.NewMemoryRepository(cfg)
	s := service.NewService(m)
	h := handlers.New(s)
	r := h.NewRouter()

	go http.ListenAndServe(":"+strings.Split(cfg.Address, ":")[1], r)

	<-quit

	m.SaveToFile()
}

func parseConfig() (cfg models.Config, err error) {
	address := flag.String("a", ADDRESS, "a address")
	restore := flag.Bool("r", RESTORE, "a restore")
	storeInterval := flag.String("i", STOREINTERVAL, "a store_interval")
	storeFile := flag.String("f", STOREFILE, "a store_file")

	flag.Parse()

	cfg.Address = *address
	cfg.Restore = *restore
	cfg.StoreInterval = *storeInterval
	cfg.StoreFile = *storeFile

	err = env.Parse(&cfg)

	return
}
