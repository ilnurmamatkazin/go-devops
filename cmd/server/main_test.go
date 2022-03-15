package main

import (
	"net"
	"os"
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/handlers"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/pg"
	"github.com/stretchr/testify/assert"
)

func Test_parseConfig(t *testing.T) {
	tests := []struct {
		name  string
		env   string
		value string
		kind  string
	}{
		{name: "positive env=RESTORE", env: "RESTORE", value: "true", kind: "positive"},
		{name: "negative env=RESTORE", env: "RESTORE", value: "12.56", kind: "negative"},
		{name: "positive env=ADDRESS", env: "ADDRESS", value: "127.0.0.1", kind: "positive"},
		{name: "negative env=ADDRESS", env: "ADDRESS", value: "12.56", kind: "negative"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.env, tt.value)
			defer os.Unsetenv(tt.env)

			_, err := parseConfig()

			if tt.env == "RESTORE" {
				if tt.kind == "positive" {
					assert.Nil(t, err)
				} else {
					assert.NotNil(t, err)
				}
			} else if tt.env == "ADDRESS" {
				ip := net.ParseIP(tt.value)

				if tt.kind == "positive" {
					assert.NotNil(t, ip)
				} else {
					assert.Nil(t, ip)
				}
			}
		})
	}
}

func Test_main(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		cfg := models.Config{
			Address:       models.Address,
			Restore:       models.Restore,
			StoreInterval: models.StoreInterval,
			StoreFile:     models.StoreFile,
			Key:           models.Key,
			Database:      models.Database,
		}

		db, err := pg.NewRepository(&cfg)
		if err == nil {
			defer func() {
				db.Close()
			}()

			repository := storage.NewStorage(&cfg, db)

			err = repository.Metric.ConnectPG()
			assert.NoError(t, err)

			service := service.NewService(&cfg, repository)
			hendler := handlers.NewHandler(service)
			_ = hendler.NewRouter()
		}

	})

}
