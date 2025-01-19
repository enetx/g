package main

import (
	"time"

	. "github.com/enetx/g"
)

func main() {
	// 1) Basic named placeholders

	// Here we have a template with two named placeholders: {foo} and {bar}.
	// The maps can be passed in any of the supported forms (map[string]any, map[String]any, etc.).
	foo := String("foo")
	bar := "bar"

	Format(String("foo: {foo}, bar: {bar}"), map[string]any{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo}, bar: {bar}", map[String]any{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo.$upper}, bar: {bar}", Map[string, any]{"foo": foo, "bar": bar}).Println()
	Format("foo: {foo}, bar: {bar}", Map[String, any]{"foo": foo, "bar": bar}).Println()

	// 2) Named placeholders with fallback and multiple modifiers

	// We define a map with "name", "age", and "city".
	// We'll illustrate fallback ( {name?noname} ), plus various modifiers
	// like $trim, $replace, $title, $substring, $date, etc.
	name := String("   john  ")
	age := 30
	city := "New York"

	named := map[string]any{
		"name":  name,
		"age":   age,
		"city":  city,
		"today": time.Now(),
	}

	// Template example:
	// {name.$trim.$replace(j,r).$title.$substring(0,-2)}ot
	// This will:
	//   - Trim whitespace around "name"
	//   - Replace 'j' with 'r'
	//   - Convert to Title
	//   - Substring from 0 to length-2 (e.g. "Jo" => "J")
	//   - Then append "ot" as literal text
	Format(
		"Hello, my name is {name.$trim.$replace(j,r).$title.$substring(0,-2).$replace(o,0)}0t. I am {age} years old and live in {city.$truncate(2)}.",
		named,
	).
		Println() // Output: Hello, my name is Root. I am 30 years old and live in Ne....

	// Another variant with fallback {unknown?name} and a date modifier
	Format("Today is {today.$date(01/02/2006)}. Name fallback example: {unknown?name.$trim.$upper}", named).
		Println() // Output: Today is 01/19/2025. Name fallback example: JOHN

	// 3) Mixing numeric placeholders with named placeholders

	// Numeric placeholders {1}, {2}, etc. use 1-based indexing for non-map arguments.
	// Named placeholders still come from any map arguments.
	Format(
		"Numeric: {1}, Named: {key}, Another numeric: {2}",
		map[string]any{"key": "named-value"},
		"positional-1", // => {1}
		"positional-2", // => {2}
	).
		Println() // Output: Numeric: positional-1, Named: named-value, Another numeric: positional-2

	// 4) Numeric-only usage

	// If you only have numeric placeholders, you can simply pass arguments
	Format("{1} + {2} + {3} + {1}", "Hello", 123, "World").
		Println() // Output: Hello + 123 + World + Hello
}
