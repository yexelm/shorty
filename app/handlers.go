package app

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yexelm/shorty/metrics"

	"github.com/gomodule/redigo/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrEmptyShortCode    = errors.New("empty short code")
	ErrEmptyRequestBody  = errors.New("empty request body")
	ErrShortCodeNotFound = errors.New("the requested short code not found")

	// for metrics
	methodToOperation = map[string]string{
		http.MethodGet:  "Get original URL by short alias",
		http.MethodPost: "Shorten the long URL",
	}
)

func (a *App) newAPI() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", a.Shorty)
	m.Handle("/metrics", promhttp.Handler())

	return m
}

// Shorty is the handler which acts according to the request method.
// GET: returns original URL by its short alias.
// POST: generates and returns short alias for the given URL.
func (a *App) Shorty(w http.ResponseWriter, r *http.Request) {
	label, ok := methodToOperation[r.Method]
	if !ok {
		label = "operation not implemented"
	}

	obs := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.LatencyHandler.WithLabelValues(label).Observe(v)
	}))
	defer obs.ObserveDuration()

	switch r.Method {
	case http.MethodGet:
		short := []byte(strings.TrimPrefix(r.URL.Path, "/"))
		if len(short) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, ErrEmptyShortCode)
			return
		}

		long, err := a.db.LongByShort(short)
		if err != nil {
			if err == redis.ErrNil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, ErrShortCodeNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}

		_, _ = w.Write(long)
	case http.MethodPost:
		longURL, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		if len(longURL) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, ErrEmptyRequestBody)
			return
		}

		short, err := a.db.ShortByLong(longURL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		buf := new(bytes.Buffer)
		buf.WriteString(r.Host + "/")
		buf.Write(short)

		_, _ = w.Write(buf.Bytes())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
