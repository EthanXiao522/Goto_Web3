package model

type Phase struct {
	ID           uint64  `json:"id"`
	PhaseNumber  uint8   `json:"phase_number"`
	Title        string  `json:"title"`
	Subtitle     string  `json:"subtitle"`
	Goal         string  `json:"goal"`
	Deliverables string  `json:"deliverables"`
	SortOrder    int     `json:"sort_order"`
	WeekCount    int     `json:"week_count,omitempty"`
	TaskCount    int     `json:"task_count,omitempty"`
	CompletedCount int   `json:"completed_count,omitempty"`
}
