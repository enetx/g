package main

import . "github.com/enetx/g"

func main() {
	b := Bytes("test rest foo bar")
	b.Println()

	b.Chars().Map(Bytes.Upper).Collect().Join().Println()
	b.Fields().Map(Bytes.Upper).Collect().Join().Println()
}
