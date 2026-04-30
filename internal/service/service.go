package service

import (
	"fmt"
	"time"

	"github.com/fourseasonsm/job_manager/internal/model"
	"github.com/fourseasonsm/job_manager/internal/repository"
	"github.com/google/uuid"
)

type CreateJobInput struct{}

type JobService struct {
	repo  repository.JobRepository
	queue chan<- string
}

func NewJobService(repo repository.JobRepository, queue chan<- string) *JobService {
	return &JobService{
		repo:  repo,
		queue: queue,
	}
}

func (s *JobService) CreateJob(input CreateJobInput) (model.Job, error) {

	job := model.Job{
		Id:        uuid.New().String(),
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(job); err != nil {
		return model.Job{}, fmt.Errorf("failed to save job: %w", err)
	}
	s.queue <- job.Id

	return job, nil
}

func (s *JobService) GetJob(id string) (model.Job, error) {
	job, err := s.repo.GetById(id)
	if err != nil {
		return model.Job{}, fmt.Errorf("failed to get job: %w", err)
	}
	return job, nil
}

func (s *JobService) ListJobs(status string) ([]model.Job, error) {
	if status == "" {
		return s.repo.GetAll()
	}
	return s.repo.GetByStatus(model.Status(status))
}
