package model

import "time"

type UserTask struct {
	ID                 uint64     `json:"id"`
	UserID             uint64     `json:"user_id"`
	TaskID             uint64     `json:"task_id"`
	IsCompleted        bool       `json:"is_completed"`
	CompletedAt        *time.Time `json:"completed_at"`
	LearningLinks      string     `json:"learning_links"`
	ImplementationPlan string     `json:"implementation_plan"`
	ImplementationCode string     `json:"implementation_code"`
	ExperienceSummary  string     `json:"experience_summary"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
