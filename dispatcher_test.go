package flash

import (
	"context"
	"testing"
)

func TestDispatcher_Stop(t *testing.T) {
	type fields struct {
		stop chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
	}{

		{
			name:   "success - close worker",
			fields: fields{stop: make(chan struct{})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDispatcher(context.TODO())
			d.Start(0)
			d.stop = tt.fields.stop
			d.Stop()
		})
	}
}

func TestDispatcher_dispatch(t *testing.T) {
	d := NewDispatcher(context.TODO())
	w := NewWorker(context.TODO(), d.pool, d.counter)
	d.workers = append(d.workers, w)

	t.Parallel()
	type fields struct {
		task *testTask
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "success - execute test task 1 - Fail false - Delay 100ms",
			fields: fields{task: &testTask{ID: 1, Fail: false, Delay: "100ms"}},
			want:   "completed",
		},
		{
			name:   "Fail - execute test task 2 - Delay 200ms - Fail true",
			fields: fields{task: &testTask{ID: 2, Delay: "200ms", Fail: true}},
			want:   "failed",
		},
		{
			name:   "success - execute test task 3 - Fail false - Delay 300ms",
			fields: fields{task: &testTask{ID: 3, Fail: false, Delay: "300ms"}},
			want:   "completed",
		},
		{
			name:   "Fail - execute test task 4 - Fail true - Delay 400ms",
			fields: fields{task: &testTask{ID: 4, Fail: true, Delay: "400ms"}},
			want:   "failed",
		},
		{
			name:   "success - execute test task 5 - Fail false - Delay 500ms",
			fields: fields{task: &testTask{ID: 5, Fail: false, Delay: "500ms"}},
			want:   "completed",
		},
		{
			name:   "Fail - execute test task 6 - Delay 600ms - Fail true",
			fields: fields{task: &testTask{ID: 6, Delay: "600ms", Fail: true}},
			want:   "failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d.Start(0)
			d.input <- tt.fields.task
			<-d.pool
		})
	}
}
