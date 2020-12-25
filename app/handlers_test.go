package app

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/yexelm/shorty/config"
	"github.com/yexelm/shorty/store"
)

var a *App

func TestMain(m *testing.M) {
	dbPort := config.GetEnv("REDIS_URL", "localhost:6379")
	appPort := config.GetEnv("APP_PORT", ":8080")
	dbNum := config.GetEnv("TEST_DB_2", "13")

	db, _ := store.New(dbPort, dbNum)
	a = New(db, appPort)
	code := m.Run()
	c := db.Pool.Get()
	c.Do("FLUSHDB")
	a.Stop()
	os.Exit(code)
}

// TestShorty checks if Shorty handler works correctly
func TestShortyPostAndGet(t *testing.T) {
	tests := []struct {
		name     string
		recorder *httptest.ResponseRecorder
		method   string
		url      string
		body     io.Reader
		wantCode int
		wantBody string
	}{
		{
			name:     "correct POST request",
			recorder: httptest.NewRecorder(),
			method:   http.MethodPost,
			url:      "",
			body:     bytes.NewReader([]byte(`ya.ru`)),
			wantCode: http.StatusOK,
			wantBody: "/b",
		},
		{
			name:     "correct GET request",
			recorder: httptest.NewRecorder(),
			method:   http.MethodGet,
			url:      "/b",
			body:     nil,
			wantCode: http.StatusOK,
			wantBody: "ya.ru",
		},
		{
			name:     "empty short link",
			recorder: httptest.NewRecorder(),
			method:   http.MethodGet,
			url:      "/",
			body:     nil,
			wantCode: http.StatusBadRequest,
			wantBody: "empty short code\n",
		},
		{
			name:     "wrong method",
			recorder: httptest.NewRecorder(),
			method:   http.MethodPut,
			url:      "/",
			body:     nil,
			wantCode: http.StatusMethodNotAllowed,
			wantBody: "",
		},
		{
			name:     "get non-existent short link",
			recorder: httptest.NewRecorder(),
			method:   http.MethodGet,
			url:      "/c",
			body:     nil,
			wantCode: http.StatusNotFound,
			wantBody: redis.ErrNil.Error() + "\n",
		},
		{
			name:     "empty POST request body",
			recorder: httptest.NewRecorder(),
			method:   http.MethodPost,
			url:      "",
			body:     bytes.NewReader([]byte("")),
			wantCode: http.StatusBadRequest,
			wantBody: "request body is empty\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, tc.url, tc.body)
			a.Shorty(tc.recorder, req)

			gotCode := tc.recorder.Code
			if gotCode != tc.wantCode {
				t.Errorf("\ngot:  %v\nwant: %v\n", gotCode, tc.wantCode)
			}
			gotBody := tc.recorder.Body.String()
			if gotBody != tc.wantBody {
				t.Errorf("\ngot:  %v\nwant: %v\n", gotBody, tc.wantBody)
			}
		})
	}
}
