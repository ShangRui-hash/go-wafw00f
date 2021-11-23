package logger

import (
	"github.com/sirupsen/logrus"
)

func Init(debug bool) error {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
	}
	return nil
}
