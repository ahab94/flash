package flash

import (
	"context"
	"testing"
)

func TestWorker_work(t *testing.T) {
	type fields struct {
		task *testTask
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "success - execute test task - fail false",
			fields: fields{task: &testTask{id: 0, fail: false, delay: "100ms"}},
			want:   "completed",
		},
		{
			name:   "fail - execute test task - fail true",
			fields: fields{task: &testTask{id: 0, fail: true, delay: "100ms"}},
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
			if tt.fields.task.status != tt.want {
				t.Errorf("worker <- task failed wanted: %s got %s", tt.want, tt.fields.task.status)
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
		task *testTask
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "success - execute test task 1 - fail false - delay 100ms",
			fields: fields{task: &testTask{id: 1, fail: false, delay: "100ms"}},
			want:   "completed",
		},
		{
			name:   "fail - execute test task 2 - delay 200ms - panic true",
			fields: fields{task: &testTask{id: 2, delay: "200ms", panic: true}},
			want:   "",
		},
		{
			name:   "success - execute test task 3 - fail false - delay 300ms",
			fields: fields{task: &testTask{id: 3, fail: false, delay: "300ms"}},
			want:   "completed",
		},
		{
			name:   "fail - execute test task 4 - fail true - delay 400ms",
			fields: fields{task: &testTask{id: 4, fail: true, delay: "400ms"}},
			want:   "failed",
		},
		{
			name:   "success - execute test task 5 - fail false - delay 500ms",
			fields: fields{task: &testTask{id: 5, fail: false, delay: "500ms"}},
			want:   "completed",
		},
		{
			name:   "fail - execute test task 6 - delay 600ms - panic true",
			fields: fields{task: &testTask{id: 6, delay: "600ms", panic: true}},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// execute task and wait for it to complete
			w.input <- tt.fields.task
			<-w.pool
			if tt.fields.task.status != tt.want {
				t.Errorf("worker <- task failed wanted: %s got %s", tt.want, tt.fields.task.status)
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
