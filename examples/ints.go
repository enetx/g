package main

import . "github.com/enetx/g"

func main() {
	NewInt(100).Random().Print()

	NewInt(-10).RandomRange(-5).Print()
	NewInt(-10).RandomRange(5).Print()

	NewInt(11).Max(10, 32, 11, 33, 908).Print()
	NewInt(11).Min(-1, 32, 11, 33, 908).Print()

	NewInt(97).Binary().Print()
	NewInt('a').Binary().Print()
	NewInt(byte('a')).Binary().Print()
}
