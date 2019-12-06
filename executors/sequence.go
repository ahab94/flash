package executors

import (
	"context"
)

// Sequence is an executor for sequential executions
type Sequence struct {
	executor
}

// NewSequence - initializes a sequence executor
func NewSequence(ctx context.Context) *Sequence {
	return &Sequence{
		executor: executor{
			name: "sequence",
			ctx:  ctx,
		},
	}
}

// Execute - executes all executables sequentially
func (s *Sequence) Execute() error {
	if err := s.executor.Execute(); err != nil {
		return err
	}

	for _, exec := range s.executables {
		if !exec.IsCompleted() {
			if err := exec.Execute(); err != nil {
				log(s.ctx, s.name).Errorf("error encountered: %+v", err)
				exec.OnFailure(err)
				return err
			}
			exec.OnSuccess()
		}
	}
	return nil
}
