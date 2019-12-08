package executors

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// Sequence is an executor for sequential executions
type Sequence struct {
	executor
}

// NewSequence - initializes a sequence executor
func NewSequence(ctx context.Context) *Sequence {
	return &Sequence{
		executor: executor{
			id:  fmt.Sprintf("%s-%s", "sequence", uuid.NewV4().String()),
			ctx: ctx,
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
				log(s.id).Errorf("error while executing: %+v", err)
				exec.OnFailure(err)
				return err
			}
			log(s.id).Infof("completed executing: %+v", exec)
			exec.OnSuccess()
		}
	}
	return nil
}
