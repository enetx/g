// Package rx provides compiled regular expressions: compile a pattern once
// with [New], then run it against g.String haystacks (or g.Bytes via the
// *Bytes method variants). The method set is IsMatch, Find, FindIter,
// Captures, CapturesIter, Index, Replace and ReplaceBy (first match),
// ReplaceAll, ReplaceAllBy, Split, SplitN.
//
//	re := rx.New(`[0-9]{4}-[0-9]{2}-[0-9]{2}`).Unwrap()
//	dates := re.FindIter("What do 1865-04-14 and 1963-11-22 have in common?").Collect().Slice()
//	// dates: Slice["1865-04-14", "1963-11-22"]
//
// The *Iter and Split* methods return lazy sequences that plug into g's
// iterator chains (Map, Filter, Take, Collect, ...). Matching runs when the
// chain is consumed, not when the method is called; the underlying scan is
// stdlib regexp's, which has no incremental API, so all matches materialize at
// that point.
//
// These operations live in their own package rather than as methods on
// g.String: as methods they would weld Slice[String] and Seq[String]
// into String's method set, dragging the whole container graph into every
// package that merely names g.String.
package rx

import (
	"regexp"

	"github.com/enetx/g"
)

// Regex is a compiled regular expression.
// The zero value is not usable; build one with [New] or [From].
type Regex struct{ re *regexp.Regexp }

// New compiles pattern into a [Regex].
//
//	re := rx.New(`\d+`).Unwrap()
func New[T ~string](pattern T) g.Result[Regex] {
	re, err := regexp.Compile(string(pattern))
	if err != nil {
		return g.Err[Regex](err)
	}

	return g.Ok(Regex{re})
}

// From wraps an already-compiled *regexp.Regexp.
func From(re *regexp.Regexp) Regex { return Regex{re} }

// Std returns the underlying *regexp.Regexp.
func (r Regex) Std() *regexp.Regexp { return r.re }

// String returns the source pattern text.
func (r Regex) String() string { return r.re.String() }

// IsMatch reports whether the pattern matches s.
func (r Regex) IsMatch(s g.String) bool { return r.re.MatchString(s.Std()) }

// Find returns the first match in s.
// A genuine empty match (e.g. patterns like `\d*`, `a*`, or `^`) returns
// Some(""), while the absence of any match returns None.
func (r Regex) Find(s g.String) g.Option[g.String] {
	loc := r.re.FindStringIndex(s.Std())
	if loc == nil {
		return g.None[g.String]()
	}

	return g.Some(s[loc[0]:loc[1]])
}

// FindIter returns a lazy sequence of all matches in s.
//
//	re.FindIter(s).Take(2).Collect().Slice()
func (r Regex) FindIter(s g.String) g.Seq[g.String] {
	return func(yield func(g.String) bool) {
		for _, m := range r.re.FindAllString(s.Std(), -1) {
			if !yield(g.String(m)) {
				return
			}
		}
	}
}

// Captures returns the first match and its capture groups: the full match at
// index 0 followed by the captured submatches.
// None if the pattern does not match.
func (r Regex) Captures(s g.String) g.Option[g.Slice[g.String]] {
	result := g.TransformSlice(r.re.FindStringSubmatch(s.Std()), g.NewString)
	if result.IsEmpty() {
		return g.None[g.Slice[g.String]]()
	}

	return g.Some(result)
}

// CapturesIter returns a lazy sequence of capture groups for every match
// in s.
func (r Regex) CapturesIter(s g.String) g.Seq[g.Slice[g.String]] {
	return func(yield func(g.Slice[g.String]) bool) {
		for _, v := range r.re.FindAllStringSubmatch(s.Std(), -1) {
			if !yield(g.TransformSlice(v, g.NewString)) {
				return
			}
		}
	}
}

// Index returns the [start, end) byte offsets of the first match in s, or None
// if the pattern does not match.
func (r Regex) Index(s g.String) g.Option[g.Slice[g.Int]] {
	result := g.TransformSlice(r.re.FindStringIndex(s.Std()), g.NewInt)
	if result.IsEmpty() {
		return g.None[g.Slice[g.Int]]()
	}

	return g.Some(result)
}

// Replace replaces only the FIRST match in s with repl, expanding
// $-references. Use [Regex.ReplaceAll] to replace every match.
func (r Regex) Replace(s, repl g.String) g.String {
	m := r.re.FindStringSubmatchIndex(s.Std())
	if m == nil {
		return s
	}

	b := make([]byte, 0, len(s))
	b = append(b, s[:m[0]]...)
	b = r.re.ExpandString(b, repl.Std(), s.Std(), m)
	b = append(b, s[m[1]:]...)

	return g.String(b)
}

// ReplaceAll replaces every match in s with repl, expanding $-references.
func (r Regex) ReplaceAll(s, repl g.String) g.String {
	return g.String(r.re.ReplaceAllString(s.Std(), repl.Std()))
}

// ReplaceBy replaces only the FIRST match in s with the result of fn applied
// to the match; the rest of the string is left untouched. The replacement is
// inserted literally — $-references are not expanded. Use [Regex.ReplaceAllBy]
// to transform every match.
func (r Regex) ReplaceBy(s g.String, fn func(match g.String) g.String) g.String {
	loc := r.re.FindStringIndex(s.Std())
	if loc == nil {
		return s
	}

	b := make([]byte, 0, len(s))
	b = append(b, s[:loc[0]]...)
	b = append(b, fn(s[loc[0]:loc[1]])...)
	b = append(b, s[loc[1]:]...)

	return g.String(b)
}

// ReplaceAllBy replaces every match in s with the result of fn applied to the
// match.
func (r Regex) ReplaceAllBy(s g.String, fn func(match g.String) g.String) g.String {
	return g.String(r.re.ReplaceAllStringFunc(s.Std(), func(m string) string { return fn(g.String(m)).Std() }))
}

// Split returns a lazy sequence of the substrings of s around every match of
// the pattern.
func (r Regex) Split(s g.String) g.Seq[g.String] { return r.SplitN(s, -1) }

// SplitN splits s around matches of the pattern into at most n substrings;
// the last substring is the unsplit remainder. If n is
// negative there is no limit; if n is zero the sequence is empty.
func (r Regex) SplitN(s g.String, n g.Int) g.Seq[g.String] {
	return func(yield func(g.String) bool) {
		for _, part := range r.re.Split(s.Std(), n.Std()) {
			if !yield(g.String(part)) {
				return
			}
		}
	}
}

// MatchAny reports whether any of the patterns matches s.
func MatchAny(s g.String, patterns ...Regex) bool {
	for _, p := range patterns {
		if p.IsMatch(s) {
			return true
		}
	}

	return false
}

// MatchAll reports whether every pattern matches s.
func MatchAll(s g.String, patterns ...Regex) bool {
	for _, p := range patterns {
		if !p.IsMatch(s) {
			return false
		}
	}

	return true
}
