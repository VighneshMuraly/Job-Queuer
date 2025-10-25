package repo

import (
	"job-queuer/app/model"

	"gorm.io/gorm"
)

type Repo interface {
	SaveJob(job *model.Job) error
	GetJobStatus(jobID string) (*model.Job, error)
    GetQueuedJobs() ([]model.Job, error)
    UpdateJobStatus(job *model.Job) error
}

type jobRepo struct {
	db *gorm.DB
}

func NewJobRepo(db *gorm.DB) Repo {
	return &jobRepo{db: db}
}

func (r *jobRepo) SaveJob(job *model.Job) error {
	return r.db.Create(job).Error
}

func (r *jobRepo) GetJobStatus(jobID string) (*model.Job, error) {
	var job model.Job
	if err := r.db.First(&job, "id = ?", jobID).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepo) GetQueuedJobs() ([]model.Job, error) {
	var jobs []model.Job
	if err := r.db.Where("status = ?", "queued").Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *jobRepo) UpdateJobStatus(job *model.Job) error {
	return r.db.Save(job).Error
}
