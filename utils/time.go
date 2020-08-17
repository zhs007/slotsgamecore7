package sgc7utils

import "time"

// TimeI - Time
type TimeI interface {
	// Now - get now time
	Now() time.Time
}

// Time - default Time
type Time struct{}

// Now - get now time
func (t Time) Now() time.Time {
	return time.Now()
}

var gTime TimeI

// FormatNow - format time
func FormatNow(t TimeI) string {
	return t.Now().Format("2006-01-02_15:04:05")
}

func init() {
	gTime = &Time{}
}
