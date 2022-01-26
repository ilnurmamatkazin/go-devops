package main

import (
	"flag"
	"fmt"
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
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Println("env.Parse", err.Error())
		os.Exit(2)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	fmt.Println("@@@@", cfg)
	repository, err := storage.New(cfg)
	if err != nil {
		log.Println("ошибка подключения к бд: ", err.Error())
		//os.Exit(2)
	} else {
		defer repository.Close()
	}

	service := service.New(cfg, repository)
	hendler := handlers.New(service)
	router := hendler.NewRouter()

	go http.ListenAndServe(":"+strings.Split(cfg.Address, ":")[1], router)

	<-quit

	repository.Save()
}

func parseConfig() (cfg models.Config, err error) {
	address := flag.String("a", models.Address, "a address")
	restore := flag.Bool("r", models.Restore, "a restore")
	storeInterval := flag.String("i", models.StoreInterval, "a store_interval")
	storeFile := flag.String("f", models.StoreFile, "a store_file")
	key := flag.String("k", models.Key, "a secret key")
	database := flag.String("d", models.Database, "a database")

	flag.Parse()

	cfg.Address = *address
	cfg.Restore = *restore
	cfg.StoreInterval = *storeInterval
	cfg.StoreFile = *storeFile
	cfg.Key = *key
	cfg.Database = *database

	err = env.Parse(&cfg)

	return
}
