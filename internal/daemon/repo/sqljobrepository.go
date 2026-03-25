package repo

import (
	"database/sql"
	"fmt"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type jobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Save(job *domain.Job) error {
	query := `
		INSERT INTO jobs (id, job_type, payload, status, attempt_count, created_at, next_retry_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, job.ID, job.JobType, job.Payload, job.Status, job.AttemptCount, job.CreatedAt, job.NextRetryAt)
	return err
}

func (r *jobRepository) Update(job *domain.Job) error {
	query := `
		UPDATE jobs
		SET status = ?, last_error = ?, attempt_count = ?, next_retry_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, job.Status, job.LastError, job.AttemptCount, job.NextRetryAt, job.ID)
	return err
}

func (r *jobRepository) GetPending(limit int) ([]domain.Job, error) {
	query := `
		SELECT id, job_type, payload, status, last_error, attempt_count, created_at, next_retry_at
		FROM jobs
		WHERE status = ? OR (status = ? AND next_retry_at <= datetime('now'))
		ORDER BY created_at ASC
		LIMIT ?
	`

	rows, err := r.db.Query(query, domain.JobStatusPending, domain.JobStatusFailed, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var jobs []domain.Job
	for rows.Next() {
		var job domain.Job
		err := rows.Scan(
			&job.ID,
			&job.JobType,
			&job.Payload,
			&job.Status,
			&job.LastError,
			&job.AttemptCount,
			&job.CreatedAt,
			&job.NextRetryAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}

		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

func (r *jobRepository) GetAll() ([]domain.Job, error) {
	query := `
		SELECT id, job_type, payload, status, last_error, attempt_count, created_at, next_retry_at
		FROM jobs
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []domain.Job
	for rows.Next() {
		var job domain.Job
		err := rows.Scan(
			&job.ID,
			&job.JobType,
			&job.Payload,
			&job.Status,
			&job.LastError,
			&job.AttemptCount,
			&job.CreatedAt,
			&job.NextRetryAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}

		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

func (r *jobRepository) Delete(jobID string) error {
	query := `DELETE FROM jobs WHERE id = ?`
	_, err := r.db.Exec(query, jobID)
	return err
}
