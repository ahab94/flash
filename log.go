package flash

import (
	logs "github.com/sirupsen/logrus"
)

var logger *logs.Logger

func init() {
	SetLogger(logs.New())
}

// SetLogger - sets custom logrus logger
func SetLogger(log *logs.Logger) {
	log.SetLevel(logs.DebugLevel)
	logger = log
}

func log(id string) *logs.Entry {
	return logger.WithFields(logs.Fields{
		"id": id,
	})
}
