package repository

import (
	"database/sql"
	"fmt"

	"github.com/xyd/web3-learning-tracker/internal/model"
)

type TaskRepo struct{ DB *sql.DB }

func (r *TaskRepo) FindByDay(dayID uint64) ([]model.Task, error) {
	rows, err := r.DB.Query(
		`SELECT id, day_id, content, estimated_hours, resource_urls, sort_order, is_checkpoint, created_at
		FROM tasks WHERE day_id = ? ORDER BY sort_order`, dayID)
	if err != nil {
		return nil, fmt.Errorf("task find by day: %w", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.DayID, &t.Content, &t.EstimatedHours, &t.ResourceURLs, &t.SortOrder, &t.IsCheckpoint, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("task scan: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepo) FindByID(id uint64) (*model.Task, error) {
	t := &model.Task{}
	err := r.DB.QueryRow(
		`SELECT id, day_id, content, estimated_hours, resource_urls, sort_order, is_checkpoint, created_at
		FROM tasks WHERE id = ?`, id,
	).Scan(&t.ID, &t.DayID, &t.Content, &t.EstimatedHours, &t.ResourceURLs, &t.SortOrder, &t.IsCheckpoint, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("task find: %w", err)
	}
	return t, nil
}

func (r *TaskRepo) FindIDsByDay(dayID uint64) ([]uint64, error) {
	rows, err := r.DB.Query(`SELECT id FROM tasks WHERE day_id = ? ORDER BY sort_order`, dayID)
	if err != nil {
		return nil, fmt.Errorf("task ids by day: %w", err)
	}
	defer rows.Close()

	var ids []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
