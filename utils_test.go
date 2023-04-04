package rfc3339

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestNsToInt(t *testing.T) {
	t.Run("converts fractions to ints", func(t *testing.T) {
		type testData struct {
			input    string
			expected int
		}
		tests := []testData{
			{input: ".005", expected: 5000000},
			{input: ".5", expected: 500000000},
			{input: ".000000005", expected: 5},
			{input: ".1234567891", expected: 123456789},
			{input: "005", expected: 5000000}, // no leading `.`
		}

		for _, test := range tests {
			result := nsToInt(test.input)
			assert.Equal(t, test.expected, result)

			// Verify that our conversion matches the standard library's conversion:
			targetDate, _ := time.Parse(time.RFC3339Nano,
				fmt.Sprintf(
					"2023-03-28T12:00:00.%s-04:00",
					strings.TrimPrefix(test.input, "."),
				),
			)
			assert.Equal(t, targetDate.Nanosecond(), result)
		}
	})

	t.Run("chops wide strings", func(t *testing.T) {
		result := nsToInt(".123456789654321")
		assert.Equal(t, 123456789, result)
	})
}

func TestToInt(t *testing.T) {
	tests := [][]interface{}{
		{"08", 8},
		{"0005", 5},
	}

	for _, test := range tests {
		result := toInt(test[0].(string))
		assert.Equal(t, test[1].(int), result)
	}
}
