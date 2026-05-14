package repository

import (
	"database/sql"
	"fmt"

	"github.com/xyd/web3-learning-tracker/internal/model"
)

type PhaseRepo struct{ DB *sql.DB }

func (r *PhaseRepo) GetAllWithProgress(userID uint64) ([]model.Phase, error) {
	rows, err := r.DB.Query(`
		SELECT p.id, p.phase_number, p.title, p.subtitle, p.goal, p.deliverables, p.sort_order,
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
		return nil, fmt.Errorf("phase get all: %w", err)
	}
	defer rows.Close()

	var phases []model.Phase
	for rows.Next() {
		var p model.Phase
		if err := rows.Scan(&p.ID, &p.PhaseNumber, &p.Title, &p.Subtitle, &p.Goal, &p.Deliverables, &p.SortOrder, &p.TaskCount, &p.CompletedCount); err != nil {
			return nil, fmt.Errorf("phase scan: %w", err)
		}
		phases = append(phases, p)
	}
	return phases, rows.Err()
}

func (r *PhaseRepo) FindByID(id uint64) (*model.Phase, error) {
	p := &model.Phase{}
	err := r.DB.QueryRow(
		`SELECT id, phase_number, title, subtitle, goal, deliverables, sort_order FROM phases WHERE id = ?`, id,
	).Scan(&p.ID, &p.PhaseNumber, &p.Title, &p.Subtitle, &p.Goal, &p.Deliverables, &p.SortOrder)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("phase find: %w", err)
	}
	return p, nil
}
