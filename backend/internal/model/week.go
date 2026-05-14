package model

type Week struct {
	ID           uint64 `json:"id"`
	PhaseID      uint64 `json:"phase_id"`
	WeekNumber   uint8  `json:"week_number"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Goal         string `json:"goal"`
	Deliverables string `json:"deliverables"`
	SortOrder    int    `json:"sort_order"`
	DayCount     int    `json:"day_count"`
	TaskCount    int    `json:"task_count"`
	CompletedCount int  `json:"completed_count"`
}
