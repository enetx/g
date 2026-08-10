package rx

import "github.com/enetx/g"

// IsMatchBytes reports whether the pattern matches bs.
func (r Regex) IsMatchBytes(bs g.Bytes) bool { return r.re.Match(bs) }

// FindBytes returns the first match in bs.
// A genuine empty match (e.g. patterns like `a*` or `^`) returns Some of an
// empty Bytes, while the absence of any match returns None.
func (r Regex) FindBytes(bs g.Bytes) g.Option[g.Bytes] {
	loc := r.re.FindIndex(bs)
	if loc == nil {
		return g.None[g.Bytes]()
	}

	return g.Some(bs[loc[0]:loc[1]])
}

// FindIterBytes returns a lazy sequence of all matches in bs.
func (r Regex) FindIterBytes(bs g.Bytes) g.Seq[g.Bytes] {
	return func(yield func(g.Bytes) bool) {
		for _, m := range r.re.FindAll(bs, -1) {
			if !yield(g.Bytes(m)) {
				return
			}
		}
	}
}

// CapturesBytes returns the first match and its capture groups: the full
// match at index 0 followed by the captured submatches. None if the pattern
// does not match.
func (r Regex) CapturesBytes(bs g.Bytes) g.Option[g.Slice[g.Bytes]] {
	result := g.TransformSlice(r.re.FindSubmatch(bs), func(b []byte) g.Bytes { return g.Bytes(b) })
	if result.IsEmpty() {
		return g.None[g.Slice[g.Bytes]]()
	}

	return g.Some(result)
}

// CapturesIterBytes returns a lazy sequence of capture groups for every match
// in bs.
func (r Regex) CapturesIterBytes(bs g.Bytes) g.Seq[g.Slice[g.Bytes]] {
	return func(yield func(g.Slice[g.Bytes]) bool) {
		for _, v := range r.re.FindAllSubmatch(bs, -1) {
			if !yield(g.TransformSlice(v, func(b []byte) g.Bytes { return g.Bytes(b) })) {
				return
			}
		}
	}
}

// IndexBytes returns the [start, end) byte offsets of the first match in bs,
// or None if the pattern does not match.
func (r Regex) IndexBytes(bs g.Bytes) g.Option[g.Slice[g.Int]] {
	result := g.TransformSlice(r.re.FindIndex(bs), g.NewInt)
	if result.IsEmpty() {
		return g.None[g.Slice[g.Int]]()
	}

	return g.Some(result)
}

// ReplaceBytes replaces only the FIRST match in bs with repl, expanding
// $-references. Use [Regex.ReplaceAllBytes] to replace every match.
func (r Regex) ReplaceBytes(bs, repl g.Bytes) g.Bytes {
	m := r.re.FindSubmatchIndex(bs)
	if m == nil {
		return bs
	}

	b := make([]byte, 0, len(bs))
	b = append(b, bs[:m[0]]...)
	b = r.re.Expand(b, repl, bs, m)
	b = append(b, bs[m[1]:]...)

	return b
}

// ReplaceAllBytes replaces every match in bs with repl, expanding
// $-references.
func (r Regex) ReplaceAllBytes(bs, repl g.Bytes) g.Bytes { return r.re.ReplaceAll(bs, repl) }

// ReplaceBytesBy replaces only the FIRST match in bs with the result of fn
// applied to the match; the rest of the slice is left untouched. The
// replacement is inserted literally — $-references are not expanded. The match
// passed to fn aliases the source backing array; clone it before mutating.
// Use [Regex.ReplaceAllBytesBy] to transform every match.
func (r Regex) ReplaceBytesBy(bs g.Bytes, fn func(match g.Bytes) g.Bytes) g.Bytes {
	loc := r.re.FindIndex(bs)
	if loc == nil {
		return bs
	}

	b := make(g.Bytes, 0, len(bs))
	b = append(b, bs[:loc[0]]...)
	b = append(b, fn(bs[loc[0]:loc[1]])...)
	b = append(b, bs[loc[1]:]...)

	return b
}

// ReplaceAllBytesBy replaces every match in bs with the result of fn applied
// to the match.
func (r Regex) ReplaceAllBytesBy(bs g.Bytes, fn func(match g.Bytes) g.Bytes) g.Bytes {
	return r.re.ReplaceAllFunc(bs, func(b []byte) []byte { return fn(g.Bytes(b)) })
}

// SplitBytes returns a lazy sequence of the subslices of bs around every match
// of the pattern.
func (r Regex) SplitBytes(bs g.Bytes) g.Seq[g.Bytes] { return r.SplitBytesN(bs, -1) }

// SplitBytesN splits bs around matches of the pattern into at most n
// subslices, mirroring the semantics of regexp.(*Regexp).Split (which only
// exists for strings): the last subslice is the unsplit remainder. If n is
// negative there is no limit; if n is zero the sequence is empty. Each yielded
// segment is a copy, so the results never alias bs's backing array.
func (r Regex) SplitBytesN(bs g.Bytes, n g.Int) g.Seq[g.Bytes] {
	return func(yield func(g.Bytes) bool) {
		if n == 0 {
			return
		}

		if len(r.re.String()) > 0 && len(bs) == 0 {
			yield(g.Bytes{})
			return
		}

		matches := r.re.FindAllIndex(bs, n.Std())

		beg := 0
		end := 0
		emitted := 0

		for _, match := range matches {
			if n > 0 && emitted+1 >= n.Std() {
				break
			}

			end = match[0]
			if match[1] != 0 {
				emitted++
				if !yield(bs[beg:end].Clone()) {
					return
				}
			}

			beg = match[1]
		}

		if end != len(bs) {
			yield(bs[beg:].Clone())
		}
	}
}

// MatchAnyBytes reports whether any of the patterns matches bs.
func MatchAnyBytes(bs g.Bytes, patterns ...Regex) bool {
	for _, p := range patterns {
		if p.IsMatchBytes(bs) {
			return true
		}
	}

	return false
}

// MatchAllBytes reports whether every pattern matches bs.
func MatchAllBytes(bs g.Bytes, patterns ...Regex) bool {
	for _, p := range patterns {
		if !p.IsMatchBytes(bs) {
			return false
		}
	}

	return true
}
