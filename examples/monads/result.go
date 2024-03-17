package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	// Example 1

	double := func(x g.Int) g.Result[g.Int] { return g.Ok(x * 2) }

	result := g.Ok(g.Int(10)).Then(double).Then(double)

	if result.IsOk() {
		fmt.Println("Result:", result.Ok()) // Output: Result: 40
	} else {
		fmt.Println("Error:", result.Err())
	}

	// Example 2

	square := func(x g.Int) g.Result[g.Int] { return g.Ok(x * x) }
	addFive := func(x g.Int) g.Result[g.Int] { return g.Ok(x + 5) }

	result = g.Ok(g.Int(10)).Then(square).Then(addFive)

	if result.IsOk() {
		fmt.Println("Result:", result.Ok()) // Output: Result: 105
	} else {
		fmt.Println("Error:", result.Err())
	}

	// Example 3
	result = g.String("15").ToInt().Then(double)

	if result.IsOk() {
		fmt.Println("Result:", result.Ok()) // Output: Result: 30
	} else {
		fmt.Println("Error:", result.Err())
	}

	// Example 4
	divideByZero := func(x float64) g.Result[float64] {
		if x == 0 {
			return g.Err[float64](fmt.Errorf("division by zero"))
		}

		return g.Ok(10.0 / float64(x))
	}

	resultf := g.Ok(0.0).Then(divideByZero)

	if resultf.IsOk() {
		fmt.Println("Result:", resultf.Ok())
	} else {
		fmt.Println("Error:", resultf.Err()) // Output: Error: division by zero
	}
}
