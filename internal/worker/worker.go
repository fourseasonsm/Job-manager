package worker

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/fourseasonsm/job_manager/internal/model"
	"github.com/fourseasonsm/job_manager/internal/repository"
)

type Worker struct {
	id    int
	repo  repository.JobRepository
	queue <-chan string
}

func NewWorker(id int, repo repository.JobRepository, queue <-chan string) *Worker {
	return &Worker{id: id, repo: repo, queue: queue}
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case jobId, ok := <-w.queue:
			if !ok {
				slog.Info("worker stopped", "worker_id", w.id)
				return
			}
			w.process(ctx, jobId)

		case <-ctx.Done():
			slog.Info("worker stopped", "worker_id", w.id)
			return
		}
	}
}

func (w *Worker) process(ctx context.Context, jobId string) {
	job, err := w.repo.GetById(jobId)
	if err != nil {
		slog.Error("worker: job not found", "job_id", jobId, "worker_id", w.id)
		return
	}
	now := time.Now()
	job.Status = model.StatusRunning
	job.StartedAt = now
	w.repo.Update(job)

	slog.Info("worker: processing job", "job_id", jobId, "worker_id", w.id)

	duration := time.Duration(1+rand.Intn(5)) * time.Second
	select {
	case <-time.After(duration):
	case <-ctx.Done():
		return
	}

	finished := time.Now()
	job.FinishedAt = finished

	if rand.Float32() < 0.1 {
		job.Status = model.StatusFailed
		slog.Warn("worker: job failed", "job_id", jobId)
	} else {
		job.Status = model.StatusDone
		slog.Info("worker: job done", "job_id", jobId)
	}

	w.repo.Update(job)
}
