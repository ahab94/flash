package xcruncher

import (
	"context"
)

func recoverPanic(ctx context.Context, executable Executable) {
	if r := recover(); r != nil {
		log(ctx).Warnf("recovered from panic while executing job: %v", executable)
	}
}
