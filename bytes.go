package g

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// NewBytes creates a new Bytes value.
func NewBytes[T ~string | ~[]byte](bs T) Bytes { return Bytes(bs) }

// Reverse returns a new Bytes with the order of its runes reversed.
func (bs Bytes) Reverse() Bytes {
	reversed := make(Bytes, bs.Len())
	i := 0

	for bs.Len() > 0 {
		r, size := utf8.DecodeLastRune(bs)
		bs = bs[:bs.Len()-size]
		i += utf8.EncodeRune(reversed[i:], r)
	}

	return reversed
}

// Replace replaces the first 'n' occurrences of 'oldB' with 'newB' in the Bytes.
func (bs Bytes) Replace(oldB, newB Bytes, n int) Bytes { return bytes.Replace(bs, oldB, newB, n) }

// ReplaceAll replaces all occurrences of 'oldB' with 'newB' in the Bytes.
func (bs Bytes) ReplaceAll(oldB, newB Bytes) Bytes { return bytes.ReplaceAll(bs, oldB, newB) }

// ReplaceRegexp replaces all occurrences of the regular expression matches in the Bytes
// with the provided newB and returns the resulting Bytes after the replacement.
func (bs Bytes) ReplaceRegexp(pattern *regexp.Regexp, newB Bytes) Bytes {
	return pattern.ReplaceAll(bs, newB)
}

// FindRegexp searches the Bytes for the first occurrence of the regular expression pattern
// and returns an Option[Bytes] containing the matched substring.
// If no match is found, the Option[Bytes] will be None.
func (bs Bytes) FindRegexp(pattern *regexp.Regexp) Option[Bytes] {
	result := Bytes(pattern.Find(bs))
	if result.Empty() {
		return None[Bytes]()
	}

	return Some(result)
}

// Trim trims the specified characters from the beginning and end of the Bytes.
func (bs Bytes) Trim(cutset String) Bytes { return bytes.Trim(bs, cutset.Std()) }

// TrimLeft trims the specified characters from the beginning of the Bytes.
func (bs Bytes) TrimLeft(cutset String) Bytes { return bytes.TrimLeft(bs, cutset.Std()) }

// TrimRight trims the specified characters from the end of the Bytes.
func (bs Bytes) TrimRight(cutset String) Bytes { return bytes.TrimRight(bs, cutset.Std()) }

// TrimPrefix trims the specified Bytes prefix from the Bytes.
func (bs Bytes) TrimPrefix(cutset Bytes) Bytes { return bytes.TrimPrefix(bs, cutset) }

// TrimSuffix trims the specified Bytes suffix from the Bytes.
func (bs Bytes) TrimSuffix(cutset Bytes) Bytes { return bytes.TrimSuffix(bs, cutset) }

// Split splits the Bytes at each occurrence of the specified Bytes separator.
func (bs Bytes) Split(sep ...Bytes) Slice[Bytes] {
	var separator []byte
	if len(sep) != 0 {
		separator = sep[0]
	}

	return sliceBytesFromStd(bytes.Split(bs, separator))
}

func sliceBytesFromStd(bb [][]byte) Slice[Bytes] {
	result := NewSlice[Bytes](0, len(bb))
	for _, v := range bb {
		result = result.Append(NewBytes(v))
	}

	return result
}

// Add appends the given Bytes to the current Bytes.
func (bs Bytes) Add(obs Bytes) Bytes { return append(bs, obs...) }

// AddPrefix prepends the given Bytes to the current Bytes.
func (bs Bytes) AddPrefix(obs Bytes) Bytes { return obs.Add(bs) }

// Std returns the Bytes as a byte slice.
func (bs Bytes) Std() []byte { return bs }

// Clone creates a new Bytes instance with the same content as the current Bytes.
func (bs Bytes) Clone() Bytes { return bytes.Clone(bs) }

// Compare compares the Bytes with another Bytes and returns an Int.
func (bs Bytes) Compare(obs Bytes) Int { return Int(bytes.Compare(bs, obs)) }

// Contains checks if the Bytes contains the specified Bytes.
func (bs Bytes) Contains(obs Bytes) bool { return bytes.Contains(bs, obs) }

// ContainsAny checks if the Bytes contains any of the specified Bytes.
func (bs Bytes) ContainsAny(obss ...Bytes) bool {
	for _, obs := range obss {
		if bs.Contains(obs) {
			return true
		}
	}

	return false
}

// ContainsAll checks if the Bytes contains all of the specified Bytes.
func (bs Bytes) ContainsAll(obss ...Bytes) bool {
	for _, obs := range obss {
		if !bs.Contains(obs) {
			return false
		}
	}

	return true
}

// ContainsAnyChars checks if the given Bytes contains any characters from the input String.
func (bs Bytes) ContainsAnyChars(chars String) bool { return bytes.ContainsAny(bs, chars.Std()) }

// ContainsRune checks if the Bytes contains the specified rune.
func (bs Bytes) ContainsRune(r rune) bool { return bytes.ContainsRune(bs, r) }

// Count counts the number of occurrences of the specified Bytes in the Bytes.
func (bs Bytes) Count(obs Bytes) int { return bytes.Count(bs, obs) }

// Empty checks if the Bytes is empty.
func (bs Bytes) Empty() bool { return bs == nil || bs.Len() == 0 }

// Eq checks if the Bytes is equal to another Bytes.
func (bs Bytes) Eq(obs Bytes) bool { return bs.Compare(obs).Eq(0) }

// EqFold compares two Bytes slices case-insensitively.
func (bs Bytes) EqFold(obs Bytes) bool { return bytes.EqualFold(bs, obs) }

// Gt checks if the Bytes is greater than another Bytes.
func (bs Bytes) Gt(obs Bytes) bool { return bs.Compare(obs).Gt(0) }

// ToString returns the Bytes as an String.
func (bs Bytes) ToString() String { return String(bs) }

// Index returns the index of the first instance of obs in bs, or -1 if bs is not present in obs.
func (bs Bytes) Index(obs Bytes) int { return bytes.Index(bs, obs) }

// IndexRegexp searches for the first occurrence of the regular expression pattern in the Bytes.
// If a match is found, it returns an Option containing an Slice with the start and end indices of the match.
// If no match is found, it returns None.
func (bs Bytes) IndexRegexp(pattern *regexp.Regexp) Option[Slice[Int]] {
	result := TransformSlice(pattern.FindIndex(bs), NewInt)
	if result.Empty() {
		return None[Slice[Int]]()
	}

	return Some(result)
}

// FindAllRegexp searches the Bytes for all occurrences of the regular expression pattern
// and returns an Option[Slice[Bytes]] containing a slice of matched substrings.
// If no matches are found, the Option[Slice[Bytes]] will be None.
func (bs Bytes) FindAllRegexp(pattern *regexp.Regexp) Option[Slice[Bytes]] {
	return bs.FindAllRegexpN(pattern, -1)
}

// FindAllRegexpN searches the Bytes for up to n occurrences of the regular expression pattern
// and returns an Option[Slice[Bytes]] containing a slice of matched substrings.
// If no matches are found, the Option[Slice[Bytes]] will be None.
// If n is negative, all occurrences will be returned.
func (bs Bytes) FindAllRegexpN(pattern *regexp.Regexp, n Int) Option[Slice[Bytes]] {
	result := sliceBytesFromStd(pattern.FindAll(bs, n.Std()))
	if result.Empty() {
		return None[Slice[Bytes]]()
	}

	return Some(result)
}

// FindSubmatchRegexp searches the Bytes for the first occurrence of the regular expression pattern
// and returns an Option[Slice[Bytes]] containing the matched substrings and submatches.
// The Option[Slice[Bytes]] will contain an Slice[Bytes] for each match,
// where each Slice[Bytes] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Bytes]] will be None.
func (bs Bytes) FindSubmatchRegexp(pattern *regexp.Regexp) Option[Slice[Bytes]] {
	result := sliceBytesFromStd(pattern.FindSubmatch(bs))
	if result.Empty() {
		return None[Slice[Bytes]]()
	}

	return Some(result)
}

// FindAllSubmatchRegexp searches the Bytes for all occurrences of the regular expression pattern
// and returns an Option[Slice[Slice[Bytes]]] containing the matched substrings and submatches.
// The Option[Slice[Slice[Bytes]]] will contain an Slice[Bytes] for each match,
// where each Slice[Bytes] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Slice[Bytes]]] will be None.
// This method is equivalent to calling SubmatchAllRegexpN with n = -1, which means it finds all occurrences.
func (bs Bytes) FindAllSubmatchRegexp(pattern *regexp.Regexp) Option[Slice[Slice[Bytes]]] {
	return bs.FindAllSubmatchRegexpN(pattern, -1)
}

// FindAllSubmatchRegexpN searches the Bytes for occurrences of the regular expression pattern
// and returns an Option[Slice[Slice[Bytes]]] containing the matched substrings and submatches.
// The Option[Slice[Slice[Bytes]]] will contain an Slice[Bytes] for each match,
// where each Slice[Bytes] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Slice[Bytes]]] will be None.
// The 'n' parameter specifies the maximum number of matches to find. If n is negative, it finds all occurrences.
func (bs Bytes) FindAllSubmatchRegexpN(pattern *regexp.Regexp, n Int) Option[Slice[Slice[Bytes]]] {
	var result Slice[Slice[Bytes]]

	for _, v := range pattern.FindAllSubmatch(bs, n.Std()) {
		result = result.Append(sliceBytesFromStd(v))
	}

	if result.Empty() {
		return None[Slice[Slice[Bytes]]]()
	}

	return Some(result)
}

// LastIndex returns the index of the last instance of obs in bs, or -1 if obs is not present in bs.
func (bs Bytes) LastIndex(obs Bytes) int { return bytes.LastIndex(bs, obs) }

// IndexByte returns the index of the first instance of the byte b in bs, or -1 if b is not
// present in bs.
func (bs Bytes) IndexByte(b byte) int { return bytes.IndexByte(bs, b) }

// LastIndexByte returns the index of the last instance of the byte b in bs, or -1 if b is not
// present in bs.
func (bs Bytes) LastIndexByte(b byte) int { return bytes.LastIndexByte(bs, b) }

// IndexRune returns the index of the first instance of the rune r in bs, or -1 if r is not
// present in bs.
func (bs Bytes) IndexRune(r rune) int { return bytes.IndexRune(bs, r) }

// Len returns the length of the Bytes.
func (bs Bytes) Len() int { return len(bs) }

// LenRunes returns the number of runes in the Bytes.
func (bs Bytes) LenRunes() int { return utf8.RuneCount(bs) }

// Lt checks if the Bytes is less than another Bytes.
func (bs Bytes) Lt(obs Bytes) bool { return bs.Compare(obs).Lt(0) }

// Map applies a function to each rune in the Bytes and returns the modified Bytes.
func (bs Bytes) Map(fn func(rune) rune) Bytes { return bytes.Map(fn, bs) }

// NormalizeNFC returns a new Bytes with its Unicode characters normalized using the NFC form.
func (bs Bytes) NormalizeNFC() Bytes { return norm.NFC.Bytes(bs) }

// Ne checks if the Bytes is not equal to another Bytes.
func (bs Bytes) Ne(obs Bytes) bool { return !bs.Eq(obs) }

// NotEmpty checks if the Bytes is not empty.
func (bs Bytes) NotEmpty() bool { return bs.Len() != 0 }

// Reader returns a *bytes.Reader initialized with the content of Bytes.
func (bs Bytes) Reader() *bytes.Reader { return bytes.NewReader(bs) }

// Repeat returns a new Bytes consisting of the current Bytes repeated 'count' times.
func (bs Bytes) Repeat(count int) Bytes { return bytes.Repeat(bs, count) }

// ToRunes returns the Bytes as a slice of runes.
func (bs Bytes) ToRunes() []rune { return bytes.Runes(bs) }

// Title converts the Bytes to title case.
func (bs Bytes) Title() Bytes { return cases.Title(language.English).Bytes(bs) }

// Lower converts the Bytes to lowercase.
func (bs Bytes) Lower() Bytes { return cases.Lower(language.English).Bytes(bs) }

// Upper converts the Bytes to uppercase.
func (bs Bytes) Upper() Bytes { return cases.Upper(language.English).Bytes(bs) }

// TrimSpace trims white space characters from the beginning and end of the Bytes.
func (bs Bytes) TrimSpace() Bytes { return bytes.TrimSpace(bs) }

// Print prints the content of the Bytes to the standard output (console)
// and returns the Bytes unchanged.
func (bs Bytes) Print() Bytes { fmt.Println(bs); return bs }
