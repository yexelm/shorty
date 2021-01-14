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
	cfg := config.New()

	storage, err := store.New(cfg.RedisURL, cfg.DbNum)
	if err != nil {
		log.Fatal(err)
	}

	application := app.New(storage, cfg.HostPort)
	go application.Serve(cfg.HostPort)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("received %v signal, shutting down", sig)
	application.Stop()

	log.Println("application stopped")
}
