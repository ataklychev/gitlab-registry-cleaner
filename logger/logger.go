package logger

import "go.uber.org/zap"

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// NewLogger make logger.
func NewLogger(production bool) (Logger, func() error) {
	var logger *zap.Logger
	if production {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}

	sugar := logger.Sugar()

	return sugar, sugar.Sync
}
