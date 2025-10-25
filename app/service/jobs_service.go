package service

import (
	"context"
	"errors"
	"job-queuer/app/model"
	repo "job-queuer/app/repository"
	"sort"
	"time"

	"math/rand"

	"github.com/google/uuid"
)

type JobService interface {
	ScheduleJob(dto model.JobDTO) (*model.Job, error)
	CheckJobStatus(jobID string) (*model.Job, error)
	GetQueuedJobs() ([]model.Job, error)
	ProcessNextJob(jobs []model.Job) (*model.Job, error)
}

type jobService struct {
	repo repo.Repo
}

func NewJobService(r repo.Repo) JobService {
	return &jobService{repo: r}
}

func (s *jobService) ScheduleJob(dto model.JobDTO) (*model.Job, error) {
	// Validate
	if dto.Type == "" || dto.Priority == "" || dto.Payload == nil {
		return nil, errors.New("missing required fields")
	}
	if dto.Priority != "low" && dto.Priority != "medium" && dto.Priority != "high" {
		return nil, errors.New("invalid priority")
	}
	job := &model.Job{
		ID:        uuid.New(),
		Type:      dto.Type,
		Payload:   dto.Payload,
		Status:    "queued",
		Priority:  dto.Priority,
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveJob(job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *jobService) CheckJobStatus(jobID string) (*model.Job, error) {
	return s.repo.GetJobStatus(jobID)
}

func (s *jobService) GetQueuedJobs() ([]model.Job, error) {
	return s.repo.GetQueuedJobs()
}

func (s *jobService) ProcessNextJob(jobs []model.Job) (*model.Job, error) {
	if len(jobs) == 0 {
		return nil, nil
	}

	// Filter out jobs that reached max retries
	validJobs := []model.Job{}
	for _, job := range jobs {
		if job.Retries < 3 {
			validJobs = append(validJobs, job)
		}
	}
	if len(validJobs) == 0 {
		return nil, nil
	}

	// Sort jobs by last update (oldest first)
	sort.SliceStable(validJobs, func(i, j int) bool {
		getTime := func(job model.Job) time.Time {
			if job.CompletedAt != nil {
				return *job.CompletedAt
			}
			if job.StartedAt != nil {
				return *job.StartedAt
			}
			return job.CreatedAt
		}
		return getTime(validJobs[i]).Before(getTime(validJobs[j]))
	})

	// Assign time factor based on hierarchy (older = higher factor)
	timeFactors := make(map[string]float64) // UUID string -> factor
	for idx, job := range validJobs {
		timeFactors[job.ID.String()] = float64(len(validJobs) - idx) // oldest = highest
	}

	// Priority factor: high=3, medium=2, low=1
	priorityFactor := func(job model.Job) float64 {
		switch job.Priority {
		case "high":
			return 3
		case "medium":
			return 2
		case "low":
			return 1
		default:
			return 1
		}
	}

	// Select the job with the highest combined factor
	var selected model.Job
	maxFactor := -1.0
	for _, job := range validJobs {
		factor := priorityFactor(job) + timeFactors[job.ID.String()]
		if factor > maxFactor {
			maxFactor = factor
			selected = job
		}
	}

	// Process only this job
	results := make(chan result, 1)
	ctx := context.Background()
	go s.processJob(ctx, selected, results)
	res := <-results
	return res.job, res.err
}

type result struct {
	job *model.Job
	err error
}

func (s *jobService) processJob(ctx context.Context, job model.Job, results chan result) {

	success := rand.Intn(100) >= 20
	now := time.Now()
	job.StartedAt = &now
	time.Sleep(5 * time.Second)
	now = time.Now()
	job.CompletedAt = &now
	if success {
		job.Status = "success"
		resultMsg := "Job completed successfully"
		job.Result = &resultMsg
		job.ErrorMessage = nil
	} else {
		job.Status = "failed"
		job.Result = nil
		errMsg := "Job failed randomly"
		job.ErrorMessage = &errMsg
		job.Retries++
		// If low fails, increase priority for next cycle
		switch job.Priority {
		case "low":
			job.Priority = "medium"
		case "medium":
			job.Priority = "high"
		}
		if job.Retries >= 3 {
			job.Status = "terminated"
			termMsg := "Job terminated after 3 retries"
			job.ErrorMessage = &termMsg
		}
	}
	err := s.repo.UpdateJobStatus(&job)
	results <- struct {
		job *model.Job
		err error
	}{job: &job, err: err}
}
