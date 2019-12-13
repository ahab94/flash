package flash

import (
	"context"
	"testing"
)

func Test_executor_IsCompleted(t *testing.T) {
	type fields struct {
		executables []Executable
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "all tasks complete",
			fields: fields{
				executables: []Executable{
					&testTask{
						Status: "completed",
					},
					&testTask{
						Status: "completed",
					},
					&testTask{
						Status: "completed",
					},
					&testTask{
						Status: "completed",
					},
				},
			},
			want: true,
		},
		{
			name: "all tasks complete",
			fields: fields{
				executables: []Executable{
					&testTask{
						Status: "completed",
					},
					&testTask{
						Status: "failed",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				id:          "0",
				ctx:         context.TODO(),
				executables: tt.fields.executables,
			}
			if got := e.IsCompleted(); got != tt.want {
				t.Errorf("IsCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}
