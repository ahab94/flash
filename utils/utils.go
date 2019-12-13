package utils

import (
	"context"

	logger "github.com/sirupsen/logrus"
)

// RecoverPanic - used to avoid crashes following unexpected Panic
func RecoverPanic(ctx context.Context) {
	if r := recover(); r != nil {
		logger.Warnf("recovered from Panic")
	}
}
