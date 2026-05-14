package database

import "fmt"

func Migrate() error {
	ddls := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(64) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS phases (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			phase_number TINYINT UNSIGNED NOT NULL,
			title VARCHAR(255) NOT NULL,
			subtitle VARCHAR(255) DEFAULT '',
			goal TEXT,
			deliverables TEXT,
			sort_order INT NOT NULL DEFAULT 0,
			INDEX idx_phases_sort (sort_order, phase_number)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS weeks (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			phase_id BIGINT UNSIGNED NOT NULL,
			week_number TINYINT UNSIGNED NOT NULL,
			title VARCHAR(255) NOT NULL,
			subtitle VARCHAR(255) DEFAULT '',
			goal TEXT,
			deliverables TEXT,
			sort_order INT NOT NULL DEFAULT 0,
			UNIQUE INDEX idx_weeks_number (week_number),
			INDEX idx_weeks_phase (phase_id, sort_order),
			FOREIGN KEY (phase_id) REFERENCES phases(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS days (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			week_id BIGINT UNSIGNED NOT NULL,
			day_number TINYINT UNSIGNED NOT NULL,
			title VARCHAR(255) NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			INDEX idx_days_week (week_id, sort_order),
			FOREIGN KEY (week_id) REFERENCES weeks(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS tasks (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			day_id BIGINT UNSIGNED NOT NULL,
			content TEXT NOT NULL,
			estimated_hours DECIMAL(4,1) DEFAULT 0,
			resource_urls JSON,
			sort_order INT NOT NULL DEFAULT 0,
			is_checkpoint TINYINT(1) DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_tasks_day (day_id, sort_order),
			FOREIGN KEY (day_id) REFERENCES days(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS user_tasks (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			user_id BIGINT UNSIGNED NOT NULL,
			task_id BIGINT UNSIGNED NOT NULL,
			is_completed TINYINT(1) DEFAULT 0,
			completed_at TIMESTAMP NULL,
			learning_links TEXT,
			implementation_plan TEXT,
			implementation_code TEXT,
			experience_summary TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uq_user_task (user_id, task_id),
			INDEX idx_user_tasks_user (user_id, is_completed),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	}

	for _, ddl := range ddls {
		if _, err := DB.Exec(ddl); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}
