package main

import (
	"fmt"
	"regexp"

	"gitlab.com/x0xO/g"
)

func main() {
	// Define a regular expression pattern
	pattern := regexp.MustCompile(`(?m)(post-)(?:\d+)`)

	// Replace the first occurrence of the pattern in the string
	g.String("post-55").ReplaceRegexp(pattern, "${1}38").Print()
	// Output: post-38

	// Find the first match of the pattern in the string
	g.String("x post-55 x").FindRegexp(pattern).Some().Print()
	// Output: post-55

	// If no match is found, provide a default value
	g.String("post-not-found").FindRegexp(pattern).UnwrapOr("post-333").Print()
	// Output: post-333

	// Find all matches of the pattern in the string
	g.String("some post-55 not found post-31 post-22").FindAllRegexp(pattern).Some().Print()
	// Output: Slice[post-55, post-31, post-22]

	// Find a specific number of matches of the pattern in the string
	g.String("some post-55 not found post-31 post-22").FindAllRegexpN(pattern, 2).Some().Print()
	// Output: Slice[post-55, post-31]

	// Get the starting indices of the first match of the pattern in the string
	g.String("post-55").IndexRegexp(pattern).Some().Print()
	// Output: Slice[0, 7]

	// Find the submatches of the first match of the pattern in the string
	g.String("some post-55 not found post-31").FindSubmatchRegexp(pattern).Some().Print()
	// Output: Slice[post-55, post-]

	// Find all submatches of the pattern in the string
	g.String("some post-55 not found post-31 post-22").FindAllSubmatchRegexp(pattern).Some().Print()
	// Output: Slice[Slice[post-55, post-], Slice[post-31, post-], Slice[post-22, post-]]

	// Find a specific number of submatches of the pattern in the string
	g.String("some post-55 not found post-31 post-22").FindAllSubmatchRegexpN(pattern, 2).Some().Print()
	// Output: Slice[Slice[post-55, post-], Slice[post-31, post-]]

	// Split the string using the regular expression pattern
	g.String("some test for split n").SplitRegexp(*regexp.MustCompile(`\s`)).Print()
	// Output: Slice[some, test, for, split, n]

	// Split the string using the regular expression pattern, limiting the number of splits
	g.String("some test for split n").SplitRegexpN(*regexp.MustCompile(`\s`), 2).Some().Print()
	// Output: Slice[some, test for split n]

	// Check if the string contains a match of the regular expression
	fmt.Println(g.String("some test").ContainsRegexp(`\s`).Unwrap())
	// Output: true

	// Check if the string contains matches for all the provided regular expressions
	fmt.Println(g.String("some test 1").ContainsRegexpAll(`\s`, `\d`).Unwrap())
	// Output: true

	// Check if the string contains a match for any of the provided regular expressions
	fmt.Println(g.String("some test").ContainsRegexpAny(`\s`, `\d`).Unwrap())
	// Output: true
}
