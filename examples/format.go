package main

import (
	"gitlab.com/x0xO/g"
)

func main() {
	foo := g.String("foo")
	bar := "bar"

	g.GSprintf(g.String("foo: {foo}\nbar: {bar}"), map[string]any{"foo": foo, "bar": bar}).Print()
	g.GSprintf("foo: {foo}\nbar: {bar}", map[g.String]any{"foo": foo, "bar": bar}).Print()
	g.GSprintf("foo: {foo}\nbar: {bar}", g.Map[string, any]{"foo": foo, "bar": bar}).Print()
	g.GSprintf("foo: {foo}\nbar: {bar}", g.Map[g.String, any]{"foo": foo, "bar": bar}).Print()

	name := "John"
	age := 30
	city := "New York"

	named := map[string]any{
		"name": name,
		"age":  age,
		"city": city,
	}

	f := g.GSprintf("Hello, my name is {name}. I am {age} years old and live in {city}.", named)
	f.Print()
}
