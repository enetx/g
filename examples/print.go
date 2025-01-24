package main

import (
	"time"

	. "github.com/enetx/g"
)

func main() {
	// Struct access
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

	Printf("Struct field: {1.Field}, Sub-Field: {1.Sub.InnerField.Upper}\n", structExample)

	// Methods
	Printf("{.SubString(0,-1,2)}\n", String("somestring"))

	Printf("Hex: {1.Hex}, Binary: {1.Binary}\n", Int(255))

	// MapOrd example
	mo := NewMapOrd[String, Map[String, Map[String, int]]]()
	mo.Set("db", Map[String, Map[String, int]]{"user": {"age": 35}})

	Printf("user {1.Get(db).Some.Get(user).Some.Get(age).Some} years old\n", mo)
	Printf("user {.db.user.age} years old\n", mo.AsAny())

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
		"numeric: {}, named: {key.Named}, another numeric: \\{{.Upper}\\}\n",
		Named{"key": struct{ Named string }{Named: "value"}},
		"positional-1",         // => {1}
		String("positional-2"), // => {2}
	)

	// Numeric-only usage
	Printf("{1} + {2} + {3} + {1}\n", "Hello", 123, "World")

	// Basic usage with a map
	mapExample := Map[string, string]{"key": "value"}
	Printf("Value from map: {.Get(key).Some}\n", mapExample)
	Printf("Value from map: {.key}\n", mapExample)

	// Nested map
	nestedMap := Map[String, Map[string, String]]{
		"outer": {"inner": "nestedValue"},
	}

	Printf("Nested value: {1.Get(outer).Some.Get(inner).Some}\n", nestedMap)
	Printf("Nested value: {.outer.inner}\n", nestedMap)

	// Map with non-string keys
	mixedKeysMap := Map[Float, String]{
		3.14: "pi",
		2.71: "e",
	}

	Printf("Float key example: {.Get(3_14).Some}\n", mixedKeysMap)
	Printf("Float key example: {.3_14}\n", mixedKeysMap)

	// Slice access
	sliceExample := Slice[String]{"first", "second", "third"}
	Printf("Slice value at index 1: {.Get(1)}\n", sliceExample)
	Printf("Slice value at index 1: {.1}\n", sliceExample)

	// Nested slice access
	nestedSlice := Slice[Slice[Int]]{{1, 2, 3}, {4, 5, 6}}
	Printf("Nested slice value: {.Get(1).Get(2)}\n", nestedSlice)
	Printf("Nested slice value: {.1.2}\n", nestedSlice)

	// Boolean keys
	boolMap := Map[bool, string]{true: "TrueValue", false: "FalseValue"}
	Printf("Boolean key true: {1.Get(true).Some}, Boolean key false: {1.Get(false).Some}\n", boolMap)
	Printf("Boolean key true: {1.true}, Boolean key false: {1.false}\n", boolMap)

	// Map with int keys
	intKeyMap := Map[int, string]{42: "Answer to everything"}
	Printf("Integer key example: {.Get(42).Some}\n", intKeyMap)
	Printf("Integer key example: {.42}\n", intKeyMap)
}
