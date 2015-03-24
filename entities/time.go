package entities

import (
	"errors"
	"time"
)

// Time is a type alias so we could override methods and add our own
type Time struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
// Time is formatted as RFC3339 instead of RFC3339Nano (original time.Time)
func (t Time) MarshalJSON() ([]byte, error) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(t.In(loc).Format(`"` + time.RFC3339 + `"`)), nil
}

// ParseTime parses date string of RFC3339 format
func ParseTime(s string) (Time, error) {
	var res Time

	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		res = Time{t}
	}

	return res, err
}
