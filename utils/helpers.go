package utils

import (
	"errors"
	"time"

	logger "github.com/sirupsen/logrus"
)

type TestTask struct {
	ID     int
	Fail   bool
	Panic  bool
	Delay  string
	Status string
}

func (t *TestTask) Execute() error {
	duration, err := time.ParseDuration(t.Delay)
	if err != nil {
		logger.Warnf("parse duration error... overriding Delay as 1 second")
		duration = time.Second
	}
	time.Sleep(duration)
	if t.Panic {
		panic(t)
	}

	if t.Fail {
		return errors.New("some error")
	}
	return nil
}

func (t *TestTask) OnFailure(err error) {
	t.Status = "failed"
}

func (t *TestTask) OnSuccess() {
	t.Status = "completed"
}

func (t *TestTask) IsCompleted() bool {
	return t.Status == "completed"
}
