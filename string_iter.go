package g

import (
	"unicode"
	"unicode/utf8"
)

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func explode(s String) SeqSlice[String] {
	return func(yield func(String) bool) {
		for len(s) > 0 {
			_, size := utf8.DecodeRuneInString(s.Std())
			if !yield(s[:size]) {
				return
			}

			s = s[size:]
		}
	}
}

func split(s, sep String, sepSave Int) SeqSlice[String] {
	if len(sep) == 0 {
		return explode(s)
	}

	return func(yield func(String) bool) {
		for {
			i := s.Index(sep)

			if i < 0 {
				break
			}

			frag := s[:i+sepSave]
			if !yield(frag) {
				return
			}

			s = s[i+sep.Len():]
		}

		yield(s)
	}
}

func fields(s String) SeqSlice[String] {
	return func(yield func(String) bool) {
		start := -1

		for i := 0; i < len(s); {
			size := 1
			r := rune(s[i])
			isSpace := asciiSpace[s[i]] != 0

			if r >= utf8.RuneSelf {
				r, size = utf8.DecodeRuneInString(s[i:].Std())
				isSpace = unicode.IsSpace(r)
			}

			if isSpace {
				if start >= 0 {
					if !yield(s[start:i]) {
						return
					}
					start = -1
				}
			} else if start < 0 {
				start = i
			}

			i += size
		}

		if start >= 0 {
			yield(s[start:])
		}
	}
}

func fieldsby(s String, fn func(rune) bool) SeqSlice[String] {
	return func(yield func(String) bool) {
		start := -1

		for i := 0; i < len(s); {
			size := 1
			r := rune(s[i])

			if r >= utf8.RuneSelf {
				r, size = utf8.DecodeRuneInString(s[i:].Std())
			}

			if fn(r) {
				if start >= 0 {
					if !yield(s[start:i]) {
						return
					}
					start = -1
				}
			} else if start < 0 {
				start = i
			}

			i += size
		}

		if start >= 0 {
			yield(s[start:])
		}
	}
}

func lines(s String) SeqSlice[String] {
	return func(yield func(String) bool) {
		for s != "" {
			var line String

			if i := s.Index("\n"); i >= 0 {
				line, s = s[:i+1].TrimRight("\r\n"), s[i+1:]
			} else {
				line, s = s, ""
			}

			if !yield(line) {
				return
			}
		}
	}
}
