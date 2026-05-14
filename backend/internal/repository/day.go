package repository

import (
	"database/sql"
	"fmt"

	"github.com/xyd/web3-learning-tracker/internal/model"
)

type DayRepo struct{ DB *sql.DB }

func (r *DayRepo) FindByID(id uint64) (*model.Day, error) {
	d := &model.Day{}
	err := r.DB.QueryRow(
		`SELECT id, week_id, day_number, title, sort_order FROM days WHERE id = ?`, id,
	).Scan(&d.ID, &d.WeekID, &d.DayNumber, &d.Title, &d.SortOrder)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("day find: %w", err)
	}
	return d, nil
}

func (r *DayRepo) FindByWeek(weekID uint64) ([]model.Day, error) {
	rows, err := r.DB.Query(
		`SELECT id, week_id, day_number, title, sort_order FROM days WHERE week_id = ? ORDER BY sort_order`, weekID)
	if err != nil {
		return nil, fmt.Errorf("day find by week: %w", err)
	}
	defer rows.Close()

	var days []model.Day
	for rows.Next() {
		var d model.Day
		if err := rows.Scan(&d.ID, &d.WeekID, &d.DayNumber, &d.Title, &d.SortOrder); err != nil {
			return nil, fmt.Errorf("day scan: %w", err)
		}
		days = append(days, d)
	}
	return days, rows.Err()
}

func (r *DayRepo) FindByWeekWithProgress(weekID, userID uint64) ([]model.Day, error) {
	rows, err := r.DB.Query(`
		SELECT d.id, d.week_id, d.day_number, d.title, d.sort_order,
		  (SELECT COUNT(*) FROM tasks t WHERE t.day_id = d.id AND t.is_checkpoint = 0) AS task_count,
		  (SELECT COUNT(*) FROM tasks t JOIN user_tasks ut ON ut.task_id = t.id
		   WHERE t.day_id = d.id AND t.is_checkpoint = 0
		   AND ut.user_id = ? AND ut.is_completed = 1) AS completed_count
		FROM days d WHERE d.week_id = ? ORDER BY d.sort_order`, userID, weekID)
	if err != nil {
		return nil, fmt.Errorf("day find with progress: %w", err)
	}
	defer rows.Close()

	var days []model.Day
	for rows.Next() {
		var d model.Day
		if err := rows.Scan(&d.ID, &d.WeekID, &d.DayNumber, &d.Title, &d.SortOrder, &d.TaskCount, &d.CompletedCount); err != nil {
			return nil, fmt.Errorf("day scan: %w", err)
		}
		days = append(days, d)
	}
	return days, rows.Err()
}
