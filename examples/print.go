package main

import (
	"time"

	. "github.com/enetx/g"
)

func main() {
	// 1) Basic named placeholders

	foo := "foo"
	bar := "bar"

	Printf("foo: {foo}, bar: {bar}\n", Named{"foo": foo, "bar": bar})

	// 2) Named placeholders with fallback and multiple modifiers

	// We define a map with "name", "age", and "city".
	// We'll illustrate fallback ( {name?noname} ), plus various modifiers
	// like $trim, $replace, $title, $substring, $date, etc.
	name := String("   john  ")
	age := 30
	city := "New York"

	named := Named{
		"name":  name,
		"age":   age,
		"city":  city,
		"today": time.Now(),
	}

	// Template example:
	// {name.$trim.$replace(j,r).$title.$substring(0,-2).$replace(o,0)}0t
	// This will:
	//   - Trim whitespace around "name"
	//   - Replace 'j' with 'r'
	//   - Convert to Title
	//   - Substring from 0 to length-2 (e.g. "Jo" => "J")
	//   - Replace 'o' with '0'
	//   - Then append "0t" as literal text
	// Output: Hello, my name is R00t. I am 30 years old and live in Ne....
	Printf(
		"Hello, my name is {name.$trim.$replace(j,r).$title.$substring(0,-2).$replace(o,0)}0t. I am {age} years old and live in {city.$truncate(2)}.\n",
		named,
	)

	// Another variant with fallback {unknown?name} and a date modifier
	// Output: Today is 01/19/2025. Name fallback example: JOHN
	Printf("Today is {today.$date(01/02/2006)}. Name fallback example: {unknown?name.$trim.$upper}\n", named)

	// 3) Mixing autoindex placeholders with named placeholders

	// Output: Numeric: positional-1, Named: {named:value}, Another numeric: {POSITIONAL-2}
	Printf(
		"Numeric: {}, Named: {key.$fmt(%+v)}, Another numeric: \\{{.$upper}\\}\n",
		Named{"key": struct{ named string }{named: "value"}},
		"positional-1", // => {1}
		"positional-2", // => {2}
	)

	// 4) Numeric-only usage

	// If you only have numeric placeholders, you can simply pass arguments
	// Output: Hello + 123 + World + Hello
	Printf("{1} + {2} + {3} + {1}", "Hello", 123, "World")
}
