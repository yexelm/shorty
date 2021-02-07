package app_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/yexelm/shorty/app"
	"github.com/yexelm/shorty/config"
	"github.com/yexelm/shorty/store"

	"github.com/alicebob/miniredis/v2"
)

var a *app.App

func TestMain(m *testing.M) {
	cfg := config.New()

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	db, _ := store.New(s.Addr(), 0)
	a = app.New(db, cfg.HostPort)
	code := m.Run()
	conn := db.Pool.Get()
	_, _ = conn.Do("FLUSHDB")
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
			wantBody: app.ErrEmptyShortCode.Error(),
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
			wantBody: app.ErrShortCodeNotFound.Error(),
		},
		{
			name:     "empty POST request body",
			recorder: httptest.NewRecorder(),
			method:   http.MethodPost,
			url:      "",
			body:     bytes.NewReader([]byte("")),
			wantCode: http.StatusBadRequest,
			wantBody: app.ErrEmptyRequestBody.Error(),
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
			resp, _ := ioutil.ReadAll(tc.recorder.Result().Body)
			gotBody := string(resp)
			if gotBody != tc.wantBody {
				t.Errorf("\ngot:  %v\nwant: %v\n", gotBody, tc.wantBody)
			}
		})
	}
}
