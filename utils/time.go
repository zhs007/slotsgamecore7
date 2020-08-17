package sgc7utils

import "time"

// TimeI - Time
type TimeI interface {
	// Now - get now time
	Now() time.Time
	// GetDefaultLocation - get DefaultLocation
	GetDefaultLocation() *time.Location
}

// Time - default Time
type Time struct {
	local *time.Location
}

// Now - get now time
func (t Time) Now() time.Time {
	return time.Now()
}

// GetDefaultLocation - get DefaultLocation
func (t Time) GetDefaultLocation() *time.Location {
	return t.local
}

var gTime TimeI

// FormatNow - format time
func FormatNow(t TimeI) string {
	return t.Now().In(t.GetDefaultLocation()).Format("2006-01-02_15:04:05")
}

func init() {
	t := &Time{}

	l, err := time.LoadLocation("")
	if err == nil {
		t.local = l
	}

	gTime = t
}
