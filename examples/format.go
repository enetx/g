package main

import (
	"gitlab.com/x0xO/g"
)

func main() {
	foo := g.String("foo")
	bar := "bar"

	g.Format(g.String("foo: {foo}, bar: {bar}"), map[string]any{"foo": foo, "bar": bar}).Print()
	g.Format("foo: {foo}, bar: {bar}", map[g.String]any{"foo": foo, "bar": bar}).Print()
	g.Format("foo: {foo}, bar: {bar}", g.Map[string, any]{"foo": foo, "bar": bar}).Print()
	g.Format("foo: {foo}, bar: {bar}", g.Map[g.String, any]{"foo": foo, "bar": bar}).Print()

	name := "John"
	age := 30
	city := "New York"

	named := map[string]any{
		"name": name,
		"age":  age,
		"city": city,
	}

	f := g.Format("Hello, my name is {name}. I am {age} years old and live in {city}.", named)
	f.Print()
}
