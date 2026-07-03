package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestBytesIsASCII(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		expected bool
	}{
		{name: "ASCII bytes", input: Bytes("Hello, World!"), expected: true},
		{name: "Empty bytes", input: Bytes(""), expected: true},
		{name: "Non-ASCII bytes", input: Bytes("héllo"), expected: false},
		{name: "High byte", input: Bytes{0x80}, expected: false},
		{name: "Boundary byte", input: Bytes{0x7f}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.IsASCII(); result != tt.expected {
				t.Errorf("Bytes(%q).IsASCII() = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBytesIsDigit(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		expected bool
	}{
		{name: "All digits", input: Bytes("1234567890"), expected: true},
		{name: "Empty bytes", input: Bytes(""), expected: false},
		{name: "Mixed digits and letters", input: Bytes("12a45"), expected: false},
		{name: "Spaces", input: Bytes("12 45"), expected: false},
		{name: "Single digit", input: Bytes("7"), expected: true},
		// Byte-wise semantics: Unicode digits are multibyte and are NOT recognized,
		// unlike the rune-aware String.IsDigit.
		{name: "Unicode digits", input: Bytes("١٢٣"), expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.IsDigit(); result != tt.expected {
				t.Errorf("Bytes(%q).IsDigit() = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBytesCut(t *testing.T) {
	tests := []struct {
		name          string
		input         Bytes
		start         Bytes
		end           Bytes
		rmtags        []bool
		wantRemainder Bytes
		wantCut       Bytes
	}{
		{
			name:          "Without rmtags",
			input:         Bytes("Hello, [world]! How are you?"),
			start:         Bytes("["),
			end:           Bytes("]"),
			wantRemainder: Bytes("Hello, [world]! How are you?"),
			wantCut:       Bytes("world"),
		},
		{
			name:          "With rmtags true",
			input:         Bytes("Hello, [world]! How are you?"),
			start:         Bytes("["),
			end:           Bytes("]"),
			rmtags:        []bool{true},
			wantRemainder: Bytes("Hello, ! How are you?"),
			wantCut:       Bytes("world"),
		},
		{
			name:          "With rmtags false",
			input:         Bytes("a<b>c"),
			start:         Bytes("<"),
			end:           Bytes(">"),
			rmtags:        []bool{false},
			wantRemainder: Bytes("a<b>c"),
			wantCut:       Bytes("b"),
		},
		{
			name:          "Start marker not found",
			input:         Bytes("Hello, world!"),
			start:         Bytes("["),
			end:           Bytes("!"),
			wantRemainder: Bytes("Hello, world!"),
			wantCut:       Bytes(""),
		},
		{
			name:          "End marker not found",
			input:         Bytes("Hello, [world!"),
			start:         Bytes("["),
			end:           Bytes("]"),
			wantRemainder: Bytes("Hello, [world!"),
			wantCut:       Bytes(""),
		},
		{
			name:          "Empty start marker",
			input:         Bytes("Hello, world!"),
			start:         Bytes(""),
			end:           Bytes("!"),
			wantRemainder: Bytes("Hello, world!"),
			wantCut:       Bytes(""),
		},
		{
			name:          "Empty end marker",
			input:         Bytes("Hello, world!"),
			start:         Bytes("!"),
			end:           Bytes(""),
			wantRemainder: Bytes("Hello, world!"),
			wantCut:       Bytes(""),
		},
		{
			name:          "Empty input",
			input:         Bytes(""),
			start:         Bytes("["),
			end:           Bytes("]"),
			wantRemainder: Bytes(""),
			wantCut:       Bytes(""),
		},
		{
			name:          "Multibyte markers with rmtags",
			input:         Bytes("пре«тело»пост"),
			start:         Bytes("«"),
			end:           Bytes("»"),
			rmtags:        []bool{true},
			wantRemainder: Bytes("препост"),
			wantCut:       Bytes("тело"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remainder, cut := tt.input.Cut(tt.start, tt.end, tt.rmtags...)
			if !remainder.Eq(tt.wantRemainder) {
				t.Errorf("Cut() remainder = %q; want %q", remainder, tt.wantRemainder)
			}

			if !cut.Eq(tt.wantCut) {
				t.Errorf("Cut() cut = %q; want %q", cut, tt.wantCut)
			}
		})
	}
}

func TestBytesCutDoesNotMutateReceiver(t *testing.T) {
	original := Bytes("Hello, [world]! How are you?")
	input := original.Clone()

	remainder, cut := input.Cut(Bytes("["), Bytes("]"), true)

	if !input.Eq(original) {
		t.Errorf("Cut() mutated the receiver: got %q; want %q", input, original)
	}

	if !remainder.Eq(Bytes("Hello, ! How are you?")) || !cut.Eq(Bytes("world")) {
		t.Errorf("Cut() = (%q, %q); want (%q, %q)", remainder, cut, "Hello, ! How are you?", "world")
	}
}

func TestBytesSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		other    Bytes
		expected Float
	}{
		{name: "Identical", input: Bytes("hello"), other: Bytes("hello"), expected: 100},
		{name: "Both empty", input: Bytes(""), other: Bytes(""), expected: 100},
		{name: "First empty", input: Bytes(""), other: Bytes("hello"), expected: 0},
		{name: "Second empty", input: Bytes("hello"), other: Bytes(""), expected: 0},
		{name: "Completely different", input: Bytes("abc"), other: Bytes("xyz"), expected: 0},
		{name: "Kitten sitting", input: Bytes("kitten"), other: Bytes("sitting"), expected: 57.14285714285714},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.Similarity(tt.other); result != tt.expected {
				t.Errorf("Bytes(%q).Similarity(%q) = %v; want %v", tt.input, tt.other, result, tt.expected)
			}
		})
	}
}

func TestBytesTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		max      Int
		expected Bytes
	}{
		{name: "Basic truncation", input: Bytes("Hello, World!"), max: 5, expected: Bytes("Hello...")},
		{name: "Shorter than max", input: Bytes("Short"), max: 10, expected: Bytes("Short")},
		{name: "Exact length", input: Bytes("Perfect"), max: 7, expected: Bytes("Perfect")},
		{name: "Zero max", input: Bytes("Hello"), max: 0, expected: Bytes("...")},
		{name: "Negative max", input: Bytes("Hello"), max: -1, expected: Bytes("Hello")},
		{name: "Empty bytes", input: Bytes(""), max: 5, expected: Bytes("")},
		// Byte-wise semantics: a 4-byte emoji is split mid-rune, unlike the
		// rune-aware String.Truncate.
		{name: "Multibyte split", input: Bytes("😊😊"), max: 5, expected: append(Bytes("😊"), 0xf0, '.', '.', '.')},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.Truncate(tt.max); !result.Eq(tt.expected) {
				t.Errorf("Bytes(%q).Truncate(%d) = %q; want %q", tt.input, tt.max, result, tt.expected)
			}
		})
	}
}

func TestBytesTruncateDoesNotMutateReceiver(t *testing.T) {
	original := Bytes("HelloWorld")
	input := original.Clone()

	result := input.Truncate(5)

	if !input.Eq(original) {
		t.Errorf("Truncate() mutated the receiver: got %q; want %q", input, original)
	}

	if !result.Eq(Bytes("Hello...")) {
		t.Errorf("Truncate() = %q; want %q", result, "Hello...")
	}
}

func TestBytesLeftJustify(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		length   Int
		pad      Bytes
		expected Bytes
	}{
		{name: "Basic padding", input: Bytes("Hello"), length: 10, pad: Bytes("."), expected: Bytes("Hello.....")},
		{name: "Repeated pad with remainder", input: Bytes("Hello"), length: 10, pad: Bytes("ab"), expected: Bytes("Helloababa")},
		{name: "Length equals len", input: Bytes("Hello"), length: 5, pad: Bytes("."), expected: Bytes("Hello")},
		{name: "Length less than len", input: Bytes("Hello"), length: 3, pad: Bytes("."), expected: Bytes("Hello")},
		{name: "Empty pad", input: Bytes("Hello"), length: 10, pad: Bytes(""), expected: Bytes("Hello")},
		{name: "Empty bytes", input: Bytes(""), length: 3, pad: Bytes("."), expected: Bytes("...")},
		// Byte-wise semantics: "€" is 3 bytes, so 5-2=3 remaining bytes fit exactly one "€".
		{name: "Multibyte pad exact", input: Bytes("ab"), length: 5, pad: Bytes("€"), expected: Bytes("ab€")},
		// 4 remaining bytes: one full "€" plus its first byte 0xe2, splitting the rune.
		{name: "Multibyte pad split", input: Bytes("ab"), length: 6, pad: Bytes("€"), expected: append(Bytes("ab€"), 0xe2)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.LeftJustify(tt.length, tt.pad); !result.Eq(tt.expected) {
				t.Errorf("Bytes(%q).LeftJustify(%d, %q) = %q; want %q", tt.input, tt.length, tt.pad, result, tt.expected)
			}
		})
	}
}

func TestBytesRightJustify(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		length   Int
		pad      Bytes
		expected Bytes
	}{
		{name: "Basic padding", input: Bytes("Hello"), length: 10, pad: Bytes("."), expected: Bytes(".....Hello")},
		{name: "Repeated pad with remainder", input: Bytes("Hello"), length: 10, pad: Bytes("ab"), expected: Bytes("ababaHello")},
		{name: "Length equals len", input: Bytes("Hello"), length: 5, pad: Bytes("."), expected: Bytes("Hello")},
		{name: "Length less than len", input: Bytes("Hello"), length: 3, pad: Bytes("."), expected: Bytes("Hello")},
		{name: "Empty pad", input: Bytes("Hello"), length: 10, pad: Bytes(""), expected: Bytes("Hello")},
		{name: "Empty bytes", input: Bytes(""), length: 3, pad: Bytes("."), expected: Bytes("...")},
		// Byte-wise semantics: 4 remaining bytes = one full "€" plus its first byte 0xe2.
		{name: "Multibyte pad split", input: Bytes("ab"), length: 6, pad: Bytes("€"), expected: append(Bytes("€"), 0xe2, 'a', 'b')},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.RightJustify(tt.length, tt.pad); !result.Eq(tt.expected) {
				t.Errorf("Bytes(%q).RightJustify(%d, %q) = %q; want %q", tt.input, tt.length, tt.pad, result, tt.expected)
			}
		})
	}
}

func TestBytesCenter(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		length   Int
		pad      Bytes
		expected Bytes
	}{
		{name: "Basic padding", input: Bytes("Hello"), length: 10, pad: Bytes("."), expected: Bytes("..Hello...")},
		{name: "Even remainder", input: Bytes("ab"), length: 6, pad: Bytes("."), expected: Bytes("..ab..")},
		{name: "Repeated pad with remainder", input: Bytes("abc"), length: 10, pad: Bytes(".."), expected: Bytes("...abc....")},
		{name: "Length equals len", input: Bytes("Hello"), length: 5, pad: Bytes("."), expected: Bytes("Hello")},
		{name: "Length less than len", input: Bytes("Hello"), length: 3, pad: Bytes("."), expected: Bytes("Hello")},
		{name: "Empty pad", input: Bytes("Hello"), length: 10, pad: Bytes(""), expected: Bytes("Hello")},
		{name: "Empty bytes", input: Bytes(""), length: 4, pad: Bytes("."), expected: Bytes("....")},
		// Byte-wise semantics: remains=6, 3 bytes per side fit exactly one "€" each.
		{name: "Multibyte pad exact", input: Bytes("ab"), length: 8, pad: Bytes("€"), expected: Bytes("€ab€")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.Center(tt.length, tt.pad); !result.Eq(tt.expected) {
				t.Errorf("Bytes(%q).Center(%d, %q) = %q; want %q", tt.input, tt.length, tt.pad, result, tt.expected)
			}
		})
	}
}

func TestBytesReplaceMulti(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		oldnew   []Bytes
		expected Bytes
	}{
		{
			name:  "Basic replacements",
			input: Bytes("Hello, world! This is a test."),
			oldnew: []Bytes{
				Bytes("Hello"), Bytes("Greetings"),
				Bytes("world"), Bytes("universe"),
				Bytes("test"), Bytes("example"),
			},
			expected: Bytes("Greetings, universe! This is a example."),
		},
		{
			name:     "No pairs",
			input:    Bytes("Hello, world!"),
			oldnew:   nil,
			expected: Bytes("Hello, world!"),
		},
		{
			name:     "Empty input",
			input:    Bytes(""),
			oldnew:   []Bytes{Bytes("a"), Bytes("b")},
			expected: Bytes(""),
		},
		{
			name:     "Replace with empty",
			input:    Bytes("banana"),
			oldnew:   []Bytes{Bytes("na"), Bytes("")},
			expected: Bytes("ba"),
		},
		{
			name:     "Earlier pair wins at same position",
			input:    Bytes("abc"),
			oldnew:   []Bytes{Bytes("ab"), Bytes("X"), Bytes("a"), Bytes("Y")},
			expected: Bytes("Xc"),
		},
		{
			name:     "Multibyte patterns",
			input:    Bytes("héllo wörld"),
			oldnew:   []Bytes{Bytes("é"), Bytes("e"), Bytes("ö"), Bytes("o")},
			expected: Bytes("hello world"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.ReplaceMulti(tt.oldnew...); !result.Eq(tt.expected) {
				t.Errorf("Bytes(%q).ReplaceMulti(%q) = %q; want %q", tt.input, tt.oldnew, result, tt.expected)
			}
		})
	}
}

func TestBytesReplaceMultiDoesNotMutateReceiver(t *testing.T) {
	original := Bytes("Hello, world!")
	input := original.Clone()

	result := input.ReplaceMulti(Bytes("world"), Bytes("universe"))

	if !input.Eq(original) {
		t.Errorf("ReplaceMulti() mutated the receiver: got %q; want %q", input, original)
	}

	if !result.Eq(Bytes("Hello, universe!")) {
		t.Errorf("ReplaceMulti() = %q; want %q", result, "Hello, universe!")
	}

	result[0] = 'X'

	if !input.Eq(original) {
		t.Errorf("ReplaceMulti() result shares memory with the receiver: got %q; want %q", input, original)
	}
}

func TestBytesRemove(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		patterns []Bytes
		expected Bytes
	}{
		{
			name:     "Basic removal",
			input:    Bytes("Hello, world! This is a test."),
			patterns: []Bytes{Bytes("Hello"), Bytes("test")},
			expected: Bytes(", world! This is a ."),
		},
		{
			name:     "No patterns",
			input:    Bytes("Hello, world!"),
			patterns: nil,
			expected: Bytes("Hello, world!"),
		},
		{
			name:     "Empty input",
			input:    Bytes(""),
			patterns: []Bytes{Bytes("a")},
			expected: Bytes(""),
		},
		{
			name:     "Pattern not found",
			input:    Bytes("Hello"),
			patterns: []Bytes{Bytes("xyz")},
			expected: Bytes("Hello"),
		},
		{
			name:     "All occurrences removed",
			input:    Bytes("banana"),
			patterns: []Bytes{Bytes("na")},
			expected: Bytes("ba"),
		},
		{
			name:     "Multibyte pattern",
			input:    Bytes("héllo héllo"),
			patterns: []Bytes{Bytes("é")},
			expected: Bytes("hllo hllo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.Remove(tt.patterns...); !result.Eq(tt.expected) {
				t.Errorf("Bytes(%q).Remove(%q) = %q; want %q", tt.input, tt.patterns, result, tt.expected)
			}
		})
	}
}

func TestBytesRemoveDoesNotMutateReceiver(t *testing.T) {
	original := Bytes("Hello, world!")
	input := original.Clone()

	result := input.Remove(Bytes("world"))

	if !input.Eq(original) {
		t.Errorf("Remove() mutated the receiver: got %q; want %q", input, original)
	}

	if !result.Eq(Bytes("Hello, !")) {
		t.Errorf("Remove() = %q; want %q", result, "Hello, !")
	}
}

func TestBytesReplaceNth(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		oldB     Bytes
		newB     Bytes
		n        Int
		expected Bytes
	}{
		{
			name:     "First occurrence",
			input:    Bytes("The quick brown dog jumped over the lazy dog."),
			oldB:     Bytes("dog"),
			newB:     Bytes("fox"),
			n:        1,
			expected: Bytes("The quick brown fox jumped over the lazy dog."),
		},
		{
			name:     "Second occurrence",
			input:    Bytes("The quick brown dog jumped over the lazy dog."),
			oldB:     Bytes("dog"),
			newB:     Bytes("fox"),
			n:        2,
			expected: Bytes("The quick brown dog jumped over the lazy fox."),
		},
		{
			name:     "Last occurrence",
			input:    Bytes("The quick brown dog jumped over the lazy dog."),
			oldB:     Bytes("dog"),
			newB:     Bytes("fox"),
			n:        -1,
			expected: Bytes("The quick brown dog jumped over the lazy fox."),
		},
		{
			name:     "Negative n (except -1)",
			input:    Bytes("banana"),
			oldB:     Bytes("na"),
			newB:     Bytes("xy"),
			n:        -2,
			expected: Bytes("banana"),
		},
		{
			name:     "Zero n",
			input:    Bytes("banana"),
			oldB:     Bytes("na"),
			newB:     Bytes("xy"),
			n:        0,
			expected: Bytes("banana"),
		},
		{
			name:     "Longer replacement",
			input:    Bytes("Hello, world!"),
			oldB:     Bytes("world"),
			newB:     Bytes("beautiful world"),
			n:        1,
			expected: Bytes("Hello, beautiful world!"),
		},
		{
			name:     "Shorter replacement",
			input:    Bytes("A wonderful day"),
			oldB:     Bytes("wonderful"),
			newB:     Bytes("nice"),
			n:        1,
			expected: Bytes("A nice day"),
		},
		{
			name:     "Not enough occurrences",
			input:    Bytes("foo bar baz"),
			oldB:     Bytes("bar"),
			newB:     Bytes("XXX"),
			n:        2,
			expected: Bytes("foo bar baz"),
		},
		{
			name:     "Empty old pattern",
			input:    Bytes("Hello"),
			oldB:     Bytes(""),
			newB:     Bytes("x"),
			n:        1,
			expected: Bytes("Hello"),
		},
		{
			name:     "Empty input",
			input:    Bytes(""),
			oldB:     Bytes("a"),
			newB:     Bytes("b"),
			n:        1,
			expected: Bytes(""),
		},
		{
			name:     "Overlapping occurrences count non-overlapped",
			input:    Bytes("aaaa"),
			oldB:     Bytes("aa"),
			newB:     Bytes("b"),
			n:        2,
			expected: Bytes("aab"),
		},
		{
			name:     "Multibyte old pattern",
			input:    Bytes("猫と犬と猫"),
			oldB:     Bytes("猫"),
			newB:     Bytes("にゃんこ"),
			n:        -1,
			expected: Bytes("猫と犬とにゃんこ"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.ReplaceNth(tt.oldB, tt.newB, tt.n); !result.Eq(tt.expected) {
				t.Errorf(
					"Bytes(%q).ReplaceNth(%q, %q, %d) = %q; want %q",
					tt.input, tt.oldB, tt.newB, tt.n, result, tt.expected,
				)
			}
		})
	}
}

func TestBytesReplaceNthDoesNotMutateReceiver(t *testing.T) {
	original := Bytes("The quick brown dog jumped over the lazy dog.")
	input := original.Clone()

	result := input.ReplaceNth(Bytes("dog"), Bytes("fox"), 1)

	if !input.Eq(original) {
		t.Errorf("ReplaceNth() mutated the receiver: got %q; want %q", input, original)
	}

	if !result.Eq(Bytes("The quick brown fox jumped over the lazy dog.")) {
		t.Errorf("ReplaceNth() = %q; want %q", result, "The quick brown fox jumped over the lazy dog.")
	}

	result[0] = 'X'

	if !input.Eq(original) {
		t.Errorf("ReplaceNth() result shares memory with the receiver: got %q; want %q", input, original)
	}
}

func TestBytesChunks(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		size     Int
		expected []Bytes
	}{
		{name: "Even split", input: Bytes("hello world!"), size: 3, expected: []Bytes{Bytes("hel"), Bytes("lo "), Bytes("wor"), Bytes("ld!")}},
		{name: "Uneven split", input: Bytes("hello"), size: 2, expected: []Bytes{Bytes("he"), Bytes("ll"), Bytes("o")}},
		{name: "Size equals length", input: Bytes("hello"), size: 5, expected: []Bytes{Bytes("hello")}},
		{name: "Size exceeds length", input: Bytes("hi"), size: 10, expected: []Bytes{Bytes("hi")}},
		{name: "Zero size", input: Bytes("hello"), size: 0, expected: nil},
		{name: "Negative size", input: Bytes("hello"), size: -1, expected: nil},
		{name: "Empty input", input: Bytes(""), size: 3, expected: nil},
		// Byte-wise semantics: "é" (0xc3 0xa9) is split across chunks, unlike the
		// rune-aware String.Chunks.
		{name: "Multibyte split mid-rune", input: Bytes("héllo"), size: 2, expected: []Bytes{Bytes("h\xc3"), Bytes("\xa9l"), Bytes("lo")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunks := tt.input.Chunks(tt.size).Collect()

			if chunks.Len() != Int(len(tt.expected)) {
				t.Fatalf("Bytes(%q).Chunks(%d) yielded %d chunks; want %d", tt.input, tt.size, chunks.Len(), len(tt.expected))
			}

			for i, chunk := range chunks {
				if !chunk.Eq(tt.expected[i]) {
					t.Errorf("Bytes(%q).Chunks(%d)[%d] = %q; want %q", tt.input, tt.size, i, chunk, tt.expected[i])
				}
			}
		})
	}
}

func TestBytesChunksEarlyBreak(t *testing.T) {
	var collected []Bytes

	for chunk := range Bytes("abcdef").Chunks(2) {
		collected = append(collected, chunk)
		if len(collected) == 2 {
			break
		}
	}

	if len(collected) != 2 || !collected[0].Eq(Bytes("ab")) || !collected[1].Eq(Bytes("cd")) {
		t.Errorf("Chunks() early break collected %q; want [%q %q]", collected, "ab", "cd")
	}
}

func TestBytesSubBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		start    Int
		end      Int
		step     []Int
		expected Bytes
	}{
		{name: "Basic subrange", input: Bytes("Hello, world!"), start: 0, end: 5, expected: Bytes("Hello")},
		{name: "With negative start", input: Bytes("Hello, world!"), start: -6, end: 12, expected: Bytes("world")},
		{name: "With negative end", input: Bytes("Hello, world!"), start: 7, end: -1, expected: Bytes("world")},
		{name: "Start exceeds end", input: Bytes("Hello, world!"), start: 5, end: 1, expected: Bytes("")},
		{name: "Step parameter", input: Bytes("abcdef"), start: 0, end: 6, step: []Int{2}, expected: Bytes("ace")},
		{name: "Negative step", input: Bytes("abcdefgh"), start: -1, end: 0, step: []Int{-1}, expected: Bytes("hgfedcb")},
		{name: "End past length", input: Bytes("Hello"), start: 0, end: 100, expected: Bytes("Hello")},
		{name: "Start past length", input: Bytes("Hello"), start: 100, end: 200, expected: Bytes("")},
		{name: "Start way negative", input: Bytes("Hello"), start: -100, end: 3, expected: Bytes("Hel")},
		{name: "Both negative oob", input: Bytes("Hello"), start: -100, end: -100, expected: Bytes("")},
		{name: "Empty input", input: Bytes(""), start: 0, end: 5, expected: Bytes("")},
		// Byte-wise semantics: indices count bytes, so the boundary falls inside
		// "é" (0xc3 0xa9) and splits the rune, unlike the rune-aware
		// String.SubString, which would return "hé" here.
		{name: "Multibyte split mid-rune", input: Bytes("héllo"), start: 0, end: 2, expected: Bytes("h\xc3")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.input.SubBytes(tt.start, tt.end, tt.step...); !result.Eq(tt.expected) {
				t.Errorf(
					"Bytes(%q).SubBytes(%d, %d, %v) = %q; want %q",
					tt.input, tt.start, tt.end, tt.step, result, tt.expected,
				)
			}
		})
	}
}

func TestBytesSubBytesDoesNotMutateReceiver(t *testing.T) {
	original := Bytes("Hello, world!")
	input := original.Clone()

	result := input.SubBytes(0, 5)

	if !input.Eq(original) {
		t.Errorf("SubBytes() mutated the receiver: got %q; want %q", input, original)
	}

	if !result.Eq(Bytes("Hello")) {
		t.Errorf("SubBytes() = %q; want %q", result, "Hello")
	}

	result[0] = 'X'

	if !input.Eq(original) {
		t.Errorf("SubBytes() result shares memory with the receiver: got %q; want %q", input, original)
	}
}

func TestBytesTextStringParity(t *testing.T) {
	inputs := []string{"", "Hello, World!", "12345", "12a45", "kitten", "abc", "Short"}

	t.Run("IsASCII", func(t *testing.T) {
		for _, s := range inputs {
			if got, want := Bytes(s).IsASCII(), String(s).IsASCII(); got != want {
				t.Errorf("Bytes(%q).IsASCII() = %v; String version = %v", s, got, want)
			}
		}
	})

	t.Run("IsDigit", func(t *testing.T) {
		for _, s := range inputs {
			if got, want := Bytes(s).IsDigit(), String(s).IsDigit(); got != want {
				t.Errorf("Bytes(%q).IsDigit() = %v; String version = %v", s, got, want)
			}
		}
	})

	t.Run("Cut", func(t *testing.T) {
		for _, rmtags := range [][]bool{nil, {false}, {true}} {
			for _, s := range []string{"Hello, [world]!", "no markers", "[unclosed", ""} {
				bRem, bCut := Bytes(s).Cut(Bytes("["), Bytes("]"), rmtags...)
				sRem, sCut := String(s).Cut("[", "]", rmtags...)

				if bRem.String() != sRem || bCut.String() != sCut {
					t.Errorf(
						"Bytes(%q).Cut(rmtags=%v) = (%q, %q); String version = (%q, %q)",
						s, rmtags, bRem, bCut, sRem, sCut,
					)
				}
			}
		}
	})

	t.Run("Similarity", func(t *testing.T) {
		pairs := [][2]string{{"kitten", "sitting"}, {"abc", "xyz"}, {"", "abc"}, {"same", "same"}, {"", ""}}
		for _, p := range pairs {
			if got, want := Bytes(p[0]).Similarity(Bytes(p[1])), String(p[0]).Similarity(String(p[1])); got != want {
				t.Errorf("Bytes(%q).Similarity(%q) = %v; String version = %v", p[0], p[1], got, want)
			}
		}
	})

	t.Run("Truncate", func(t *testing.T) {
		for _, s := range inputs {
			for _, max := range []Int{-1, 0, 3, 5, 100} {
				if got, want := Bytes(s).Truncate(max), String(s).Truncate(max); got.String() != want {
					t.Errorf("Bytes(%q).Truncate(%d) = %q; String version = %q", s, max, got, want)
				}
			}
		}
	})

	t.Run("ReplaceMulti", func(t *testing.T) {
		for _, s := range inputs {
			got := Bytes(s).ReplaceMulti(Bytes("l"), Bytes("L"), Bytes("o"), Bytes("0"))
			want := String(s).ReplaceMulti("l", "L", "o", "0")

			if got.String() != want {
				t.Errorf("Bytes(%q).ReplaceMulti() = %q; String version = %q", s, got, want)
			}
		}
	})

	t.Run("Remove", func(t *testing.T) {
		for _, s := range inputs {
			got := Bytes(s).Remove(Bytes("l"), Bytes("o"))
			want := String(s).Remove("l", "o")

			if got.String() != want {
				t.Errorf("Bytes(%q).Remove() = %q; String version = %q", s, got, want)
			}
		}
	})

	t.Run("ReplaceNth", func(t *testing.T) {
		for _, s := range inputs {
			for _, n := range []Int{-2, -1, 0, 1, 2, 100} {
				got := Bytes(s).ReplaceNth(Bytes("l"), Bytes("LL"), n)
				want := String(s).ReplaceNth("l", "LL", n)

				if got.String() != want {
					t.Errorf("Bytes(%q).ReplaceNth(n=%d) = %q; String version = %q", s, n, got, want)
				}
			}
		}
	})

	t.Run("Chunks", func(t *testing.T) {
		for _, s := range inputs {
			for _, size := range []Int{-1, 0, 1, 2, 3, 100} {
				bChunks := Bytes(s).Chunks(size).Collect()
				sChunks := String(s).Chunks(size).Collect()

				if bChunks.Len() != sChunks.Len() {
					t.Errorf(
						"Bytes(%q).Chunks(%d) yielded %d chunks; String version yielded %d",
						s, size, bChunks.Len(), sChunks.Len(),
					)

					continue
				}

				for i, chunk := range bChunks {
					if chunk.String() != sChunks[i] {
						t.Errorf("Bytes(%q).Chunks(%d)[%d] = %q; String version = %q", s, size, i, chunk, sChunks[i])
					}
				}
			}
		}
	})

	t.Run("SubBytes", func(t *testing.T) {
		bounds := [][2]Int{{0, 5}, {-6, -1}, {7, -1}, {5, 1}, {0, 100}, {-100, 3}, {100, 200}, {-100, -100}, {-1, 0}}

		for _, s := range inputs {
			for _, se := range bounds {
				for _, step := range [][]Int{nil, {2}, {-1}} {
					got := Bytes(s).SubBytes(se[0], se[1], step...)
					want := String(s).SubString(se[0], se[1], step...)

					if got.String() != want {
						t.Errorf(
							"Bytes(%q).SubBytes(%d, %d, %v) = %q; String version = %q",
							s, se[0], se[1], step, got, want,
						)
					}
				}
			}
		}
	})

	t.Run("Justify", func(t *testing.T) {
		for _, s := range inputs {
			for _, pad := range []string{"", ".", "ab"} {
				for _, length := range []Int{0, 3, 10, 20} {
					if got, want := Bytes(s).LeftJustify(length, Bytes(pad)), String(s).LeftJustify(length, String(pad)); got.String() != want {
						t.Errorf("Bytes(%q).LeftJustify(%d, %q) = %q; String version = %q", s, length, pad, got, want)
					}

					if got, want := Bytes(s).RightJustify(length, Bytes(pad)), String(s).RightJustify(length, String(pad)); got.String() != want {
						t.Errorf("Bytes(%q).RightJustify(%d, %q) = %q; String version = %q", s, length, pad, got, want)
					}

					if got, want := Bytes(s).Center(length, Bytes(pad)), String(s).Center(length, String(pad)); got.String() != want {
						t.Errorf("Bytes(%q).Center(%d, %q) = %q; String version = %q", s, length, pad, got, want)
					}
				}
			}
		}
	})
}
