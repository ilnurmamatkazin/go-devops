// Сервис сохранения системных метрик в системе.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/pg"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/transport/grpc"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/transport/http/handlers"
	"github.com/ilnurmamatkazin/go-devops/internal/model"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	build := model.NewBuild(buildVersion, buildDate, buildCommit)
	build.Print()

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
		db.Conn = nil
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
	hendler := handlers.NewHandler(&cfg, service)
	router := hendler.NewRouter()

	go http.ListenAndServe(":"+strings.Split(cfg.Address, ":")[1], router)
	go grpc.StartGRPC(cfg, service)

	<-quit

	if db.Conn != nil {
		repository.Metric.Save()
	}

}

// parseConfig функция получения значений флагов и переменных среды.
func parseConfig() (cfg models.Config, err error) {
	if !flag.Parsed() {
		config := flag.String("c", models.JSONConfig, "a json config")

		if *config != "" {
			data, err := os.ReadFile(*config)

			if err != nil {
				log.Printf("os.ReadFile error: %s", err.Error())
			} else {
				err = json.Unmarshal(data, &cfg)
				if err != nil {
					log.Printf("json.Unmarshal error: %s", err.Error())
				}
			}
		}

		address := flag.String("a", models.Address, "a address")
		restore := flag.Bool("r", models.Restore, "a restore")
		storeInterval := flag.String("i", models.StoreInterval, "a store_interval")
		storeFile := flag.String("f", models.StoreFile, "a store_file")
		key := flag.String("k", models.Key, "a secret key")
		database := flag.String("d", models.Database, "a database")
		privateKey := flag.String("crypto-key", models.PrivateKey, "a crypto key")
		trustedSubnet := flag.String("t", models.TrustedSubnet, "a CIDR")

		flag.Parse()

		cfg.Address = *address
		cfg.Restore = *restore
		cfg.StoreInterval = *storeInterval
		cfg.StoreFile = *storeFile
		cfg.Key = *key
		cfg.Database = *database
		cfg.PrivateKey = *privateKey
		cfg.TrustedSubnet = *trustedSubnet

	}

	err = env.Parse(&cfg)

	return
}
