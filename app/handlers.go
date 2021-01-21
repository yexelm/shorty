package app

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yexelm/shorty/metrics"

	"github.com/gomodule/redigo/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	errEmptyShortCode    = errors.New("empty short code")
	errEmptyRequestBody  = errors.New("empty request body")
	errShortCodeNotFound = errors.New("the requested short code not found")

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

// Shorty is the handler which does different things based on request method.
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
			http.Error(w, errEmptyShortCode.Error(), http.StatusBadRequest)
			return
		}

		long, err := a.db.LongByShort(short)
		if err != nil {
			if err == redis.ErrNil {
				http.Error(w, errShortCodeNotFound.Error(), http.StatusNotFound)
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
			http.Error(w, errEmptyRequestBody.Error(), http.StatusBadRequest)
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
