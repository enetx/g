package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Define a regular expression pattern
	pattern := String(`(?m)(post-)(?:\d+)`).Regexp().Compile().Unwrap()

	// Replace the first occurrence of the pattern in the string
	String("post-55").Regexp().Replace(pattern, "${1}38").Println()
	// Output: post-38

	// Find the first match of the pattern in the string
	o := String("x post-55 x").Regexp().Find(pattern)
	// Output: post-55

	// switch pattern
	switch {
	case o.IsSome():
		o.Some().Println()
	default:
		fmt.Println("not found")
	}

	// If no match is found, provide a default value
	String("post-not-found").Regexp().Find(pattern).UnwrapOr("post-333").Println()
	// Output: post-333

	// Find all matches of the pattern in the string
	String("some post-55 not found post-31 post-22").Regexp().FindAll(pattern).Unwrap().Println()
	// Output: Slice[post-55, post-31, post-22]

	// Find a specific number of matches of the pattern in the string
	String("some post-55 not found post-31 post-22").Regexp().FindAllN(pattern, 2).Some().Println()
	// Output: Slice[post-55, post-31]

	// Get the starting indices of the first match of the pattern in the string
	String("post-55").Regexp().Index(pattern).Some().Println()
	// Output: Slice[0, 7]

	// Find the submatches of the first match of the pattern in the string
	String("some post-55 not found post-31").Regexp().FindSubmatch(pattern).Some().Println()
	// Output: Slice[post-55, post-]

	// Find all submatches of the pattern in the string
	String("some post-55 not found post-31 post-22").Regexp().FindAllSubmatch(pattern).Some().Println()
	// Output: Slice[Slice[post-55, post-], Slice[post-31, post-], Slice[post-22, post-]]

	// Find a specific number of submatches of the pattern in the string
	String("some post-55 not found post-31 post-22").Regexp().FindAllSubmatchN(pattern, 2).Some().Println()
	// Output: Slice[Slice[post-55, post-], Slice[post-31, post-]]

	patterns := String(`\s`).Regexp().Compile().Unwrap()
	patternd := String(`\d`).Regexp().Compile().Unwrap()

	// Split the string using the regular expression pattern
	String("some test for split n").Regexp().Split(patterns).Println()
	// Output: Slice[some, test, for, split, n]

	// Split the string using the regular expression pattern, limiting the number of splits
	String("some test for split n").Regexp().SplitN(patterns, 2).Some().Println()
	// Output: Slice[some, test for split n]

	// Check if the string contains a match of the regular expression
	fmt.Println(String("some test").Regexp().Match(patterns))
	// Output: true

	// Check if the string contains matches for all the provided regular expressions
	fmt.Println(String("some test 1").Regexp().MatchAll(patterns, patternd))
	// Output: true

	// Check if the string contains a match for any of the provided regular expressions
	fmt.Println(String("some test").Regexp().MatchAny(patterns, patternd))
	// Output: true
}
