package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/inodaf/neoman/internal/daemon/worker"
)

func NewScheduler(worker *worker.Worker) *Scheduler {
	return &Scheduler{
		worker: worker,
		ticker: time.NewTicker(5 * time.Second),
	}
}

type Scheduler struct {
	worker *worker.Worker
	ticker *time.Ticker
}

func (s *Scheduler) Run(ctx context.Context) error {
	defer s.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.ticker.C:
			err := s.worker.ProcessJobs(ctx, worker.ProcessJobsInput{Limit: 4})
			if err != nil {
				slog.Error("scheduler error processing jobs", "error", err)
			}
		}
	}
}
