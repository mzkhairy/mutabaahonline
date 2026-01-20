package email

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// Message represents an email with optional text and HTML bodies.
type Message struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Text    string // plain text body
	HTML    string // HTML body
	Headers map[string]string
}

func (m *Message) Validate() error {
	if _, err := mail.ParseAddress(m.From); err != nil {
		return fmt.Errorf("invalid from: %w", err)
	}
	if len(allRecipients(m)) == 0 {
		return errors.New("at least one recipient is required")
	}
	return nil
}

func allRecipients(m *Message) []string {
	var rcpt []string
	rcpt = append(rcpt, m.To...)
	rcpt = append(rcpt, m.Cc...)
	rcpt = append(rcpt, m.Bcc...)
	return rcpt
}

// buildMIME builds a minimal MIME message. Uses multipart/alternative when HTML is present.
func buildMIME(m *Message) []byte {
	// Common headers
	headers := map[string]string{
		"From":         m.From,
		"To":           strings.Join(m.To, ", "),
		"Cc":           strings.Join(m.Cc, ", "),
		"Subject":      encodeHeader(m.Subject),
		"MIME-Version": "1.0",
		"Date":         time.Now().Format(time.RFC1123Z),
	}
	// Merge custom headers (allow override if needed)
	for k, v := range m.Headers {
		headers[k] = v
	}

	var sb strings.Builder
	// Write headers (skip empty values like empty Cc)
	for k, v := range headers {
		if strings.TrimSpace(v) == "" {
			continue
		}
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v)
		sb.WriteString("\r\n")
	}

	// Body
	if m.HTML != "" {
		boundary := fmt.Sprintf("mixed_%d", time.Now().UnixNano())
		sb.WriteString("Content-Type: multipart/alternative; boundary=\"")
		sb.WriteString(boundary)
		sb.WriteString("\"\r\n\r\n")

		// text part (optional, recommended for clients without HTML)
		if m.Text != "" {
			sb.WriteString("--" + boundary + "\r\n")
			sb.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
			sb.WriteString(m.Text)
			sb.WriteString("\r\n")
		}
		// html part
		sb.WriteString("--" + boundary + "\r\n")
		sb.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		sb.WriteString(m.HTML)
		sb.WriteString("\r\n--" + boundary + "--\r\n")
	} else {
		// plain text only
		sb.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		sb.WriteString(m.Text)
		sb.WriteString("\r\n")
	}
	return []byte(sb.String())
}

// naive header encoding: keep it simple for ASCII subjects; real-world should use RFC 2047 encoding
func encodeHeader(s string) string { return s }
