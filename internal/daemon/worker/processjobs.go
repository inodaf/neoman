package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/inodaf/neoman/internal/daemon/domain"
)

type ProcessJobsInput struct {
	Limit int
}

func (w *Worker) ProcessJobs(ctx context.Context, input ProcessJobsInput) error {
	jobs, err := w.jobRepository.GetPending(input.Limit)
	if err != nil {
		slog.Error("failed to get pending jobs", "error", err)
		return err
	}

	for _, job := range jobs {
		w.processJob(ctx, job)
	}

	return nil
}

func (w *Worker) processJob(ctx context.Context, job domain.Job) {
	job.MarkInProgress()
	err := w.jobRepository.Update(&job)
	if err != nil {
		slog.Error("failed to mark job in_progress", "job_id", job.ID, "error", err)
		return
	}

	var processErr error
	switch job.JobType {
	case domain.JobTypeIndexPages:
		processErr = w.processIndexPagesJob(ctx, job)
	default:
		processErr = fmt.Errorf("unknown job type: %s", job.JobType)
	}

	if processErr == nil {
		err := w.jobRepository.Delete(job.ID)
		if err != nil {
			slog.Error("failed to delete completed job", "job_id", job.ID, "error", err)
		}
		slog.Info("job completed successfully", "job_id", job.ID)
		return
	}

	if job.HasMaxRetriesExceeded() {
		job.MarkFailed(processErr.Error())
		slog.Error("job failed after max retries", "job_id", job.ID, "error", processErr)
	} else {
		job.ScheduleRetry()
		slog.Warn("job will be retried", "job_id", job.ID, "attempt", job.AttemptCount, "next_retry_at", job.NextRetryAt, "error", processErr)
	}

	err = w.jobRepository.Update(&job)
	if err != nil {
		slog.Error("failed to update job after failure", "job_id", job.ID, "error", err)
	}
}

func (w *Worker) processIndexPagesJob(ctx context.Context, job domain.Job) error {
	var payload struct {
		Author     string `json:"author"`
		Repository string `json:"repository"`
		Source     string `json:"source"`
	}

	err := json.Unmarshal([]byte(job.Payload), &payload)
	if err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	return w.IndexPages(IndexPagesInput{
		Author:     payload.Author,
		Repository: payload.Repository,
		Source:     domain.RemoteSource(payload.Source),
	})
}
