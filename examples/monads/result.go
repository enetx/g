package main

import (
	"fmt"
	"strconv"

	. "github.com/enetx/g"
)

func main() {
	// Example 1: Chaining operations on a Result using Then

	// Define a function to double an integer
	double := func(x Int) Result[Int] { return Ok(x * 2) }

	// Create a Result containing the integer 10, then double it twice using Then
	result := Ok(Int(10)).Then(double).Then(double)

	// Check if the result is an Ok value or an error
	if result.IsOk() {
		// Print the value if it's an Ok result
		fmt.Println("Result:", result.Ok()) // Output: Result: 40
	} else {
		// Print the error if it's an error result
		fmt.Println("Error:", result.Err())
	}

	// Example 2: Chaining multiple operations on a Result using Then

	// Define a function to square an integer
	square := func(x Int) Result[Int] { return Ok(x * x) }

	// Define a function to add five to an integer
	addFive := func(x Int) Result[Int] { return Ok(x + 5) }

	// Chain square and addFive operations on a Result containing the integer 10
	result = Ok(Int(10)).Then(square).Then(addFive)

	// Check if the result is an Ok value or an error
	if result.IsOk() {
		// Print the value if it's an Ok result
		fmt.Println("Result:", result.Ok()) // Output: Result: 105
	} else {
		// Print the error if it's an error result
		fmt.Println("Error:", result.Err())
	}

	// Example 3: Converting a string to an integer and then doubling the result

	// Convert the string "15" to an integer and then double it
	result = String("15").TryInt().Then(double)

	// Check if the result is an Ok value or an error
	if result.IsOk() {
		// Print the value if it's an Ok result
		fmt.Println("Result:", result.Ok()) // Output: Result: 30
	} else {
		// Print the error if it's an error result
		fmt.Println("Error:", result.Err())
	}

	// Example 4: Handling division by zero

	// Define a function to divide 10.0 by a float64 value, handling division by zero
	divideByZero := func(x float64) Result[float64] {
		if x == 0 {
			return Err[float64](fmt.Errorf("division by zero"))
		}
		return Ok(10.0 / x)
	}

	// Attempt to divide 10.0 by 0.0
	resultFloat := Ok(0.0).Then(divideByZero)

	// Check if the result is an Ok value or an error
	if resultFloat.IsOk() {
		// Print the value if it's an Ok result
		fmt.Println("Result:", resultFloat.Ok())
	} else {
		// Print the error if it's an error result
		fmt.Println("Error:", resultFloat.Err()) // Output: Error: division by zero
	}

	// Example 5: Converting a string to an integer using MapResult and MapToResult

	// Define a string containing a valid integer
	str := "123"

	// Create a Result containing the string
	strResult := Ok(str)

	// Use ResultMap to convert the string to an integer
	intResult := TransformResult(strResult, func(s string) Result[int] { return ResultOf(strconv.Atoi(s)) })

	// Alternatively, use ResultOfMap to convert the string to an integer
	// This simplifies the process by directly passing strconv.Atoi
	intResult = TransformResultOf(strResult, strconv.Atoi)

	// Check if the intResult is an error or contains a value
	if intResult.IsErr() {
		fmt.Println("Error:", intResult.Err())
	} else {
		fmt.Println("Integer Value:", intResult.Ok())
	}
}
