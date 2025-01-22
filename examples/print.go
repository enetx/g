package main

import (
	"time"

	. "github.com/enetx/g"
)

func main() {
	// MapOrd example
	mo := NewMapOrd[any, any]()
	mo.
		Set("db", map[String]Map[String, int]{"user": {"age": 35}})

	Printf("user {.$get(db.user.age)} years old\n", mo)

	// Map example
	data := map[string]any{
		"user": map[string]string{
			"email": "user@example.com",
		},
	}

	Printf("Email: {.$get(user.email)}\n", data)

	// Basic named placeholders
	foo := "foo"
	bar := "bar"

	Printf("foo: {foo}, bar: {bar}\n", Named{"foo": foo, "bar": bar})

	// Named placeholders with fallback and multiple modifiers
	name := String("   john  ")
	age := 30
	city := "New York"

	named := Named{
		"name":  name,
		"age":   age,
		"city":  city,
		"today": time.Now(),
	}

	Printf(
		"Hello, my name is {name.$trim.$replace(j,r).$title.$substring(0,-2).$replace(o,0)}0t. I am {age} years old and live in {city.$truncate(2)}.\n",
		named,
	)

	Printf("Today is {today.$date(01/02/2006)}. Name fallback example: {unknown?name.$trim.$upper}\n", named)

	// Mixing autoindex placeholders with named placeholders
	Printf(
		"Numeric: {}, Named: {key.$fmt(%+v)}, Another numeric: \\{{.$upper}\\}\n",
		Named{"key": struct{ named string }{named: "value"}},
		"positional-1", // => {1}
		"positional-2", // => {2}
	)

	// Numeric-only usage
	Printf("{1} + {2} + {3} + {1}", "Hello", 123, "World")

	// Basic $get usage with a map
	mapExample := map[string]string{"key": "value"}
	Printf("Value from map: {1.$get(key)}\n", mapExample)

	// Nested map with $get
	nestedMap := map[String]map[string]String{
		"outer": {"inner": "nestedValue"},
	}

	Printf("Nested value: {1.$get(outer.inner)}\n", nestedMap)

	// Map with non-string keys
	mixedKeysMap := Map[Float, String]{
		3.14: "pi",
		2.71: "e",
	}

	Printf("Float key example: {.$get(3_14)}\n", mixedKeysMap)

	// Slice access with $get
	sliceExample := Slice[String]{"first", "second", "third"}
	Printf("Slice value at index 1: {.$get(1)}\n", sliceExample)

	// Nested slice access with $get
	nestedSlice := Slice[Slice[Int]]{{1, 2, 3}, {4, 5, 6}}
	Printf("Nested slice value: {1.$get(1.2)}\n", nestedSlice)

	// Struct access with $get
	type MyStruct struct {
		Field string
		Sub   struct {
			InnerField string
		}
	}

	structExample := MyStruct{
		Field: "fieldValue",
		Sub:   struct{ InnerField string }{InnerField: "innerValue"},
	}

	Printf("Struct field: {1.$get(Field)}, Sub-Field: {1.$get(Sub.InnerField)}\n", structExample)

	// Combination of map, slice, and struct
	complexExample := map[string]map[string][]struct {
		Key   string
		Value int
	}{
		"outer": {
			"middle": {
				{Key: "exampleKey", Value: 42},
			},
		},
	}

	Printf("Complex example: {1.$get(outer.middle.0.Key)} => {1.$get(outer.middle.0.Value)}\n", complexExample)

	// Full complexity with $get
	fullComplex := map[string]map[string]map[string][]map[string]string{
		"level1": {
			"level2": {
				"level3": {
					{"finalKey": "finalValue"},
				},
			},
		},
	}

	Printf("Full complexity: {1.$get(level1.level2.level3.0.finalKey)}\n", fullComplex)

	// Boolean keys
	boolMap := map[bool]string{true: "TrueValue", false: "FalseValue"}
	Printf("Boolean key true: {1.$get(true)}, Boolean key false: {1.$get(false)}\n", boolMap)

	// Map with int keys
	intKeyMap := map[int]string{42: "Answer to everything"}
	Printf("Integer key example: {1.$get(42)}\n", intKeyMap)

	// Complex nested structures
	complexNested := struct {
		Map   Map[String, Slice[string]]
		Array [2]map[int]string
	}{
		Map: Map[String, Slice[string]]{
			"list": {"item1", "item2"},
		},
		Array: [2]map[int]string{
			{1: "first", 2: "second"},
			{3: "third", 4: "fourth"},
		},
	}

	Printf("Nested struct: {1.$get(Map.list.1)}, Array item: {1.$get(Array.1.3)}\n", complexNested)
}
