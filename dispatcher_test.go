package flash

import (
	"context"
	"testing"

	"github.com/ahab94/flash/utils"
)

func TestDispatcher_Stop(t *testing.T) {
	badChan := make(chan struct{})
	close(badChan)

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
		{
			name:   "success - close worker on bad channel with recovery",
			fields: fields{stop: badChan},
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
		task *utils.TestTask
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "success - execute test task 1 - Fail false - Delay 100ms",
			fields: fields{task: &utils.TestTask{ID: 1, Fail: false, Delay: "100ms"}},
			want:   "completed",
		},
		{
			name:   "Fail - execute test task 2 - Delay 200ms - Panic true",
			fields: fields{task: &utils.TestTask{ID: 2, Delay: "200ms", Panic: true}},
			want:   "",
		},
		{
			name:   "success - execute test task 3 - Fail false - Delay 300ms",
			fields: fields{task: &utils.TestTask{ID: 3, Fail: false, Delay: "300ms"}},
			want:   "completed",
		},
		{
			name:   "Fail - execute test task 4 - Fail true - Delay 400ms",
			fields: fields{task: &utils.TestTask{ID: 4, Fail: true, Delay: "400ms"}},
			want:   "failed",
		},
		{
			name:   "success - execute test task 5 - Fail false - Delay 500ms",
			fields: fields{task: &utils.TestTask{ID: 5, Fail: false, Delay: "500ms"}},
			want:   "completed",
		},
		{
			name:   "Fail - execute test task 6 - Delay 600ms - Panic true",
			fields: fields{task: &utils.TestTask{ID: 6, Delay: "600ms", Panic: true}},
			want:   "",
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
