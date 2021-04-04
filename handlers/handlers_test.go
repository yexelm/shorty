package handlers

import (
	"errors"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func initCtx(method, URI string, body []byte) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI(URI)
	ctx.Request.SetHostBytes(ctx.Request.URI().Host())
	ctx.Request.SetBody(body)
	ctx.Request.Header.SetMethod(method)
	return ctx
}

func Test_Handle(t *testing.T) {
	t.Parallel()
	ao := assert.New(t)
	mockEnv, env := loadMockEnv(t)
	defer mockEnv.Ctrl.Finish()

	type testData struct {
		tCase  string
		method string

		expectedCode int
	}

	testTable := []testData{
		{
			tCase:        "method not allowed",
			method:       "PUT",
			expectedCode: fasthttp.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.tCase, func(t *testing.T) {
			ctx := initCtx(tc.method, "", nil)

			env.Handle(ctx)
			ao.Equal(tc.expectedCode, ctx.Response.StatusCode())
		})
	}
}

func Test_longer(t *testing.T) {
	t.Parallel()
	ao := assert.New(t)
	mockEnv, env := loadMockEnv(t)
	defer mockEnv.Ctrl.Finish()

	type testData struct {
		tCase        string
		ctx          *fasthttp.RequestCtx
		URI          string
		expectedFunc func()
		expectedBody string
		expectedCode int
	}

	testTable := []testData{
		{
			tCase:        "empty short code",
			ctx:          nil,
			URI:          "",
			expectedFunc: func() {},

			expectedBody: ErrEmptyShortCode.Error(),
			expectedCode: fasthttp.StatusBadRequest,
		},
		{
			tCase: "not in cache",
			ctx:   nil,
			URI:   "shortcode",
			expectedFunc: func() {
				mockEnv.Cache.EXPECT().Longer([]byte("shortcode")).Return(nil, redis.ErrNil)
			},

			expectedBody: ErrShortCodeNotFound.Error(),
			expectedCode: fasthttp.StatusNotFound,
		},

		{
			tCase: "cache error",
			ctx:   nil,
			URI:   "shortcode",
			expectedFunc: func() {
				mockEnv.Cache.EXPECT().Longer([]byte("shortcode")).Return(nil, errors.New("some cache error"))
			},
			expectedBody: "some cache error",
			expectedCode: fasthttp.StatusInternalServerError,
		},
		{
			tCase: "success",
			ctx:   nil,
			URI:   "shortcode",
			expectedFunc: func() {
				mockEnv.Cache.EXPECT().Longer([]byte("shortcode")).Return([]byte("fullURL"), nil)
			},
			expectedBody: "fullURL",
			expectedCode: fasthttp.StatusOK,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.tCase, func(t *testing.T) {
			tc.ctx = initCtx("GET", tc.URI, nil)
			tc.expectedFunc()
			env.Handle(tc.ctx)

			ao.Equal(tc.expectedCode, tc.ctx.Response.StatusCode())
			ao.Equal(tc.expectedBody, string(tc.ctx.Response.Body()))
		})
	}
}

func Test_shorter(t *testing.T) {
	t.Parallel()
	ao := assert.New(t)
	mockEnv, env := loadMockEnv(t)
	defer mockEnv.Ctrl.Finish()

	type testData struct {
		tCase        string
		ctx          *fasthttp.RequestCtx
		body         []byte
		expectedFunc func()

		expectedBody string
		expectedCode int
	}

	testTable := []testData{
		{
			tCase:        "empty request body",
			ctx:          nil,
			body:         nil,
			expectedFunc: func() {},

			expectedBody: ErrEmptyRequestBody.Error(),
			expectedCode: fasthttp.StatusBadRequest,
		},
		{
			tCase: "error while getting shorter",
			ctx:   nil,
			body:  []byte("originalURL"),
			expectedFunc: func() {
				mockEnv.Cache.EXPECT().Shorter([]byte("originalURL")).Return(nil, errors.New("some error"))
			},
			expectedBody: "some error",
			expectedCode: fasthttp.StatusInternalServerError,
		},
		{
			tCase: "success",
			ctx:   nil,
			body:  []byte("originalURL"),
			expectedFunc: func() {
				mockEnv.Cache.EXPECT().Shorter([]byte("originalURL")).Return([]byte("shortcode"), nil)
			},
			expectedBody: "host.com/shortcode",
			expectedCode: fasthttp.StatusOK,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.tCase, func(t *testing.T) {
			tc.ctx = initCtx("POST", "http://host.com", tc.body)
			tc.expectedFunc()
			env.Handle(tc.ctx)

			ao.Equal(tc.expectedCode, tc.ctx.Response.StatusCode())
			ao.Equal(tc.expectedBody, string(tc.ctx.Response.Body()))
		})
	}
}
