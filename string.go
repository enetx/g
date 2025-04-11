package g

import (
	"math/big"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// NewString creates a new String from the provided string.
func NewString[T ~string | rune | byte | ~[]rune | ~[]byte](str T) String { return String(str) }

// Ptr returns a pointer to the current String value.
func (s String) Ptr() *String { return &s }

// Transform applies a transformation function to the String and returns the result.
func (s String) Transform(fn func(String) String) String { return fn(s) }

// Builder returns a new Builder initialized with the content of the String.
func (s String) Builder() *Builder { return NewBuilder().Write(s) }

// Min returns the minimum of Strings.
func (s String) Min(b ...String) String { return cmp.Min(append(b, s)...) }

// Max returns the maximum of Strings.
func (s String) Max(b ...String) String { return cmp.Max(append(b, s)...) }

// Random generates a random String of the specified length, selecting characters from predefined sets.
// If additional character sets are provided, only those will be used; the default set (ASCII_LETTERS and DIGITS)
// is excluded unless explicitly provided.
//
// Parameters:
// - count (Int): Length of the random String to generate.
// - letters (...String): Additional character sets to consider for generating the random String (optional).
//
// Returns:
// - String: Randomly generated String with the specified length.
//
// Example usage:
//
//	randomString := g.String.Random(10)
//	randomString contains a random String with 10 characters.
func (String) Random(length Int, letters ...String) String {
	var chars Slice[rune]

	if len(letters) != 0 {
		chars = letters[0].Runes()
	} else {
		chars = (ASCII_LETTERS + DIGITS).Runes()
	}

	result := NewBuilder()

	for range length {
		result.WriteRune(chars.Random())
	}

	return result.String()
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

// ToBigInt attempts to convert the String receiver into an Option containing a *big.Int.
// This function assumes the string represents a numerical value, which can be in decimal,
// hexadecimal (prefixed with "0x"), or octal (prefixed with "0") format. The function
// leverages the SetString method of the math/big package, automatically detecting the
// numeric base when set to 0.
//
// If the string is correctly formatted and represents a valid number, ToBigInt returns
// a Some containing the *big.Int parsed from the string. If the string is empty, contains
// invalid characters, or does not conform to a recognizable numeric format, ToBigInt
// returns a None, indicating that the conversion was unsuccessful.
//
// Returns:
//   - An Option[*big.Int] encapsulating the conversion result. It returns Some[*big.Int]
//     with the parsed value if successful, otherwise None[*big.Int] if the parsing fails.
func (s String) ToBigInt() Option[*big.Int] {
	if bigInt, ok := new(big.Int).SetString(s.Std(), 0); ok {
		return Some(bigInt)
	}

	return None[*big.Int]()
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
func (s String) Title() String { return String(cases.Title(language.English).String(s.Std())) }

// Lower returns the String in lowercase.
func (s String) Lower() String { return String(cases.Lower(language.English).String(s.Std())) }

// Upper returns the String in uppercase.
func (s String) Upper() String { return String(cases.Upper(language.English).String(s.Std())) }

// Trim removes leading and trailing white space from the String.
func (s String) Trim() String { return String(strings.TrimSpace(s.Std())) }

// TrimStart removes leading white space from the String.
func (s String) TrimStart() String { return trimStringStart(s) }

// TrimEnd removes trailing white space from the String.
func (s String) TrimEnd() String { return trimStringEnd(s) }

// TrimSet removes the specified set of characters from both the beginning and end of the String.
func (s String) TrimSet(cutset String) String {
	return String(strings.Trim(s.Std(), cutset.Std()))
}

// TrimStartSet removes the specified set of characters from the beginning of the String.
func (s String) TrimStartSet(cutset String) String {
	return String(strings.TrimLeft(s.Std(), cutset.Std()))
}

// TrimEndSet removes the specified set of characters from the end of the String.
func (s String) TrimEndSet(cutset String) String {
	return String(strings.TrimRight(s.Std(), cutset.Std()))
}

// StripPrefix trims the specified prefix from the String.
func (s String) StripPrefix(prefix String) String {
	return String(strings.TrimPrefix(s.Std(), prefix.Std()))
}

// StripSuffix trims the specified suffix from the String.
func (s String) StripSuffix(suffix String) String {
	return String(strings.TrimSuffix(s.Std(), suffix.Std()))
}

// Replace replaces the 'oldS' String with the 'newS' String for the specified number of
// occurrences.
func (s String) Replace(oldS, newS String, n Int) String {
	return String(strings.Replace(s.Std(), oldS.Std(), newS.Std(), n.Std()))
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
	on := Slice[String](oldnew).ToStringSlice()
	return String(strings.NewReplacer(on...).Replace(s.Std()))
}

// Remove removes all occurrences of specified substrings from the String.
//
// Parameters:
//
// - matches ...String: Substrings to be removed from the string. Specify as many substrings as needed.
//
// Returns:
//
// - String: A new string with all specified substrings removed.
//
// Example usage:
//
//	original := g.String("Hello, world! This is a test.")
//	modified := original.Remove(
//	    "Hello",
//	    "test",
//	)
//	// modified contains ", world! This is a ."
func (s String) Remove(matches ...String) String {
	for _, match := range matches {
		s = s.ReplaceAll(match, "")
	}

	return s
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
func (s String) ReplaceNth(oldS, newS String, n Int) String {
	if n < -1 || len(oldS) == 0 {
		return s
	}

	count, i := Int(0), Int(0)

	for {
		pos := s[i:].Index(oldS)
		if pos == -1 {
			break
		}

		pos += i
		count++

		if count == n || (n == -1 && s[pos+oldS.Len():].Index(oldS) == -1) {
			return s[:pos] + newS + s[pos+oldS.Len():]
		}

		i = pos + oldS.Len()
	}

	return s
}

// Contains checks if the String contains the specified substring.
func (s String) Contains(substr String) bool { return f.Contains(substr)(s) }

// ContainsAny checks if the String contains any of the specified substrings.
func (s String) ContainsAny(substrs ...String) bool {
	return Slice[String](substrs).
		Iter().
		Any(func(substr String) bool { return s.Contains(substr) })
}

// ContainsAll checks if the given String contains all the specified substrings.
func (s String) ContainsAll(substrs ...String) bool {
	return Slice[String](substrs).
		Iter().
		All(func(substr String) bool { return s.Contains(substr) })
}

// ContainsAnyChars checks if the String contains any characters from the specified String.
func (s String) ContainsAnyChars(chars String) bool { return f.ContainsAnyChars(chars)(s) }

// StartsWith checks if the String starts with the specified prefix.
// It uses a higher-order function to perform the check.
func (s String) StartsWith(prefix String) bool { return f.StartsWith(prefix)(s) }

// StartsWithAny checks if the String starts with any of the provided prefixes.
// The method accepts a variable number of arguments, allowing for checking against multiple
// prefixes at once. It iterates over the provided prefixes and uses the HasPrefix function from
// the strings package to check if the String starts with each prefix.
// The function returns true if the String starts with any of the prefixes, and false otherwise.
//
// Example usage:
//
//	s := g.String("http://example.com")
//	if s.StartsWithAny("http://", "https://") {
//	   // do something
//	}
func (s String) StartsWithAny(prefixes ...String) bool {
	return Slice[String](prefixes).
		Iter().
		Any(func(prefix String) bool { return s.StartsWith(prefix) })
}

// EndsWith checks if the String ends with the specified suffix.
// It uses a higher-order function to perform the check.
func (s String) EndsWith(suffix String) bool { return f.EndsWith(suffix)(s) }

// EndsWithAny checks if the String ends with any of the provided suffixes.
// The method accepts a variable number of arguments, allowing for checking against multiple
// suffixes at once. It iterates over the provided suffixes and uses the HasSuffix function from
// the strings package to check if the String ends with each suffix.
// The function returns true if the String ends with any of the suffixes, and false otherwise.
//
// Example usage:
//
//	s := g.String("example.com")
//	if s.EndsWithAny(".com", ".net") {
//	   // do something
//	}
func (s String) EndsWithAny(suffixes ...String) bool {
	return Slice[String](suffixes).
		Iter().
		Any(func(suffix String) bool { return s.EndsWith(suffix) })
}

// Lines splits the String by lines and returns the iterator.
func (s String) Lines() SeqSlice[String] { return linesString(s) }

// Fields splits the String into a slice of substrings, removing any whitespace, and returns the iterator.
func (s String) Fields() SeqSlice[String] { return fieldsString(s) }

// FieldsBy splits the String into a slice of substrings using a custom function to determine the field boundaries,
// and returns the iterator.
func (s String) FieldsBy(fn func(r rune) bool) SeqSlice[String] { return fieldsbyString(s, fn) }

// Split splits the String by the specified separator and returns the iterator.
func (s String) Split(sep ...String) SeqSlice[String] {
	var separator String
	if len(sep) != 0 {
		separator = sep[0]
	}

	return splitString(s, separator, 0)
}

// SplitAfter splits the String after each instance of the specified separator and returns the iterator.
func (s String) SplitAfter(sep String) SeqSlice[String] { return splitString(s, sep, sep.Len()) }

// SplitN splits the String into substrings using the provided separator and returns an Slice[String] of the results.
// The n parameter controls the number of substrings to return:
// - If n is negative, there is no limit on the number of substrings returned.
// - If n is zero, an empty Slice[String] is returned.
// - If n is positive, at most n substrings are returned.
func (s String) SplitN(sep String, n Int) Slice[String] {
	return TransformSlice(strings.SplitN(s.Std(), sep.Std(), n.Std()), NewString)
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
// - size (Int): The size of the chunks to split the String into.
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
func (s String) Chunks(size Int) Slice[String] {
	if size <= 0 || s.Empty() {
		return nil
	}

	if size >= s.Len() {
		return Slice[String]{s}
	}

	return TransformSlice(s.Split().Chunks(size).Collect(), func(ch Slice[String]) String { return ch.Join() })
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

	endIndex := s[startIndex+start.Len():].Index(end)
	if endIndex == -1 {
		return s, ""
	}

	cut := s[startIndex+start.Len() : startIndex+start.Len()+endIndex]

	startCutIndex := startIndex
	endCutIndex := startIndex + start.Len() + endIndex

	if len(rmtags) != 0 && !rmtags[0] {
		startCutIndex += start.Len()
	} else {
		endCutIndex += end.Len()
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

	s1 := s.Runes()
	s2 := str.Runes()

	lenS1 := s.LenRunes()
	lenS2 := str.LenRunes()

	if lenS1 > lenS2 {
		s1, s2, lenS1, lenS2 = s2, s1, lenS2, lenS1
	}

	distance := NewSlice[Int](lenS1 + 1)

	for i, r2 := range s2 {
		prev := Int(i) + 1

		for j, r1 := range s1 {
			current := distance[j]
			if r2 != r1 {
				current = distance[j].Add(1).Min(prev + 1).Min(distance[j+1] + 1)
			}

			distance[j], prev = prev, current
		}

		distance[lenS1] = prev
	}

	return Float(1).Sub(distance[lenS1].Float() / lenS1.Max(lenS2).Float()).Mul(100)
}

// Cmp compares two Strings and returns an cmp.Ordering indicating their relative order.
// The result will be cmp.Equal if s==str, cmp.Less if s < str, and cmp.Greater if s > str.
func (s String) Cmp(str String) cmp.Ordering { return cmp.Cmp(s, str) }

// Append appends the specified String to the current String.
func (s String) Append(str String) String { return s + str }

// Prepend prepends the specified String to the current String.
func (s String) Prepend(str String) String { return str + s }

// ContainsRune checks if the String contains the specified rune.
func (s String) ContainsRune(r rune) bool { return strings.ContainsRune(s.Std(), r) }

// Count returns the number of non-overlapping instances of the substring in the String.
func (s String) Count(substr String) Int { return Int(strings.Count(s.Std(), substr.Std())) }

// Empty checks if the String is empty.
func (s String) Empty() bool { return len(s) == 0 }

// Eq checks if two Strings are equal.
func (s String) Eq(str String) bool { return s == str }

// EqFold compares two String strings case-insensitively.
func (s String) EqFold(str String) bool { return strings.EqualFold(s.Std(), str.Std()) }

// Gt checks if the String is greater than the specified String.
func (s String) Gt(str String) bool { return s > str }

// Gte checks if the String is greater than or equal to the specified String.
func (s String) Gte(str String) bool { return s >= str }

// Bytes returns the String as an Bytes.
func (s String) Bytes() Bytes { return Bytes(s) }

// Index returns the index of the first instance of the specified substring in the String, or -1
// if substr is not present in s.
func (s String) Index(substr String) Int { return Int(strings.Index(s.Std(), substr.Std())) }

// LastIndex returns the index of the last instance of the specified substring in the String, or -1
// if substr is not present in s.
func (s String) LastIndex(substr String) Int { return Int(strings.LastIndex(s.Std(), substr.Std())) }

// IndexRune returns the index of the first instance of the specified rune in the String.
func (s String) IndexRune(r rune) Int { return Int(strings.IndexRune(s.Std(), r)) }

// Len returns the length of the String.
func (s String) Len() Int { return Int(len(s)) }

// LenRunes returns the number of runes in the String.
func (s String) LenRunes() Int { return Int(utf8.RuneCountInString(s.Std())) }

// Lt checks if the String is less than the specified String.
func (s String) Lt(str String) bool { return s < str }

// Lte checks if the String is less than or equal to the specified String.
func (s String) Lte(str String) bool { return s <= str }

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
func (s String) Repeat(count Int) String { return String(strings.Repeat(s.Std(), count.Std())) }

// Reverse reverses the String.
func (s String) Reverse() String { return s.Bytes().Reverse().String() }

// Runes returns the String as a slice of runes.
func (s String) Runes() Slice[rune] { return []rune(s) }

// Chars splits the String into individual characters and returns the iterator.
func (s String) Chars() SeqSlice[String] { return s.Split() }

// SubString extracts a substring from the String starting at the 'start' index and ending before the 'end' index.
// The function also supports an optional 'step' parameter to define the increment between indices in the substring.
// If 'start' or 'end' index is negative, they represent positions relative to the end of the String:
// - A negative 'start' index indicates the position from the end of the String, moving backward.
// - A negative 'end' index indicates the position from the end of the String.
// The function ensures that indices are adjusted to fall within the valid range of the String's length.
// If indices are out of bounds or if 'start' exceeds 'end', the function returns the original String unmodified.
func (s String) SubString(start, end Int, step ...Int) String {
	return String(s.Runes().SubSlice(start, end, step...))
}

// Std returns the String as a string.
func (s String) Std() string { return string(s) }

// Sprintf applies a specified format to the String object.
func (s String) Sprintf(template String) String { return Sprintf(template, s) }

// Truncate shortens the String to the specified maximum length. If the String exceeds the
// specified length, it is truncated, and an ellipsis ("...") is appended to indicate the truncation.
//
// If the length of the String is less than or equal to the specified maximum length, the
// original String is returned unchanged.
//
// The method respects Unicode characters and truncates based on the number of runes,
// not bytes.
//
// Parameters:
//   - max: The maximum number of runes allowed in the resulting String.
//
// Returns:
//   - A new String truncated to the specified maximum length with "..." appended
//     if truncation occurs. Otherwise, returns the original String.
//
// Example usage:
//
//	s := g.String("Hello, World!")
//	result := s.Truncate(5)
//	// result: "Hello..."
//
//	s2 := g.String("Short")
//	result2 := s2.Truncate(10)
//	// result2: "Short"
//
//	s3 := g.String("😊😊😊😊😊")
//	result3 := s3.Truncate(3)
//	// result3: "😊😊😊..."
func (s String) Truncate(max Int) String {
	if max.IsNegative() || s.LenRunes().Lte(max) {
		return s
	}

	return String(s.Runes().SubSlice(0, max)).Append("...")
}

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
func (s String) LeftJustify(length Int, pad String) String {
	if s.LenRunes() >= length || pad.Eq("") {
		return s
	}

	output := NewBuilder()

	_ = output.Write(s)
	writePadding(output, pad, pad.LenRunes(), length-s.LenRunes())

	return output.String()
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
func (s String) RightJustify(length Int, pad String) String {
	if s.LenRunes() >= length || pad.Empty() {
		return s
	}

	output := NewBuilder()

	writePadding(output, pad, pad.LenRunes(), length-s.LenRunes())
	_ = output.Write(s)

	return output.String()
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
func (s String) Center(length Int, pad String) String {
	if s.LenRunes() >= length || pad.Empty() {
		return s
	}

	output := NewBuilder()

	remains := length - s.LenRunes()

	writePadding(output, pad, pad.LenRunes(), remains/2)
	_ = output.Write(s)
	writePadding(output, pad, pad.LenRunes(), (remains+1)/2)

	return output.String()
}

// writePadding writes the padding String to the output Builder to fill the remaining length.
// It repeats the padding String as necessary and appends any remaining runes from the padding
// String.
func writePadding(output *Builder, pad String, padlen, remains Int) {
	if repeats := remains / padlen; repeats > 0 {
		_ = output.Write(pad.Repeat(repeats))
	}

	padrunes := pad.Runes()
	for i := range remains % padlen {
		_ = output.WriteRune(padrunes[i])
	}
}

// Print writes the content of the String to the standard output (console)
// and returns the String unchanged.
func (s String) Print() String { Print(s); return s }

// Println writes the content of the String to the standard output (console) with a newline
// and returns the String unchanged.
func (s String) Println() String { Println(s); return s }

// trimStringStart returns a slice of the string s, with all leading
// white space removed, as defined by Unicode.
func trimStringStart(s String) String {
	start := 0

	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			return trimStringFuncStart(s[start:], unicode.IsSpace)
		}

		if asciiSpace[c] == 0 {
			break
		}
	}

	return s[start:]
}

// trimStringEnd returns a slice of the string s, with all trailing
// white space removed, as defined by Unicode.
func trimStringEnd(s String) String {
	stop := len(s)

	for ; stop > 0; stop-- {
		c := s[stop-1]
		if c >= utf8.RuneSelf {
			return trimStringFuncEnd(s[:stop], unicode.IsSpace)
		}

		if asciiSpace[c] == 0 {
			break
		}
	}

	return s[:stop]
}

// Helper function to trim leading characters using a unicode function.
func trimStringFuncStart(s String, fn func(rune) bool) String {
	start := 0

	for start < len(s) {
		r, size := utf8.DecodeRuneInString(s[start:].Std())
		if !fn(r) {
			break
		}

		start += size
	}

	return s[start:]
}

// Helper function to trim trailing characters using a unicode function.
func trimStringFuncEnd(s String, fn func(rune) bool) String {
	stop := len(s)

	for stop > 0 {
		r, size := utf8.DecodeLastRuneInString(s[:stop].Std())
		if !fn(r) {
			break
		}

		stop -= size
	}

	return s[:stop]
}
