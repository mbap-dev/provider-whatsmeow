package log

import (
	"fmt"
	"log"
)

// At package initialisation configure the standard logger to include
// timestamps and file/line information.  This aids in debugging while
// remaining light weight.  If you wish to integrate with a more
// sophisticated structured logger such as Zap or Logrus you can do so
// here.
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
	msg := fmt.Sprintf(format, args...)
	log.Printf("level=info provider=whatsmeow session_id=%s message_id=%s %s", e.SessionID, e.MessageID, msg)
}

// Error emits an error log message.  The provider, session and
// message identifiers are automatically included as structured fields.
func (e *Entry) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("level=error provider=whatsmeow session_id=%s message_id=%s %s", e.SessionID, e.MessageID, msg)
}
