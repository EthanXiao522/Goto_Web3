package service

import (
	"database/sql"
	"fmt"
)

type ProgressService struct {
	DB *sql.DB
}

type ProgressOverview struct {
	TotalTasks     int `json:"total_tasks"`
	CompletedTasks int `json:"completed_tasks"`
	TotalPhases    int `json:"total_phases"`
	CompletedPhases int `json:"completed_phases"`
	TotalWeeks     int `json:"total_weeks"`
	CompletedWeeks int `json:"completed_weeks"`
}

type PhaseProgress struct {
	PhaseID       uint64  `json:"phase_id"`
	PhaseNumber   uint8   `json:"phase_number"`
	Title         string  `json:"title"`
	TaskCount     int     `json:"task_count"`
	CompletedCount int    `json:"completed_count"`
	Percentage    float64 `json:"percentage"`
}

type WeekProgress struct {
	WeekID        uint64  `json:"week_id"`
	WeekNumber    uint8   `json:"week_number"`
	Title         string  `json:"title"`
	TaskCount     int     `json:"task_count"`
	CompletedCount int    `json:"completed_count"`
}

type DashboardData struct {
	Overview       ProgressOverview `json:"overview"`
	PhaseProgress  []PhaseProgress  `json:"phase_progress"`
	WeekProgress   []WeekProgress   `json:"week_progress"`
	RecentTasks    []RecentTask     `json:"recent_tasks"`
}

type RecentTask struct {
	TaskID      uint64  `json:"task_id"`
	Content     string  `json:"content"`
	PhaseTitle  string  `json:"phase_title"`
	WeekNumber  uint8   `json:"week_number"`
	IsCompleted bool    `json:"is_completed"`
	CompletedAt *string `json:"completed_at"`
}

func NewProgressService(db *sql.DB) *ProgressService {
	return &ProgressService{DB: db}
}

func (s *ProgressService) GetOverview(userID uint64) (*ProgressOverview, error) {
	o := &ProgressOverview{}

	err := s.DB.QueryRow(`
		SELECT COUNT(*) FROM tasks WHERE is_checkpoint = 0
	`).Scan(&o.TotalTasks)
	if err != nil {
		return nil, fmt.Errorf("progress total tasks: %w", err)
	}

	err = s.DB.QueryRow(`
		SELECT COUNT(*) FROM user_tasks WHERE user_id = ? AND is_completed = 1
	`, userID).Scan(&o.CompletedTasks)
	if err != nil {
		return nil, fmt.Errorf("progress completed tasks: %w", err)
	}

	err = s.DB.QueryRow(`SELECT COUNT(*) FROM phases`).Scan(&o.TotalPhases)
	if err != nil {
		return nil, fmt.Errorf("progress total phases: %w", err)
	}

	err = s.DB.QueryRow(`SELECT COUNT(*) FROM weeks`).Scan(&o.TotalWeeks)
	if err != nil {
		return nil, fmt.Errorf("progress total weeks: %w", err)
	}

	err = s.DB.QueryRow(`
		SELECT COUNT(DISTINCT w.phase_id)
		FROM weeks w
		WHERE NOT EXISTS (
			SELECT 1 FROM days d
			JOIN tasks t ON t.day_id = d.id
			LEFT JOIN user_tasks ut ON ut.task_id = t.id AND ut.user_id = ? AND ut.is_completed = 1
			WHERE d.week_id = w.id AND t.is_checkpoint = 0 AND ut.id IS NULL
		)
	`, userID).Scan(&o.CompletedPhases)
	if err != nil {
		return nil, fmt.Errorf("progress completed phases: %w", err)
	}

	err = s.DB.QueryRow(`
		SELECT COUNT(*) FROM weeks w
		WHERE NOT EXISTS (
			SELECT 1 FROM days d
			JOIN tasks t ON t.day_id = d.id
			LEFT JOIN user_tasks ut ON ut.task_id = t.id AND ut.user_id = ? AND ut.is_completed = 1
			WHERE d.week_id = w.id AND t.is_checkpoint = 0 AND ut.id IS NULL
		)
	`, userID).Scan(&o.CompletedWeeks)
	if err != nil {
		return nil, fmt.Errorf("progress completed weeks: %w", err)
	}

	return o, nil
}

func (s *ProgressService) GetDashboard(userID uint64) (*DashboardData, error) {
	overview, err := s.GetOverview(userID)
	if err != nil {
		return nil, err
	}

	data := &DashboardData{Overview: *overview}

	phaseRows, err := s.DB.Query(`
		SELECT p.id, p.phase_number, p.title,
		  (SELECT COUNT(*) FROM weeks w JOIN days d ON d.week_id = w.id
		   JOIN tasks t ON t.day_id = d.id
		   WHERE w.phase_id = p.id AND t.is_checkpoint = 0) AS task_count,
		  (SELECT COUNT(*) FROM weeks w JOIN days d ON d.week_id = w.id
		   JOIN tasks t ON t.day_id = d.id
		   JOIN user_tasks ut ON ut.task_id = t.id
		   WHERE w.phase_id = p.id AND t.is_checkpoint = 0
		   AND ut.user_id = ? AND ut.is_completed = 1) AS completed_count
		FROM phases p ORDER BY p.sort_order`, userID)
	if err != nil {
		return nil, fmt.Errorf("dashboard phases: %w", err)
	}
	defer phaseRows.Close()

	for phaseRows.Next() {
		var pp PhaseProgress
		if err := phaseRows.Scan(&pp.PhaseID, &pp.PhaseNumber, &pp.Title, &pp.TaskCount, &pp.CompletedCount); err != nil {
			return nil, fmt.Errorf("dashboard phase scan: %w", err)
		}
		if pp.TaskCount > 0 {
			pp.Percentage = float64(pp.CompletedCount) / float64(pp.TaskCount) * 100
		}
		data.PhaseProgress = append(data.PhaseProgress, pp)
	}

	weekRows, err := s.DB.Query(`
		SELECT w.id, w.week_number, w.title,
		  (SELECT COUNT(*) FROM days d JOIN tasks t ON t.day_id = d.id
		   WHERE d.week_id = w.id AND t.is_checkpoint = 0) AS task_count,
		  (SELECT COUNT(*) FROM days d JOIN tasks t ON t.day_id = d.id
		   JOIN user_tasks ut ON ut.task_id = t.id
		   WHERE d.week_id = w.id AND t.is_checkpoint = 0
		   AND ut.user_id = ? AND ut.is_completed = 1) AS completed_count
		FROM weeks w ORDER BY w.sort_order`, userID)
	if err != nil {
		return nil, fmt.Errorf("dashboard weeks: %w", err)
	}
	defer weekRows.Close()

	for weekRows.Next() {
		var wp WeekProgress
		if err := weekRows.Scan(&wp.WeekID, &wp.WeekNumber, &wp.Title, &wp.TaskCount, &wp.CompletedCount); err != nil {
			return nil, fmt.Errorf("dashboard week scan: %w", err)
		}
		data.WeekProgress = append(data.WeekProgress, wp)
	}

	recentRows, err := s.DB.Query(`
		SELECT t.id, t.content, p.title, w.week_number,
		       ut.is_completed, ut.completed_at
		FROM user_tasks ut
		JOIN tasks t ON t.id = ut.task_id
		JOIN days d ON d.id = t.day_id
		JOIN weeks w ON w.id = d.week_id
		JOIN phases p ON p.id = w.phase_id
		WHERE ut.user_id = ?
		ORDER BY ut.updated_at DESC LIMIT 5`, userID)
	if err != nil {
		return nil, fmt.Errorf("dashboard recent: %w", err)
	}
	defer recentRows.Close()

	for recentRows.Next() {
		var rt RecentTask
		var completedAt sql.NullTime
		if err := recentRows.Scan(&rt.TaskID, &rt.Content, &rt.PhaseTitle, &rt.WeekNumber, &rt.IsCompleted, &completedAt); err != nil {
			return nil, fmt.Errorf("dashboard recent scan: %w", err)
		}
		if completedAt.Valid {
			s := completedAt.Time.Format("2006-01-02 15:04")
			rt.CompletedAt = &s
		}
		data.RecentTasks = append(data.RecentTasks, rt)
	}

	return data, nil
}
