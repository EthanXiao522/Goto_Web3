package importer

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

func Seed(db *sql.DB, data *ParsedData) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("seeder: begin tx: %w", err)
	}
	defer tx.Rollback()

	phaseStmt, _ := tx.Prepare(`INSERT IGNORE INTO phases (phase_number, title, subtitle, goal, deliverables, sort_order) VALUES (?, ?, ?, ?, ?, ?)`)
	weekStmt, _ := tx.Prepare(`INSERT IGNORE INTO weeks (phase_id, week_number, title, subtitle, goal, deliverables, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	dayStmt, _ := tx.Prepare(`INSERT IGNORE INTO days (week_id, day_number, title, sort_order) VALUES (?, ?, ?, ?)`)
	taskStmt, _ := tx.Prepare(`INSERT IGNORE INTO tasks (day_id, content, estimated_hours, resource_urls, sort_order, is_checkpoint) VALUES (?, ?, ?, ?, ?, ?)`)

	var phaseCount, weekCount, dayCount, taskCount int

	for _, phase := range data.Phases {
		res, err := phaseStmt.Exec(phase.PhaseNumber, phase.Title, phase.Subtitle, phase.Goal, phase.Deliverables, phase.SortOrder)
		if err != nil {
			return fmt.Errorf("seeder: phase %d: %w", phase.PhaseNumber, err)
		}
		phaseID, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("seeder: phase %d last insert: %w", phase.PhaseNumber, err)
		}
		if phaseID > 0 {
			phaseCount++
		}

		for _, week := range phase.Weeks {
			res, err := weekStmt.Exec(phaseID, week.WeekNumber, week.Title, week.Subtitle, week.Goal, week.Deliverables, week.SortOrder)
			if err != nil {
				return fmt.Errorf("seeder: week %d: %w", week.WeekNumber, err)
			}
			weekID, err := res.LastInsertId()
			if err != nil {
				return fmt.Errorf("seeder: week %d last insert: %w", week.WeekNumber, err)
			}
			if weekID > 0 {
				weekCount++
			}

			for _, day := range week.Days {
				res, err := dayStmt.Exec(weekID, day.DayNumber, day.Title, day.SortOrder)
				if err != nil {
					return fmt.Errorf("seeder: day %d: %w", day.DayNumber, err)
				}
				dayID, err := res.LastInsertId()
				if err != nil {
					return fmt.Errorf("seeder: day %d last insert: %w", day.DayNumber, err)
				}
				if dayID > 0 {
					dayCount++
				}

				for _, task := range day.Tasks {
					urlsJSON := encodeURLs(task.ResourceURLs)
					res, err := taskStmt.Exec(dayID, task.Content, task.EstimatedHours, urlsJSON, task.SortOrder, task.IsCheckpoint)
					if err != nil {
						return fmt.Errorf("seeder: task %q: %w", truncate(task.Content, 50), err)
					}
					if id, _ := res.LastInsertId(); id > 0 {
						taskCount++
					}
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("seeder: commit: %w", err)
	}

	fmt.Printf("Imported: %d phases, %d weeks, %d days, %d tasks\n", phaseCount, weekCount, dayCount, taskCount)
	return nil
}

func encodeURLs(urls []string) string {
	if len(urls) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(urls)
	return string(b)
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n]) + "..."
	}
	return s
}
