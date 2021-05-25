package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/private-square/bkst-users-api/utils/dateutils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

var (
	logger *zap.Logger
)

func init() {
	logConfig := zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "caller",
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths: []string{"stdout"},
	}

	var err error
	if logger, err = logConfig.Build(); err != nil {
		log.Fatalln("Unable to initialize logger")
	}
}

func Info(msg string, tags ...zapcore.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, tags...)
}

func Warn(msg string, tags ...zapcore.Field) {
	logger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, tags...)
}

func Error(msg string, err error, tags ...zapcore.Field) {
	if err != nil {
		tags = append(tags, zap.NamedError("error", err))
	}
	logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, tags...)
}

func Fatal(msg string) {
	logger.WithOptions(zap.AddCallerSkip(1)).Fatal(msg)
}

// GinZap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
func GinZap() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := dateutils.GetDateTimeNow()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		end := dateutils.GetDateTimeNow()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.WithOptions(zap.AddCallerSkip(1)).Error(e)
			}
		} else {
			logger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("latency", latency.String()),
			)
		}
	}
}
