package logger

import (
	"fmt"

	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type zeroLogger struct {
	logger zerolog.Logger
	tags   map[string]string
}

// NewLogger new simple logger
func NewLogger() Log {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return &zeroLogger{
		logger: log.With().Logger(),
	}
}

// NewLogByConfig new zerolog impl
func NewLogByConfig(cfg *config.Config) Log {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	tags := map[string]string{
		"app": cfg.App,
		"env": cfg.Env,
	}

	subContext := log.With()
	for k, v := range tags {
		subContext = subContext.Str(k, v)
	}

	sublogger := subContext.Logger()

	return &zeroLogger{
		logger: sublogger,
		tags:   tags,
	}
}

func (zl zeroLogger) sendErrorToSentry(err error) {
	scope := sentry.NewScope()
	scope.SetTags(zl.tags)
}

func (zl zeroLogger) Print(v ...any) {
	zl.logger.Print(v...)
}

func (zl zeroLogger) Printf(format string, v ...any) {
	zl.logger.Printf(format, v...)
}

func (zl zeroLogger) Println(v ...any) {
	zl.logger.Print(v...)
	zl.logger.Print("\n")
}

func (zl zeroLogger) Debug(v ...any) {
	zl.logger.Debug().Msg(fmt.Sprint(v...))
}

func (zl zeroLogger) Debugf(format string, v ...any) {
	zl.logger.Debug().Msg(fmt.Sprintf(format, v...))
}

func (zl zeroLogger) Info(v ...any) {
	zl.logger.Info().Msg(fmt.Sprint(v...))
}

func (zl zeroLogger) Infof(format string, v ...any) {
	zl.logger.Info().Msg(fmt.Sprintf(format, v...))
}

func (zl zeroLogger) Warn(v ...any) {
	zl.logger.Warn().Msg(fmt.Sprint(v...))
}

func (zl zeroLogger) Warnf(format string, v ...any) {
	zl.logger.Warn().Msg(fmt.Sprintf(format, v...))
}

func (zl zeroLogger) Error(err error, v ...any) {
	zl.logger.Error().Stack().Err(err).Msg(fmt.Sprint(v...))
}

func (zl zeroLogger) Errorf(err error, format string, v ...any) {
	zl.logger.Error().Stack().Err(err).Msg(fmt.Sprintf(format, v...))
}

func (zl zeroLogger) Fatal(err error, v ...any) {
	zl.logger.Fatal().Stack().Err(err).Msg(fmt.Sprint(v...))
	go zl.sendErrorToSentry(err)
}

func (zl zeroLogger) Fatalf(err error, format string, v ...any) {
	zl.logger.Fatal().Stack().Err(err).Msg(fmt.Sprintf(format, v...))
	go zl.sendErrorToSentry(err)
}
