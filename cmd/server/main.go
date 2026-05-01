package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fourseasonsm/job_manager/internal/handler"
	"github.com/fourseasonsm/job_manager/internal/repository"
	"github.com/fourseasonsm/job_manager/internal/service"
	"github.com/fourseasonsm/job_manager/internal/worker"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	queue := make(chan string, 10)
	repo := repository.NewInMemoryRepository()
	svc := service.NewJobService(repo, queue)
	h := handler.NewJobHandler(svc)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", h.CreateJob)
	mux.HandleFunc("GET /jobs", h.ListJobs)
	mux.HandleFunc("GET /jobs/{id}", h.GetJob)

	ctx, cancel := context.WithCancel(context.Background())
	pool := worker.NewPool(3, repo, queue)
	pool.Start(ctx)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		slog.Info("shutting down...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		srv.Shutdown(shutdownCtx)
	}()

	slog.Info("server started", "addr", ":8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}

	pool.Wait()
	slog.Info("server stopped")
}
