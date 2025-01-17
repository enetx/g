package main

import . "github.com/enetx/g"

func main() {
	foo := String("foo")
	bar := "bar"

	Format(String("foo: {foo}, bar: {bar}"), map[string]any{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo}, bar: {bar}", map[String]any{"foo": foo, "bar": bar}).Println()
	Format("foo: {$upper:foo}, bar: {bar}", Map[string, any]{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo}, bar: {bar}", Map[String, any]{"foo": foo, "bar": bar}).Println()

	name := String("   john  ")
	age := 30
	city := "New York"

	named := map[string]any{
		"name":   name,
		"age":    age,
		"city":   city,
		"noname": "       no name         ",
	}

	f := Format(
		"Hello, my name is {$trim:name}. I am {age} years old and live in {city}.",
		// "Hello, my name is {$trim.$upper:name}. I am {age} years old and live in {city}.",
		// "Hello, my name is {$trim.$rot13.$upper:name?noname}. I am {age} years old and live in {city}.",
		named,
	)

	f.Println()
}
