package timestamp

import (
	"fmt"
	"time"
)

// Timestamp is a resource that returns the current or set time in a specified'
// format.
type Timestamp struct {
	// format is a textual representation of the time value formatted
	// according to the layout defined by the argument.  It follows
	// the go representation for time formats.  There are several mapped
	// values that will convert to the go equivalent including: rfc3339,
	// rfc3999nano, unix, and unixnano.
	format    string
	timestamp string
	provider  Provider
}

// New returns a new timestamp resource.
func New(cfg Config) *Timestamp {
	return &Timestamp{
		format:    cfg.Format,
		timestamp: cfg.Timestamp,
		provider:  RealTime{},
	}
}

// WithProvider sets the time provider for the timestamp resource. This is
// currently only used in testing.
func (t *Timestamp) WithProvider(provider Provider) *Timestamp {
	t.provider = provider
	return t
}

// Get implements the resource interface and returns the current time in the
// specified format.
func (t *Timestamp) Get() string {
	if t.timestamp != "" {
		return t.timestamp
	}

	now := t.provider.Now()
	var format string
	switch t.format {
	case "unix":
		return fmt.Sprint(now.Unix())
	case "unixnano":
		return fmt.Sprint(now.UnixNano())
	case "rfc3339":
		format = time.RFC3339
	case "rfc3339nano":
		format = time.RFC3339Nano
	case "rfc1123":
		format = time.RFC1123
	case "rfc1123z":
		format = time.RFC1123Z
	case "common_log":
		format = CommonLogFormat
	default:
		format = t.format
	}

	formatted := now.Format(format)
	// if the format is not valid, it just returns the format and we'll just return
	// the default format.
	if formatted == format {
		return now.Format("")
	}

	return formatted
}
