# RFC3339

A library that provides simplicity around parsing [RFC3339][3339] `date-time`
and `full-date` representations. The standard `time.Time` library is a bit
rough in that formats must be known ahead of time, and supplied when serializing
to strings. The types provided by this library extend the standard `time.Time`
type.

[3339]: https://www.rfc-editor.org/rfc/rfc3339

## Install

```sh
$ go get github.com/jsumners/go-rfc3339
```

## Example

```go
package main

import (
	"fmt"
	"github.com/jsumners/go-rfc3339"
)

func main() {
	dt, err := rfc3339.NewDateTimeFromString("2023-04-04T12:45:00-04:00")
	if err != nil {
		panic(err)
	}

	fmt.Printf("timestamp: %s", dt.ToString())
	// Output:
	// timestamp: 2023-04-04T12:45:00-04:00
}
```
