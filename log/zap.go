package log

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapWrapper struct {
	logger *zap.SugaredLogger
}

func NewLogger() *ZapWrapper {
	return &ZapWrapper{
		logger: newZapLogger(),
	}
}

func newZapLogger() *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:       "_m",
		NameKey:          "logger",
		LevelKey:         "_l",
		EncodeLevel:      zapcore.LowercaseLevelEncoder,
		TimeKey:          "_t",
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		CallerKey:        "",
		FunctionKey:      "",
		StacktraceKey:    "",
		LineEnding:       "",
		EncodeDuration:   func(_ time.Duration, _ zapcore.PrimitiveArrayEncoder) {},
		EncodeCaller:     func(_ zapcore.EntryCaller, _ zapcore.PrimitiveArrayEncoder) {},
		EncodeName:       func(_ string, _ zapcore.PrimitiveArrayEncoder) {},
		ConsoleSeparator: "",
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), os.Stdout, zapcore.DebugLevel)

	return zap.New(core).Sugar()
}

func (l *ZapWrapper) Info(msg string, args ...interface{}) {
	l.logger.Infow(msg, args...)
}

func (l *ZapWrapper) Error(msg string, args ...interface{}) {
	l.logger.Errorw(msg, args...)
}

func (l *ZapWrapper) Close() error {
	err := l.logger.Sync()

	if err == nil || isSyncInvalidError(err) {
		return nil
	}

	return fmt.Errorf("sync zap log: %w", err)
}

func isSyncInvalidError(err error) bool {
	var pathErr *os.PathError

	if !errors.As(err, &pathErr) {
		return false
	}

	switch {
	case errors.Is(pathErr.Err, syscall.ENOTTY):
	case errors.Is(pathErr.Err, syscall.EINVAL):
	default:
		return false
	}

	return true
}
