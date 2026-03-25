package domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"time"
)

type JobType string

const (
	JobTypeIndexPages JobType = "index_pages"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusInProgress JobStatus = "in_progress"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID           string
	JobType      JobType
	Payload      string // JSON string
	Status       JobStatus
	LastError    *string
	AttemptCount int
	CreatedAt    time.Time
	NextRetryAt  *time.Time
}

func NewJob(jobType JobType, payload string) (*Job, error) {
	if jobType == "" {
		return nil, errors.New("job_type cannot be empty")
	}
	if payload == "" {
		return nil, errors.New("payload cannot be empty")
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	jobID := hex.EncodeToString(b)

	return &Job{
		ID:        jobID,
		JobType:   jobType,
		Payload:   payload,
		Status:    JobStatusPending,
		CreatedAt: time.Now(),
	}, nil
}

func (j *Job) MarkInProgress() {
	j.Status = JobStatusInProgress
}

func (j *Job) MarkFailed(errorMsg string) {
	j.Status = JobStatusFailed
	j.LastError = &errorMsg
}

func (j *Job) ScheduleRetry() {
	j.AttemptCount++
	backoffSeconds := max(int(math.Pow(2, float64(j.AttemptCount))), 1800)
	nextRetry := time.Now().Add(time.Duration(backoffSeconds) * time.Second)
	j.NextRetryAt = &nextRetry
}

func (j *Job) HasMaxRetriesExceeded() bool {
	return j.AttemptCount >= 3
}
