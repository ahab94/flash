package flash

import (
	"context"
	"testing"
)

func TestConcurrent_Execute(t *testing.T) {
	d := NewDispatcher(context.TODO())
	d.Start(10)

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
		{
			name: "success - work all tasks - expect incomplete",
			fields: fields{
				executables: []Executable{
					&testTask{
						ID:    1,
						Fail:  true,
						Delay: "2s",
					}, &testTask{
						ID:    2,
						Fail:  false,
						Delay: "100ms",
					},
				},
				completion: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConcurrent(context.TODO(), d)
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
	d := NewDispatcher(context.TODO())
	d.Start(100)
	tasks := nTasks(1000)
	c := NewConcurrent(context.TODO(), d)
	for _, task := range tasks {
		c.Add(task)
	}

	b.ResetTimer()
	if err := c.Execute(); err != nil {
		b.Errorf("Execute() error = %v", err)
	}
}
