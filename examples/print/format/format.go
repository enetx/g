package main

import (
	"math"
	"strconv"
	"time"

	. "github.com/enetx/g"
)

// Price demonstrates custom formatting without reflection. FormatValue receives
// the text after ':'; an empty spec is used by {}.
type Price float64

func (p Price) FormatValue(spec String) String {
	value := strconv.FormatFloat(float64(p), 'f', 2, 64)

	switch spec {
	case "", "decimal":
		return String(value)
	case "currency":
		return "$" + String(value)
	case "cents":
		return Int(math.Round(float64(p) * 100)).String()
	default:
		return "<unsupported price format: " + spec + ">"
	}
}

func main() {
	// --- Auto-index placeholders ---
	// {} consumes the next positional argument in order.
	Println("{} {} {}", "one", "two", "three") // one two three

	// --- Positional placeholders ---
	// {N} references argument N (1-based). The same argument can be reused.
	Println("{2} comes before {1}", "world", "hello") // hello comes before world

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
	Println("template syntax: {{{value}}}", Named{"value": "example"}) // template syntax: {example}

	// --- Additional representations ---
	Println("decimal={:d} char={:c} quoted={:q} unicode={:U}", 42, 65, "go\n", 'A')

	// FormatTo appends into an existing reusable builder.
	var output Builder
	FormatTo(&output, "{}={:04d}", "answer", 42)
	output.String().Println()

	// TryFormat validates braces, argument resolution, modifiers, and format specs.
	TryFormat("value={:d}", 42).Inspect(func(s String) { s.Println() })

	// Formattable gives a type its own format specs without reflection.
	price := Price(19.95)
	Println("default={} currency={:currency} cents={:cents}", price, price, price)
	// default=19.95 currency=$19.95 cents=1995

	// --- Transform, then format ---
	// Placeholders resolve data only (fields, keys, indices) — methods are
	// never invoked. Apply transformations to the value
	// before formatting.
	Println("{}", String("  hello world  ").Trim().Title()) // Hello World
	Println("{}", String("gopher").Upper().Reverse())       // REHPOG
	Println("{}", String("ha").Repeat(3))                   // hahaha
	Println("{}", String("Hello, World!").Truncate(5))      // Hello...

	// --- Numeric type formatters ---
	Println("dec={} hex={} bin={} oct={}", Int(42), Int(42).Hex(), Int(42).Binary(), Int(42).Octal())

	// --- Float modifiers ---
	Println("{} / {}", Float(3.14159).Round(), Float(3.14159).RoundDecimal(2)) // 3 / 3.14

	// --- Struct field access ---
	type Point struct{ X, Y int }
	p := Point{X: 10, Y: 20}
	Println("x={1.X} y={1.Y}", p) // x=10 y=20

	// --- Map access ---
	m := Map[String, String]{"key": "value"}
	Println("map[key] = {.key}", m) // map[key] = value

	// --- Slice access ---
	s := Slice[String]{"alpha", "beta", "gamma"}
	Println("index 0={.0} index 2={.2}", s, s) // index 0=alpha index 2=gamma

	// --- Date formatting: format the value first, then hand it to Format ---
	today := Named{"today": time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC).Format("2006-01-02")}
	Println("today is {today}", today) // today is 2025-06-15

	// --- type and debug format specs ---
	// {:T} → fmt.Sprintf("%T", v), {:?} → fmt.Sprintf("%#v", v)
	Println("{:T}  {:?}", Int(99), Int(99)) // g.Int  99

	// --- Mixing named and positional ---
	Println("lang={lang}, arg={}, arg={}",
		Named{"lang": "Go"}, "first", "second") // lang=Go, arg=first, arg=second
}
