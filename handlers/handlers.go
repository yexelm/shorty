package handlers

import (
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"

	"github.com/yexelm/shorty/metrics"
)

//go:generate mockgen -source=handlers.go -destination=handlers_mocks.go -package=handlers -self_package=shorty/handlers

var (
	ErrEmptyShortCode    = errors.New("empty short code")
	ErrEmptyRequestBody  = errors.New("empty request body")
	ErrShortCodeNotFound = errors.New("the requested short code not found")

	// for metrics
	methodToOperation = map[string]string{
		fasthttp.MethodGet:  "Get original URI by short alias",
		fasthttp.MethodPost: "Shorten the long URI",
	}
)

type LongerShorter interface {
	Longer(short []byte) ([]byte, error)
	Shorter(long []byte) ([]byte, error)
}

func (env *Environment) Handle(ctx *fasthttp.RequestCtx) {
	m, ok := methodToOperation[string(ctx.Method())]
	if !ok {
		m = "operation not implemented"
	}

	obs := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.LatencyHandler.WithLabelValues(m).Observe(v)
	}))
	defer obs.ObserveDuration()

	switch {
	case ctx.IsGet():
		env.longer(ctx)
	case ctx.IsPost():
		env.shorter(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

// longer returns the original URI for the given short code.
func (env *Environment) longer(ctx *fasthttp.RequestCtx) {
	short := []byte(strings.TrimPrefix(string(ctx.Path()), "/"))

	if len(short) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(ErrEmptyShortCode.Error())
		return
	}

	originalURL, err := env.Cache.Longer(short)
	if err != nil && err == redis.ErrNil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.WriteString(ErrShortCodeNotFound.Error())
		return
	}

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.Write(originalURL)
}

// shorter converts the original URI into the short alias and returns it.
func (env *Environment) shorter(ctx *fasthttp.RequestCtx) {
	longURL, err := io.ReadAll(bytes.NewReader(ctx.Request.Body()))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	if len(longURL) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(ErrEmptyRequestBody.Error())
		return
	}

	short, err := env.Cache.Shorter(longURL)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	buf := &bytes.Buffer{}
	buf.Write(append(ctx.URI().Host(), '/'))
	buf.Write(short)

	ctx.Write(buf.Bytes())
}
