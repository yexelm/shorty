package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yexelm/shorty/app"
	"github.com/yexelm/shorty/config"
	"github.com/yexelm/shorty/store"
)

func main() {

	redisURL := config.GetEnv("REDIS_URL", "localhost:6379")
	appPort := config.GetEnv("APP_PORT", ":8080")
	dbNum := config.GetEnv("DB_NUM", "13")

	storage, err := store.New(redisURL, dbNum)
	if err != nil {
		log.Fatal(err)
	}

	a := app.New(storage, appPort)
	go a.Serve(appPort)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("received %v signal, shutting down", sig)
	a.Stop()

	log.Println("application stopped")
}
