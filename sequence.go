package flash

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// Sequence is an executor for sequential executions
type Sequence struct {
	executor
	successHandler func()
	failHandler    func(err error)
}

type SequenceOption func(*Sequence)

// NewSequence - initializes a sequence executor
func NewSequence(ctx context.Context, opts ...SequenceOption) *Sequence {
	seq := &Sequence{
		executor: executor{
			id:  fmt.Sprintf("%s-%s", "sequence", uuid.NewV4().String()),
			ctx: ctx,
		},
	}

	for _, opt := range opts {
		opt(seq)
	}

	return seq
}

// SequenceFailHandler - inits fail handler
func SequenceFailHandler(fail func(err error)) SequenceOption {
	return func(s *Sequence) {
		s.failHandler = fail
	}
}

// SequenceSuccessHandler - inits success handler
func SequenceSuccessHandler(success func()) SequenceOption {
	return func(s *Sequence) {
		s.successHandler = success
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

// OnSuccess - handles completion callback
func (s *Sequence) OnSuccess() {
	s.executor.OnSuccess()
	if s.successHandler != nil {
		s.successHandler()
	}
}

// OnFailure - handles failure callback
func (s *Sequence) OnFailure(err error) {
	s.executor.OnFailure(err)
	if s.failHandler != nil {
		s.failHandler(err)
	}
}
