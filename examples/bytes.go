package main

import . "github.com/enetx/g"

func main() {
	b := NewBytes("test rest foo bar")

	b.Split().Map(Bytes.Upper).Collect().Println()
	b.Fields().Map(Bytes.Upper).Collect().Println()
}
