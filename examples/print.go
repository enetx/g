package main

import (
	"time"

	"github.com/enetx/g"
	. "github.com/enetx/g"
	"github.com/enetx/g/ref"
)

type user struct{ id *int }

func (u *user) GetID() int {
	return *u.id
}

func main() {
	a := 1

	Println("{}00% loaded", a)
	g.Println("{}00% loaded", a)

	Println("%s", a)
	g.Println("%d", a)

	user := &user{id: ref.Of(19)}
	Println("{.GetID}", user)

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

	Print("Struct field: {1.Field}, Sub-Field: {1.Sub.InnerField.Upper}, {1.type}, {1.debug}\n", structExample)

	// Methods
	Print("{.SubString(0,-1,2).Center(9,=)}\n", String("somestring"))

	Print("Hex: {.Hex}, Binary: {1.Binary}\n", Int(255))

	// MapOrd example
	mo := NewMapOrd[String, Map[String, Map[String, int]]]()
	mo.Set("db", Map[String, Map[String, int]]{"user": {"age": 35}})

	Print("user {1.Get(db).Some.Get(user).Some.Get(age).Some} years old\n", mo)
	Print("user {.db.user.age} years old\n", mo.AsAny())

	// Basic named placeholders
	foo := String("foo")
	bar := "bar"

	Print("foo: {foo.Upper}, bar: {bar}\n", Named{"foo": foo, "bar": bar})

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

	Println(
		"Hello, my name is {name.Trim.ReplaceAll(j,r).Title.SubString(0,-2).ReplaceAll(o,0)}0t. I am {age} years old and live in {city.Truncate(2)}.",
		named,
	)

	Println("Today is {today.Format(01/02/2006)}. Name fallback example: {unknown?name.Trim.Upper}", named)

	Print("Name fallback example: {unknown?name.Trim.Upper}\n", named)

	// Mixing autoindex placeholders with named placeholders
	Print(
		"numeric: {}, named: {key.Named}, another numeric: \\{{.Upper}\\}\n",
		Named{"key": struct{ Named string }{Named: "value"}},
		"positional-1",         // => {1}
		String("positional-2"), // => {2}
	)

	// Numeric-only usage
	Print("{1} + {2} + {3} + {1}\n", "Hello", 123, "World")

	// Basic usage with a map
	mapExample := Map[string, string]{"key": "value"}
	Print("Value from map: {.Get(key).Some}\n", mapExample)
	Print("Value from map: {.key}\n", mapExample)

	// Nested map
	nestedMap := Map[String, Map[string, String]]{
		"outer": {"inner": "nestedValue"},
	}

	Print("Nested value: {1.Get(outer).Some.Get(inner).Some}\n", nestedMap)
	Print("Nested value: {.outer.inner}\n", nestedMap)

	// Map with non-string keys
	mixedKeysMap := Map[Float, String]{
		3.14: "pi",
		2.71: "e",
	}

	Print("Float key example: {.Get(3_14).Some}\n", mixedKeysMap)
	Print("Float key example: {.3_14}\n", mixedKeysMap)

	// Slice access
	sliceExample := Slice[String]{"first", "second", "third"}
	Print("Slice value at index 1: {.Get(1).Some}\n", sliceExample)
	Print("Slice value at index 1: {.1}\n", sliceExample)

	// Nested slice access
	nestedSlice := Slice[Slice[Int]]{{1, 2, 3}, {4, 5, 6}}
	Print("Nested slice value: {.Get(1).Some.Get(2).Some}\n", nestedSlice)
	Print("Nested slice value: {.1.2}\n", nestedSlice)

	// Boolean keys
	boolMap := Map[bool, string]{true: "TrueValue", false: "FalseValue"}
	Print("Boolean key true: {1.Get(true).Some}, Boolean key false: {1.Get(false).Some}\n", boolMap)
	Print("Boolean key true: {1.true}, Boolean key false: {1.false}\n", boolMap)

	// Map with int keys
	intKeyMap := Map[int, string]{42: "Answer to everything"}
	Print("Integer key example: {.Get(42).Some}\n", intKeyMap)
	Print("Integer key example: {.42}\n", intKeyMap)
}
