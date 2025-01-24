package main

import (
	"time"

	. "github.com/enetx/g"
)

func main() {
	// Methods
	Printf("{.SubString(0,-1,2)}\n", String("somestring"))

	Printf("Hex: {1.Hex}, Binary: {1.Binary}\n", Int(255))

	// MapOrd example
	mo := NewMapOrd[String, Map[String, Map[String, int]]]()
	mo.Set("db", Map[String, Map[String, int]]{"user": {"age": 35}})

	Printf("user {1.Get(db).Some.Get(user).Some.Get(age).Some} years old\n", mo)

	// Basic named placeholders
	foo := String("foo")
	bar := "bar"

	Printf("foo: {foo.Upper}, bar: {bar}\n", Named{"foo": foo, "bar": bar})

	// Named placeholders with fallback and multiple modifiers
	name := String("   john  ")
	age := 3
	city := String("New York")

	named := Named{
		"name":  name,
		"age":   age,
		"city":  city,
		"today": time.Now(),
	}

	Printf(
		"Hello, my name is {name.Trim.ReplaceAll(j,r).Title.SubString(0,-2).ReplaceAll(o,0)}0t. I am {age} years old and live in {city.Truncate(2)}.\n",
		named,
	)

	Printf("Today is {today.Format(01/02/2006)}. Name fallback example: {unknown?name.Trim.Upper}\n", named)

	Printf("Name fallback example: {unknown?name.Trim.Upper}\n", named)

	// Mixing autoindex placeholders with named placeholders
	Printf(
		"Numeric: {}, Named: {key}, Another numeric: \\{{.Upper}\\}\n",
		Named{"key": struct{ named string }{named: "value"}},
		"positional-1",         // => {1}
		String("positional-2"), // => {2}
	)

	// Numeric-only usage
	Printf("{1} + {2} + {3} + {1}", "Hello", 123, "World")

	// Basic usage with a map
	mapExample := Map[string, string]{"key": "value"}
	Printf("Value from map: {1.key.Get(key).Some}\n", mapExample)

	// Nested map
	nestedMap := Map[String, Map[string, String]]{
		"outer": {"inner": "nestedValue"},
	}

	Printf("Nested value: {1.Get(outer).Some.Get(inner).Some}\n", nestedMap)

	// Map with non-string keys
	mixedKeysMap := Map[Float, String]{
		3.14: "pi",
		2.71: "e",
	}

	Printf("Float key example: {.Get(3_14).Some}\n", mixedKeysMap)

	// Slice access with $get
	sliceExample := Slice[String]{"first", "second", "third"}
	Printf("Slice value at index 1: {.Get(1)}\n", sliceExample)

	// Nested slice access with $get
	nestedSlice := Slice[Slice[Int]]{{1, 2, 3}, {4, 5, 6}}
	Printf("Nested slice value: {1.Get(1).Get(2)}\n", nestedSlice)

	// Boolean keys
	boolMap := Map[bool, string]{true: "TrueValue", false: "FalseValue"}
	Printf("Boolean key true: {1.Get(true).Some}, Boolean key false: {1.Get(false).Some}\n", boolMap)

	// Map with int keys
	intKeyMap := Map[int, string]{42: "Answer to everything"}
	Printf("Integer key example: {1.Get(42).Some}\n", intKeyMap)
}
