package logging

import (
	"time"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Z *zap.Logger
	S *zap.SugaredLogger
}

// NewZap создает zap логгер.
// mode: "dev" (читаемый вывод) или "prod" (JSON).
func NewZap(mode string) (*Logger, error) {
	var cfg zap.Config
	switch mode {
	case "dev":
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	case "prod":
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "ts"
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	default:
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.OutputPaths = []string{"stdout"}      // можно добавить файл
	cfg.ErrorOutputPaths = []string{"stderr"} // можно добавить файл

	z, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{Z: z, S: z.Sugar()}, nil
}

func LogRequest(log *Logger, r *http.Request, status int, duration time.Duration) {
	log.Z.Info("request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Int("status", status),
		zap.Duration("duration", duration),
	)
}	