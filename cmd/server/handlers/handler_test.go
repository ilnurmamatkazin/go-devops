package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/models"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/service"
	"github.com/ilnurmamatkazin/go-devops/cmd/server/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ADDRESS       = "localhost:8080"
	STOREINTERVAL = 300
	STOREFILE     = "./tmp/devops-metrics-db.json"
	RESTORE       = true
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
		Address:       ADDRESS,
		StoreInterval: STOREINTERVAL,
		StoreFile:     STOREFILE,
		Restore:       RESTORE,
	}

	if err := env.Parse(&cfg); err != nil {
		os.Exit(2)
	}

	m := memory.NewMemoryRepository(cfg)
	s := service.NewService(m)
	h := New(s)
	r := h.NewRouter()

	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, _ := testRequest(t, ts, "POST", "/update/counter/testCounter/100")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// assert.Equal(t, "brand:renault", body)
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
