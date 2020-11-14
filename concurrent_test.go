package flash

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/ahab94/engine"
)

func TestConcurrent_Execute(t *testing.T) {
	e := engine.NewEngine(context.TODO())
	e.Start(10)

	t.Parallel()
	type fields struct {
		executables []Executable
		completion  bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success - Work 100 tasks - expect complete",
			fields: fields{
				executables: nTasks(100),
				completion:  true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConcurrent(context.TODO(), e, true)
			for _, task := range tt.fields.executables {
				c.Add(task)
			}
			if err := c.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.fields.completion != c.IsCompleted() {
				t.Errorf("Execute() tasks expected to be completed but incomplete %+v", tt.fields.executables)
			}
		})
	}
}

func BenchmarkConcurrent_Execute(b *testing.B) {
	e := engine.NewEngine(context.TODO())
	e.Start(100)
	tasks := nTasks(1000)

	b.Logf("goroutines before adding tasks %d", runtime.NumGoroutine())
	c := NewConcurrent(context.TODO(), e, true)
	for _, task := range tasks {
		c.Add(task)
	}

	b.ResetTimer()
	b.Logf("starting goroutines %d", runtime.NumGoroutine())
	if err := c.Execute(); err != nil {
		b.Errorf("Execute() error = %v", err)
	}

	time.Sleep(5 * time.Second)
	b.Logf("ending goroutines %d", runtime.NumGoroutine())
}
