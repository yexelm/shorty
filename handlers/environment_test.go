package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/yexelm/shorty/config"
)

type MockEnv struct {
	Ctrl  *gomock.Controller
	Cache *MockLongerShorter
}

func loadMockEnv(t *testing.T) (*MockEnv, *Environment) {
	ctrl := gomock.NewController(t)

	cache := NewMockLongerShorter(ctrl)

	mockEnv := &MockEnv{
		Ctrl:  ctrl,
		Cache: cache,
	}

	env := &Environment{
		Config: config.New(),
		Cache:  cache,
	}

	return mockEnv, env
}
