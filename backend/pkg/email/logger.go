package email

import (
	"context"
	"fmt"
	"io"
	"os"
)

// LoggerSender writes messages to an io.Writer; useful for dev/tests.
type LoggerSender struct{ w io.Writer }

func NewLoggerSender(w io.Writer) *LoggerSender {
	if w == nil {
		w = os.Stdout
	}
	return &LoggerSender{w: w}
}

func (s *LoggerSender) Send(ctx context.Context, msg Message) error {
	if err := msg.Validate(); err != nil {
		return err
	}
	_, err := fmt.Fprintf(s.w, "[email] to=%v subject=%q\n", append(append([]string{}, msg.To...), msg.Cc...), msg.Subject)
	return err
}
