package flash

import (
	"errors"
	"time"
)

type testTask struct {
	id     int
	fail   bool
	panic  bool
	delay  string
	status string
}

func (t *testTask) Execute() error {
	duration, err := time.ParseDuration(t.delay)
	if err != nil {
		logger.Warnf("parse duration error... overriding delay as 1 second")
		duration = time.Second
	}
	time.Sleep(duration)
	if t.panic {
		panic(t)
	}

	if t.fail {
		return errors.New("some error")
	}
	return nil
}

func (t *testTask) OnFailure(err error) {
	t.status = "failed"
}

func (t *testTask) OnSuccess() {
	t.status = "completed"
}

func (t *testTask) IsCompleted() bool {
	return t.status == "completed"
}
