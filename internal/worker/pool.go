package worker

import (
	"context"
	"sync"

	"github.com/fourseasonsm/job_manager/internal/repository"
)

type Pool struct {
	workers []*Worker
	wg      sync.WaitGroup
}

func NewPool(size int, repo repository.JobRepository, queue <-chan string) *Pool {
	workers := make([]*Worker, size)
	for i := range workers {
		workers[i] = NewWorker(i+1, repo, queue)
	}
	return &Pool{workers: workers}
}

func (p *Pool) Start(ctx context.Context) {
	for _, w := range p.workers {
		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done()
			w.Run(ctx)
		}(w)
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
