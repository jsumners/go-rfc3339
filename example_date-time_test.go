package rfc3339_test

import (
	"encoding/json"
	"fmt"
	"github.com/jsumners/go-rfc3339"
)

type MyJson struct {
	Created rfc3339.DateTime `json:"created"`
}

// This example shows how a [rfc3339.DateTime] can be used when
// deserializing JSON data.
func Example_datetime() {
	input := `{"created":"2023-04-04T12:30:00-04:00"}`

	var myJson MyJson
	json.Unmarshal([]byte(input), &myJson)

	fmt.Printf("%v\n", myJson)
	// Output:
	// {2023-04-04 12:30:00 -0400 UTC-04:00}
}
