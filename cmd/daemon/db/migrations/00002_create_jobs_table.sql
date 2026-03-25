-- +goose Up
CREATE TABLE IF NOT EXISTS jobs (
    id TEXT PRIMARY KEY,
    job_type TEXT NOT NULL,
    payload TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    last_error TEXT,
    attempt_count INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    next_retry_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_jobs_status_next_retry 
    ON jobs(status, next_retry_at) 
    WHERE status IN ('pending', 'failed');

-- +goose Down
DROP INDEX IF EXISTS idx_jobs_status_next_retry;
DROP TABLE IF EXISTS jobs;
