package main

import (
	"fmt"

	"github.com/enetx/g/rx"
)

func main() {
	// Compile a pattern once, then run it against haystacks.
	re := rx.New(`(?m)(post-)(?:\d+)`).Unwrap()

	// Replace the first match; $-references expand capture groups.
	re.Replace("post-55 post-66", "${1}38").Println() // post-38 post-66

	// ReplaceAll replaces every match.
	re.ReplaceAll("post-55 post-66", "${1}38").Println() // post-38 post-38

	// Find returns the first match as an Option.
	re.Find("x post-55 x").Some().Println() // post-55

	// A miss is None; UnwrapOr supplies a default.
	re.Find("post-not-found").UnwrapOr("post-333").Println() // post-333

	// FindIter is a lazy sequence over all matches.
	re.FindIter("some post-55 not found post-31 post-22").
		Collect().
		Slice().
		Println()
		// Slice[post-55, post-31, post-22]

	// Iterator adapters compose before collection.
	re.FindIter("some post-55 not found post-31 post-22").Take(2).Collect().Slice().Println() // Slice[post-55, post-31]

	// Index reports the [start, end) offsets of the first match.
	re.Index("post-55").Some().Println() // Slice[0, 7]

	// Captures returns the full match followed by capture groups.
	re.Captures("some post-55 not found post-31").Some().Println() // Slice[post-55, post-]

	// CapturesIter yields capture groups for every match.
	re.CapturesIter("some post-55 not found post-31 post-22").
		Collect().
		Slice().
		Println()
		// Slice[Slice[post-55, post-], ...]

	// Split slices the haystack around matches.
	spaces := rx.New(`\s`).Unwrap()
	digits := rx.New(`\d`).Unwrap()

	spaces.Split("some test for split n").Collect().Slice().Println() // Slice[some, test, for, split, n]

	// SplitN keeps the unsplit remainder in the last element.
	spaces.SplitN("some test for split n", 2).Collect().Slice().Println() // Slice[some, test for split n]

	// IsMatch reports whether the pattern matches.
	fmt.Println(spaces.IsMatch("some test")) // true

	// MatchAll / MatchAny take several compiled patterns.
	fmt.Println(rx.MatchAll("some test 1", spaces, digits)) // true
	fmt.Println(rx.MatchAny("some_test", spaces, digits))   // false
}
