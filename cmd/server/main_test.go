package main

import (
	"bytes"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilnurmamatkazin/go-devops/cmd/server/handlers"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

func TestParseMetric(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "simple test #1 - good",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
			},
			request: "/update/counter/testCounter/100",
		},
		{
			name: "simple test #2 - invalid value",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
			request: "/update/counter/testCounter/invalid_value",
		},
		{
			name: "simple test #3 - without_id",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
			request: "/update/counter/",
		},
		{
			name: "simple test #4 - unknown",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  501,
			},
			request: "/update/unknown/testCounter/100",
		},
		// {
		// 	name: "simple test #5 - bad request",
		// 	want: want{
		// 		contentType: "text/json; charset=utf-8",
		// 		statusCode:  400,
		// 	},
		// 	request: "/update/counter/testCounter/100",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, 100)

			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlers.ParseCounterMetric)
			h.ServeHTTP(w, request)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			// assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			// userResult, err := ioutil.ReadAll(result.Body)
			// require.NoError(t, err)
			// err = result.Body.Close()
			// require.NoError(t, err)

			// var user User
			// err = json.Unmarshal(userResult, &user)
			// require.NoError(t, err)

			// assert.Equal(t, tt.want.user, user)
		})
	}
}
