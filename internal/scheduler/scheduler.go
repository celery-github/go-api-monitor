package scheduler

import (
	"context"
	"time"
)

type Runner func(ctx context.Context) error

type Scheduler struct {
	interval time.Duration
	run      Runner
}

func New(interval time.Duration, run Runner) *Scheduler {
	if interval <= 0 {
		interval = 60 * time.Second
	}
	return &Scheduler{interval: interval, run: run}
}

func (s *Scheduler) Run(ctx context.Context) error {
	// run immediately once
	if err := s.run(ctx); err != nil {
		return err
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.run(ctx); err != nil {
				return err
			}
		}
	}
}
