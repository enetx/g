package g

import (
	"unicode"
	"unicode/utf8"
)

func explodeBytes(bs Bytes) SeqSlice[Bytes] {
	return func(yield func(Bytes) bool) {
		for len(bs) > 0 {
			_, size := utf8.DecodeRune(bs)
			if !yield(bs[:size]) {
				return
			}

			bs = bs[size:]
		}
	}
}

func splitBytes(bs, sep Bytes, sepSave Int) SeqSlice[Bytes] {
	if len(sep) == 0 {
		return explodeBytes(bs)
	}

	return func(yield func(Bytes) bool) {
		for {
			i := bs.Index(sep)
			if i < 0 {
				break
			}

			frag := bs[:i+sepSave]
			if !yield(frag) {
				return
			}

			bs = bs[i+sep.Len():]
		}

		yield(bs)
	}
}

func fieldsBytes(bs Bytes) SeqSlice[Bytes] {
	return func(yield func(Bytes) bool) {
		start := -1

		for i := 0; i < len(bs); {
			size := 1
			r := rune(bs[i])
			isSpace := asciiSpace[bs[i]] != 0

			if r >= utf8.RuneSelf {
				r, size = utf8.DecodeRune(bs[i:])
				isSpace = unicode.IsSpace(r)
			}

			if isSpace {
				if start >= 0 {
					if !yield(bs[start:i]) {
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
			yield(bs[start:])
		}
	}
}

func fieldsbyBytes(bs Bytes, f func(rune) bool) SeqSlice[Bytes] {
	return func(yield func(Bytes) bool) {
		start := -1

		for i := 0; i < len(bs); {
			size := 1
			r := rune(bs[i])

			if r >= utf8.RuneSelf {
				r, size = utf8.DecodeRune(bs[i:])
			}

			if f(r) {
				if start >= 0 {
					if !yield(bs[start:i]) {
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
			yield(bs[start:])
		}
	}
}
