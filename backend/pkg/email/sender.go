package email

import "context"

// Sender is a transport-agnostic email sender.
// Implementations: SMTPSender, LoggerSender (dev/test).
type Sender interface {
	Send(ctx context.Context, msg Message) error
}
