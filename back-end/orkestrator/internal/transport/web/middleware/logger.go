package middleware

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerConfig struct {
	gin.LoggerConfig
	SkipPaths []string
}

func Logger() gin.HandlerFunc {
	return LoggerWithConfig(gin.LoggerConfig{}, log.Logger)
}

func LoggerWithWriter(out io.Writer, notlogged ...string) gin.HandlerFunc {
	return LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: notlogged,
	}, log.Logger)
}

func LoggerWithConfig(conf gin.LoggerConfig, logger zerolog.Logger) gin.HandlerFunc {
	return newLoggerMiddleware(conf, logger)
}

func newLoggerMiddleware(conf gin.LoggerConfig, logger zerolog.Logger) gin.HandlerFunc {
	skip := computeSkip(conf)

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		// Process request
		c.Next()
		// Log only when path is not being skipped
		if _, ok := skip[path]; ok {
			return
		}

		logger.Info().
			Str("StartTimestamp", fmt.Sprintf("%d", start.Unix())).
			Str("ClientIP", c.ClientIP()).
			Str("Method", c.Request.Method).
			Str("Status", fmt.Sprintf("%d", c.Writer.Status())).
			Str("BodySize", fmt.Sprintf("%d", c.Writer.Size())).
			Str("ErrorMessage", c.Errors.ByType(gin.ErrorTypePrivate).String()).
			Str("Path", path).
			Str("Query", raw).
			Msg(" ")
	}
}

func computeSkip(conf gin.LoggerConfig) map[string]struct{} {
	notlogged := conf.SkipPaths

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return skip
}
