package rfc3339

import (
	"fmt"
	"github.com/jsumners/go-reggie"
	"strings"
	"time"
)

var fullDateRegex = reggie.MustCompile(
	`^(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})$`,
)

// IsFullDateString verifies if an input string matches the format of
// an RFC 3339 `full-date` representation.
func IsFullDateString(input string) bool {
	return fullDateRegex.MatchString(input)
}

// NewFullDateFromString creates a new [FullDate] instance from an RFC 3339
// `full-date` string representation. Note that the time parts will be set to
// 00:00:00.000 at the UTC (+00:00) offset.
func NewFullDateFromString(input string) (FullDate, error) {
	matches := fullDateRegex.FindStringSubmatch(input)

	if matches == nil {
		return FullDate{}, fmt.Errorf("`%s` is not a full-date string", input)
	}

	year := fullDateRegex.SubmatchWithName("year")
	month := fullDateRegex.SubmatchWithName("month")
	day := fullDateRegex.SubmatchWithName("day")

	d := time.Date(
		toInt(year), time.Month(toInt(month)), toInt(day),
		0, 0, 0, 0,
		time.FixedZone("UTC", 0),
	)

	return FullDate{Time: d}, nil
}

// ToString serializes the [FullDate] instance to an RFC 3339 full-date
// string representation.
func (fd FullDate) ToString() string {
	return fmt.Sprintf("%04d-%02d-%02d", fd.Year(), fd.Month(), fd.Day())
}

func (fd FullDate) MarshalJSON() ([]byte, error) {
	if fd.IsZero() {
		return []byte("null"), nil
	}
	serialized := fmt.Sprintf(`"%s"`, fd.ToString())
	return []byte(serialized), nil
}

func (fd *FullDate) UnmarshalJSON(data []byte) error {
	timeStr := strings.Trim(string(data), `"`)
	if timeStr == "null" || timeStr == "" {
		return nil
	}

	d, err := NewFullDateFromString(timeStr)
	if err != nil {
		return err
	}

	fd.Time = d.Time

	return nil
}
