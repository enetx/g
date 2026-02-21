package main

import (
	"time"

	. "github.com/enetx/g"
)

func main() {
	// --- Auto-index placeholders ---
	// {} consumes the next positional argument in order.
	Println("{} {} {}", "one", "two", "three") // one two three

	// --- Positional placeholders ---
	// {N} references argument N (1-based). The same argument can be reused.
	Println("{2} comes before {1}", "world", "hello") // hello comes before world
	Println("{1.Hex}, {1.Binary}", Int(255))          // ff, 11111111

	// --- Named placeholders ---
	// Pass a Named map alongside positional args.
	named := Named{
		"lang": "Go",
		"year": 2009,
	}
	Println("{lang} was released in {year}.", named) // Go was released in 2009.

	// --- Fallback ---
	// {key?fallback} uses fallback key when the primary key is missing.
	Println("Hello, {user?guest}!", Named{"guest": "stranger"}) // Hello, stranger!

	// --- Escape literal braces ---
	Println("template syntax: \\{{value}\\}", Named{"value": "example"}) // template syntax: {example}

	// --- Modifier chains ---
	// Modifiers call methods on the value via reflection — any method works.
	Println("{.Trim.Title}", String("  hello world  ")) // Hello World
	Println("{.Upper.Reverse}", String("gopher"))       // REHPOG
	Println("{.Repeat(3)}", String("ha"))               // hahaha
	Println("{.Truncate(5)}", String("Hello, World!"))  // Hello...

	// --- Numeric type formatters ---
	Println("dec={} hex={.Hex} bin={.Binary} oct={.Octal}", Int(42), Int(42), Int(42), Int(42))

	// --- Float modifiers ---
	Println("{.Round} / {.RoundDecimal(2)}", Float(3.14159), Float(3.14159)) // 3 / 3.14

	// --- Struct field access ---
	type Point struct{ X, Y int }
	p := Point{X: 10, Y: 20}
	Println("x={1.X} y={1.Y}", p) // x=10 y=20

	// --- Map access ---
	m := Map[String, String]{"key": "value"}
	Println("map[key] = {.key}", m)           // map[key] = value
	Println("map[key] = {.Get(key).Some}", m) // same via method

	// --- Slice access ---
	s := Slice[String]{"alpha", "beta", "gamma"}
	Println("index 0={.0} index 2={.2}", s, s) // index 0=alpha index 2=gamma

	// --- Date formatting via method call ---
	today := Named{"today": time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)}
	Println("today is {today.Format(2006-01-02)}", today) // today is 2025-06-15

	// --- type and debug specifiers ---
	// type → fmt.Sprintf("%T", v), debug → fmt.Sprintf("%#v", v)
	Println("{.type}  {.debug}", Int(99), Int(99)) // g.Int  99

	// --- Mixing named and positional ---
	Println("lang={lang}, arg={}, arg={}",
		Named{"lang": "Go"}, "first", "second") // lang=Go, arg=first, arg=second
}
