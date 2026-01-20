package logx

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type ctxKey int

const loggerKey ctxKey = iota

var defaultLogger zerolog.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func SetDefault(l zerolog.Logger) { defaultLogger = l }

func NewContext(ctx context.Context, l zerolog.Logger) context.Context { return context.WithValue(ctx, loggerKey, l) }

func WithContext(ctx context.Context) zerolog.Logger {
    if v := ctx.Value(loggerKey); v != nil {
        if l, ok := v.(zerolog.Logger); ok {
            return l
        }
    }
    return defaultLogger
}
