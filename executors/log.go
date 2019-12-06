package executors

import (
	"context"

	logs "github.com/sirupsen/logrus"
)

var logger *logs.Logger

func init() {
	logger = logs.New()

	logger.SetFormatter(&logs.TextFormatter{
		FullTimestamp: true,
	})
}

func log(ctx context.Context, executor string) *logs.Entry {
	if ctx != nil {
		logger.WithContext(ctx)
	}

	return logger.WithFields(logs.Fields{
		"package":  "executors",
		"executor": executor,
	})
}
