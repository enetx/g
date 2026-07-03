package g

import "strings"

// IsASCII checks if all bytes in the Bytes are ASCII bytes.
func (bs Bytes) IsASCII() bool {
	for i := range bs {
		if bs[i] >= 0x80 {
			return false
		}
	}

	return true
}

// IsDigit checks if the Bytes is non-empty and all bytes are ASCII digits ('0'-'9').
// Unlike String.IsDigit, which is rune-aware and accepts Unicode digits, this method
// operates byte-wise and only recognizes ASCII digits.
func (bs Bytes) IsDigit() bool {
	if bs.IsEmpty() {
		return false
	}

	for i := range bs {
		if bs[i] < '0' || bs[i] > '9' {
			return false
		}
	}

	return true
}

// ReplaceMulti performs multiple replacements within the Bytes.
//
// The replacements are provided as pairs of old and new Bytes, in the same
// order as in String.ReplaceMulti. Replacements are performed in a single pass:
// at each position the earliest, longest match among the old patterns wins, and
// matches do not overlap. The number of arguments must be even; otherwise the
// method panics.
//
// Both this method and String.ReplaceMulti match patterns as raw byte
// sequences, so the two versions behave identically on the same data.
//
// Parameters:
//
// - oldnew ...Bytes: Pairs of Bytes to be replaced. Specify as many pairs as needed.
//
// Returns:
//
// - Bytes: A new Bytes with replacements applied. The receiver is not modified.
//
// Example usage:
//
//	original := g.Bytes("Hello, world! This is a test.")
//	replaced := original.ReplaceMulti(
//	    g.Bytes("Hello"), g.Bytes("Greetings"),
//	    g.Bytes("world"), g.Bytes("universe"),
//	    g.Bytes("test"), g.Bytes("example"),
//	)
//	// replaced contains "Greetings, universe! This is a example."
func (bs Bytes) ReplaceMulti(oldnew ...Bytes) Bytes {
	pairs := make([]string, len(oldnew))
	for i, b := range oldnew {
		pairs[i] = string(b)
	}

	return Bytes(strings.NewReplacer(pairs...).Replace(bs.StringUnsafe().Std()))
}

// Remove removes all occurrences of the specified patterns from the Bytes.
//
// Both this method and String.Remove match patterns as raw byte sequences,
// so the two versions behave identically on the same data.
//
// Parameters:
//
// - patterns ...Bytes: Patterns to be removed from the Bytes. Specify as many patterns as needed.
//
// Returns:
//
//   - Bytes: A new Bytes with all specified patterns removed. The receiver is
//     not modified. If no patterns are given, the original Bytes is returned.
//
// Example usage:
//
//	original := g.Bytes("Hello, world! This is a test.")
//	modified := original.Remove(
//	    g.Bytes("Hello"),
//	    g.Bytes("test"),
//	)
//	// modified contains ", world! This is a ."
func (bs Bytes) Remove(patterns ...Bytes) Bytes {
	if len(patterns) == 0 {
		return bs
	}

	pairs := make([]string, len(patterns)*2)
	for i, pattern := range patterns {
		pairs[i*2] = string(pattern)
		pairs[i*2+1] = ""
	}

	return Bytes(strings.NewReplacer(pairs...).Replace(bs.StringUnsafe().Std()))
}

// ReplaceNth returns a new Bytes with the nth occurrence of oldB
// replaced with newB. If there aren't enough occurrences of oldB, the
// original Bytes is returned. If n is less than -1, the original Bytes
// is also returned. If n is -1, the last occurrence of oldB is replaced with newB.
//
// Both this method and String.ReplaceNth match patterns as raw byte sequences,
// so the two versions behave identically on the same data.
//
// Returns:
//
//   - Bytes: A new Bytes with the nth occurrence of oldB replaced with newB.
//     The receiver is not modified.
//
// Example usage:
//
//	bs := g.Bytes("The quick brown dog jumped over the lazy dog.")
//	result := bs.ReplaceNth(g.Bytes("dog"), g.Bytes("fox"), 2)
//	fmt.Println(result)
//
// Output: "The quick brown dog jumped over the lazy fox.".
func (bs Bytes) ReplaceNth(oldB, newB Bytes, n Int) Bytes {
	if n < -1 || oldB.IsEmpty() {
		return bs
	}

	count, i := Int(0), Int(0)

	for {
		pos := bs[i:].Index(oldB)
		if pos == -1 {
			break
		}

		pos += i
		count++

		if count == n || (n == -1 && bs[pos+oldB.Len():].Index(oldB) == -1) {
			result := make(Bytes, 0, bs.Len()+newB.Len()-oldB.Len())
			result = append(result, bs[:pos]...)
			result = append(result, newB...)
			result = append(result, bs[pos+oldB.Len():]...)

			return result
		}

		i = pos + oldB.Len()
	}

	return bs
}

// Chunks splits the Bytes into chunks of the specified size.
//
// This function iterates through the Bytes, yielding chunks of the specified size.
// If size is less than or equal to 0 or the Bytes is empty,
// it returns an empty SeqSlice[Bytes].
// If size is greater than or equal to the length of the Bytes,
// it returns a SeqSlice[Bytes] containing the original Bytes.
//
// Unlike String.Chunks, which counts runes, the chunk size is measured in
// bytes, so multibyte UTF-8 sequences may be split across chunks.
//
// The yielded chunks are subslices sharing memory with the original Bytes,
// as with Split; clone them if independent copies are needed.
//
// Parameters:
//
// - size (Int): The size of the chunks to split the Bytes into.
//
// Returns:
//
// - SeqSlice[Bytes]: A sequence of Bytes chunks of the specified size.
//
// Example usage:
//
//	bs := g.Bytes("Hello, World!")
//	chunks := bs.Chunks(4).Collect()
//
// chunks contains {"Hell", "o, W", "orld", "!"}.
func (bs Bytes) Chunks(size Int) SeqSlice[Bytes] {
	if size.Lte(0) || bs.IsEmpty() {
		return func(func(Bytes) bool) {}
	}

	n := size.Std()
	l := len(bs)

	if n >= l {
		return func(yield func(Bytes) bool) { yield(bs) }
	}

	return func(yield func(Bytes) bool) {
		for i := 0; i < l; i += n {
			if !yield(bs[i:min(i+n, l)]) {
				return
			}
		}
	}
}

// Cut returns two Bytes values. The first Bytes contains the remainder of the
// original Bytes after the cut. The second Bytes contains the bytes between the
// first occurrences of the 'start' and 'end' markers, with tags removed if specified.
//
// The function searches for the 'start' and 'end' markers within the Bytes.
// If both are found, it returns the first Bytes containing the remainder of the
// original Bytes after the cut, followed by the second Bytes containing the bytes
// between the first occurrences of 'start' and 'end' with tags removed if specified.
//
// If either 'start' or 'end' is empty or not found in the Bytes, it returns the
// original Bytes as the first Bytes, and an empty Bytes as the second.
//
// Parameters:
//
// - start (Bytes): The Bytes marking the beginning of the segment to be cut.
//
// - end (Bytes): The Bytes marking the end of the segment to be cut.
//
//   - rmtags (bool, optional): An optional boolean parameter indicating whether
//     to remove the matched region (including the 'start' and 'end' tags) from the
//     remainder. Defaults to false, in which case the remainder equals the original
//     Bytes and only the cut content is extracted.
//
// Returns:
//
//   - Bytes: The first Bytes containing the remainder of the original Bytes.
//     When rmtags is true the matched region is removed from it; otherwise it is
//     the original Bytes. Returns the original Bytes if 'start' or 'end' is empty
//     or not found.
//
//   - Bytes: The second Bytes containing the bytes between the first occurrences of
//     'start' and 'end', or an empty Bytes if 'start' or 'end' is empty or not found.
//
// Example usage:
//
//	bs := g.Bytes("Hello, [world]! How are you?")
//	remainder, cut := bs.Cut(g.Bytes("["), g.Bytes("]"), true)
//	// remainder: "Hello, ! How are you?"
//	// cut: "world"
func (bs Bytes) Cut(start, end Bytes, rmtags ...bool) (Bytes, Bytes) {
	if start.IsEmpty() || end.IsEmpty() {
		return bs, Bytes("")
	}

	startIndex := bs.Index(start)
	if startIndex == -1 {
		return bs, Bytes("")
	}

	startEnd := startIndex + start.Len()
	endIndex := bs[startEnd:].Index(end)
	if endIndex == -1 {
		return bs, Bytes("")
	}

	cut := bs[startEnd : startEnd+endIndex]

	if len(rmtags) == 0 || !rmtags[0] {
		return bs, cut
	}

	tail := startEnd + endIndex + end.Len()

	remainder := make(Bytes, 0, startIndex+(bs.Len()-tail))
	remainder = append(remainder, bs[:startIndex]...)
	remainder = append(remainder, bs[tail:]...)

	return remainder, cut
}

// SubBytes extracts a subrange from the Bytes starting at the 'start' index and ending before the 'end' index.
// The function also supports an optional 'step' parameter to define the increment between indices in the result.
// If 'start' or 'end' index is negative, they represent positions relative to the end of the Bytes:
// - A negative 'start' index indicates the position from the end of the Bytes, moving backward.
// - A negative 'end' index indicates the position from the end of the Bytes.
// The function ensures that indices are adjusted to fall within the valid range of the Bytes' length.
// Out-of-bounds indices are clamped to the Bytes' bounds instead of panicking;
// if 'start' exceeds 'end' (for a positive step) the result is an empty Bytes.
//
// Unlike String.SubString, which indexes runes, all indices and the step are
// measured in bytes, so a boundary that falls inside a multibyte UTF-8 sequence
// splits the rune and the result may not be valid UTF-8. A negative step
// reverses bytes, not runes.
//
// The result is a newly allocated Bytes; the receiver is not modified.
func (bs Bytes) SubBytes(start, end Int, step ...Int) Bytes {
	n := bs.Len()

	clamp := func(i Int) Int {
		if i < 0 {
			i += n
		}

		if i < 0 {
			return 0
		}

		if i > n {
			return n
		}

		return i
	}

	return Bytes(Slice[byte](bs).SubSlice(clamp(start), clamp(end), step...))
}

// Similarity calculates the similarity between two Bytes using the
// Levenshtein distance algorithm and returns the similarity percentage as a Float.
//
// The function compares two Bytes using the Levenshtein distance,
// which measures the difference between two sequences by counting the number
// of single-byte edits required to change one sequence into the other.
// The similarity is then calculated by normalizing the distance by the maximum
// length of the two input Bytes.
//
// Unlike String.Similarity, which compares runes, this method operates byte-wise,
// so multibyte UTF-8 sequences are compared byte by byte.
//
// Parameters:
//
// - obs (Bytes): The Bytes to compare with bs.
//
// Returns:
//
// - Float: The similarity percentage between the two Bytes as a value between 0 and 100.
//
// Example usage:
//
//	b1 := g.Bytes("kitten")
//	b2 := g.Bytes("sitting")
//	similarity := b1.Similarity(b2) // 57.14285714285714
func (bs Bytes) Similarity(obs Bytes) Float {
	if bs.Eq(obs) {
		return 100
	}

	if bs.IsEmpty() || obs.IsEmpty() {
		return 0
	}

	s1, s2 := bs, obs

	n1, n2 := len(s1), len(s2)

	if n1 > n2 {
		s1, s2, n1, n2 = s2, s1, n2, n1
	}

	distance := make([]int, n1+1)

	for i, b2 := range s2 {
		prev := i + 1

		for j, b1 := range s1 {
			current := distance[j]
			if b2 != b1 {
				current = min(distance[j]+1, min(prev+1, distance[j+1]+1))
			}

			distance[j], prev = prev, current
		}

		distance[n1] = prev
	}

	return Float(1-float64(distance[n1])/float64(max(n1, n2))) * 100
}

// Truncate shortens the Bytes to the specified maximum length. If the Bytes exceeds the
// specified length, it is truncated, and an ellipsis ("...") is appended to indicate the truncation.
//
// If the length of the Bytes is less than or equal to the specified maximum length, the
// original Bytes is returned unchanged.
//
// Unlike String.Truncate, which is rune-aware, this method truncates based on the number
// of bytes, so multibyte UTF-8 sequences may be split.
//
// Parameters:
//   - max: The maximum number of bytes allowed in the resulting Bytes.
//
// Returns:
//   - A new Bytes truncated to the specified maximum length with "..." appended
//     if truncation occurs. Otherwise, returns the original Bytes.
//
// Example usage:
//
//	bs := g.Bytes("Hello, World!")
//	result := bs.Truncate(5)
//	// result: "Hello..."
//
//	bs2 := g.Bytes("Short")
//	result2 := bs2.Truncate(10)
//	// result2: "Short"
func (bs Bytes) Truncate(max Int) Bytes {
	if max.IsNegative() {
		return bs
	}

	if bs.Len() <= max {
		return bs
	}

	return append(bs[:max:max], "..."...)
}

// LeftJustify justifies the Bytes to the left by adding padding to the right, up to the
// specified length. If the length of the Bytes is already greater than or equal to the specified
// length, or the pad is empty, the original Bytes is returned.
//
// The padding Bytes is repeated as necessary to fill the remaining length.
// The padding is added to the right of the Bytes.
//
// Unlike String.LeftJustify, which counts runes, both length and padding are
// measured in bytes.
//
// Parameters:
//   - length: The desired length of the resulting justified Bytes.
//   - pad: The Bytes used as padding.
//
// Example usage:
//
//	bs := g.Bytes("Hello")
//	result := bs.LeftJustify(10, g.Bytes("..."))
//	// result: "Hello....."
func (bs Bytes) LeftJustify(length Int, pad Bytes) Bytes {
	if bs.Len() >= length || pad.IsEmpty() {
		return bs
	}

	buf := make(Bytes, 0, length)
	buf = append(buf, bs...)
	buf = appendPadding(buf, pad, length-bs.Len())

	return buf
}

// RightJustify justifies the Bytes to the right by adding padding to the left, up to the
// specified length. If the length of the Bytes is already greater than or equal to the specified
// length, or the pad is empty, the original Bytes is returned.
//
// The padding Bytes is repeated as necessary to fill the remaining length.
// The padding is added to the left of the Bytes.
//
// Unlike String.RightJustify, which counts runes, both length and padding are
// measured in bytes.
//
// Parameters:
//   - length: The desired length of the resulting justified Bytes.
//   - pad: The Bytes used as padding.
//
// Example usage:
//
//	bs := g.Bytes("Hello")
//	result := bs.RightJustify(10, g.Bytes("..."))
//	// result: ".....Hello"
func (bs Bytes) RightJustify(length Int, pad Bytes) Bytes {
	if bs.Len() >= length || pad.IsEmpty() {
		return bs
	}

	buf := make(Bytes, 0, length)
	buf = appendPadding(buf, pad, length-bs.Len())
	buf = append(buf, bs...)

	return buf
}

// Center justifies the Bytes by adding padding on both sides, up to the specified length.
// If the length of the Bytes is already greater than or equal to the specified length, or the
// pad is empty, the original Bytes is returned.
//
// The padding Bytes is repeated as necessary to evenly distribute the remaining length on both
// sides.
// The padding is added to the left and right of the Bytes.
//
// Unlike String.Center, which counts runes, both length and padding are
// measured in bytes.
//
// Parameters:
//   - length: The desired length of the resulting justified Bytes.
//   - pad: The Bytes used as padding.
//
// Example usage:
//
//	bs := g.Bytes("Hello")
//	result := bs.Center(10, g.Bytes("..."))
//	// result: "..Hello..."
func (bs Bytes) Center(length Int, pad Bytes) Bytes {
	slen := bs.Len()
	if slen >= length || pad.IsEmpty() {
		return bs
	}

	remains := length - slen

	buf := make(Bytes, 0, length)
	buf = appendPadding(buf, pad, remains/2)
	buf = append(buf, bs...)
	buf = appendPadding(buf, pad, (remains+1)/2)

	return buf
}

// appendPadding appends the padding Bytes to buf to fill the remaining length.
// It repeats the padding Bytes as necessary and appends any remaining bytes from
// the padding Bytes.
func appendPadding(buf, pad Bytes, remains Int) Bytes {
	padlen := pad.Len()

	for range remains / padlen {
		buf = append(buf, pad...)
	}

	if rem := remains % padlen; rem != 0 {
		buf = append(buf, pad[:rem]...)
	}

	return buf
}
