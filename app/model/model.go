package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type JobDTO struct {
	Type     string          `json:"type"`
	Payload  json.RawMessage `json:"payload"`
	Priority string          `json:"priority"` // "low", "medium", "high"
}

func (JobDTO) TableName() string {
	return "jobs"
}

type Job struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey"`
	Type         string          `json:"type" gorm:"type:varchar(50);not null"`
	Payload      json.RawMessage `json:"payload" gorm:"type:jsonb;not null"`
	Status       string          `json:"status" gorm:"type:varchar(20);not null;check:status IN ('queued','processing','success','failed')"`
	Priority     string          `json:"priority" gorm:"type:varchar(10);not null;check:priority IN ('low','medium','high')"`
	CreatedAt    time.Time       `json:"created_at" gorm:"not null;autoCreateTime"`
	StartedAt    *time.Time      `json:"started_at,omitempty" gorm:""`
	CompletedAt  *time.Time      `json:"completed_at,omitempty" gorm:""`
	Result       *string         `json:"result,omitempty" gorm:"type:text"`
	ErrorMessage *string         `json:"error_message,omitempty" gorm:"type:text"`
	Retries      int             `json:"retries" gorm:"type:int;default:0"`
}

func (Job) TableName() string {
	return "jobs"
}
