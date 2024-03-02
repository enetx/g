package main

import (
	"gitlab.com/x0xO/g"
)

func main() {
	foo := g.String("foo")
	bar := "bar"

	g.GSprintf("foo: {foo}\nbar: {bar}", map[string]any{"foo": foo, "bar": bar}).Print()
	g.GSprintf("foo: {foo}\nbar: {bar}", map[g.String]any{"foo": foo, "bar": bar}).Print()
	g.GSprintf("foo: {foo}\nbar: {bar}", g.Map[string, any]{"foo": foo, "bar": bar}).Print()
	g.GSprintf("foo: {foo}\nbar: {bar}", g.Map[g.String, any]{"foo": foo, "bar": bar}).Print()

	values := map[string]any{
		"Name": "John",
		"Age":  30,
		"City": "New York",
	}

	formatString := "Hello, my name is {Name}. I am {Age} years old and live in {City}."
	formattedString := g.GSprintf(formatString, values)
	formattedString.Print()
}
