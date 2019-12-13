package flash

import (
	"context"
	"testing"

	"github.com/ahab94/flash/utils"
)

func TestWorker_work(t *testing.T) {
	type fields struct {
		task *utils.TestTask
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "success - execute test task - Fail false",
			fields: fields{task: &utils.TestTask{ID: 0, Fail: false, Delay: "100ms"}},
			want:   "completed",
		},
		{
			name:   "Fail - execute test task - Fail true",
			fields: fields{task: &utils.TestTask{ID: 0, Fail: true, Delay: "100ms"}},
			want:   "failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Worker{
				id:      "0",
				ctx:     context.TODO(),
				counter: &counter{},
				stop:    make(chan struct{}),
				input:   make(chan Executable),
				pool:    make(chan chan Executable),
			}
			// start worker
			w.Start()

			// execute task and wait for it to complete
			w.input <- tt.fields.task
			<-w.pool
			if tt.fields.task.Status != tt.want {
				t.Errorf("worker <- task failed wanted: %s got %s", tt.want, tt.fields.task.Status)
			}
		})
	}
}

func TestWorker_workParallel(t *testing.T) {
	w := &Worker{
		id:      "0",
		ctx:     context.TODO(),
		counter: &counter{},
		stop:    make(chan struct{}),
		input:   make(chan Executable),
		pool:    make(chan chan Executable),
	}
	// start worker
	w.Start()

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
			// execute task and wait for it to complete
			w.input <- tt.fields.task
			<-w.pool
			if tt.fields.task.Status != tt.want {
				t.Errorf("worker <- task failed wanted: %s got %s", tt.want, tt.fields.task.Status)
			}
		})
	}
}

func TestWorker_Stop(t *testing.T) {
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
			w := &Worker{
				id:      "0",
				ctx:     context.TODO(),
				counter: &counter{},
				stop:    tt.fields.stop,
				input:   make(chan Executable),
				pool:    make(chan chan Executable),
			}
			// start worker...
			w.Start()
			w.Stop()
		})
	}
}
