package model

type Day struct {
	ID             uint64 `json:"id"`
	WeekID         uint64 `json:"week_id"`
	DayNumber      uint8  `json:"day_number"`
	Title          string `json:"title"`
	SortOrder      int    `json:"sort_order"`
	TaskCount      int    `json:"task_count,omitempty"`
	CompletedCount int    `json:"completed_count,omitempty"`
}
