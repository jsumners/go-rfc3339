package rfc3339

import (
	"github.com/spf13/cast"
	"strings"
)

// nsToInt converts a fractional second string, e.g. `.005`, to an integer that
// is acceptable by [time.Date]. In short, [time.Date] requires nanoseconds to
// be an integer up to 9 digits wide. Any input string that is wider than
// 10 characters, including the leading `.`, will be chopped to 10 characters.
//
// See the following snippets for the core implementation:
// + https://cs.opensource.google/go/go/+/refs/tags/go1.20.2:src/time/format_go;l=118-126
// + https://cs.opensource.google/go/go/+/refs/tags/go1.20.2:src/time/format.go;l=1494-1517;drc=06264b740e3bfe619f5e90359d8f0d521bd47806
func nsToInt(input string) int {
	toConvert := input
	if len(input) > 10 {
		toConvert = toConvert[0:10]
	}
	toConvert = strings.TrimPrefix(toConvert, ".")

	result := toInt(toConvert)
	for i := 0; i < 10-len(toConvert)-1; i += 1 {
		result = result * 10
	}
	return result
}

// toInt converts a string to an integer. Leading `0` characters will be
// trimmed before converting.
func toInt(input string) int {
	return cast.ToInt(strings.TrimPrefix(input, "0"))
}
