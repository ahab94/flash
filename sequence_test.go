package flash

import (
	"context"
	"testing"
)

func TestSequence_Execute(t *testing.T) {
	type fields struct {
		executables []Executable
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success - execute all tasks",
			fields: fields{executables: []Executable{
				&testTask{
					Fail:  false,
					Delay: "100ms",
				},
				&testTask{
					Fail:  false,
					Delay: "200ms",
				},
			}},
			wantErr: false,
		},
		{
			name: "fail - execute all tasks - fail true",
			fields: fields{executables: []Executable{
				&testTask{
					Fail:  true,
					Delay: "100ms",
				},
				&testTask{
					Fail:  false,
					Delay: "200ms",
				},
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewSequence(context.TODO())
			for _, task := range tt.fields.executables {
				p.Add(task)
			}
			if err := p.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
