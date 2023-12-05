package main

import "gitlab.com/x0xO/g"

func main() {
	d := g.NewDir(".").Copy("copy").Unwrap()

	d.Path().Unwrap().Print()
}
