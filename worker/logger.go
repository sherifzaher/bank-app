package worker

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type RedisLogger struct{}

func NewRedisLogger() *RedisLogger {
	return &RedisLogger{}
}

func (logger *RedisLogger) Print(level zerolog.Level, opts ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(opts...))
}

func (logger *RedisLogger) Debug(args ...interface{}) {
	logger.Print(zerolog.DebugLevel, args...)
}
func (logger *RedisLogger) Info(args ...interface{}) {
	logger.Print(zerolog.InfoLevel, args...)
}
func (logger *RedisLogger) Warn(args ...interface{}) {
	logger.Print(zerolog.WarnLevel, args...)
}
func (logger *RedisLogger) Error(args ...interface{}) {
	logger.Print(zerolog.ErrorLevel, args...)
}
func (logger *RedisLogger) Fatal(args ...interface{}) {
	logger.Print(zerolog.FatalLevel, args...)
}
