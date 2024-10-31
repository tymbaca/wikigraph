package logger

import (
	"github.com/rs/zerolog"
)

func init() {
	Init(zerolog.DebugLevel)
}

var logger *zerolog.Logger

// Init initializes a default logger.
func Init(level zerolog.Level) {
	l := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger().Level(level)

	logger = &l
}

// Debug prints log in debug level
func Debug(msg string) {
	logger.Debug().Msg(msg)
}

// Debugf prints formatted message in debug level
func Debugf(fmt string, v ...interface{}) {
	logger.Debug().Msgf(fmt, v...)
}

// Error prints log in error level
func Error(msg string) {
	logger.Error().Msg(msg)
}

// Errorf prints formatted message in error level
func Errorf(fmt string, v ...interface{}) {
	logger.Error().Msgf(fmt, v...)
}

// ErrorWithFields prints fields in error level
func ErrorWithFields(fields map[string]interface{}) {
	logger.Error().Fields(fields).Send()
}

// Info prints log in info level
func Info(msg string) {
	logger.Info().Msg(msg)
}

// Infof prints formatted message in info level
func Infof(fmt string, v ...interface{}) {
	logger.Info().Msgf(fmt, v...)
}

// InfoWithFields prints fields in info level
func InfoWithFields(fields map[string]interface{}) {
	logger.Info().Fields(fields).Send()
}

// Warn prints log in warn level
func Warn(msg string) {
	logger.Warn().Msg(msg)
}

// Warnf prints formatted message in warn level
func Warnf(fmt string, v ...interface{}) {
	logger.Warn().Msgf(fmt, v...)
}

// Fatal prints log in fatal level
func Fatal(msg string) {
	logger.Fatal().Msg(msg)
}

// Fatalf prints formatted message in fatal level
func Fatalf(fmt string, v ...interface{}) {
	logger.Fatal().Msgf(fmt, v...)
}
