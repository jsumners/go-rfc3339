# RFC3339

A library that provides simplicity around parsing [RFC3339][3339] `date-time`
and `full-date` representations. The standard `time.Time` library is a bit
rough in that formats must be known ahead of time, and supplied when serializing
to strings. The types provided by this library extend the standard `time.Time`
type.

The `date-time` and `full-date` types implement the [Scanner][scanner] and
[Valuer][valuer] interfaces so that they can be stored as strings in a
database.

[3339]: https://www.rfc-editor.org/rfc/rfc3339
[scanner]: https://pkg.go.dev/database/sql#Scanner
[valuer]: https://pkg.go.dev/database/sql/driver#Valuer

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

## Example (sql)

```go
package main

import (
	"database/sql"
	"fmt"

	"github.com/jsumners/go-rfc3339"
	_ "modernc.org/sqlite"
)

func main() {
	db, _ := sql.Open("sqlite", "./db.sqlite")
	defer db.Close()

	dt, _ := rfc3339.NewDateTimeFromString("2023-09-28T08:30:00-04:00")
	db.Exec(
		"insert into example (date_time) values (@date_time)",
		sql.Named("date_time", dt),
	)

	rows, _ := db.Query("select date_time from values")
	for rows.Next() {
		var dt rfc3339.DateTime
		rows.Scan(&dt)
		fmt.Printf("dt = %+v\n", dt)
	}
}
```
