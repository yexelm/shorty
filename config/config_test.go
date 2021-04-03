package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setStringField(t *testing.T) {
	t.Parallel()
	ao := assert.New(t)
	os.Setenv("foo", "bar")

	type testData struct {
		tCase        string
		key          string
		defaultValue string
		expected     string
	}

	testTable := []testData{
		{
			tCase:        "success",
			key:          "foo",
			defaultValue: "bar",
			expected:     "bar",
		},
		{
			tCase:        "default value",
			key:          "bar",
			defaultValue: "baz",
			expected:     "baz",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.tCase, func(t *testing.T) {
			ao.Equal(tc.expected, setStringField(tc.key, tc.defaultValue))
		})
	}
}

func Test_setIntField(t *testing.T) {
	t.Parallel()
	ao := assert.New(t)
	os.Setenv("foo", "125")
	os.Setenv("baz", "this will fail strconv")

	type testData struct {
		tCase        string
		key          string
		defaultValue int
		expected     int
	}

	testTable := []testData{
		{
			tCase:        "success",
			key:          "foo",
			defaultValue: 0,
			expected:     125,
		},
		{
			tCase:        "default value",
			key:          "bar",
			defaultValue: 12,
			expected:     12,
		},
		{
			tCase:        "failed to parse value from env",
			key:          "baz",
			defaultValue: 13,
			expected:     13,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.tCase, func(t *testing.T) {
			ao.Equal(tc.expected, setIntField(tc.key, tc.defaultValue))
		})
	}
}

func Test_New(t *testing.T) {
	t.Parallel()
	ao := assert.New(t)

	type testData struct {
		tcase    string
		expected *Config
	}

	testTable := []testData{
		{
			tcase: "default values",
			expected: &Config{
				RedisURL:      defaultRedisURL,
				HostPort:      defaultHostPort,
				ContainerPort: defaultContainerPort,
				DbNum:         defaultDbNum,
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.tcase, func(t *testing.T) {
			ao.Equal(tc.expected, New())
		})
	}
}
