package types

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

func log(ctx context.Context) *logs.Entry {
	if ctx != nil {
		logger.WithContext(ctx)
	}

	return logger.WithFields(logs.Fields{
		"package": "flash/types",
	})
}
