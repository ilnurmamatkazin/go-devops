package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/pg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestNewRouter(t *testing.T) {
	cfg := models.Config{
		Address:       models.Address,
		StoreInterval: models.StoreInterval,
		StoreFile:     models.StoreFile,
		Restore:       models.Restore,
		Key:           models.Key,
		Database:      models.Database,
	}

	if err := env.Parse(&cfg); err != nil {
		log.Println("Ошибка чтения конфигурации")
		os.Exit(2)
	}

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
	hendler := NewHandler(service)
	router := hendler.NewRouter()

	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, _ := testRequest(t, ts, "POST", "/update/counter/testCounter/100")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	resp, _ = testRequest(t, ts, "POST", "/update/counter/testCounter/invalid_value")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	resp, _ = testRequest(t, ts, "POST", "/update/counter/")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	defer resp.Body.Close()

	resp, _ = testRequest(t, ts, "POST", "/update/unknown/testCounter/100")
	assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)
	defer resp.Body.Close()

}
