package log

import (
	"fmt"
	stdlog "log"
	"os"
	"strings"
)

// At package initialisation configure the standard logger to include
// timestamps and file/line information.  This aids in debugging while
// remaining light weight.  If you wish to integrate with a more
// sophisticated structured logger such as Zap or Logrus you can do so
// here.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelError
)

var minLevel = LevelError

func init() {
	stdlog.SetFlags(stdlog.LstdFlags | stdlog.Lshortfile)
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL"))) {
	case "debug":
		minLevel = LevelDebug
	case "info":
		minLevel = LevelInfo
	default:
		minLevel = LevelError
	}
}

// Entry provides a minimal structure for logging contextual
// information.  Each Entry carries a session identifier and a message
// identifier which are included in log output.  The provider name is
// hard coded as "whatsmeow".
type Entry struct {
	SessionID string
	MessageID string
}

// WithSession constructs a new Entry with the given session ID.  Use
// this helper when logging events associated with a particular
// WhatsApp session.
func WithSession(sessionID string) *Entry {
	return &Entry{SessionID: sessionID}
}

// WithMessageID returns a copy of the current entry with the
// supplied message ID set.  This is useful when logging actions
// related to a specific outbound or inbound message.
func (e *Entry) WithMessageID(msgID string) *Entry {
	return &Entry{SessionID: e.SessionID, MessageID: msgID}
}

// Info emits an informational log message.  The provider, session and
// message identifiers are automatically included as structured fields.
func (e *Entry) Info(format string, args ...interface{}) {
	if minLevel > LevelInfo {
		return
	}
	msg := fmt.Sprintf(format, args...)
	stdlog.Printf("level=info provider=whatsmeow session_id=%s message_id=%s %s", e.SessionID, e.MessageID, msg)
}

// Error emits an error log message.  The provider, session and
// message identifiers are automatically included as structured fields.
func (e *Entry) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	stdlog.Printf("level=error provider=whatsmeow session_id=%s message_id=%s %s", e.SessionID, e.MessageID, msg)
}

// Debug emits a debug log message gated by LOG_LEVEL.
func (e *Entry) Debug(format string, args ...interface{}) {
	if minLevel > LevelDebug {
		return
	}
	msg := fmt.Sprintf(format, args...)
	stdlog.Printf("level=debug provider=whatsmeow session_id=%s message_id=%s %s", e.SessionID, e.MessageID, msg)
}

// Package-level helpers for logs not tied to a particular Entry/session
func Debugf(format string, args ...interface{}) {
	if minLevel > LevelDebug {
		return
	}
	stdlog.Printf("level=debug provider=whatsmeow %s", fmt.Sprintf(format, args...))
}

func Infof(format string, args ...interface{}) {
	if minLevel > LevelInfo {
		return
	}
	stdlog.Printf("level=info provider=whatsmeow %s", fmt.Sprintf(format, args...))
}

func Errorf(format string, args ...interface{}) {
	stdlog.Printf("level=error provider=whatsmeow %s", fmt.Sprintf(format, args...))
}
