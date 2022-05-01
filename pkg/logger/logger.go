package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "errors":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3

	logger := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger()

	return &Logger{
		logger: &logger,
	}
}

func (l *Logger) log(level, message string, args ...interface{}) {
	if len(args) == 0 {
		switch level {
		case "errors":
			l.logger.Error().Msg(message)
		case "warn":
			l.logger.Warn().Msg(message)
		case "debug":
			l.logger.Debug().Msg(message)
		case "fatal":
			l.logger.Fatal().Msg(message)
		default:
			l.logger.Info().Msg(message)
		}
	} else {
		switch level {
		case "errors":
			l.logger.Error().Msgf(message, args...)
		case "warn":
			l.logger.Warn().Msgf(message, args...)
		case "debug":
			l.logger.Debug().Msgf(message, args...)
		case "fatal":
			l.logger.Fatal().Msgf(message, args...)
		default:
			l.logger.Info().Msgf(message, args...)
		}
	}
}

func (l *Logger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(level, msg.Error(), args...)
	case string:
		l.log(level, msg, args...)
	default:
		l.log(level, fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg("debug", message, args)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.log("info", message, args)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.log("warn", message, args)
}

func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.msg("errors", message, args)
}

func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg("fatal", message, args...)

	os.Exit(1)
}
