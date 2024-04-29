package main

import "github.com/enetx/g"

func main() {
	b := g.NewBytes("test rest foo bar")

	b.Split().Map(g.Bytes.Upper).Collect().Print()
	b.Fields().Map(g.Bytes.Upper).Collect().Print()
}
