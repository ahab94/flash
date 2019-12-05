package flash

import (
	"context"
)

// RecoverPanic - used to avoid crashes following unexpected panic
func RecoverPanic(ctx context.Context, executable Executable) {
	if r := recover(); r != nil {
		log(ctx).Warnf("recovered from panic while executing job: %v", executable)
	}
}
