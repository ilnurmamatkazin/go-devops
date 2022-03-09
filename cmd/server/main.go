// Сервис сохранения системных метрик в системе.
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
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/pg"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Println("env.Parse", err.Error())
		os.Exit(2)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	db, err := pg.NewRepository(&cfg)
	if err != nil {
		log.Println("ошибка подключения к бд: ", err.Error())
	} else {
		defer func() {
			db.Close()
		}()
	}
	repository := storage.NewStorage(&cfg, db)

	if err = repository.Metric.ConnectPG(); err != nil {
		log.Println("ошибка загрузки сохраненых данных", err.Error())
		os.Exit(2)
	}

	service := service.NewService(&cfg, repository)
	hendler := handlers.NewHandler(service)
	router := hendler.NewRouter()

	go http.ListenAndServe(":"+strings.Split(cfg.Address, ":")[1], router)

	<-quit
	repository.Metric.Save()
}

// parseConfig функция получения значений флагов и переменных среды.
func parseConfig() (cfg models.Config, err error) {
	if !flag.Parsed() {
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
	}

	err = env.Parse(&cfg)

	return
}
