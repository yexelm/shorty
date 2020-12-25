package app

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gomodule/redigo/redis"
)

func (a *App) newAPI() http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/", a.Shorty)

	return m
}

// Shorty is the handler which does different things based on request method.
// GET: returns original URL by its short alias.
// POST: generates and returns short alias for the given URL.
func (a *App) Shorty(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		short := []byte(strings.TrimPrefix(r.URL.Path, "/"))
		if len(short) == 0 {
			http.Error(w, "empty short code", http.StatusBadRequest)
			return
		}

		long, err := a.db.LongByShort(short)
		if err != nil {
			if err == redis.ErrNil {
				http.Error(w, redis.ErrNil.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(long)
	case http.MethodPost:
		longURL, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(longURL) == 0 {
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}

		short, err := a.db.ShortByLong(longURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf := new(bytes.Buffer)
		buf.WriteString(r.Host + "/")
		buf.Write(short)

		_, _ = w.Write(buf.Bytes())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
