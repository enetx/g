package g

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// NewBytes creates a new Bytes value.
func NewBytes[T ~string | ~[]byte](bs T) Bytes { return Bytes(bs) }

// Transform applies a transformation function to the Bytes and returns the result.
func (bs Bytes) Transform(fn func(Bytes) Bytes) Bytes { return fn(bs) }

// Reverse returns a new Bytes with the order of its runes reversed.
func (bs Bytes) Reverse() Bytes {
	reversed := make(Bytes, bs.Len())
	i := 0

	for bs.Len() > 0 {
		r, size := utf8.DecodeLastRune(bs)
		bs = bs[:bs.Len().Std()-size]
		i += utf8.EncodeRune(reversed[i:], r)
	}

	return reversed
}

// Replace replaces the first 'n' occurrences of 'oldB' with 'newB' in the Bytes.
func (bs Bytes) Replace(oldB, newB Bytes, n Int) Bytes { return bytes.Replace(bs, oldB, newB, n.Std()) }

// ReplaceAll replaces all occurrences of 'oldB' with 'newB' in the Bytes.
func (bs Bytes) ReplaceAll(oldB, newB Bytes) Bytes { return bytes.ReplaceAll(bs, oldB, newB) }

// RxReplace replaces all occurrences of the regular expression matches in the Bytes
// with the provided newB and returns the resulting Bytes after the replacement.
func (bs Bytes) RxReplace(pattern *regexp.Regexp, newB Bytes) Bytes {
	return pattern.ReplaceAll(bs, newB)
}

// RxFind searches the Bytes for the first occurrence of the regular expression pattern
// and returns an Option[Bytes] containing the matched substring.
// If no match is found, the Option[Bytes] will be None.
func (bs Bytes) RxFind(pattern *regexp.Regexp) Option[Bytes] {
	result := Bytes(pattern.Find(bs))
	if result.Empty() {
		return None[Bytes]()
	}

	return Some(result)
}

// Trim trims leading and trailing white space from the Bytes.
func (bs Bytes) Trim() Bytes { return bytes.TrimSpace(bs) }

// TrimStart removes leading white space from the Bytes.
func (bs Bytes) TrimStart() Bytes { return trimBytesStart(bs) }

// TrimEnd removes trailing white space from the Bytes.
func (bs Bytes) TrimEnd() Bytes { return trimBytesEnd(bs) }

// TrimSet trims the specified set of characters from both the beginning and end of the Bytes.
func (bs Bytes) TrimSet(cutset String) Bytes { return bytes.Trim(bs, cutset.Std()) }

// TrimStartSet removes the specified set of characters from the beginning of the Bytes.
func (bs Bytes) TrimStartSet(cutset String) Bytes { return bytes.TrimLeft(bs, cutset.Std()) }

// TrimEndSet removes the specified set of characters from the end of the Bytes.
func (bs Bytes) TrimEndSet(cutset String) Bytes { return bytes.TrimRight(bs, cutset.Std()) }

// StripPrefix trims the specified Bytes prefix from the Bytes.
func (bs Bytes) StripPrefix(cutset Bytes) Bytes { return bytes.TrimPrefix(bs, cutset) }

// StripSuffix trims the specified Bytes suffix from the Bytes.
func (bs Bytes) StripSuffix(cutset Bytes) Bytes { return bytes.TrimSuffix(bs, cutset) }

// Split splits the Bytes by the specified separator and returns the iterator.
func (bs Bytes) Split(sep ...Bytes) SeqSlice[Bytes] {
	var separator []byte
	if len(sep) != 0 {
		separator = sep[0]
	}

	return splitBytes(bs, separator, 0)
}

// SplitAfter splits the Bytes after each instance of the specified separator and returns the iterator.
func (bs Bytes) SplitAfter(sep Bytes) SeqSlice[Bytes] { return splitBytes(bs, sep, sep.Len()) }

// Fields splits the Bytes into a slice of substrings, removing any whitespace, and returns the iterator.
func (bs Bytes) Fields() SeqSlice[Bytes] { return fieldsBytes(bs) }

// FieldsBy splits the Bytes into a slice of substrings using a custom function to determine the field boundaries,
// and returns the iterator.
func (bs Bytes) FieldsBy(fn func(r rune) bool) SeqSlice[Bytes] { return fieldsbyBytes(bs, fn) }

// Add appends the given Bytes to the current Bytes.
func (bs Bytes) Add(obs Bytes) Bytes { return append(bs, obs...) }

// AddPrefix prepends the given Bytes to the current Bytes.
func (bs Bytes) AddPrefix(obs Bytes) Bytes { return obs.Add(bs) }

// Std returns the Bytes as a byte slice.
func (bs Bytes) Std() []byte { return bs }

// Clone creates a new Bytes instance with the same content as the current Bytes.
func (bs Bytes) Clone() Bytes { return bytes.Clone(bs) }

// Cmp compares the Bytes with another Bytes and returns an cmp.Ordering.
func (bs Bytes) Cmp(obs Bytes) cmp.Ordering { return cmp.Ordering(bytes.Compare(bs, obs)) }

// Contains checks if the Bytes contains the specified Bytes.
func (bs Bytes) Contains(obs Bytes) bool { return bytes.Contains(bs, obs) }

// ContainsAny checks if the Bytes contains any of the specified Bytes.
func (bs Bytes) ContainsAny(obss ...Bytes) bool {
	return Slice[Bytes](obss).
		Iter().
		Any(func(obs Bytes) bool { return bs.Contains(obs) })
}

// ContainsAll checks if the Bytes contains all of the specified Bytes.
func (bs Bytes) ContainsAll(obss ...Bytes) bool {
	return Slice[Bytes](obss).
		Iter().
		All(func(obs Bytes) bool { return bs.Contains(obs) })
}

// ContainsAnyChars checks if the given Bytes contains any characters from the input String.
func (bs Bytes) ContainsAnyChars(chars String) bool { return bytes.ContainsAny(bs, chars.Std()) }

// ContainsRune checks if the Bytes contains the specified rune.
func (bs Bytes) ContainsRune(r rune) bool { return bytes.ContainsRune(bs, r) }

// Count counts the number of occurrences of the specified Bytes in the Bytes.
func (bs Bytes) Count(obs Bytes) Int { return Int(bytes.Count(bs, obs)) }

// Empty checks if the Bytes is empty.
func (bs Bytes) Empty() bool { return len(bs) == 0 }

// Eq checks if the Bytes is equal to another Bytes.
func (bs Bytes) Eq(obs Bytes) bool { return bs.Cmp(obs).IsEq() }

// EqFold compares two Bytes slices case-insensitively.
func (bs Bytes) EqFold(obs Bytes) bool { return bytes.EqualFold(bs, obs) }

// Gt checks if the Bytes is greater than another Bytes.
func (bs Bytes) Gt(obs Bytes) bool { return bs.Cmp(obs).IsGt() }

// String returns the Bytes as an String.
func (bs Bytes) String() String { return String(bs) }

// Index returns the index of the first instance of obs in bs, or -1 if bs is not present in obs.
func (bs Bytes) Index(obs Bytes) Int { return Int(bytes.Index(bs, obs)) }

// RxMatch checks if the Bytes contains a match for the specified regular expression pattern.
func (bs Bytes) RxMatch(pattern *regexp.Regexp) bool { return f.RxMatch[Bytes](pattern)(bs) }

// RxMatchAny checks if the Bytes contains a match for any of the specified regular
// expression patterns.
func (bs Bytes) RxMatchAny(patterns ...*regexp.Regexp) bool {
	return Slice[*regexp.Regexp](patterns).
		Iter().
		Any(func(pattern *regexp.Regexp) bool { return bs.RxMatch(pattern) })
}

// RxMatchAll checks if the Bytes contains a match for all of the specified regular expression patterns.
func (bs Bytes) RxMatchAll(patterns ...*regexp.Regexp) bool {
	return Slice[*regexp.Regexp](patterns).
		Iter().
		All(func(pattern *regexp.Regexp) bool { return bs.RxMatch(pattern) })
}

// RxIndex searches for the first occurrence of the regular expression pattern in the Bytes.
// If a match is found, it returns an Option containing an Slice with the start and end indices of the match.
// If no match is found, it returns None.
func (bs Bytes) RxIndex(pattern *regexp.Regexp) Option[Slice[Int]] {
	result := TransformSlice(pattern.FindIndex(bs), NewInt)
	if result.Empty() {
		return None[Slice[Int]]()
	}

	return Some(result)
}

// RxFindAll searches the Bytes for all occurrences of the regular expression pattern
// and returns an Option[Slice[Bytes]] containing a slice of matched substrings.
// If no matches are found, the Option[Slice[Bytes]] will be None.
func (bs Bytes) RxFindAll(pattern *regexp.Regexp) Option[Slice[Bytes]] {
	return bs.RxFindAllN(pattern, -1)
}

// RxFindAllN searches the Bytes for up to n occurrences of the regular expression pattern
// and returns an Option[Slice[Bytes]] containing a slice of matched substrings.
// If no matches are found, the Option[Slice[Bytes]] will be None.
// If n is negative, all occurrences will be returned.
func (bs Bytes) RxFindAllN(pattern *regexp.Regexp, n Int) Option[Slice[Bytes]] {
	result := TransformSlice(pattern.FindAll(bs, n.Std()), NewBytes)
	if result.Empty() {
		return None[Slice[Bytes]]()
	}

	return Some(result)
}

// RxFindSubmatch searches the Bytes for the first occurrence of the regular expression pattern
// and returns an Option[Slice[Bytes]] containing the matched substrings and submatches.
// The Option[Slice[Bytes]] will contain an Slice[Bytes] for each match,
// where each Slice[Bytes] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Bytes]] will be None.
func (bs Bytes) RxFindSubmatch(pattern *regexp.Regexp) Option[Slice[Bytes]] {
	result := TransformSlice(pattern.FindSubmatch(bs), NewBytes)
	if result.Empty() {
		return None[Slice[Bytes]]()
	}

	return Some(result)
}

// RxFindAllSubmatch searches the Bytes for all occurrences of the regular expression pattern
// and returns an Option[Slice[Slice[Bytes]]] containing the matched substrings and submatches.
// The Option[Slice[Slice[Bytes]]] will contain an Slice[Bytes] for each match,
// where each Slice[Bytes] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Slice[Bytes]]] will be None.
// This method is equivalent to calling SubmatchAllRegexpN with n = -1, which means it finds all occurrences.
func (bs Bytes) RxFindAllSubmatch(pattern *regexp.Regexp) Option[Slice[Slice[Bytes]]] {
	return bs.RxFindAllSubmatchN(pattern, -1)
}

// RxFindAllSubmatchN searches the Bytes for occurrences of the regular expression pattern
// and returns an Option[Slice[Slice[Bytes]]] containing the matched substrings and submatches.
// The Option[Slice[Slice[Bytes]]] will contain an Slice[Bytes] for each match,
// where each Slice[Bytes] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Slice[Bytes]]] will be None.
// The 'n' parameter specifies the maximum number of matches to find. If n is negative, it finds all occurrences.
func (bs Bytes) RxFindAllSubmatchN(pattern *regexp.Regexp, n Int) Option[Slice[Slice[Bytes]]] {
	var result Slice[Slice[Bytes]]

	for _, v := range pattern.FindAllSubmatch(bs, n.Std()) {
		result = append(result, TransformSlice(v, NewBytes))
	}

	if result.Empty() {
		return None[Slice[Slice[Bytes]]]()
	}

	return Some(result)
}

// LastIndex returns the index of the last instance of obs in bs, or -1 if obs is not present in bs.
func (bs Bytes) LastIndex(obs Bytes) Int { return Int(bytes.LastIndex(bs, obs)) }

// IndexByte returns the index of the first instance of the byte b in bs, or -1 if b is not
// present in bs.
func (bs Bytes) IndexByte(b byte) Int { return Int(bytes.IndexByte(bs, b)) }

// LastIndexByte returns the index of the last instance of the byte b in bs, or -1 if b is not
// present in bs.
func (bs Bytes) LastIndexByte(b byte) Int { return Int(bytes.LastIndexByte(bs, b)) }

// IndexRune returns the index of the first instance of the rune r in bs, or -1 if r is not
// present in bs.
func (bs Bytes) IndexRune(r rune) Int { return Int(bytes.IndexRune(bs, r)) }

// Len returns the length of the Bytes.
func (bs Bytes) Len() Int { return Int(len(bs)) }

// LenRunes returns the number of runes in the Bytes.
func (bs Bytes) LenRunes() Int { return Int(utf8.RuneCount(bs)) }

// Lt checks if the Bytes is less than another Bytes.
func (bs Bytes) Lt(obs Bytes) bool { return bs.Cmp(obs).IsLt() }

// Map applies a function to each rune in the Bytes and returns the modified Bytes.
func (bs Bytes) Map(fn func(rune) rune) Bytes { return bytes.Map(fn, bs) }

// NormalizeNFC returns a new Bytes with its Unicode characters normalized using the NFC form.
func (bs Bytes) NormalizeNFC() Bytes { return norm.NFC.Bytes(bs) }

// Ne checks if the Bytes is not equal to another Bytes.
func (bs Bytes) Ne(obs Bytes) bool { return !bs.Eq(obs) }

// NotEmpty checks if the Bytes is not empty.
func (bs Bytes) NotEmpty() bool { return !bs.Empty() }

// Reader returns a *bytes.Reader initialized with the content of Bytes.
func (bs Bytes) Reader() *bytes.Reader { return bytes.NewReader(bs) }

// Repeat returns a new Bytes consisting of the current Bytes repeated 'count' times.
func (bs Bytes) Repeat(count Int) Bytes { return bytes.Repeat(bs, count.Std()) }

// Runes returns the Bytes as a slice of runes.
func (bs Bytes) Runes() []rune { return bytes.Runes(bs) }

// Title converts the Bytes to title case.
func (bs Bytes) Title() Bytes { return cases.Title(language.English).Bytes(bs) }

// Lower converts the Bytes to lowercase.
func (bs Bytes) Lower() Bytes { return cases.Lower(language.English).Bytes(bs) }

// Upper converts the Bytes to uppercase.
func (bs Bytes) Upper() Bytes { return cases.Upper(language.English).Bytes(bs) }

// Print prints the content of the Bytes to the standard output (console)
// and returns the Bytes unchanged.
func (bs Bytes) Print() Bytes { fmt.Println(bs); return bs }

// trimBytesStart trims the leading whitespace characters from the byte slice.
func trimBytesStart(s []byte) []byte {
	start := 0

	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			return trimBytesFuncStart(s[start:], unicode.IsSpace)
		}

		if asciiSpace[c] == 0 {
			break
		}
	}

	if start == len(s) {
		return nil
	}

	return s[start:]
}

// trimBytesEnd trims the trailing whitespace characters from the byte slice.
func trimBytesEnd(s []byte) []byte {
	stop := len(s)

	for ; stop > 0; stop-- {
		c := s[stop-1]
		if c >= utf8.RuneSelf {
			return trimBytesFuncEnd(s[:stop], unicode.IsSpace)
		}

		if asciiSpace[c] == 0 {
			break
		}
	}

	if stop == 0 {
		return nil
	}

	return s[:stop]
}

// Helper function to trim leading characters using a unicode function
func trimBytesFuncStart(s []byte, fn func(rune) bool) []byte {
	start := 0

	for start < len(s) {
		r, size := utf8.DecodeRune(s[start:])
		if !fn(r) {
			break
		}

		start += size
	}

	return s[start:]
}

// Helper function to trim trailing characters using a unicode function
func trimBytesFuncEnd(s []byte, fn func(rune) bool) []byte {
	stop := len(s)

	for stop > 0 {
		r, size := utf8.DecodeLastRune(s[:stop])
		if !fn(r) {
			break
		}

		stop -= size
	}

	return s[:stop]
}
