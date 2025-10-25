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
	ProcessTopJob(jobs []model.Job) (*model.Job, error)
	ProcessJobsWithQuota(jobs []model.Job) ([]*model.Job, error)
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

func (s *jobService) ProcessTopJob(jobs []model.Job) (*model.Job, error) {
	if len(jobs) == 0 {
		return nil, nil
	}
	job := jobs[0]
	// Emulate job running
	time.Sleep(2 * time.Second)
	success := rand.Intn(100) >= 20 // 80% success
	now := time.Now()
	job.StartedAt = &now
	job.CompletedAt = &now
	if success {
		job.Status = "success"
		result := "Job completed successfully"
		job.Result = &result
		job.ErrorMessage = nil
	} else {
		job.Status = "failed"
		job.Result = nil
		errMsg := "Job failed randomly"
		job.ErrorMessage = &errMsg
		job.Retries++
		// If low fails, increase priority for next cycle
		if job.Priority == "low" {
			job.Priority = "medium"
		} else if job.Priority == "medium" {
			job.Priority = "high"
		}
		// If retries >= 3, terminate job
		if job.Retries >= 3 {
			job.Status = "terminated"
			termMsg := "Job terminated after 3 retries"
			job.ErrorMessage = &termMsg
		}
	}
	if err := s.repo.UpdateJobStatus(&job); err != nil {
		return nil, err
	}
	return &job, nil
}

// ProcessJobsWithQuota processes jobs based on priority quotas
func (s *jobService) ProcessJobsWithQuota(jobs []model.Job) ([]*model.Job, error) {
	quota := map[string]int{"high": 3, "medium": 2, "low": 1}
	processed := []*model.Job{}
	// used := map[string]int{"high": 0, "medium": 0, "low": 0}

	// Remove jobs with retries >= 3
	filtered := []model.Job{}
	for _, job := range jobs {
		if job.Retries < 3 {
			filtered = append(filtered, job)
		}
	}

	// Split jobs by priority
	highJobs, medJobs, lowJobs := []model.Job{}, []model.Job{}, []model.Job{}
	for _, job := range filtered {
		switch job.Priority {
		case "high":
			highJobs = append(highJobs, job)
		case "medium":
			medJobs = append(medJobs, job)
		case "low":
			lowJobs = append(lowJobs, job)
		}
	}

	// Sort each by oldest updated timestamp (CompletedAt, StartedAt, CreatedAt)
	sortJobs := func(jobs []model.Job) {
		sort.SliceStable(jobs, func(i, j int) bool {
			getTime := func(job model.Job) time.Time {
				if job.CompletedAt != nil {
					return *job.CompletedAt
				}
				if job.StartedAt != nil {
					return *job.StartedAt
				}
				return job.CreatedAt
			}
			return getTime(jobs[i]).Before(getTime(jobs[j]))
		})
	}
	sortJobs(highJobs)
	sortJobs(medJobs)
	sortJobs(lowJobs)

	// Allocate jobs by quota
	toProcess := []model.Job{}
	toProcess = append(toProcess, highJobs[:min(quota["high"], len(highJobs))]...)
	toProcess = append(toProcess, medJobs[:min(quota["medium"], len(medJobs))]...)
	toProcess = append(toProcess, lowJobs[:min(quota["low"], len(lowJobs))]...)

	results := make(chan result, len(toProcess))
	ctx := context.Background()
	for _, job := range toProcess {
		go s.processJob(ctx, job, results)
	}
	for i := 0; i < len(toProcess); i++ {
		res := <-results
		if res.err != nil {
			return processed, res.err
		}
		processed = append(processed, res.job)
	}
	return processed, nil
}

type result struct {
	job *model.Job
	err error
}

func (s *jobService) processJob(ctx context.Context, job model.Job, results chan result) {

	success := rand.Intn(100) >= 20
	now := time.Now()
	job.StartedAt = &now
	time.Sleep(2 * time.Second)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
