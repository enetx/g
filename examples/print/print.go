package main

import (
	. "github.com/enetx/g"
)

func main() {
	// Struct access: placeholders resolve DATA — struct fields, map keys and
	// slice indices. Methods are never invoked;
	// transform values before formatting.
	type MyStruct struct {
		Field string
		Sub   struct {
			InnerField String
		}
	}

	structExample := MyStruct{
		Field: "fieldValue",
		Sub:   struct{ InnerField String }{InnerField: "innerValue"},
	}

	Print("Struct field: {1.Field}, Sub-Field: {1.Sub.InnerField}, {1:T}, {1:?}\n", structExample)

	// Transform first, then format.
	Print("{}\n", String("somestring").SubString(0, -1, 2).Center(9, "="))

	Print("Hex: {}, Binary: {}\n", Int(255).Hex(), Int(255).Binary())

	// MapOrd example: dotted path walks Key/Value pairs.
	mo := NewMapOrd[String, Map[String, Map[String, int]]]()
	mo.Insert("db", Map[String, Map[String, int]]{"user": {"age": 35}})

	Print("user {.db.user.age} years old\n", mo)

	// Basic named placeholders
	foo := String("foo")
	bar := "bar"

	Print("foo: {foo}, bar: {bar}\n", Named{"foo": foo.Upper(), "bar": bar})

	// Named placeholders with fallback
	named := Named{
		"name": String("   john  ").Trim().Title(),
		"age":  3,
		"city": String("New York"),
	}

	Println("Hello, my name is {name}. I am {age} years old and live in {city}.", named)

	Println("Name fallback example: {unknown?name}", named)

	// Mixing autoindex placeholders with named placeholders
	Print(
		"numeric: {}, named: {key.Named}, literal braces: {{{}}}\n",
		Named{"key": struct{ Named string }{Named: "value"}},
		"positional-1",         // => {1}
		String("positional-2"), // => {2}
	)

	// Numeric-only usage
	Print("{1} + {2} + {3} + {1}\n", "Hello", 123, "World")

	// Basic usage with a map
	mapExample := Map[string, string]{"key": "value"}
	Print("Value from map: {.key}\n", mapExample)

	// Nested map
	nestedMap := Map[String, Map[string, String]]{
		"outer": {"inner": "nestedValue"},
	}

	Print("Nested value: {.outer.inner}\n", nestedMap)

	// Map with non-string keys ("_" reads as "." for float keys)
	mixedKeysMap := Map[Float, String]{
		3.14: "pi",
		2.71: "e",
	}

	Print("Float key example: {.3_14}\n", mixedKeysMap)

	// Slice access
	sliceExample := Slice[String]{"first", "second", "third"}
	Print("Slice value at index 1: {.1}\n", sliceExample)

	// Nested slice access
	nestedSlice := Slice[Slice[Int]]{{1, 2, 3}, {4, 5, 6}}
	Print("Nested slice value: {.1.2}\n", nestedSlice)

	// Boolean keys
	boolMap := Map[bool, string]{true: "TrueValue", false: "FalseValue"}
	Print("Boolean key true: {1.true}, Boolean key false: {1.false}\n", boolMap)

	// Map with int keys
	intKeyMap := Map[int, string]{42: "Answer to everything"}
	Print("Integer key example: {.42}\n", intKeyMap)
}
