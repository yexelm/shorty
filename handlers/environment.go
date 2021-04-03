package handlers

import (
	"log"

	"github.com/yexelm/shorty/config"
	"github.com/yexelm/shorty/store"
)

type Environment struct {
	Config *config.Config
	Cache  LongerShorter
}

func LoadEnvironment() *Environment {
	cfg := config.New()
	cache, err := store.New(cfg.RedisURL, cfg.DbNum)
	if err != nil {
		log.Fatal(err)
	}

	env := Environment{
		Config: cfg,
		Cache:  cache,
	}

	return &env
}
