package logrus_zap_hook

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type ZapHook struct {
	Logger *zap.Logger
}

func NewZapHook(logger *zap.Logger) (*ZapHook, error) {
	return &ZapHook{
		Logger: logger,
	}, nil
}

func (hook *ZapHook) Fire(entry *logrus.Entry) error {
	fields := make([]zap.Field, 0, 10)

	for key, value := range entry.Data {
		if key == logrus.ErrorKey {
			fields = append(fields, zap.Error(value.(error)))
		} else {
			fields = append(fields, zap.Any(key, value))
		}
	}

	switch entry.Level {
	case logrus.PanicLevel:
		hook.Logger.Panic(entry.Message, fields...)
	case logrus.FatalLevel:
		hook.Logger.Fatal(entry.Message, fields...)
	case logrus.ErrorLevel:
		hook.Logger.Error(entry.Message, fields...)
	case logrus.WarnLevel:
		hook.Logger.Warn(entry.Message, fields...)
	case logrus.InfoLevel:
		hook.Logger.Info(entry.Message, fields...)
	case logrus.DebugLevel, logrus.TraceLevel:
		hook.Logger.Debug(entry.Message, fields...)
	}

	return nil
}

func (hook *ZapHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
