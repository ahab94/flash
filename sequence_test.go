package flash

import (
	"context"
	"testing"
)

func TestSequence_Execute(t *testing.T) {
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
			name: "success - Work 10 tasks - expect complete",
			fields: fields{
				executables: nTasks(10),
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSequence(context.TODO())
			for _, task := range tt.fields.executables {
				s.Add(task)
			}
			if err := s.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.fields.completion != s.IsCompleted() {
				t.Errorf("Execute() tasks expected to be completed but incomplete %+v", tt.fields.executables)
			}
		})
	}
}
