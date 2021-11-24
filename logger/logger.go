package logger

import (
	"github.com/sirupsen/logrus"
)

func Init(debug, silent bool) error {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if silent {
		logrus.SetLevel(logrus.ErrorLevel)
	}
	if !silent && !debug {
		logrus.SetLevel(logrus.InfoLevel)
	}
	return nil
}
