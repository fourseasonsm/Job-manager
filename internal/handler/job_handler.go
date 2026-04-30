package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fourseasonsm/job_manager/internal/repository"
	"github.com/fourseasonsm/job_manager/internal/response"
	"github.com/fourseasonsm/job_manager/internal/service"
)

type JobHandler struct {
	service *service.JobService
}

func NewJobHandler(svc *service.JobService) *JobHandler {
	return &JobHandler{service: svc}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var input service.CreateJobInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	job, err := h.service.CreateJob(input)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, job)
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, err := h.service.GetJob(id)
	if err != nil {
		if errors.Is(err, repository.ErrJobNotFound) {
			response.Error(w, http.StatusNotFound, "job not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, job)
}

func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	jobs, err := h.service.ListJobs(status)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, jobs)
}
