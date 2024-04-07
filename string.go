package g

import (
	"cmp"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/enetx/g/pkg/minmax"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// NewString creates a new String from the provided string.
func NewString[T ~string | rune | byte | ~[]rune | ~[]byte](str T) String { return String(str) }

// Builder returns a new Builder initialized with the content of the String.
func (s String) Builder() *Builder { return NewBuilder().Write(s) }

// Min returns the minimum of Strings.
func (s String) Min(b ...String) String { return minmax.Min(s, b...) }

// Max returns the maximum of Strings.
func (s String) Max(b ...String) String { return minmax.Max(s, b...) }

// Random generates a random String of the specified length, selecting characters from predefined sets.
// If additional character sets are provided, only those will be used; the default set (ASCII_LETTERS and DIGITS)
// is excluded unless explicitly provided.
//
// Parameters:
// - count (int): Length of the random String to generate.
// - letters (...String): Additional character sets to consider for generating the random String (optional).
//
// Returns:
// - String: Randomly generated String with the specified length.
//
// Example usage:
//
//	randomString := g.String.Random(10)
//	randomString contains a random String with 10 characters.
func (String) Random(count int, letters ...String) String {
	var chars Slice[rune]

	if len(letters) != 0 {
		chars = letters[0].ToRunes()
	} else {
		chars = (ASCII_LETTERS + DIGITS).ToRunes()
	}

	var result strings.Builder

	for range count {
		result.WriteRune(chars.Random())
	}

	return String(result.String())
}

// IsASCII checks if all characters in the String are ASCII bytes.
func (s String) IsASCII() bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}

	return true
}

// IsDigit checks if all characters in the String are digits.
func (s String) IsDigit() bool {
	if s.Empty() {
		return false
	}

	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}

// ToInt tries to parse the String as an int and returns an Int.
func (s String) ToInt() Result[Int] {
	hint, err := strconv.ParseInt(s.Std(), 0, 32)
	if err != nil {
		return Err[Int](err)
	}

	return Ok(Int(hint))
}

// ToFloat tries to parse the String as a float64 and returns an Float.
func (s String) ToFloat() Result[Float] {
	float, err := strconv.ParseFloat(s.Std(), 64)
	if err != nil {
		return Err[Float](err)
	}

	return Ok(Float(float))
}

// Title converts the String to title case.
func (s String) Title() String {
	return String(cases.Title(language.English).String(s.Std()))
}

// Lower returns the String in lowercase.
func (s String) Lower() String {
	return String(cases.Lower(language.English).String(s.Std()))
}

// Upper returns the String in uppercase.
func (s String) Upper() String {
	return String(cases.Upper(language.English).String(s.Std()))
}

// Trim trims characters in the cutset from the beginning and end of the String.
func (s String) Trim(cutset String) String {
	return String(strings.Trim(s.Std(), cutset.Std()))
}

// TrimLeft trims characters in the cutset from the beginning of the String.
func (s String) TrimLeft(cutset String) String {
	return String(strings.TrimLeft(s.Std(), cutset.Std()))
}

// TrimRight trims characters in the cutset from the end of the String.
func (s String) TrimRight(cutset String) String {
	return String(strings.TrimRight(s.Std(), cutset.Std()))
}

// TrimPrefix trims the specified prefix from the String.
func (s String) TrimPrefix(prefix String) String {
	return String(strings.TrimPrefix(s.Std(), prefix.Std()))
}

// TrimSuffix trims the specified suffix from the String.
func (s String) TrimSuffix(suffix String) String {
	return String(strings.TrimSuffix(s.Std(), suffix.Std()))
}

// Replace replaces the 'oldS' String with the 'newS' String for the specified number of
// occurrences.
func (s String) Replace(oldS, newS String, n int) String {
	return String(strings.Replace(s.Std(), oldS.Std(), newS.Std(), n))
}

// ReplaceAll replaces all occurrences of the 'oldS' String with the 'newS' String.
func (s String) ReplaceAll(oldS, newS String) String {
	return String(strings.ReplaceAll(s.Std(), oldS.Std(), newS.Std()))
}

// ReplaceMulti creates a custom replacer to perform multiple string replacements.
//
// Parameters:
//
// - oldnew ...String: Pairs of strings to be replaced. Specify as many pairs as needed.
//
// Returns:
//
// - String: A new string with replacements applied using the custom replacer.
//
// Example usage:
//
//	original := g.String("Hello, world! This is a test.")
//	replaced := original.ReplaceMulti(
//	    "Hello", "Greetings",
//	    "world", "universe",
//	    "test", "example",
//	)
//	// replaced contains "Greetings, universe! This is an example."
func (s String) ReplaceMulti(oldnew ...String) String {
	on := SliceOf(oldnew...).ToStringSlice()
	return String(strings.NewReplacer(on...).Replace(s.Std()))
}

// ReplaceRegexp replaces all occurrences of the regular expression matches in the String
// with the provided newS (as a String) and returns the resulting String after the replacement.
func (s String) ReplaceRegexp(pattern *regexp.Regexp, newS String) String {
	return String(pattern.ReplaceAllString(s.Std(), newS.Std()))
}

// FindRegexp searches the String for the first occurrence of the regulare xpression pattern
// and returns an Option[String] containing the matched substring.
// If no match is found, it returns None.
func (s String) FindRegexp(pattern *regexp.Regexp) Option[String] {
	result := String(pattern.FindString(s.Std()))
	if result.Empty() {
		return None[String]()
	}

	return Some(result)
}

// ReplaceNth returns a new String instance with the nth occurrence of oldS
// replaced with newS. If there aren't enough occurrences of oldS, the
// original String is returned. If n is less than -1, the original String
// is also returned. If n is -1, the last occurrence of oldS is replaced with newS.
//
// Returns:
//
// - A new String instance with the nth occurrence of oldS replaced with newS.
//
// Example usage:
//
//	s := g.String("The quick brown dog jumped over the lazy dog.")
//	result := s.ReplaceNth("dog", "fox", 2)
//	fmt.Println(result)
//
// Output: "The quick brown dog jumped over the lazy fox.".
func (s String) ReplaceNth(oldS, newS String, n int) String {
	if n < -1 || len(oldS) == 0 {
		return s
	}

	count, i := 0, 0

	for {
		pos := s[i:].Index(oldS)
		if pos == -1 {
			break
		}

		pos += i
		count++

		if count == n || (n == -1 && s[pos+len(oldS):].Index(oldS) == -1) {
			return s[:pos] + newS + s[pos+len(oldS):]
		}

		i = pos + len(oldS)
	}

	return s
}

// ContainsRegexp checks if the String contains a match for the specified regular expression pattern.
func (s String) ContainsRegexp(pattern String) Result[bool] {
	return ResultOf(regexp.MatchString(pattern.Std(), s.Std()))
}

// ContainsRegexpAny checks if the String contains a match for any of the specified regular
// expression patterns.
func (s String) ContainsRegexpAny(patterns ...String) Result[bool] {
	for _, pattern := range patterns {
		if r := s.ContainsRegexp(pattern); r.IsErr() || r.Ok() {
			return r
		}
	}

	return Ok(false)
}

// ContainsRegexpAll checks if the String contains a match for all of the specified regular expression patterns.
func (s String) ContainsRegexpAll(patterns ...String) Result[bool] {
	for _, pattern := range patterns {
		if r := s.ContainsRegexp(pattern); r.IsErr() || !r.Ok() {
			return r
		}
	}

	return Ok(true)
}

// Contains checks if the String contains the specified substring.
func (s String) Contains(substr String) bool {
	return strings.Contains(s.Std(), substr.Std())
}

// ContainsAny checks if the String contains any of the specified substrings.
func (s String) ContainsAny(substrs ...String) bool {
	for _, substr := range substrs {
		if s.Contains(substr) {
			return true
		}
	}

	return false
}

// ContainsAll checks if the given String contains all the specified substrings.
func (s String) ContainsAll(substrs ...String) bool {
	for _, substr := range substrs {
		if !s.Contains(substr) {
			return false
		}
	}

	return true
}

// ContainsAnyChars checks if the String contains any characters from the specified String.
func (s String) ContainsAnyChars(chars String) bool {
	return strings.ContainsAny(s.Std(), chars.Std())
}

// StartsWith checks if the String starts with any of the provided prefixes.
// The method accepts a variable number of arguments, allowing for checking against multiple
// prefixes at once. It iterates over the provided prefixes and uses the HasPrefix function from
// the strings package to check if the String starts with each prefix.
// The function returns true if the String starts with any of the prefixes, and false otherwise.
//
// Example usage:
//
//	s := g.String("http://example.com")
//	if s.StartsWith("http://", "https://") {
//	   // do something
//	}
func (s String) StartsWith(prefixes ...String) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(string(s), prefix.Std()) {
			return true
		}
	}

	return false
}

// EndsWith checks if the String ends with any of the provided suffixes.
// The method accepts a variable number of arguments, allowing for checking against multiple
// suffixes at once. It iterates over the provided suffixes and uses the HasSuffix function from
// the strings package to check if the String ends with each suffix.
// The function returns true if the String ends with any of the suffixes, and false otherwise.
//
// Example usage:
//
//	s := g.String("example.com")
//	if s.EndsWith(".com", ".net") {
//	   // do something
//	}
func (s String) EndsWith(suffixes ...String) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(string(s), suffix.Std()) {
			return true
		}
	}

	return false
}

// Split splits the String by the specified separator.
func (s String) Split(sep ...String) Slice[String] {
	var separator string
	if len(sep) != 0 {
		separator = sep[0].Std()
	}

	return SliceMap(strings.Split(s.Std(), separator), NewString)
}

// SplitLines splits the String by lines.
func (s String) SplitLines() Slice[String] { return s.TrimSpace().Split("\n") }

// SplitN splits the String into substrings using the provided separator and returns an Slice[String] of the results.
// The n parameter controls the number of substrings to return:
// - If n is negative, there is no limit on the number of substrings returned.
// - If n is zero, an empty Slice[String] is returned.
// - If n is positive, at most n substrings are returned.
func (s String) SplitN(sep String, n int) Slice[String] {
	return SliceMap(strings.SplitN(s.Std(), sep.Std(), n), NewString)
}

// SplitRegexp splits the String into substrings using the provided regular expression pattern and returns an Slice[String] of the results.
// The regular expression pattern is provided as a regexp.Regexp parameter.
func (s String) SplitRegexp(pattern regexp.Regexp) Slice[String] {
	return SliceMap(pattern.Split(s.Std(), -1), NewString)
}

// SplitRegexpN splits the String into substrings using the provided regular expression pattern and returns an Slice[String] of the results.
// The regular expression pattern is provided as a regexp.Regexp parameter.
// The n parameter controls the number of substrings to return:
// - If n is negative, there is no limit on the number of substrings returned.
// - If n is zero, an empty Slice[String] is returned.
// - If n is positive, at most n substrings are returned.
func (s String) SplitRegexpN(pattern regexp.Regexp, n int) Option[Slice[String]] {
	result := SliceMap(pattern.Split(s.Std(), n), NewString)
	if result.Empty() {
		return None[Slice[String]]()
	}

	return Some(result)
}

// Fields splits the String into a slice of substrings, removing any whitespace.
func (s String) Fields() Slice[String] {
	return SliceMap(strings.Fields(s.Std()), NewString)
}

// Chunks splits the String into chunks of the specified size.
//
// This function iterates through the String, creating new String chunks of the specified size.
// If size is less than or equal to 0 or the String is empty,
// it returns an empty Slice[String].
// If size is greater than or equal to the length of the String,
// it returns an Slice[String] containing the original String.
//
// Parameters:
//
// - size (int): The size of the chunks to split the String into.
//
// Returns:
//
// - Slice[String]: A slice of String chunks of the specified size.
//
// Example usage:
//
//	text := g.String("Hello, World!")
//	chunks := text.Chunks(4)
//
// chunks contains {"Hell", "o, W", "orld", "!"}.
func (s String) Chunks(size int) Slice[String] {
	if size <= 0 || s.Empty() {
		return nil
	}

	if size >= len(s) {
		return Slice[String]{s}
	}

	return SliceMap(s.Split().Iter().Chunks(size).Collect(), func(ch Slice[String]) String { return ch.Join() })
}

// Cut returns two String values. The first String contains the remainder of the
// original String after the cut. The second String contains the text between the
// first occurrences of the 'start' and 'end' strings, with tags removed if specified.
//
// The function searches for the 'start' and 'end' strings within the String.
// If both are found, it returns the first String containing the remainder of the
// original String after the cut, followed by the second String containing the text
// between the first occurrences of 'start' and 'end' with tags removed if specified.
//
// If either 'start' or 'end' is empty or not found in the String, it returns the
// original String as the second String, and an empty String as the first.
//
// Parameters:
//
// - start (String): The String marking the beginning of the text to be cut.
//
// - end (String): The String marking the end of the text to be cut.
//
//   - rmtags (bool, optional): An optional boolean parameter indicating whether
//     to remove 'start' and 'end' tags from the cut text. Defaults to false.
//
// Returns:
//
//   - String: The first String containing the remainder of the original String
//     after the cut, with tags removed if specified,
//     or an empty String if 'start' or 'end' is empty or not found.
//
//   - String: The second String containing the text between the first occurrences of
//     'start' and 'end', or the original String if 'start' or 'end' is empty or not found.
//
// Example usage:
//
//	s := g.String("Hello, [world]! How are you?")
//	remainder, cut := s.Cut("[", "]")
//	// remainder: "Hello, ! How are you?"
//	// cut: "world"
func (s String) Cut(start, end String, rmtags ...bool) (String, String) {
	if start.Empty() || end.Empty() {
		return s, ""
	}

	startIndex := s.Index(start)
	if startIndex == -1 {
		return s, ""
	}

	endIndex := s[startIndex+len(start):].Index(end)
	if endIndex == -1 {
		return s, ""
	}

	cut := s[startIndex+len(start) : startIndex+len(start)+endIndex]

	startCutIndex := startIndex
	endCutIndex := startIndex + len(start) + endIndex

	if len(rmtags) != 0 && !rmtags[0] {
		startCutIndex += len(start)
	} else {
		endCutIndex += len(end)
	}

	remainder := s[:startCutIndex] + s[endCutIndex:]

	return remainder, cut
}

// Similarity calculates the similarity between two Strings using the
// Levenshtein distance algorithm and returns the similarity percentage as an Float.
//
// The function compares two Strings using the Levenshtein distance,
// which measures the difference between two sequences by counting the number
// of single-character edits required to change one sequence into the other.
// The similarity is then calculated by normalizing the distance by the maximum
// length of the two input Strings.
//
// Parameters:
//
// - str (String): The String to compare with s.
//
// Returns:
//
// - Float: The similarity percentage between the two Strings as a value between 0 and 100.
//
// Example usage:
//
//	s1 := g.String("kitten")
//	s2 := g.String("sitting")
//	similarity := s1.Similarity(s2) // 57.14285714285714
func (s String) Similarity(str String) Float {
	if s.Eq(str) {
		return 100
	}

	if s.Empty() || str.Empty() {
		return 0
	}

	s1 := s.ToRunes()
	s2 := str.ToRunes()

	lenS1 := s.LenRunes()
	lenS2 := str.LenRunes()

	if lenS1 > lenS2 {
		s1, s2, lenS1, lenS2 = s2, s1, lenS2, lenS1
	}

	distance := NewSlice[Int](lenS1 + 1)

	for i, r2 := range s2 {
		prev := Int(i).Add(1)

		for j, r1 := range s1 {
			current := distance[j]
			if r2 != r1 {
				current = distance[j].Add(1).Min(prev.Add(1)).Min(distance[j+1].Add(1))
			}

			distance[j], prev = prev, current
		}

		distance[lenS1] = prev
	}

	return Float(1).Sub(distance[lenS1].ToFloat().Div(Int(lenS1).Max(Int(lenS2)).ToFloat())).Mul(100)
}

// Compare compares two Strings and returns an Int indicating their relative order.
// The result will be 0 if s==str, -1 if s < str, and +1 if s > str.
func (s String) Compare(str String) Int { return Int(cmp.Compare(s, str)) }

// Append appends the specified String to the current String.
func (s String) Append(str String) String { return s + str }

// Prepend prepends the specified String to the current String.
func (s String) Prepend(str String) String { return str + s }

// ContainsRune checks if the String contains the specified rune.
func (s String) ContainsRune(r rune) bool { return strings.ContainsRune(s.Std(), r) }

// Count returns the number of non-overlapping instances of the substring in the String.
func (s String) Count(substr String) int { return strings.Count(s.Std(), substr.Std()) }

// Empty checks if the String is empty.
func (s String) Empty() bool { return len(s) == 0 }

// Eq checks if two Strings are equal.
func (s String) Eq(str String) bool { return s.Compare(str).Eq(0) }

// EqFold compares two String strings case-insensitively.
func (s String) EqFold(str String) bool { return strings.EqualFold(s.Std(), str.Std()) }

// Gt checks if the String is greater than the specified String.
func (s String) Gt(str String) bool { return s.Compare(str).Gt(0) }

// Bytes returns the String as an Bytes.
func (s String) ToBytes() Bytes { return Bytes(s) }

// Index returns the index of the first instance of the specified substring in the String, or -1
// if substr is not present in s.
func (s String) Index(substr String) int { return strings.Index(s.Std(), substr.Std()) }

// IndexRegexp searches for the first occurrence of the regular expression pattern in the String.
// If a match is found, it returns an Option containing an Slice with the start and end indices of the match.
// If no match is found, it returns None.
func (s String) IndexRegexp(pattern *regexp.Regexp) Option[Slice[Int]] {
	result := SliceMap(pattern.FindStringIndex(s.Std()), NewInt)
	if result.Empty() {
		return None[Slice[Int]]()
	}

	return Some(result)
}

// FindAllRegexp searches the String for all occurrences of the regular expression pattern
// and returns an Option[Slice[String]] containing a slice of matched substrings.
// If no matches are found, the Option[Slice[String]] will be None.
func (s String) FindAllRegexp(pattern *regexp.Regexp) Option[Slice[String]] {
	return s.FindAllRegexpN(pattern, -1)
}

// FindAllRegexpN searches the String for up to n occurrences of the regular expression pattern
// and returns an Option[Slice[String]] containing a slice of matched substrings.
// If no matches are found, the Option[Slice[String]] will be None.
// If n is negative, all occurrences will be returned.
func (s String) FindAllRegexpN(pattern *regexp.Regexp, n int) Option[Slice[String]] {
	result := SliceMap(pattern.FindAllString(s.Std(), n), NewString)
	if result.Empty() {
		return None[Slice[String]]()
	}

	return Some(result)
}

// FindSubmatchRegexp searches the String for the first occurrence of the regular expression pattern
// and returns an Option[Slice[String]] containing the matched substrings and submatches.
// The Option will contain an Slice[String] with the full match at index 0, followed by any captured submatches.
// If no match is found, it returns None.
func (s String) FindSubmatchRegexp(pattern *regexp.Regexp) Option[Slice[String]] {
	result := SliceMap(pattern.FindStringSubmatch(s.Std()), NewString)
	if result.Empty() {
		return None[Slice[String]]()
	}

	return Some(result)
}

// FindAllSubmatchRegexp searches the String for all occurrences of the regular expression pattern
// and returns an Option[Slice[Slice[String]]] containing the matched substrings and submatches.
// The Option[Slice[Slice[String]]] will contain an Slice[String] for each match,
// where each Slice[String] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Slice[String]]] will be None.
// This method is equivalent to calling SubmatchAllRegexpN with n = -1, which means it finds all occurrences.
func (s String) FindAllSubmatchRegexp(pattern *regexp.Regexp) Option[Slice[Slice[String]]] {
	return s.FindAllSubmatchRegexpN(pattern, -1)
}

// FindAllSubmatchRegexpN searches the String for occurrences of the regular expression pattern
// and returns an Option[Slice[Slice[String]]] containing the matched substrings and submatches.
// The Option[Slice[Slice[String]]] will contain an Slice[String] for each match,
// where each Slice[String] will contain the full match at index 0, followed by any captured submatches.
// If no match is found, the Option[Slice[Slice[String]]] will be None.
// The 'n' parameter specifies the maximum number of matches to find. If n is negative, it finds all occurrences.
func (s String) FindAllSubmatchRegexpN(pattern *regexp.Regexp, n int) Option[Slice[Slice[String]]] {
	var result Slice[Slice[String]]

	for _, v := range pattern.FindAllStringSubmatch(s.Std(), n) {
		result = append(result, SliceMap(v, NewString))
	}

	if result.Empty() {
		return None[Slice[Slice[String]]]()
	}

	return Some(result)
}

// LastIndex returns the index of the last instance of the specified substring in the String, or -1
// if substr is not present in s.
func (s String) LastIndex(substr String) int {
	return strings.LastIndex(s.Std(), substr.Std())
}

// IndexRune returns the index of the first instance of the specified rune in the String.
func (s String) IndexRune(r rune) int { return strings.IndexRune(s.Std(), r) }

// Len returns the length of the String.
func (s String) Len() int { return len(s) }

// LenRunes returns the number of runes in the String.
func (s String) LenRunes() int { return utf8.RuneCountInString(s.Std()) }

// Lt checks if the String is less than the specified String.
func (s String) Lt(str String) bool { return s.Compare(str).Lt(0) }

// Map applies the provided function to all runes in the String and returns the resulting String.
func (s String) Map(fn func(rune) rune) String { return String(strings.Map(fn, s.Std())) }

// NormalizeNFC returns a new String with its Unicode characters normalized using the NFC form.
func (s String) NormalizeNFC() String { return String(norm.NFC.String(s.Std())) }

// Ne checks if two Strings are not equal.
func (s String) Ne(str String) bool { return !s.Eq(str) }

// NotEmpty checks if the String is not empty.
func (s String) NotEmpty() bool { return s.Len() != 0 }

// Reader returns a *strings.Reader initialized with the content of String.
func (s String) Reader() *strings.Reader { return strings.NewReader(s.Std()) }

// Repeat returns a new String consisting of the specified count of the original String.
func (s String) Repeat(count int) String { return String(strings.Repeat(s.Std(), count)) }

// Reverse reverses the String.
func (s String) Reverse() String { return s.ToBytes().Reverse().ToString() }

// ToRunes returns the String as a slice of runes.
func (s String) ToRunes() Slice[rune] { return []rune(s) }

// Chars returns the individual characters of the String as a slice of Strings.
// Each element in the returned slice represents a single character in the original String.
func (s String) Chars() Slice[String] { return s.Split() }

// Std returns the String as a string.
func (s String) Std() string { return string(s) }

// TrimSpace trims whitespace from the beginning and end of the String.
func (s String) TrimSpace() String { return String(strings.TrimSpace(s.Std())) }

// Format applies a specified format to the String object.
func (s String) Format(format String) String { return Sprintf(format, s) }

// LeftJustify justifies the String to the left by adding padding to the right, up to the
// specified length. If the length of the String is already greater than or equal to the specified
// length, or the pad is empty, the original String is returned.
//
// The padding String is repeated as necessary to fill the remaining length.
// The padding is added to the right of the String.
//
// Parameters:
//   - length: The desired length of the resulting justified String.
//   - pad: The String used as padding.
//
// Example usage:
//
//	s := g.String("Hello")
//	result := s.LeftJustify(10, "...")
//	// result: "Hello....."
func (s String) LeftJustify(length int, pad String) String {
	if s.LenRunes() >= length || pad.Eq("") {
		return s
	}

	var output strings.Builder

	_, _ = output.WriteString(s.Std())
	writePadding(&output, pad, pad.LenRunes(), length-s.LenRunes())

	return String(output.String())
}

// RightJustify justifies the String to the right by adding padding to the left, up to the
// specified length. If the length of the String is already greater than or equal to the specified
// length, or the pad is empty, the original String is returned.
//
// The padding String is repeated as necessary to fill the remaining length.
// The padding is added to the left of the String.
//
// Parameters:
//   - length: The desired length of the resulting justified String.
//   - pad: The String used as padding.
//
// Example usage:
//
//	s := g.String("Hello")
//	result := s.RightJustify(10, "...")
//	// result: ".....Hello"
func (s String) RightJustify(length int, pad String) String {
	if s.LenRunes() >= length || pad.Empty() {
		return s
	}

	var output strings.Builder

	writePadding(&output, pad, pad.LenRunes(), length-s.LenRunes())
	_, _ = output.WriteString(s.Std())

	return String(output.String())
}

// Center justifies the String by adding padding on both sides, up to the specified length.
// If the length of the String is already greater than or equal to the specified length, or the
// pad is empty, the original String is returned.
//
// The padding String is repeated as necessary to evenly distribute the remaining length on both
// sides.
// The padding is added to the left and right of the String.
//
// Parameters:
//   - length: The desired length of the resulting justified String.
//   - pad: The String used as padding.
//
// Example usage:
//
//	s := g.String("Hello")
//	result := s.Center(10, "...")
//	// result: "..Hello..."
func (s String) Center(length int, pad String) String {
	if s.LenRunes() >= length || pad.Empty() {
		return s
	}

	var output strings.Builder

	remains := length - s.LenRunes()
	writePadding(&output, pad, pad.LenRunes(), remains/2)
	_, _ = output.WriteString(s.Std())
	writePadding(&output, pad, pad.LenRunes(), (remains+1)/2)

	return String(output.String())
}

// writePadding writes the padding String to the output Builder to fill the remaining length.
// It repeats the padding String as necessary and appends any remaining runes from the padding
// String.
func writePadding(output *strings.Builder, pad String, padlen, remains int) {
	if repeats := remains / padlen; repeats > 0 {
		_, _ = output.WriteString(pad.Repeat(repeats).Std())
	}

	padrunes := pad.ToRunes()
	for i := range remains % padlen {
		_, _ = output.WriteRune(padrunes[i])
	}
}

// Print prints the content of the String to the standard output (console)
// and returns the String unchanged.
func (s String) Print() String { fmt.Println(s); return s }
