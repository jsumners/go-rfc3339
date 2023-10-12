package rfc3339

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jsumners/go-reggie"
)

var fullDateRegex = reggie.MustCompile(
	`^(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})$`,
)

// IsFullDateString verifies if an input string matches the format of
// an RFC 3339 `full-date` representation.
func IsFullDateString(input string) bool {
	return fullDateRegex.MatchString(input)
}

// MustParseDateString wraps [NewFullDateFromString] such that if an error
// happens it generates a panic.
func MustParseDateString(input string) FullDate {
	fd, err := NewFullDateFromString(input)
	if err != nil {
		panic(err)
	}
	return fd
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

// Value implements the [driver.Valuer] interface to facilitate
// storing [FullDate] values as strings in a database.
func (fd FullDate) Value() (driver.Value, error) {
	return fd.ToString(), nil
}

// Scan implements the [sql.Scanner] interface to facilitate reading
// [FullDate] strings stored in a database.
func (fd *FullDate) Scan(value any) error {
	if value == nil {
		return nil
	}

	rv := reflect.TypeOf(value)
	if rv.Name() != "string" {
		return fmt.Errorf("value must be a string, got: %s", rv.Name())
	}

	// ConvertValue always coerces to a string and does not return
	// an error.
	str, _ := driver.String.ConvertValue(value)
	if str == "" {
		*fd = FullDate{}
		return nil
	}

	parsed, err := NewFullDateFromString(str.(string))
	if err != nil {
		return err
	}

	*fd = parsed
	return nil
}
