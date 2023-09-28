package rfc3339

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFullDate_IsFullDateString(t *testing.T) {
	t.Run("returns true for date string", func(t *testing.T) {
		result := IsFullDateString("2023-04-01")
		assert.True(t, result)
	})

	t.Run("returns false if not a date string", func(t *testing.T) {
		result := IsFullDateString("04/01/2023")
		assert.False(t, result)
	})
}

func TestFullDate_NewFromString(t *testing.T) {
	t.Run("returns error for bad input", func(t *testing.T) {
		fd, err := NewFullDateFromString("2023/04/03")
		assert.Empty(t, fd)
		assert.ErrorContains(t, err, "is not a full-date string")
	})

	t.Run("returns a full-time", func(t *testing.T) {
		fd, err := NewFullDateFromString("2023-04-03")
		assert.NoError(t, err)

		expected := time.Date(
			2023, 4, 3,
			0, 0, 0, 0,
			time.FixedZone("UTC", 0),
		)
		assert.Equal(t, expected, fd.Time)
	})
}

func TestFullDate_ToString(t *testing.T) {
	t.Run("formats a current full-date", func(t *testing.T) {
		fd, err := NewFullDateFromString("2023-04-03")
		assert.NoError(t, err)
		assert.Equal(t, "2023-04-03", fd.ToString())
	})

	t.Run("formats the zero date", func(t *testing.T) {
		fd, err := NewFullDateFromString("0001-01-01")
		assert.NoError(t, err)
		assert.Equal(t, "0001-01-01", fd.ToString())
	})
}

func TestFullDate_MarshalJSON(t *testing.T) {
	type j struct {
		Created FullDate `json:"created"`
	}

	t.Run("returns null for empty value", func(t *testing.T) {
		js := j{Created: FullDate{}}
		result, err := json.Marshal(js)
		assert.NoError(t, err)
		assert.Equal(t, `{"created":null}`, string(result))
	})

	t.Run("serializes to expected strings", func(t *testing.T) {
		input, _ := NewFullDateFromString("2023-04-04")
		expected := `{"created":"2023-04-04"}`
		js := j{Created: input}
		result, err := json.Marshal(js)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(result))
	})
}

func TestFullDate_UnmarshalJSON(t *testing.T) {
	type testJson struct {
		Created FullDate `json:"created"`
	}

	t.Run("returns nil for null", func(t *testing.T) {
		input := `{"created":null}`

		var result testJson
		err := json.Unmarshal([]byte(input), &result)

		assert.NoError(t, err)
		assert.True(t, result.Created.IsZero())
	})

	t.Run("returns nil for empty string", func(t *testing.T) {
		input := `{"created":""}`

		var result testJson
		err := json.Unmarshal([]byte(input), &result)

		assert.NoError(t, err)
		assert.True(t, result.Created.IsZero())
	})

	t.Run("returns error for bad input", func(t *testing.T) {
		input := `{"created":"2023/04/01"}`

		var result testJson
		err := json.Unmarshal([]byte(input), &result)

		assert.ErrorContains(t, err, "is not a full-date string")
		assert.True(t, result.Created.IsZero())
	})

	t.Run("unmarshals full string", func(t *testing.T) {
		input := `{"created":"2023-04-01"}`

		var result testJson
		err := json.Unmarshal([]byte(input), &result)

		assert.NoError(t, err)
		assert.Equal(
			t,
			"2023-04-01T00:00:00Z",
			result.Created.Format(time.RFC3339),
		)
	})
}

func Test_FDValue(t *testing.T) {
	fd, _ := NewFullDateFromString("2023-09-28")
	str, err := fd.Value()
	assert.Nil(t, err)
	assert.Equal(t, "2023-09-28", str)
}

func Test_FDScan(t *testing.T) {
	t.Run("handles nil input", func(t *testing.T) {
		fd := FullDate{}
		err := fd.Scan(nil)
		assert.Nil(t, err)
	})

	t.Run("only scans strings", func(t *testing.T) {
		fd := FullDate{}
		err := fd.Scan(42)
		assert.ErrorContains(t, err, "value must be a string, got: int")
	})

	t.Run("empty instance is empty", func(t *testing.T) {
		expected := FullDate{}
		source, _ := NewFullDateFromString("2023-09-28")

		err := source.Scan("")
		assert.Nil(t, err)
		assert.Equal(t, expected, source)
	})

	t.Run("returns full-date parse error", func(t *testing.T) {
		fd := FullDate{}
		err := fd.Scan("2023/09/28")
		assert.ErrorContains(t, err, "is not a full-date string")
	})

	t.Run("scans correctly", func(t *testing.T) {
		fd := FullDate{}
		err := fd.Scan("2023-09-28")
		assert.Nil(t, err)
		assert.Equal(t, "2023-09-28", fd.ToString())
	})
}
