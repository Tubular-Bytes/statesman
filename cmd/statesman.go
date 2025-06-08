package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Tubular-Bytes/statesman/pkg/router"
	"github.com/Tubular-Bytes/statesman/pkg/router/middleware"
	"github.com/gorilla/mux"
)

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	r := mux.NewRouter()
	r.Use(middleware.Logging)

	r.HandleFunc("/state", router.HandleState).Methods(
		http.MethodGet,
		http.MethodPost,
		router.MethodLock,
		router.MethodUnlock,
	)
	r.HandleFunc("/health", router.HandleHealth).Methods(http.MethodGet)

	daemon := &http.Server{
		Addr:    ":3111",
		Handler: r,
	}

	go func() {
		if err := daemon.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("failed to start server", "error", err)
		}
	}()

	slog.Info("server started", "address", daemon.Addr)

	// Wait for a signal
	<-sigchan
	daemon.Shutdown(context.Background())
}
