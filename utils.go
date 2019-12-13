package flash

import (
	"context"
)

// RecoverPanic - used to avoid crashes following unexpected panic
func RecoverPanic(ctx context.Context) {
	if r := recover(); r != nil {
		logger.Warnf("recovered from panic")
	}
}
