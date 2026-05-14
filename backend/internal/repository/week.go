package repository

import (
	"database/sql"
	"fmt"

	"github.com/xyd/web3-learning-tracker/internal/model"
)

type WeekRepo struct{ DB *sql.DB }

func (r *WeekRepo) FindByPhase(phaseID uint64) ([]model.Week, error) {
	rows, err := r.DB.Query(
		`SELECT id, phase_id, week_number, title, subtitle, goal, deliverables, sort_order
		FROM weeks WHERE phase_id = ? ORDER BY sort_order`, phaseID)
	if err != nil {
		return nil, fmt.Errorf("week find by phase: %w", err)
	}
	defer rows.Close()

	var weeks []model.Week
	for rows.Next() {
		var w model.Week
		if err := rows.Scan(&w.ID, &w.PhaseID, &w.WeekNumber, &w.Title, &w.Subtitle, &w.Goal, &w.Deliverables, &w.SortOrder); err != nil {
			return nil, fmt.Errorf("week scan: %w", err)
		}
		weeks = append(weeks, w)
	}
	return weeks, rows.Err()
}

func (r *WeekRepo) FindByID(id uint64) (*model.Week, error) {
	w := &model.Week{}
	err := r.DB.QueryRow(
		`SELECT id, phase_id, week_number, title, subtitle, goal, deliverables, sort_order
		FROM weeks WHERE id = ?`, id,
	).Scan(&w.ID, &w.PhaseID, &w.WeekNumber, &w.Title, &w.Subtitle, &w.Goal, &w.Deliverables, &w.SortOrder)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("week find: %w", err)
	}
	return w, nil
}

func (r *WeekRepo) FindByPhaseWithProgress(phaseID, userID uint64) ([]model.Week, error) {
	rows, err := r.DB.Query(`
		SELECT w.id, w.phase_id, w.week_number, w.title, w.subtitle, w.goal, w.deliverables, w.sort_order,
		  (SELECT COUNT(*) FROM days d JOIN tasks t ON t.day_id = d.id
		   WHERE d.week_id = w.id AND t.is_checkpoint = 0) AS task_count,
		  (SELECT COUNT(*) FROM days d JOIN tasks t ON t.day_id = d.id
		   JOIN user_tasks ut ON ut.task_id = t.id
		   WHERE d.week_id = w.id AND t.is_checkpoint = 0
		   AND ut.user_id = ? AND ut.is_completed = 1) AS completed_count
		FROM weeks w WHERE w.phase_id = ? ORDER BY w.sort_order`, userID, phaseID)
	if err != nil {
		return nil, fmt.Errorf("week find with progress: %w", err)
	}
	defer rows.Close()

	var weeks []model.Week
	for rows.Next() {
		var w model.Week
		if err := rows.Scan(&w.ID, &w.PhaseID, &w.WeekNumber, &w.Title, &w.Subtitle, &w.Goal, &w.Deliverables, &w.SortOrder, &w.TaskCount, &w.CompletedCount); err != nil {
			return nil, fmt.Errorf("week scan: %w", err)
		}
		weeks = append(weeks, w)
	}
	return weeks, rows.Err()
}
