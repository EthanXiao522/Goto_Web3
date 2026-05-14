package model

import "time"

type Task struct {
	ID             uint64    `json:"id"`
	DayID          uint64    `json:"day_id"`
	Content        string    `json:"content"`
	EstimatedHours float64   `json:"estimated_hours"`
	ResourceURLs   string    `json:"resource_urls"`
	SortOrder      int       `json:"sort_order"`
	IsCheckpoint   bool      `json:"is_checkpoint"`
	CreatedAt      time.Time `json:"created_at"`
}
