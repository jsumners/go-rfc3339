package rfc3339

import "time"

type RFC3339 interface {
	ToString() string
}

// DateTime represents an RFC 3339 `date-time`. It is a wrapper for
// [time.Time].
type DateTime struct {
	time.Time
}

// FullDate represents an RFC 3339 `full-date`. It is a wrapper for
// [time.Time]. FullDate objects set the time parts to midnight (00:00:00) at
// the UTC (+00:00) offset.
type FullDate struct {
	time.Time
}
