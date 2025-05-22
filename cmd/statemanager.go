package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Tubular-Bytes/statesman/pkg/router"
	"github.com/gorilla/mux"
)

func main() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	r := mux.NewRouter()

	r.HandleFunc("/state", router.HandleGetState).Methods(http.MethodGet)
	r.HandleFunc("/state", router.HandlePostState).Methods(http.MethodPost)
	r.HandleFunc("/lock", router.HandleLock)
	r.HandleFunc("/unlock", router.HandleUnlock)
	r.HandleFunc("/health", router.HandleHealth).Methods(http.MethodGet)

	daemon := &http.Server{
		Addr:    ":3111",
		Handler: r,
	}

	go func() {
		if err := daemon.ListenAndServe(); err != nil {
			slog.Error("failed to start server", "error", err)
		}
	}()

	// Wait for a signal
	<-sigchan
	daemon.Shutdown(context.Background())
}
