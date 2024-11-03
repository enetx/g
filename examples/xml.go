package main

import (
	"encoding/xml"
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	type Plant struct {
		XMLName xml.Name      `xml:"plant"`
		ID      Int           `xml:"id,attr"`
		Name    String        `xml:"name"`
		Origin  Slice[string] `xml:"origin"`
	}

	coffee := &Plant{ID: 27, Name: "Coffee"}
	coffee.Origin = SliceOf("Ethiopia", "Brazil")

	s := NewString("").Encode().XML(coffee, "", "  ").Unwrap().Append("\n")
	fmt.Println(s)

	var coffee2 Plant

	s.Decode().XML(&coffee2)
	fmt.Println(coffee2.Origin.Get(-1))
}
