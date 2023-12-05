package main

import (
	"fmt"
	"regexp"

	"gitlab.com/x0xO/g"
)

func main() {
	pattern := regexp.MustCompile(`(?m)(post-)(?:\d+)`)

	g.String("post-55").ReplaceRegexp(pattern, "${1}38").Print()
	// post-38

	g.String("x post-55 x").FindRegexp(pattern).Some().Print()
	// post-55

	g.String("post-not-found").FindRegexp(pattern).UnwrapOr("post-333").Print()
	// post-333

	g.String("some post-55 not found post-31 post-22").FindAllRegexp(pattern).Some().Print()
	// Slice[post-55, post-31, post-22]

	g.String("some post-55 not found post-31 post-22").FindAllRegexpN(pattern, 2).Some().Print()
	// Slice[post-55, post-31]

	g.String("post-55").IndexRegexp(pattern).Some().Print()
	// Slice[0, 7]

	g.String("some post-55 not found post-31").FindSubmatchRegexp(pattern).Some().Print()
	// Slice[post-55, post-]

	g.String("some post-55 not found post-31 post-22").FindAllSubmatchRegexp(pattern).Some().Print()
	// Slice[Slice[post-55, post-], Slice[post-31, post-], Slice[post-22, post-]]

	g.String("some post-55 not found post-31 post-22").FindAllSubmatchRegexpN(pattern, 2).Some().Print()
	// Slice[Slice[post-55, post-], Slice[post-31, post-]]

	g.String("some test for split n").SplitRegexp(*regexp.MustCompile(`\s`)).Print()
	// Slice[some, test, for, split, n]

	g.String("some test for split n").SplitRegexpN(*regexp.MustCompile(`\s`), 2).Some().Print()
	// Slice[some, test for split n]

	fmt.Println(g.String("some test").ContainsRegexp(`\s`).Unwrap())
	fmt.Println(g.String("some test 1").ContainsRegexpAll(`\s`, `\d`).Unwrap())
	fmt.Println(g.String("some test").ContainsRegexpAny(`\s`, `\d`).Unwrap())
}
