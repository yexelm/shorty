package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"

	"github.com/yexelm/shorty/handlers"
)

func main() {
	env := handlers.LoadEnvironment()

	srv := &fasthttp.Server{
		Handler:      env.Handle,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	go metrics()

	go stop(srv)

	err := srv.ListenAndServe(":" + strconv.Itoa(env.Config.HostPort))
	if err != nil {
		log.Fatal(err)
	}
}

func metrics() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func stop(s *fasthttp.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("received %v signal, shutting down", sig)
	err := s.Shutdown()
	if err != nil {
		log.Printf("failed to gracefully shutdown the server due to %v", err)
	}

	log.Println("application stopped")
}
