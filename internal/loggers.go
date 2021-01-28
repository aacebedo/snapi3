package internal

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	NormalLogger  *logrus.Entry //nolint:gochecknoglobals //Need the main logger to be global
	VerboseLogger *logrus.Entry //nolint:gochecknoglobals //Need the main logger to be global
)

func InitLoggers() {
	log := logrus.New()
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
	NormalLogger = log.WithFields(logrus.Fields{"emitter": "snapi3"})
	VerboseLogger = log.WithFields(logrus.Fields{"emitter": "snapi3"})
}
