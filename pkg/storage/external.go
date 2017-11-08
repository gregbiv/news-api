package storage

import (
	"encoding/json"
	"fmt"
	"time"
)

// JSONDate a simple type to ensure json date is correctly formatted
type JSONDate time.Time

// MarshalJSON time to make this into json
func (t JSONDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02"))
	return []byte(stamp), nil
}

// UnmarshalJSON time to get a date from json
func (t *JSONDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	// so it accepts an empty string as an empty time.Time
	if s == "" {
		return nil
	}

	date, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*t = JSONDate(date)

	return nil
}

//Time transforms the string representation into a time.Time object
func (t JSONDate) Time() time.Time {
	return time.Time(t)
}
