package config

import (
	"os"
	"strconv"
)

const (
	redisURL, defaultRedisURL             = "REDIS_URL", "redis:6379"
	hostPort, defaultHostPort             = "HOST_PORT", 8080
	containerPort, defaultContainerPort   = "CONTAINER_PORT", 8080
	dbNum, defaultDbNum                   = "DB_NUM", 0
	testHandlersDb, defaultTestHandlersDb = "TEST_HANDLERS_DB", 12
	testStorageDb, defaultTestStorageDb   = "TEST_STORAGE_DB", 13
)

// Config contains app configuration
type Config struct {
	RedisURL       string
	HostPort       int
	ContainerPort  int
	DbNum          int
	TestHandlersDb int
	TestStorageDb  int
}

// New returns a new instance of Config
func New() *Config {
	c := Config{}
	c.RedisURL = setStringField(redisURL, defaultRedisURL)

	c.HostPort = setIntField(hostPort, defaultHostPort)
	c.ContainerPort = setIntField(containerPort, defaultContainerPort)
	c.DbNum = setIntField(dbNum, defaultDbNum)
	c.TestHandlersDb = setIntField(testHandlersDb, defaultTestHandlersDb)
	c.TestStorageDb = setIntField(testStorageDb, defaultTestStorageDb)

	return &c
}

func setIntField(key string, defaultValue int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	intV, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue
	}

	return intV
}

func setStringField(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return v
}
