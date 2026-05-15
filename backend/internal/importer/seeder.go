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

	var phaseCount, weekCount, dayCount, taskCount int

	for _, phase := range data.Phases {
		res, _ := tx.Exec(`INSERT IGNORE INTO phases (phase_number, title, subtitle, goal, deliverables, sort_order) VALUES (?, ?, ?, ?, ?, ?)`,
			phase.PhaseNumber, phase.Title, phase.Subtitle, phase.Goal, phase.Deliverables, phase.SortOrder)
		phaseID := lastInsertID(res)
		if phaseID == 0 {
			tx.QueryRow("SELECT id FROM phases WHERE phase_number = ?", phase.PhaseNumber).Scan(&phaseID)
		} else {
			phaseCount++
		}

		for _, week := range phase.Weeks {
			res, _ := tx.Exec(`INSERT IGNORE INTO weeks (phase_id, week_number, title, subtitle, goal, deliverables, sort_order) VALUES (?, ?, ?, ?, ?, ?, ?)`,
				phaseID, week.WeekNumber, week.Title, week.Subtitle, week.Goal, week.Deliverables, week.SortOrder)
			weekID := lastInsertID(res)
			if weekID == 0 {
				tx.QueryRow("SELECT id FROM weeks WHERE week_number = ?", week.WeekNumber).Scan(&weekID)
			} else {
				weekCount++
			}

			for _, day := range week.Days {
				res, _ := tx.Exec(`INSERT IGNORE INTO days (week_id, day_number, title, sort_order) VALUES (?, ?, ?, ?)`,
					weekID, day.DayNumber, day.Title, day.SortOrder)
				dayID := lastInsertID(res)
				if dayID == 0 {
					tx.QueryRow("SELECT id FROM days WHERE week_id = ? AND day_number = ?", weekID, day.DayNumber).Scan(&dayID)
				} else {
					dayCount++
				}

				for _, task := range day.Tasks {
					urlsJSON := encodeURLs(task.ResourceURLs)
					res, _ := tx.Exec(`INSERT IGNORE INTO tasks (day_id, content, estimated_hours, resource_urls, sort_order, is_checkpoint) VALUES (?, ?, ?, ?, ?, ?)`,
						dayID, task.Content, task.EstimatedHours, urlsJSON, task.SortOrder, task.IsCheckpoint)
					if lastInsertID(res) > 0 {
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

func lastInsertID(res sql.Result) uint64 {
	id, _ := res.LastInsertId()
	return uint64(id)
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
