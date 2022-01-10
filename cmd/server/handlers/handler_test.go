package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// )

// func TestNewRouter(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want *chi.Mux
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewRouter(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewRouter() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

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
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, _ := testRequest(t, ts, "POST", "/update/counter/testCounter/100")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// assert.Equal(t, "brand:renault", body)

	resp, _ = testRequest(t, ts, "POST", "/update/counter/testCounter/invalid_value")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, _ = testRequest(t, ts, "POST", "/update/counter/")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	resp, _ = testRequest(t, ts, "POST", "/update/unknown/testCounter/100")
	assert.Equal(t, http.StatusNotImplemented, resp.StatusCode)

}
