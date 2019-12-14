package flash

import (
	"context"
	"testing"
)

func TestParallel_Execute(t *testing.T) {
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
			fields: fields{executables: []Executable{
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
			p := NewParallel(context.TODO())
			for _, task := range tt.fields.executables {
				p.Add(task)
			}
			if err := p.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.fields.completion != p.IsCompleted() {
				t.Errorf("Execute() tasks expected to be completed but incomplete %+v", tt.fields.executables)
			}
		})
	}
}

func BenchmarkParallel_Execute(b *testing.B) {
	tasks := nTasks(1000)
	p := NewParallel(context.TODO())
	for _, task := range tasks {
		p.Add(task)
	}

	b.ResetTimer()
	if err := p.Execute(); err != nil {
		b.Errorf("Execute() error = %v", err)
	}
}
