package main

import (
	"fmt"
	"strconv"

	"github.com/enetx/g"
)

func main() {
	// Creating Some and None Options
	someOption := g.Some(42)
	noneOption := g.None[int]()

	// Checking if Option is Some or None
	fmt.Println(someOption.IsSome()) // Output: true
	fmt.Println(noneOption.IsNone()) // Output: true

	// Unwrapping Options
	fmt.Println(someOption.Unwrap())     // Output: 42
	fmt.Println(noneOption.UnwrapOr(10)) // Output: 10

	// Mapping over Options
	doubledOption := someOption.Then(func(val int) g.Option[int] {
		return g.Some(val * 2)
	})

	fmt.Println(doubledOption.Unwrap()) // Output: 84

	// Using OptionMap to transform the value inside Option
	addTwoOption := g.OptionMap(someOption, func(val int) g.Option[string] {
		return g.Some("result: " + strconv.Itoa(val+2))
	})

	fmt.Println(addTwoOption.Unwrap()) // Output: "result: 44"

	// Using UnwrapOrDefault to handle None Option with default value
	defaultValue := noneOption.UnwrapOrDefault()
	fmt.Println(defaultValue) // Output: 0 (default value for int)

	// Using Then to chain operations on Option
	resultOption := someOption.
		Then(
			func(val int) g.Option[int] {
				if val > 10 {
					return g.Some(val * 2)
				}
				return g.None[int]()
			}).
		Then(
			func(val int) g.Option[int] {
				return g.Some(val + 5)
			})

	fmt.Println(resultOption.Unwrap()) // Output: 89

	// Using Expect to handle None Option
	noneOption.Expect("This is None")
	// The above line will panic with message "This is None"
}
