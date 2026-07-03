package g

import (
	"regexp"

	"github.com/enetx/g/f"
)

// regexps struct wraps a Bytes and provides regex-related methods.
type regexpb struct{ bytes Bytes }

// Regexp wraps a Bytes into an re struct to provide regex-related methods.
func (bs Bytes) Regexp() regexpb { return regexpb{bs} }

// Find searches the Bytes for the first occurrence of the regular expression pattern
// and returns an Option[Bytes] containing the matched substring.
// If no match is found, the Option[Bytes] will be None.
// A genuine empty match (e.g. patterns like `a*` or `^`) returns Some of an empty Bytes.
func (r regexpb) Find(pattern *regexp.Regexp) Option[Bytes] {
	loc := pattern.FindIndex(r.bytes)
	if loc == nil {
		return None[Bytes]()
	}

	return Some(Bytes(r.bytes[loc[0]:loc[1]]))
}

// Match checks if the Bytes contains a match for the specified regular expression pattern.
func (r regexpb) Match(pattern *regexp.Regexp) bool { return f.Match[[]byte](pattern)(r.bytes) }

// MatchAny checks if the Bytes contains a match for any of the specified regular
// expression patterns.
func (r regexpb) MatchAny(patterns ...*regexp.Regexp) bool {
	return Slice[*regexp.Regexp](patterns).
		Iter().
		Any(func(pattern *regexp.Regexp) bool { return r.Match(pattern) })
}

// MatchAll checks if the Bytes contains a match for all of the specified regular expression patterns.
func (r regexpb) MatchAll(patterns ...*regexp.Regexp) bool {
	return Slice[*regexp.Regexp](patterns).
		Iter().
		All(func(pattern *regexp.Regexp) bool { return r.Match(pattern) })
}

// Index searches for the first occurrence of the regular expression pattern in the Bytes.
// If a match is found, it returns an Option containing an Slice with the start and end indices of the match.
// If no match is found, it returns None.
func (r regexpb) Index(pattern *regexp.Regexp) Option[Slice[Int]] {
	result := transformSlice(pattern.FindIndex(r.bytes), NewInt)
	if result.IsEmpty() {
		return None[Slice[Int]]()
	}

	return Some(result)
}

// FindAll searches the Bytes for all occurrences of the regular expression pattern
// and returns an Option[Slice[Bytes]] containing a slice of matched substrings.
// If no matches are found, the Option[Slice[Bytes]] will be None.
func (r regexpb) FindAll(pattern *regexp.Regexp) Option[Slice[Bytes]] {
	return r.FindAllN(pattern, -1)
}

// FindAllN searches the Bytes for up to n occurrences of the regular expression pattern
// and returns an Option[Slice[Bytes]] containing a slice of matched substrings.
// If no matches are found, the Option[Slice[Bytes]] will be None.
// If n is negative, all occurrences will be returned.
func (r regexpb) FindAllN(pattern *regexp.Regexp, n Int) Option[Slice[Bytes]] {
	result := transformSlice(pattern.FindAll(r.bytes, n.Std()), func(bs []byte) Bytes { return Bytes(bs) })
	if result.IsEmpty() {
		return None[Slice[Bytes]]()
	}

	return Some(result)
}

// FindSubmatch searches the Bytes for the first occurrence of the regular expression pattern
// and returns an Option[Slice[Bytes]] containing the matched substrings and submatches.
// The Option[Slice[Bytes]] will contain an Slice[Bytes] for each match,
// where each Slice[Bytes] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Bytes]] will be None.
func (r regexpb) FindSubmatch(pattern *regexp.Regexp) Option[Slice[Bytes]] {
	result := transformSlice(pattern.FindSubmatch(r.bytes), func(bs []byte) Bytes { return Bytes(bs) })
	if result.IsEmpty() {
		return None[Slice[Bytes]]()
	}

	return Some(result)
}

// FindAllSubmatch searches the Bytes for all occurrences of the regular expression pattern
// and returns an Option[Slice[Slice[Bytes]]] containing the matched substrings and submatches.
// The Option[Slice[Slice[Bytes]]] will contain an Slice[Bytes] for each match,
// where each Slice[Bytes] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Slice[Bytes]]] will be None.
// This method is equivalent to calling SubmatchAllRegexpN with n = -1, which means it finds all occurrences.
func (r regexpb) FindAllSubmatch(pattern *regexp.Regexp) Option[Slice[Slice[Bytes]]] {
	return r.FindAllSubmatchN(pattern, -1)
}

// FindAllSubmatchN searches the Bytes for occurrences of the regular expression pattern
// and returns an Option[Slice[Slice[Bytes]]] containing the matched substrings and submatches.
// The Option[Slice[Slice[Bytes]]] will contain an Slice[Bytes] for each match,
// where each Slice[Bytes] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Slice[Bytes]]] will be None.
// The 'n' parameter specifies the maximum number of matches to find. If n is negative, it finds all occurrences.
func (r regexpb) FindAllSubmatchN(pattern *regexp.Regexp, n Int) Option[Slice[Slice[Bytes]]] {
	var result Slice[Slice[Bytes]]

	for _, v := range pattern.FindAllSubmatch(r.bytes, n.Std()) {
		result = append(result, transformSlice(v, func(bs []byte) Bytes { return Bytes(bs) }))
	}

	if result.IsEmpty() {
		return None[Slice[Slice[Bytes]]]()
	}

	return Some(result)
}

// Replace replaces all occurrences of the regular expression matches in the Bytes
// with the provided newB and returns the resulting Bytes after the replacement.
func (r regexpb) Replace(pattern *regexp.Regexp, newB Bytes) Bytes {
	return pattern.ReplaceAll(r.bytes, newB)
}

// ReplaceBy replaces all occurrences of the regular expression matches in the Bytes
// by applying a custom transformation function to each match.
// The function `fn` takes a Bytes representing a match and returns a Bytes that will replace it.
func (r regexpb) ReplaceBy(pattern *regexp.Regexp, fn func(match Bytes) Bytes) Bytes {
	return pattern.ReplaceAllFunc(r.bytes, func(b []byte) []byte { return fn(Bytes(b)) })
}

// Split splits the Bytes into substrings using the provided regular expression pattern and returns an Slice[Bytes] of the results.
// The regular expression pattern is provided as a regexp.Regexp parameter.
func (r regexpb) Split(pattern *regexp.Regexp) Slice[Bytes] {
	return r.splitN(pattern, -1)
}

// SplitN splits the Bytes into substrings using the provided regular expression pattern and returns an Slice[Bytes] of the results.
// The regular expression pattern is provided as a regexp.Regexp parameter.
// The n parameter controls the number of substrings to return:
// - If n is negative, there is no limit on the number of substrings returned.
// - If n is zero, an empty Slice[Bytes] is returned.
// - If n is positive, at most n substrings are returned.
func (r regexpb) SplitN(pattern *regexp.Regexp, n Int) Option[Slice[Bytes]] {
	result := r.splitN(pattern, n.Std())
	if result.IsEmpty() {
		return None[Slice[Bytes]]()
	}

	return Some(result)
}

// splitN slices r.bytes around the matches of pattern, mirroring the semantics of
// regexp.(*Regexp).Split (which only exists for strings). Each returned segment is a
// copy, so the result never aliases the receiver's backing array.
func (r regexpb) splitN(pattern *regexp.Regexp, n int) Slice[Bytes] {
	if n == 0 {
		return nil
	}

	if len(pattern.String()) > 0 && len(r.bytes) == 0 {
		return Slice[Bytes]{Bytes{}}
	}

	matches := pattern.FindAllIndex(r.bytes, n)
	result := make(Slice[Bytes], 0, len(matches)+1)

	beg := 0
	end := 0

	for _, match := range matches {
		if n > 0 && len(result)+1 >= n {
			break
		}

		end = match[0]
		if match[1] != 0 {
			result = append(result, Bytes(r.bytes[beg:end]).Clone())
		}

		beg = match[1]
	}

	if end != len(r.bytes) {
		result = append(result, Bytes(r.bytes[beg:]).Clone())
	}

	return result
}
