package main

import (
	"time"

	. "github.com/enetx/g"
)

func main() {
	foo := String("foo")
	bar := "bar"

	Format(String("foo: {foo}, bar: {bar}"), map[string]any{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo}, bar: {bar}", map[String]any{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo.$upper}, bar: {bar}", Map[string, any]{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo}, bar: {bar}", Map[String, any]{"foo": foo, "bar": bar}).Println()

	name := String("   john  ")
	age := 30
	city := "New York"

	named := map[string]any{
		"name":   name,
		"age":    age,
		"city":   city,
		"noname": "       no name         ",
		"today":  time.Now(),
	}

	f := Format(
		// "Hello, my name is {name.$trim}. I am {age} years old and live in {city}.",
		// "Hello, my name is {name.$trim.$upper}. I am {age} years old and live in {city}.",
		"Hello, my name is {name?noname.$trim.$upper}. I am {age} years old and live in {city}. Today {today.$format(01/02/2006)}.",
		named,
	)

	f.Println()

	///////////////////////////////////////////////////////////////////////////

	handlers := Map[String, func(v any, args ...String) any]{
		"$double": func(v any, _ ...String) any { return (v.(Int) * 2).String() },
		"$prefix": func(v any, _ ...String) any { return "prefix_" + v.(String) },
	}

	args := map[string]any{
		"value": Int(42),
		"text":  String("example"),
		"date":  time.Now(),
	}

	format := "{value.$double} and {text.$upper} {date.$format(01-02-2006)}"
	result := Format(format, args, handlers)
	result.Println() // Output: "84 and EXAMPLE 01-17-2025"
}
