package service

import (
	"encoding/json"
	"testing"
	"time"

	"job-queuer/app/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SaveJob(job *model.Job) error {
	args := m.Called(job)
	return args.Error(0)
}
func (m *MockRepo) GetJobStatus(jobID string) (*model.Job, error) {
	args := m.Called(jobID)
	return args.Get(0).(*model.Job), args.Error(1)
}
func (m *MockRepo) GetQueuedJobs() ([]model.Job, error) {
	args := m.Called()
	return args.Get(0).([]model.Job), args.Error(1)
}
func (m *MockRepo) UpdateJobStatus(job *model.Job) error {
	args := m.Called(job)
	return args.Error(0)
}

func TestScheduleJob(t *testing.T) {
	repo := new(MockRepo)
	service := NewJobService(repo)
	dto := model.JobDTO{
		Type:     "email",
		Payload:  json.RawMessage(`{"to":"user@example.com"}`),
		Priority: "high",
	}
	repo.On("SaveJob", mock.AnythingOfType("*model.Job")).Return(nil)
	job, err := service.ScheduleJob(dto)
	assert.NoError(t, err)
	assert.Equal(t, "email", job.Type)
	assert.Equal(t, "high", job.Priority)
	assert.Equal(t, "queued", job.Status)
	repo.AssertExpectations(t)
}

func TestScheduleJob_Invalid(t *testing.T) {
	repo := new(MockRepo)
	service := NewJobService(repo)
	dto := model.JobDTO{Type: "", Payload: nil, Priority: ""}
	job, err := service.ScheduleJob(dto)
	assert.Error(t, err)
	assert.Nil(t, job)
}

func TestCheckJobStatus(t *testing.T) {
	repo := new(MockRepo)
	service := NewJobService(repo)
	jobID := uuid.New().String()
	job := &model.Job{ID: uuid.MustParse(jobID), Status: "queued"}
	repo.On("GetJobStatus", jobID).Return(job, nil)
	res, err := service.CheckJobStatus(jobID)
	assert.NoError(t, err)
	assert.Equal(t, jobID, res.ID.String())
	repo.AssertExpectations(t)
}

func TestGetQueuedJobs(t *testing.T) {
	repo := new(MockRepo)
	service := NewJobService(repo)
	jobs := []model.Job{{ID: uuid.New(), Status: "queued"}}
	repo.On("GetQueuedJobs").Return(jobs, nil)
	res, err := service.GetQueuedJobs()
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	repo.AssertExpectations(t)
}

func TestProcessNextJob_Success(t *testing.T) {
	repo := new(MockRepo)
	service := NewJobService(repo)
	job := model.Job{
		ID:        uuid.New(),
		Type:      "email",
		Priority:  "high",
		Status:    "queued",
		CreatedAt: time.Now().Add(-10 * time.Minute),
	}
	jobs := []model.Job{job}
	repo.On("UpdateJobStatus", mock.AnythingOfType("*model.Job")).Return(nil)
	res, err := service.ProcessNextJob(jobs)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	repo.AssertExpectations(t)
}

func TestProcessNextJob_MaxRetries(t *testing.T) {
	repo := new(MockRepo)
	service := NewJobService(repo)
	job := model.Job{
		ID:        uuid.New(),
		Type:      "email",
		Priority:  "high",
		Status:    "queued",
		Retries:   3,
		CreatedAt: time.Now().Add(-10 * time.Minute),
	}
	jobs := []model.Job{job}
	res, err := service.ProcessNextJob(jobs)
	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestProcessNextJob_Empty(t *testing.T) {
	repo := new(MockRepo)
	service := NewJobService(repo)
	jobs := []model.Job{}
	res, err := service.ProcessNextJob(jobs)
	assert.NoError(t, err)
	assert.Nil(t, res)
}
