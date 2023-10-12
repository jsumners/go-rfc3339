package rfc3339

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jsumners/go-reggie"
)

var regexParts = []string{
	`^(?P<year>\d{4})`, "-",
	`(?P<month>\d{2})`, "-",
	`(?P<day>\d{2})`,
	`(?P<separator>[tT])`,
	`(?P<hour>\d{2})`, ":",
	`(?P<minute>\d{2})`, ":",
	`(?P<second>\d{2})`,
	`(?P<secfrac>\.\d+)?`,
	"(",
	`(?P<offsetTime>[+-]\d{2}:\d{2})`,
	"|",
	`(?P<offsetZ>[zZ])`,
	")$",
}

var dateTimeRegex = reggie.MustCompile(
	strings.Join(regexParts, ""),
)

// IsDateTimeString verifies if an input string matches the format
// of an RFC 3339 `date-time` representation.
func IsDateTimeString(input string) bool {
	return dateTimeRegex.MatchString(input)
}

// NewDateTimeFromString creates a new [DateTime] instance from an RFC 3339
// `date-time`  string representation. Note that the maximum precision of
// fractional seconds is limited to 9 places. This is due to [time.Date]'s
// implementation of fractional seconds (basically, it supports a floating-point
// exponent of 10^9).
func NewDateTimeFromString(input string) (DateTime, error) {
	matches := dateTimeRegex.FindStringSubmatch(input)

	if matches == nil {
		return DateTime{}, fmt.Errorf("input is not a date-time string: %s", input)
	}

	year := dateTimeRegex.SubmatchWithName("year")
	month := dateTimeRegex.SubmatchWithName("month")
	day := dateTimeRegex.SubmatchWithName("day")
	hour := dateTimeRegex.SubmatchWithName("hour")
	minute := dateTimeRegex.SubmatchWithName("minute")
	second := dateTimeRegex.SubmatchWithName("second")

	var secFrac = 0
	secFracString := dateTimeRegex.SubmatchWithName("secfrac")
	if secFracString != "" {
		secFrac = nsToInt(secFracString)
	}

	offsetTime := dateTimeRegex.SubmatchWithName("offsetTime")
	if offsetTime == "" {
		offsetZString := dateTimeRegex.SubmatchWithName("offsetZ")
		if offsetZString == "Z" {
			offsetTime = "+00:00"
		}
	}

	// Calculate the zone UTC offset in seconds from the provided
	// offset portion of the string, and convert it into a UTC offset zone
	// name (https://en.wikipedia.org/wiki/UTC_offset).
	zHour := toInt(offsetTime[1:3])
	zMin := toInt(offsetTime[4:6])
	zOffset := (zHour*60 + zMin) * 60
	if offsetTime[0] == '-' {
		zOffset = -1 * zOffset
	}
	zoneName := fmt.Sprintf("UTC%s", offsetTime)

	// Use the UTC offset zone name to generate a [time.Location] for use
	// in creating the new instance.
	var location *time.Location
	if zoneName == "UTC+00:00" {
		location = time.UTC
	} else {
		location = time.FixedZone(zoneName, zOffset)
	}

	date := time.Date(
		toInt(year), time.Month(toInt(month)), toInt(day),
		toInt(hour), toInt(minute), toInt(second), secFrac,
		location,
	)
	dt := DateTime{
		Time: date,
	}

	return dt, nil
}

// NewFromTime wraps the `time` instance as an RFC 3339 [DateTime].
func NewFromTime(time time.Time) DateTime {
	return DateTime{Time: time}
}

// ToString serializes the [DateTime] instance to a full RFC 3339 date-time
// string representation.
func (dt DateTime) ToString() string {
	return dt.Format(time.RFC3339Nano)
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return []byte("null"), nil
	}
	serialized := dt.Format(time.RFC3339Nano)
	serialized = fmt.Sprintf(`"%s"`, serialized)
	return []byte(serialized), nil
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	timeStr := strings.Trim(string(data), `"`)
	if timeStr == "null" || timeStr == "" {
		return nil
	}

	d, err := NewDateTimeFromString(timeStr)
	if err != nil {
		return err
	}

	dt.Time = d.Time

	return nil
}

// Value implements the [driver.Valuer] interface to facilitate
// storing [DateTime] values as strings in a database.
func (dt DateTime) Value() (driver.Value, error) {
	return dt.ToString(), nil
}

// Scan implements the [sql.Scanner] interface to facilitate reading
// [DateTime] strings stored in a database.
func (dt *DateTime) Scan(value any) error {
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
		*dt = DateTime{}
		return nil
	}

	parsed, err := NewDateTimeFromString(str.(string))
	if err != nil {
		return err
	}

	*dt = parsed
	return nil
}
