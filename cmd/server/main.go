package main

import (
	"flag"
	"log"
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

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Println("env.Parse", err.Error())
		os.Exit(2)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	repository := memory.NewMemoryRepository(cfg)
	service := service.NewService(repository)
	hendler := handlers.New(service)
	router := hendler.NewRouter()

	go http.ListenAndServe(":"+strings.Split(cfg.Address, ":")[1], router)

	<-quit

	repository.SaveToFile()
}

func parseConfig() (cfg models.Config, err error) {
	address := flag.String("a", models.Address, "a address")
	restore := flag.Bool("r", models.Restore, "a restore")
	storeInterval := flag.String("i", models.StoreInterval, "a store_interval")
	storeFile := flag.String("f", models.StoreFile, "a store_file")

	flag.Parse()

	cfg.Address = *address
	cfg.Restore = *restore
	cfg.StoreInterval = *storeInterval
	cfg.StoreFile = *storeFile

	err = env.Parse(&cfg)

	return
}
