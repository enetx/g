package main

import "gitlab.com/x0xO/g"

func main() {
	g.NewInt(100).Random().Print()

	g.NewInt(-10).RandomRange(-5).Print()
	g.NewInt(-10).RandomRange(5).Print()

	g.NewInt(11).Max(10, 32, 11, 33, 908).Print()
	g.NewInt(11).Min(-1, 32, 11, 33, 908).Print()

	g.NewInt(97).ToBinary().Print()
	g.NewInt('a').ToBinary().Print()
	g.NewInt(byte('a')).ToBinary().Print()
}
