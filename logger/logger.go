package logger

import (
	"github.com/samber/do"
	"go.uber.org/zap"
)

var (
	Log *zap.Logger
)

func init() {

	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.DisableStacktrace = true
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	log, err := zapConfig.Build()
	if err != nil {
		panic(err.Error())
	}

	Log = log
}

// NewWithDI creates a new logger instance
func New(i *do.Injector) (*zap.Logger, error) {
	return Log, nil
}
