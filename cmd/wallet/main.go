package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mux := http.NewServeMux()
	mux.HandleFunc("/GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Health Check", "remote", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"Status":"Ok"}` + "\n"))
	})
	srv := http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	go func() {
		logger.Info("Wallet Service listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server Failed ", "err", err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	logger.Info("Shutting down gracefuullly")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Shutdown failed", "err", err)
	}
	logger.Info("Server stopped")

}

//if server stops before 10 second
// t=0.0   SIGTERM → ctx.Done() closes → main unblocks
// t=0.0   main creates shutdownCtx (10s), calls srv.Shutdown(shutdownCtx)
// t=0.0   Shutdown closes listener
// t=0.0   ListenAndServe returns ErrServerClosed  ← NOW, not at 10s
// t=0.0   goroutine sees ErrServerClosed → ignores it, goroutine exits quietly
// t=0.0→  Shutdown waits for the in-flight handler to finish
// t=0.3   handler finishes, connection drained
// t=0.3   Shutdown returns nil
// t=0.3   log "Server stopped"
// t=0.3   main returns → deferred cancels run → process exits

//if server didnt complete stopping in 10 seconds
// t=0.0   ListenAndServe already returned ErrServerClosed (goroutine gone)
// t=10.0  Shutdown gives up → returns DeadlineExceeded
// t=10.0  log "Shutdown failed"
// t=10.0  main returns → process exits → stuck handler goroutine dies with process
