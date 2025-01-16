package main

import . "github.com/enetx/g"

func main() {
	foo := String("foo")
	bar := "bar"

	Format(String("foo: {foo}, bar: {bar}"), map[string]any{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo}, bar: {bar}", map[String]any{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo}, bar: {bar}", Map[string, any]{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo}, bar: {bar}", Map[String, any]{"foo": foo, "bar": bar}).Println()

	name := "John"
	age := 30
	city := "New York"

	named := map[string]any{
		"name": name,
		"age":  age,
		"city": city,
	}

	f := Format("Hello, my name is {name}. I am {age} years old and live in {city}.", named)
	f.Println()
}
