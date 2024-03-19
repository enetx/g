package main

import (
	"encoding/xml"
	"fmt"

	"github.com/enetx/g"
)

func main() {
	type Plant struct {
		XMLName xml.Name        `xml:"plant"`
		ID      g.Int           `xml:"id,attr"`
		Name    g.String        `xml:"name"`
		Origin  g.Slice[string] `xml:"origin"`
	}

	coffee := &Plant{ID: 27, Name: "Coffee"}
	coffee.Origin = g.SliceOf("Ethiopia", "Brazil")

	s := g.NewString("").Enc().XML(coffee, "", "  ").Unwrap().Add("\n")
	fmt.Println(s)

	var coffee2 Plant

	s.Dec().XML(&coffee2)
	fmt.Println(coffee2.Origin.Get(-1))
}
