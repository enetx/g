package main

import . "github.com/enetx/g"

func main() {
	b := Bytes("test rest foo bar")
	b.Println()

	SliceOf(b.Chars()...).Iter().Map(Bytes.Upper).Collect().Slice().Join().Println()
	SliceOf(b.Fields()...).Iter().Map(Bytes.Upper).Collect().Slice().Join().Println()
}
