package log

import (
	"os"
	"strings"
	"time"

	appcfg "mutabaahapi/internal/infra/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func New(cfg appcfg.Config) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Determine log level: from config (LOG_LEVEL) or APP_ENV fallback
	levelStr := cfg.LogLevel
	var level zerolog.Level
	if levelStr != "" {
		if parsed, err := zerolog.ParseLevel(strings.ToLower(levelStr)); err == nil {
			level = parsed
		} else {
			level = zerolog.InfoLevel
		}
	} else {
		// fallback: dev/local -> debug, else info
		if e := strings.ToLower(cfg.Env); e == "debug" || e == "dev" || e == "local" {
			level = zerolog.DebugLevel
		} else {
			level = zerolog.InfoLevel
		}
	}
	zerolog.SetGlobalLevel(level)

	// Choose output: pretty console if dev/local or LOG_PRETTY=true
	pretty := cfg.LogPretty
	if !pretty {
		if e := strings.ToLower(cfg.Env); e == "dev" || e == "local" {
			pretty = true
		}
	}

	if pretty {
		cw := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		return zerolog.New(cw).With().Timestamp().Logger()
	}
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
