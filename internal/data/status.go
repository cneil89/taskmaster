package data

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Status int

const (
	DEFINING Status = iota
	READY
	INPROGRESS
	UNDERREVIEW
	COMPLETED
)

var stringToStatus = map[string]Status{
	"defining":     DEFINING,
	"ready":        READY,
	"in progress":  INPROGRESS,
	"under review": UNDERREVIEW,
	"completed":    COMPLETED,
}

var statusToString = [...]string{"Defining", "Ready", "In Progress", "Under Review", "Completed"}

func ParseStatus(v string) (Status, error) {
	s, ok := stringToStatus[strings.ToLower(strings.TrimSpace(v))]
	if !ok {
		return 0, fmt.Errorf("invalid status %q", v)
	}
	return s, nil
}

func (s Status) IsValid() bool {
	return s >= DEFINING && s <= COMPLETED
}

func (s Status) String() string {
	if !s.IsValid() {
		return "unknown"
	}
	return statusToString[s]
}

// For database/sql writes
func (s Status) Value() (driver.Value, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid status: %d", s)
	}
	return s.String(), nil
}

// For database/sql reads
func (s *Status) Scan(value any) error {
	if value == nil {
		return fmt.Errorf("status cannot be null")
	}

	switch v := value.(type) {
	case string:
		parsed, err := ParseStatus(v)
		if err != nil {
			return err
		}
		*s = parsed
		return nil

	case []byte:
		parsed, err := ParseStatus(string(v))
		if err != nil {
			return err
		}
		*s = parsed
		return nil

	default:
		return fmt.Errorf("cannot scan %T into Status", value)
	}
}
