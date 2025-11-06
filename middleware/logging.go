package middleware

import (
	"day2/logging"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

type Logger struct {
	Z *zap.Logger
	S *zap.SugaredLogger
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(p []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(p)
	r.size += n
	return n, err
}


func LoggingZap(lg *logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &responseRecorder{ResponseWriter: w}
			next.ServeHTTP(rec, r)
			lg.Z.Info("http",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rec.status),
				zap.Int("size", rec.size),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

func NewZap(mode string) (*Logger, error) {
	var cfg zap.Config
	switch mode {
	case "prod":
		cfg = zap.NewProductionConfig()
	default:
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	// stdout/stderr по умолчанию — не меняем, при желании добавь файлы

	z, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{Z: z, S: z.Sugar()}, nil
}