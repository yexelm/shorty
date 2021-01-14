package app

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/yexelm/shorty/store"
)

// App keeps Redis state and the server.
type App struct {
	db     *store.Storage
	server *http.Server
}

// New returns an instance of the application server.
func New(storage *store.Storage, appPort int) *App {
	s := App{db: storage}

	s.server = &http.Server{
		Addr:         ":" + strconv.Itoa(appPort),
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		Handler:      s.newAPI(),
	}

	return &s
}

// Serve starts the application server.
func (a *App) Serve(port int) {
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("app server error: %v", err)
		}
	}()
	log.Printf("App started at port %v", port)
}

// Stop closes the storage and stops the application server. It logs any encountered errors.
func (a *App) Stop() {
	a.db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Printf("failed to gracefully shut down api server due to: %v", err)
	}

	log.Println("successfully shut down api server")
}
