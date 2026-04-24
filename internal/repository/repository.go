package repository

import (
	"errors"
	"sync"

	"github.com/fourseasonsm/job_manager/internal/model"
)

var (
	ErrJobNotFound  = errors.New("Job not found")
	ErrJobAlrExists = errors.New("Job already exists")
)

type JobRepository interface {
	Save(job model.Job) error
	GetById(id string) (model.Job, error)
	GetAll()
	GetByStatus(status model.Status) ([]model.Job, error)
	Update(job model.Job) error
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	jobs map[string]model.Job
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		jobs: make(map[string]model.Job),
	}
}

func (r *InMemoryRepository) Save(job model.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[job.Id]; exists {
		return ErrJobAlrExists
	}
	r.jobs[job.Id] = job
	return nil
}

func (r *InMemoryRepository) GetById(id string) (model.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, ok := r.jobs[id]
	if !ok {
		return model.Job{}, ErrJobNotFound
	}
	return job, nil
}

func (r *InMemoryRepository) GetAll() ([]model.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobs := make([]model.Job, 0, len(r.jobs))
	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *InMemoryRepository) GetByStatus(status model.Status) ([]model.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var jobs []model.Job
	for _, job := range r.jobs {
		if job.Status == status {
			jobs = append(jobs, job)
		}
	}
	return jobs, nil
}

func (r *InMemoryRepository) Update(job model.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[job.Id]; !exists {
		return ErrJobNotFound
	}
	r.jobs[job.Id] = job
	return nil
}
