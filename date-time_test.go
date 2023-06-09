package rfc3339

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDateTime_IsDateTimeString(t *testing.T) {
	t.Run("returns true", func(t *testing.T) {
		result := IsDateTimeString("2023-03-31T16:30:00-04:00")
		assert.True(t, result)
	})

	t.Run("returns false", func(t *testing.T) {
		result := IsDateTimeString("2023-03-31T16:30-04:00")
		assert.False(t, result)
	})
}

func TestDateTime_NewFromString(t *testing.T) {
	t.Run("returns error if string does not conform", func(t *testing.T) {
		dt, err := NewDateTimeFromString("2023-03-24")
		assert.Empty(t, dt)
		assert.Errorf(t, err, "input is not a date-time string: %s", "2023-03-24")
	})

	t.Run("returns object at correct offset", func(t *testing.T) {
		// This test covers a case where zero padded integers trigger
		// a bad path in the conversion.
		// See https://github.com/spf13/cast/issues/182.
		dt, err := NewDateTimeFromString("2023-04-01T08:30:00-04:00")
		expected := time.Date(
			2023,
			time.April,
			1,
			8,
			30,
			0,
			0,
			time.FixedZone("UTC-04:00", -14400),
		)

		assert.NoError(t, err)
		assert.Equal(t, expected, dt.Time)
	})

	t.Run("parses string without fractional second", func(t *testing.T) {
		dt, err := NewDateTimeFromString("2023-03-24T22:30:00-04:00")
		require.NoError(t, err)

		assert := assert.New(t)
		assert.Equal(2023, dt.Year())
		assert.Equal(time.March, dt.Month())
		assert.Equal(24, dt.Day())
		assert.Equal(22, dt.Hour())
		assert.Equal(30, dt.Minute())
		assert.Equal(0, dt.Second())
		assert.Equal(0, dt.Nanosecond())
		_, offset := dt.Zone()
		assert.Equal(-14400, offset)
	})

	t.Run("parses string at Z offset", func(t *testing.T) {
		dt, err := NewDateTimeFromString("2023-03-24T22:30:00.005Z")
		require.NoError(t, err)

		assert := assert.New(t)
		assert.Equal(2023, dt.Year())
		assert.Equal(time.March, dt.Month())
		assert.Equal(24, dt.Day())
		assert.Equal(22, dt.Hour())
		assert.Equal(30, dt.Minute())
		assert.Equal(0, dt.Second())
		assert.Equal(5000000, dt.Nanosecond())
		_, offset := dt.Zone()
		assert.Equal(0, offset)
	})
}

func TestDateTime_ToString(t *testing.T) {
	expected := "2023-03-24T22:30:00.005Z"
	dt, err := NewDateTimeFromString(expected)
	require.NoError(t, err)
	assert.Equal(t, expected, dt.ToString())
}

func TestDateTime_MarshalJSON(t *testing.T) {
	type j struct {
		Created DateTime `json:"created"`
	}

	t.Run("returns null for empty value", func(t *testing.T) {
		js := j{Created: DateTime{}}
		result, err := json.Marshal(js)
		assert.NoError(t, err)
		assert.Equal(t, `{"created":null}`, string(result))
	})

	t.Run("serializes to expected strings", func(t *testing.T) {
		type testData struct {
			input    DateTime
			expected string
		}

		tests := []testData{
			{
				NewFromTime(time.Date(2023, 4, 1, 11, 45, 0, 0, time.FixedZone("", -14400))),
				"2023-04-01T11:45:00-04:00",
			},
			{
				NewFromTime(time.Date(2023, 4, 1, 11, 45, 0, 5000000, time.FixedZone("", -14400))),
				"2023-04-01T11:45:00.005-04:00",
			},
		}

		for _, test := range tests {
			expected := fmt.Sprintf(`{"created":"%s"}`, test.expected)
			js := j{Created: test.input}
			result, err := json.Marshal(js)
			assert.NoError(t, err)
			assert.Equal(t, expected, string(result))
		}
	})
}

func TestDateTime_UnmarshalJSON(t *testing.T) {
	type testJson struct {
		Created DateTime `json:"created"`
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
		input := `{"created":"2023-04-01 12:00:00-04:00"}`

		var result testJson
		err := json.Unmarshal([]byte(input), &result)

		assert.ErrorContains(t, err, "input is not a date-time string")
		assert.True(t, result.Created.IsZero())
	})

	t.Run("unmarshals full string", func(t *testing.T) {
		input := `{"created":"2023-04-01T08:30:00-04:00"}`

		var result testJson
		err := json.Unmarshal([]byte(input), &result)

		assert.NoError(t, err)
		assert.Equal(
			t,
			"2023-04-01T08:30:00-04:00",
			result.Created.Format(time.RFC3339),
		)
	})
}
