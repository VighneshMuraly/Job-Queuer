package apihandler

import (
	"encoding/json"
	"job-queuer/app/model"
	"job-queuer/app/service"
	"net/http"
)

type JobHandler interface {
	ScheduleJob(w http.ResponseWriter, r *http.Request)
	CheckJobStatus(w http.ResponseWriter, r *http.Request)
}

type jobHandler struct {
	service service.JobService
}

func NewJobHandler(s service.JobService) JobHandler {
	return &jobHandler{service: s}
}

func (h *jobHandler) ScheduleJob(w http.ResponseWriter, r *http.Request) {
	var dto model.JobDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	job, err := h.service.ScheduleJob(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (h *jobHandler) CheckJobStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		http.Error(w, "Missing job ID", http.StatusBadRequest)
		return
	}
	job, err := h.service.CheckJobStatus(jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}
