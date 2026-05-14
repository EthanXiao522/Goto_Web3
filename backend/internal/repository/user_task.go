package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/xyd/web3-learning-tracker/internal/model"
)

type UserTaskRepo struct{ DB *sql.DB }

func (r *UserTaskRepo) LazyCreate(userID, taskID uint64) (*model.UserTask, error) {
	_, err := r.DB.Exec(
		`INSERT IGNORE INTO user_tasks (user_id, task_id) VALUES (?, ?)`, userID, taskID)
	if err != nil {
		return nil, fmt.Errorf("user_task lazy create: %w", err)
	}
	return r.FindByUserAndTask(userID, taskID)
}

func (r *UserTaskRepo) FindByUserAndTask(userID, taskID uint64) (*model.UserTask, error) {
	ut := &model.UserTask{}
	err := r.DB.QueryRow(
		`SELECT id, user_id, task_id, is_completed, completed_at,
		        learning_links, implementation_plan, implementation_code, experience_summary,
		        created_at, updated_at
		FROM user_tasks WHERE user_id = ? AND task_id = ?`, userID, taskID,
	).Scan(&ut.ID, &ut.UserID, &ut.TaskID, &ut.IsCompleted, &ut.CompletedAt,
		&ut.LearningLinks, &ut.ImplementationPlan, &ut.ImplementationCode, &ut.ExperienceSummary,
		&ut.CreatedAt, &ut.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user_task find: %w", err)
	}
	return ut, nil
}

func (r *UserTaskRepo) UpdateComplete(userID, taskID uint64, completed bool) error {
	var completedAt interface{}
	if completed {
		now := time.Now()
		completedAt = &now
	} else {
		completedAt = nil
	}
	_, err := r.DB.Exec(
		`UPDATE user_tasks SET is_completed = ?, completed_at = ? WHERE user_id = ? AND task_id = ?`,
		completed, completedAt, userID, taskID)
	if err != nil {
		return fmt.Errorf("user_task update complete: %w", err)
	}
	return nil
}

func (r *UserTaskRepo) UpdateFields(userID, taskID uint64, fields map[string]string) error {
	query := `UPDATE user_tasks SET `
	var args []interface{}
	for col, val := range fields {
		if val == "" {
			continue
		}
		query += col + " = ?, "
		args = append(args, val)
	}
	if len(args) == 0 {
		return nil
	}
	query = query[:len(query)-2] + ` WHERE user_id = ? AND task_id = ?`
	args = append(args, userID, taskID)
	_, err := r.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("user_task update fields: %w", err)
	}
	return nil
}

func (r *UserTaskRepo) FindByUserAndTaskIDs(userID uint64, taskIDs []uint64) (map[uint64]*model.UserTask, error) {
	if len(taskIDs) == 0 {
		return map[uint64]*model.UserTask{}, nil
	}
	query := `SELECT id, user_id, task_id, is_completed, completed_at,
		learning_links, implementation_plan, implementation_code, experience_summary,
		created_at, updated_at FROM user_tasks WHERE user_id = ? AND task_id IN (`
	args := []interface{}{userID}
	for i, id := range taskIDs {
		if i > 0 {
			query += ", "
		}
		query += "?"
		args = append(args, id)
	}
	query += ")"

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("user_task batch find: %w", err)
	}
	defer rows.Close()

	result := make(map[uint64]*model.UserTask)
	for rows.Next() {
		ut := &model.UserTask{}
		if err := rows.Scan(&ut.ID, &ut.UserID, &ut.TaskID, &ut.IsCompleted, &ut.CompletedAt,
			&ut.LearningLinks, &ut.ImplementationPlan, &ut.ImplementationCode, &ut.ExperienceSummary,
			&ut.CreatedAt, &ut.UpdatedAt); err != nil {
			return nil, fmt.Errorf("user_task batch scan: %w", err)
		}
		result[ut.TaskID] = ut
	}
	return result, rows.Err()
}
