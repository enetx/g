package main

import (
	"fmt"
	"strconv"

	. "github.com/enetx/g"
)

func main() {
	// Creating Some and None Options
	someOption := Some(42)
	noneOption := None[int]()

	// Checking if Option is Some or None
	fmt.Println(someOption.IsSome()) // Output: true
	fmt.Println(noneOption.IsNone()) // Output: true

	// Unwrapping Options
	fmt.Println(someOption.Unwrap())     // Output: 42
	fmt.Println(noneOption.UnwrapOr(10)) // Output: 10

	// Mapping over Options
	doubledOption := someOption.Then(func(val int) Option[int] {
		return Some(val * 2)
	})

	fmt.Println(doubledOption.Unwrap()) // Output: 84

	// Using Then to transform the value inside Option into a different type
	addTwoOption := someOption.Then(func(val int) Option[string] {
		return Some("result: " + strconv.Itoa(val+2))
	})

	fmt.Println(addTwoOption.Unwrap()) // Output: "result: 44"

	// Using UnwrapOrDefault to handle None Option with default value
	defaultValue := noneOption.UnwrapOrDefault()
	fmt.Println(defaultValue) // Output: 0 (default value for int)

	// Using Then to chain operations on Option
	resultOption := someOption.
		Then(
			func(val int) Option[int] {
				if val > 10 {
					return Some(val * 2)
				}
				return None[int]()
			}).
		Then(
			func(val int) Option[int] {
				return Some(val + 5)
			})

	fmt.Println(resultOption.Unwrap()) // Output: 89

	// Using OptionOf to create Option based on a condition
	m := map[string]int{"one": 100}

	v, ok := m["one"]

	optionFromCondition := OptionOf(v, ok)
	fmt.Println(optionFromCondition.Unwrap()) // Output: 100

	v, ok = m["two"]

	optionFromConditionFalse := OptionOf(v, ok)
	fmt.Println(optionFromConditionFalse.UnwrapOr(50)) // Output: 50 (default value when None)

	// MapOr collapses the Option into a plain value of any type
	fmt.Println(Some(21).MapOr("none", func(v int) string { return fmt.Sprintf("v=%d", v*2) })) // v=42
	fmt.Println(None[int]().MapOr("none", func(int) string { return "x" }))                     // none

	// ThenOf bridges Go's comma-ok functions into the chain, changing type
	lookup := map[string]int{"a": 1}
	fmt.Println(Some("a").
		ThenOf(func(k string) (int, bool) { v, ok := lookup[k]; return v, ok })) // Some(1)

	// IsNoneOr — complement of IsSomeAnd
	fmt.Println(None[int]().IsNoneOr(func(int) bool { return false })) // true

	// Using Expect to handle None Option
	noneOption.Expect("This is None")
	// The above line will panic with message "This is None"

}
