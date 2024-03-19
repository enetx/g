package g_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/enetx/g"
)

func TestStringMinMax(t *testing.T) {
	// Test cases for Min method
	minTestCases := []struct {
		inputs   []g.String
		expected g.String
	}{
		{[]g.String{"apple", "banana", "orange"}, "apple"},
		{[]g.String{"cat", "dog", "elephant"}, "cat"},
		{[]g.String{"123", "456", "789"}, "123"},
	}

	for _, testCase := range minTestCases {
		result := g.NewString(testCase.inputs[0]).Min(testCase.inputs[1:]...)
		if result != testCase.expected {
			t.Errorf("Min test failed. Expected: %s, Got: %s", testCase.expected, result)
		}
	}

	// Test cases for Max method
	maxTestCases := []struct {
		inputs   []g.String
		expected g.String
	}{
		{[]g.String{"apple", "banana", "orange"}, "orange"},
		{[]g.String{"cat", "dog", "elephant"}, "elephant"},
		{[]g.String{"123", "456", "789"}, "789"},
	}

	for _, testCase := range maxTestCases {
		result := g.NewString(testCase.inputs[0]).Max(testCase.inputs[1:]...)
		if result != testCase.expected {
			t.Errorf("Max test failed. Expected: %s, Got: %s", testCase.expected, result)
		}
	}
}

func TestChars(t *testing.T) {
	// Testing on a regular string with ASCII characters
	asciiStr := g.String("Hello")
	asciiExpected := g.Slice[g.String]{"H", "e", "l", "l", "o"}
	asciiResult := asciiStr.Chars()
	if !reflect.DeepEqual(asciiExpected, asciiResult) {
		t.Errorf("Expected %v, but got %v", asciiExpected, asciiResult)
	}

	// Testing on a string with Unicode characters (Russian)
	unicodeStr := g.String("Привет")
	unicodeExpected := g.Slice[g.String]{"П", "р", "и", "в", "е", "т"}
	unicodeResult := unicodeStr.Chars()
	if !reflect.DeepEqual(unicodeExpected, unicodeResult) {
		t.Errorf("Expected %v, but got %v", unicodeExpected, unicodeResult)
	}

	// Testing on a string with Unicode characters (Chinese)
	chineseStr := g.String("你好")
	chineseExpected := g.Slice[g.String]{"你", "好"}
	chineseResult := chineseStr.Chars()
	if !reflect.DeepEqual(chineseExpected, chineseResult) {
		t.Errorf("Expected %v, but got %v", chineseExpected, chineseResult)
	}

	// Additional test with a mix of ASCII and Unicode characters
	mixedStr := g.String("Hello 你好")
	mixedExpected := g.Slice[g.String]{"H", "e", "l", "l", "o", " ", "你", "好"}
	mixedResult := mixedStr.Chars()
	if !reflect.DeepEqual(mixedExpected, mixedResult) {
		t.Errorf("Expected %v, but got %v", mixedExpected, mixedResult)
	}

	// Testing on a string with special characters and symbols
	specialStr := g.String("Hello, 你好! How are you today? こんにちは")
	specialExpected := g.Slice[g.String]{
		"H",
		"e",
		"l",
		"l",
		"o",
		",",
		" ",
		"你",
		"好",
		"!",
		" ",
		"H",
		"o",
		"w",
		" ",
		"a",
		"r",
		"e",
		" ",
		"y",
		"o",
		"u",
		" ",
		"t",
		"o",
		"d",
		"a",
		"y",
		"?",
		" ",
		"こ",
		"ん",
		"に",
		"ち",
		"は",
	}
	specialResult := specialStr.Chars()
	if !reflect.DeepEqual(specialExpected, specialResult) {
		t.Errorf("Expected %v, but got %v", specialExpected, specialResult)
	}

	// Testing on a string with emojis
	emojiStr := g.String("Hello, 😊🌍🚀")
	emojiExpected := g.Slice[g.String]{"H", "e", "l", "l", "o", ",", " ", "😊", "🌍", "🚀"}
	emojiResult := emojiStr.Chars()
	if !reflect.DeepEqual(emojiExpected, emojiResult) {
		t.Errorf("Expected %v, but got %v", emojiExpected, emojiResult)
	}

	// Testing on an empty string
	emptyStr := g.String("")
	emptyExpected := g.Slice[g.String]{}
	emptyResult := emptyStr.Chars()
	if !reflect.DeepEqual(emptyExpected, emptyResult) {
		t.Errorf("Expected %v, but got %v", emptyExpected, emptyResult)
	}
}

func TestStringIsDigit(t *testing.T) {
	tests := []struct {
		name string
		str  g.String
		want bool
	}{
		{"empty", g.String(""), false},
		{"one", g.String("1"), true},
		{"nine", g.String("99999"), true},
		{"non-digit", g.String("1111a"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.str.IsDigit(); got != tt.want {
				t.Errorf("String.IsDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToInt(t *testing.T) {
	tests := []struct {
		name string
		str  g.String
		want g.Int
	}{
		{
			name: "empty",
			str:  g.String(""),
			want: 0,
		},
		{
			name: "one digit",
			str:  g.String("1"),
			want: 1,
		},
		{
			name: "two digits",
			str:  g.String("12"),
			want: 12,
		},
		{
			name: "one letter",
			str:  g.String("a"),
			want: 0,
		},
		{
			name: "one digit and one letter",
			str:  g.String("1a"),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.str.ToInt().UnwrapOrDefault(); got != tt.want {
				t.Errorf("String.ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToTitle(t *testing.T) {
	tests := []struct {
		name string
		str  g.String
		want g.String
	}{
		{"empty", "", ""},
		{"one word", "hello", "Hello"},
		{"two words", "hello world", "Hello World"},
		{"three words", "hello world, how are you?", "Hello World, How Are You?"},
		{"multiple hyphens", "foo-bar-baz", "Foo-Bar-Baz"},
		{"non-ascii letters", "こんにちは, 世界!", "こんにちは, 世界!"},
		{"all whitespace", "   \t\n   ", "   \t\n   "},
		{"numbers", "12345 67890", "12345 67890"},
		{"arabic", "مرحبا بالعالم", "مرحبا بالعالم"},
		{"chinese", "你好世界", "你好世界"},
		{"czech", "ahoj světe", "Ahoj Světe"},
		{"danish", "hej verden", "Hej Verden"},
		{"dutch", "hallo wereld", "Hallo Wereld"},
		{"french", "bonjour tout le monde", "Bonjour Tout Le Monde"},
		{"german", "hallo welt", "Hallo Welt"},
		{"hebrew", "שלום עולם", "שלום עולם"},
		{"hindi", "नमस्ते दुनिया", "नमस्ते दुनिया"},
		{"hungarian", "szia világ", "Szia Világ"},
		{"italian", "ciao mondo", "Ciao Mondo"},
		{"japanese", "こんにちは世界", "こんにちは世界"},
		{"korean", "안녕하세요 세상", "안녕하세요 세상"},
		{"norwegian", "hei verden", "Hei Verden"},
		{"polish", "witaj świecie", "Witaj Świecie"},
		{"portuguese", "olá mundo", "Olá Mundo"},
		{"russian", "привет мир", "Привет Мир"},
		{"spanish", "hola mundo", "Hola Mundo"},
		{"swedish", "hej världen", "Hej Världen"},
		{"turkish", "merhaba dünya", "Merhaba Dünya"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.str.Title(); got.Ne(tt.want) {
				t.Errorf("String.ToTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringAdd(t *testing.T) {
	tests := []struct {
		name string
		str  g.String
		s    g.String
		want g.String
	}{
		{
			name: "empty",
			str:  g.String(""),
			s:    g.String(""),
			want: g.String(""),
		},
		{
			name: "empty_hs",
			str:  g.String(""),
			s:    g.String("test"),
			want: g.String("test"),
		},
		{
			name: "empty_s",
			str:  g.String("test"),
			s:    g.String(""),
			want: g.String("test"),
		},
		{
			name: "not_empty",
			str:  g.String("test"),
			s:    g.String("test"),
			want: g.String("testtest"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.str.Add(tt.s); got != tt.want {
				t.Errorf("String.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringAddPrefix(t *testing.T) {
	tests := []struct {
		name string
		str  g.String
		s    g.String
		want g.String
	}{
		{
			name: "empty",
			str:  g.String(""),
			s:    g.String(""),
			want: g.String(""),
		},
		{
			name: "empty_hs",
			str:  g.String(""),
			s:    g.String("test"),
			want: g.String("test"),
		},
		{
			name: "empty_s",
			str:  g.String("test"),
			s:    g.String(""),
			want: g.String("test"),
		},
		{
			name: "not_empty",
			str:  g.String("rest"),
			s:    g.String("test"),
			want: g.String("testrest"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.str.AddPrefix(tt.s); got != tt.want {
				t.Errorf("String.AddPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringRandom(t *testing.T) {
	for i := range 100 {
		random := g.NewString("").Random(i)

		if random.Len() != i {
			t.Errorf("Random string length %d is not equal to %d", random.Len(), i)
		}
	}
}

func TestStringChunks(t *testing.T) {
	str := g.String("")
	chunks := str.Chunks(3)

	if chunks.Len() != 0 {
		t.Errorf("Expected empty slice, but got %v", chunks)
	}

	str = g.String("hello")
	chunks = str.Chunks(10)

	if chunks.Len() != 1 {
		t.Errorf("Expected 1 chunk, but got %v", chunks.Len())
	}

	if chunks[0] != str {
		t.Errorf("Expected chunk to be %v, but got %v", str, chunks.Get(0))
	}

	str = g.String("hello")
	chunks = str.Chunks(2)

	if chunks.Len() != 3 {
		t.Errorf("Expected 3 chunks, but got %v", chunks.Len())
	}

	expectedChunks := g.Slice[g.String]{"he", "ll", "o"}

	for i, c := range chunks {
		if c != expectedChunks.Get(i) {
			t.Errorf("Expected chunk %v to be %v, but got %v", i, expectedChunks.Get(i), c)
		}
	}

	str = g.String("hello world")
	chunks = str.Chunks(3)

	if chunks.Len() != 4 {
		t.Errorf("Expected 4 chunks, but got %v", chunks.Len())
	}

	expectedChunks = g.Slice[g.String]{"hel", "lo ", "wor", "ld"}

	for i, c := range chunks {
		if c != expectedChunks.Get(i) {
			t.Errorf("Expected chunk %v to be %v, but got %v", i, expectedChunks.Get(i), c)
		}
	}

	str = g.String("hello")
	chunks = str.Chunks(5)

	if chunks.Len() != 1 {
		t.Errorf("Expected 1 chunk, but got %v", chunks.Len())
	}

	if chunks.Get(0) != str {
		t.Errorf("Expected chunk to be %v, but got %v", str, chunks.Get(0))
	}

	str = g.String("hello")
	chunks = str.Chunks(-1)

	if chunks.Len() != 0 {
		t.Errorf("Expected empty slice, but got %v", chunks)
	}
}

func TestStringCut(t *testing.T) {
	tests := []struct {
		name   string
		input  g.String
		start  g.String
		end    g.String
		output g.String
	}{
		{"Basic", "Hello [start]world[end]!", "[start]", "[end]", "world"},
		{"No start", "Hello world!", "[start]", "[end]", ""},
		{"No end", "Hello [start]world!", "[start]", "[end]", ""},
		{"Start equals end", "Hello [tag]world[tag]!", "[tag]", "[tag]", "world"},
		{
			"Multiple instances",
			"A [start]first[end] B [start]second[end] C",
			"[start]",
			"[end]",
			"first",
		},
		{"Empty input", "", "[start]", "[end]", ""},
		{"Empty start and end", "Hello world!", "", "", ""},
		{
			"Nested tags",
			"A [start]first [start]nested[end] value[end] B",
			"[start]",
			"[end]",
			"first [start]nested",
		},
		{
			"Overlapping tags",
			"A [start]first[end][start]second[end] B",
			"[start]",
			"[end]",
			"first",
		},
		{"Unicode characters", "Привет [начало]мир[конец]!", "[начало]", "[конец]", "мир"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, result := test.input.Cut(test.start, test.end)
			if result != test.output {
				t.Errorf("Expected '%s', got '%s'", test.output, result)
			}
		})
	}
}

func TestURLEncode(t *testing.T) {
	testCases := []struct {
		input    g.String
		safe     g.String
		expected g.String
	}{
		{
			input:    "https://www.test.com/?query=test&param=value",
			safe:     "",
			expected: "https%3A%2F%2Fwww.test.com%2F%3Fquery%3Dtest%26param%3Dvalue",
		},
		{
			input:    "https://www.test.com/?query=test&param=value",
			safe:     ":",
			expected: "https:%2F%2Fwww.test.com%2F%3Fquery%3Dtest%26param%3Dvalue",
		},
		{
			input:    "https://www.test.com/?query=test&param=value",
			safe:     ":/&",
			expected: "https://www.test.com/%3Fquery%3Dtest&param%3Dvalue",
		},
		{
			input:    "https://www.test.com/?query=test&param=value",
			safe:     ":/&?",
			expected: "https://www.test.com/?query%3Dtest&param%3Dvalue",
		},
		{
			input:    "https://www.test.com/?query=test&param=value",
			safe:     ":/&?=",
			expected: "https://www.test.com/?query=test&param=value",
		},
	}

	for _, tc := range testCases {
		encoded := tc.input.Enc().URL(tc.safe)
		if encoded != tc.expected {
			t.Errorf(
				"For input: %s and safe: %s, expected: %s, but got: %s",
				tc.input,
				tc.safe,
				tc.expected,
				encoded.Std(),
			)
		}
	}
}

func TestURLDecode(t *testing.T) {
	tests := []struct {
		input    g.String
		expected g.String
	}{
		{
			input:    "hello+world",
			expected: "hello world",
		},
		{
			input:    "hello%20world",
			expected: "hello world",
		},
		{
			input:    "a%2Bb%3Dc%2Fd",
			expected: "a+b=c/d",
		},
		{
			input:    "foo%3Fbar%3Dbaz%26abc%3D123",
			expected: "foo?bar=baz&abc=123",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for _, test := range tests {
		actual := test.input.Dec().URL().Unwrap()
		if actual != test.expected {
			t.Errorf("UnEscape(%s): expected %s, but got %s", test.input, test.expected, actual)
		}
	}
}

func TestStringCompare(t *testing.T) {
	testCases := []struct {
		str1     g.String
		str2     g.String
		expected g.Int
	}{
		{"apple", "banana", -1},
		{"banana", "apple", 1},
		{"banana", "banana", 0},
		{"apple", "Apple", 1},
		{"", "", 0},
	}

	for _, tc := range testCases {
		result := tc.str1.Compare(tc.str2)
		if !result.Eq(tc.expected) {
			t.Errorf("Compare(%q, %q): expected %d, got %d", tc.str1, tc.str2, tc.expected, result)
		}
	}
}

func TestStringEq(t *testing.T) {
	testCases := []struct {
		str1     g.String
		str2     g.String
		expected bool
	}{
		{"apple", "banana", false},
		{"banana", "banana", true},
		{"Apple", "apple", false},
		{"", "", true},
	}

	for _, tc := range testCases {
		result := tc.str1.Eq(tc.str2)
		if result != tc.expected {
			t.Errorf("Eq(%q, %q): expected %t, got %t", tc.str1, tc.str2, tc.expected, result)
		}
	}
}

func TestStringNe(t *testing.T) {
	testCases := []struct {
		str1     g.String
		str2     g.String
		expected bool
	}{
		{"apple", "banana", true},
		{"banana", "banana", false},
		{"Apple", "apple", true},
		{"", "", false},
	}

	for _, tc := range testCases {
		result := tc.str1.Ne(tc.str2)
		if result != tc.expected {
			t.Errorf("Ne(%q, %q): expected %t, got %t", tc.str1, tc.str2, tc.expected, result)
		}
	}
}

func TestStringGt(t *testing.T) {
	testCases := []struct {
		str1     g.String
		str2     g.String
		expected bool
	}{
		{"apple", "banana", false},
		{"banana", "apple", true},
		{"Apple", "apple", false},
		{"banana", "banana", false},
		{"", "", false},
	}

	for _, tc := range testCases {
		result := tc.str1.Gt(tc.str2)
		if result != tc.expected {
			t.Errorf("Gt(%q, %q): expected %t, got %t", tc.str1, tc.str2, tc.expected, result)
		}
	}
}

func TestStringLt(t *testing.T) {
	testCases := []struct {
		str1     g.String
		str2     g.String
		expected bool
	}{
		{"apple", "banana", true},
		{"banana", "apple", false},
		{"Apple", "apple", true},
		{"banana", "banana", false},
		{"", "", false},
	}

	for _, tc := range testCases {
		result := tc.str1.Lt(tc.str2)
		if result != tc.expected {
			t.Errorf("Lt(%q, %q): expected %t, got %t", tc.str1, tc.str2, tc.expected, result)
		}
	}
}

func TestStringReverse(t *testing.T) {
	testCases := []struct {
		in      g.String
		wantOut g.String
	}{
		{in: "", wantOut: ""},
		{in: " ", wantOut: " "},
		{in: "a", wantOut: "a"},
		{in: "ab", wantOut: "ba"},
		{in: "abc", wantOut: "cba"},
		{in: "abcdefg", wantOut: "gfedcba"},
		{in: "ab丂d", wantOut: "d丂ba"},
		{in: "abåd", wantOut: "dåba"},

		{in: "世界", wantOut: "界世"},
		{in: "🙂🙃", wantOut: "🙃🙂"},
		{in: "こんにちは", wantOut: "はちにんこ"},

		// Punctuation and whitespace
		{in: "Hello, world!", wantOut: "!dlrow ,olleH"},
		{in: "Hello\tworld!", wantOut: "!dlrow\tolleH"},
		{in: "Hello\nworld!", wantOut: "!dlrow\nolleH"},

		// Mixed languages and scripts
		{in: "Hello, 世界!", wantOut: "!界世 ,olleH"},
		{in: "Привет, мир!", wantOut: "!рим ,тевирП"},
		{in: "안녕하세요, 세계!", wantOut: "!계세 ,요세하녕안"},

		// Palindromes
		{in: "racecar", wantOut: "racecar"},
		{in: "A man, a plan, a canal: Panama", wantOut: "amanaP :lanac a ,nalp a ,nam A"},

		{
			in:      "The quick brown fox jumps over the lazy dog.",
			wantOut: ".god yzal eht revo spmuj xof nworb kciuq ehT",
		},
		{in: "A man a plan a canal panama", wantOut: "amanap lanac a nalp a nam A"},
		{in: "Was it a car or a cat I saw?", wantOut: "?was I tac a ro rac a ti saW"},
		{in: "Never odd or even", wantOut: "neve ro ddo reveN"},
		{in: "Do geese see God?", wantOut: "?doG ees eseeg oD"},
		{in: "A Santa at NASA", wantOut: "ASAN ta atnaS A"},
		{in: "Yo, Banana Boy!", wantOut: "!yoB ananaB ,oY"},
		{in: "Madam, in Eden I'm Adam", wantOut: "madA m'I nedE ni ,madaM"},
		{in: "Never odd or even", wantOut: "neve ro ddo reveN"},
		{in: "Was it a car or a cat I saw?", wantOut: "?was I tac a ro rac a ti saW"},
		{in: "Do geese see God?", wantOut: "?doG ees eseeg oD"},
		{in: "No 'x' in Nixon", wantOut: "noxiN ni 'x' oN"},
		{in: "A Santa at NASA", wantOut: "ASAN ta atnaS A"},
		{in: "Yo, Banana Boy!", wantOut: "!yoB ananaB ,oY"},
	}

	for _, tc := range testCases {
		result := tc.in.Reverse()
		if result.Ne(tc.wantOut) {
			t.Errorf("Reverse(%s): expected %s, got %s", tc.in, result, tc.wantOut)
		}
	}
}

func TestStringNormalizeNFC(t *testing.T) {
	testCases := []struct {
		input    g.String
		expected g.String
	}{
		{input: "Mëtàl Hëàd", expected: "Mëtàl Hëàd"},
		{input: "Café", expected: "Café"},
		{input: "Ĵūņě", expected: "Ĵūņě"},
		{input: "𝓽𝓮𝓼𝓽 𝓬𝓪𝓼𝓮", expected: "𝓽𝓮𝓼𝓽 𝓬𝓪𝓼𝓮"},
		{input: "ḀϊṍẀṙṧ", expected: "ḀϊṍẀṙṧ"},
		{input: "えもじ れんしゅう", expected: "えもじ れんしゅう"},
		{input: "Научные исследования", expected: "Научные исследования"},
		{input: "🌟Unicode✨", expected: "🌟Unicode✨"},
		{input: "A\u0308", expected: "Ä"},
		{input: "o\u0308", expected: "ö"},
		{input: "u\u0308", expected: "ü"},
		{input: "O\u0308", expected: "Ö"},
		{input: "U\u0308", expected: "Ü"},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			output := tc.input.NormalizeNFC()
			if output != tc.expected {
				t.Errorf("Normalize(%q) = %q; want %q", tc.input, output, tc.expected)
			}
		})
	}
}

func TestStringSimilarity(t *testing.T) {
	testCases := []struct {
		str1     g.String
		str2     g.String
		expected g.Float
	}{
		{"hello", "hello", 100},
		{"hello", "world", 20},
		{"hello", "", 0},
		{"", "", 100},
		{"cat", "cats", 75},
		{"kitten", "sitting", 57.14},
		{"good", "bad", 25},
		{"book", "back", 50},
		{"abcdef", "azced", 50},
		{"tree", "three", 80},
		{"house", "horse", 80},
		{"language", "languish", 62.50},
		{"programming", "programmer", 72.73},
		{"algorithm", "logarithm", 77.78},
		{"software", "hardware", 50},
		{"tea", "ate", 33.33},
		{"pencil", "pen", 50},
		{"information", "informant", 63.64},
		{"coffee", "toffee", 83.33},
		{"developer", "develop", 77.78},
		{"distance", "difference", 50},
		{"similar", "similarity", 70},
		{"apple", "apples", 83.33},
		{"internet", "internets", 88.89},
		{"education", "dedication", 80},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			output := tc.str1.Similarity(tc.str2)
			if output.RoundDecimal(2).Ne(tc.expected) {
				t.Errorf(
					"g.String(\"%s\").SimilarText(\"%s\") = %.2f%% but want %.2f%%\n",
					tc.str1,
					tc.str2,
					output,
					tc.expected,
				)
			}
		})
	}
}

func TestStringReader(t *testing.T) {
	tests := []struct {
		name     string
		str      g.String
		expected string
	}{
		{"Empty String", "", ""},
		{"Single character String", "a", "a"},
		{"Multiple characters String", "hello world", "hello world"},
		{"String with special characters", "こんにちは、世界！", "こんにちは、世界！"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := test.str.Reader()
			resultBytes, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("Error reading from *strings.Reader: %v", err)
			}

			result := string(resultBytes)

			if result != test.expected {
				t.Errorf("Reader() content = %s, expected %s", result, test.expected)
			}
		})
	}
}

func TestStringContainsAny(t *testing.T) {
	testCases := []struct {
		name    string
		input   g.String
		substrs g.Slice[g.String]
		want    bool
	}{
		{
			name:    "ContainsAny_OneSubstringMatch",
			input:   "This is an example",
			substrs: []g.String{"This", "missing"},
			want:    true,
		},
		{
			name:    "ContainsAny_NoSubstringMatch",
			input:   "This is an example",
			substrs: []g.String{"notfound", "missing"},
			want:    false,
		},
		{
			name:    "ContainsAny_EmptySubstrings",
			input:   "This is an example",
			substrs: []g.String{},
			want:    false,
		},
		{
			name:    "ContainsAny_EmptyInput",
			input:   "",
			substrs: []g.String{"notfound", "missing"},
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.ContainsAny(tc.substrs...)
			if got != tc.want {
				t.Errorf("ContainsAny() = %v; want %v", got, tc.want)
			}
		})
	}
}

func TestStringContainsAll(t *testing.T) {
	testCases := []struct {
		name    string
		input   g.String
		substrs g.Slice[g.String]
		want    bool
	}{
		{
			name:    "ContainsAll_AllSubstringsMatch",
			input:   "This is an example",
			substrs: []g.String{"This", "example"},
			want:    true,
		},
		{
			name:    "ContainsAll_NotAllSubstringsMatch",
			input:   "This is an example",
			substrs: []g.String{"This", "missing"},
			want:    false,
		},
		{
			name:    "ContainsAll_EmptySubstrings",
			input:   "This is an example",
			substrs: []g.String{},
			want:    true,
		},
		{
			name:    "ContainsAll_EmptyInput",
			input:   "",
			substrs: []g.String{"notfound", "missing"},
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.ContainsAll(tc.substrs...)
			if got != tc.want {
				t.Errorf("ContainsAll() = %v; want %v", got, tc.want)
			}
		})
	}
}

func TestStringReplaceNth(t *testing.T) {
	tests := []struct {
		name     string
		str      g.String
		oldS     g.String
		newS     g.String
		n        int
		expected g.String
	}{
		{
			"First occurrence",
			"The quick brown dog jumped over the lazy dog.",
			"dog",
			"fox",
			1,
			"The quick brown fox jumped over the lazy dog.",
		},
		{
			"Second occurrence",
			"The quick brown dog jumped over the lazy dog.",
			"dog",
			"fox",
			2,
			"The quick brown dog jumped over the lazy fox.",
		},
		{
			"Last occurrence",
			"The quick brown dog jumped over the lazy dog.",
			"dog",
			"fox",
			-1,
			"The quick brown dog jumped over the lazy fox.",
		},
		{
			"Negative n (except -1)",
			"The quick brown dog jumped over the lazy dog.",
			"dog",
			"fox",
			-2,
			"The quick brown dog jumped over the lazy dog.",
		},
		{
			"Zero n",
			"The quick brown dog jumped over the lazy dog.",
			"dog",
			"fox",
			0,
			"The quick brown dog jumped over the lazy dog.",
		},
		{
			"Longer replacement",
			"Hello, world!",
			"world",
			"beautiful world",
			1,
			"Hello, beautiful world!",
		},
		{
			"Shorter replacement",
			"A wonderful day, isn't it?",
			"wonderful",
			"nice",
			1,
			"A nice day, isn't it?",
		},
		{
			"Replace entire string",
			"Hello, world!",
			"Hello, world!",
			"Greetings, world!",
			1,
			"Greetings, world!",
		},
		{"No replacement", "Hello, world!", "x", "y", 1, "Hello, world!"},
		{"Nonexistent substring", "Hello, world!", "foobar", "test", 1, "Hello, world!"},
		{"Replace empty string", "Hello, world!", "", "x", 1, "Hello, world!"},
		{"Multiple identical substrings", "banana", "na", "xy", 1, "baxyna"},
		{"Multiple identical substrings, last", "banana", "na", "xy", -1, "banaxy"},
		{"Replace with empty string", "Hello, world!", "world", "", 1, "Hello, !"},
		{"Empty input string", "", "world", "test", 1, ""},
		{"Empty input, empty oldS, empty newS", "", "", "", 1, ""},
		{"Replace multiple spaces", "Hello    world!", "    ", " ", 1, "Hello world!"},
		{"Unicode characters", "こんにちは世界！", "世界", "World", 1, "こんにちはWorld！"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.str.ReplaceNth(test.oldS, test.newS, test.n)
			if result != test.expected {
				t.Errorf("ReplaceNth() got %q, want %q", result, test.expected)
			}
		})
	}
}

func TestStringIsASCII(t *testing.T) {
	testCases := []struct {
		input    g.String
		expected bool
	}{
		{"Hello, world!", true},
		{"こんにちは", false},
		{"", true},
		{"1234567890", true},
		{"ABCabc", true},
		{"~`!@#$%^&*()-_+={[}]|\\:;\"'<,>.?/", true},
		{"áéíóú", false},
		{"Привет", false},
	}

	for _, tc := range testCases {
		result := tc.input.IsASCII()
		if result != tc.expected {
			t.Errorf("IsASCII(%q) returned %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestStringReplaceMulti(t *testing.T) {
	// Test case 1: Basic replacements
	input1 := g.String("Hello, world! This is a test.")
	replaced1 := input1.ReplaceMulti("Hello", "Greetings", "world", "universe", "test", "example")
	expected1 := g.String("Greetings, universe! This is a example.")
	if replaced1 != expected1 {
		t.Errorf("Test case 1 failed: Expected '%s', got '%s'", expected1, replaced1)
	}

	// Test case 2: Replacements with special characters
	input2 := g.String("The price is $100.00, not $200.00!")
	replaced2 := input2.ReplaceMulti("$100.00", "$50.00", "$200.00", "$80.00")
	expected2 := g.String("The price is $50.00, not $80.00!")
	if replaced2 != expected2 {
		t.Errorf("Test case 2 failed: Expected '%s', got '%s'", expected2, replaced2)
	}

	// Test case 3: No replacements
	input3 := g.String("No replacements here.")
	replaced3 := input3.ReplaceMulti("Hello", "Greetings", "world", "universe")
	if replaced3 != input3 {
		t.Errorf("Test case 3 failed: Expected '%s', got '%s'", input3, replaced3)
	}

	// Test case 4: Empty string
	input4 := g.String("")
	replaced4 := input4.ReplaceMulti("Hello", "Greetings")
	if replaced4 != input4 {
		t.Errorf("Test case 4 failed: Expected '%s', got '%s'", input4, replaced4)
	}
}

func TestStringToFloat(t *testing.T) {
	// Test cases for valid float strings
	validFloatCases := []struct {
		input    g.String
		expected g.Float
	}{
		{"3.14", g.Float(3.14)},
		{"-123.456", g.Float(-123.456)},
		{"0.0", g.Float(0.0)},
	}

	for _, testCase := range validFloatCases {
		result := testCase.input.ToFloat()
		if result.IsErr() {
			t.Errorf("ToFloat test failed for %s. Unexpected error: %v", testCase.input, result.Err())
		}

		if result.Ok().Ne(testCase.expected) {
			t.Errorf("ToFloat test failed for %s. Expected: %v, Got: %v",
				testCase.input,
				testCase.expected,
				result.Ok(),
			)
		}
	}

	// Test cases for invalid float strings
	invalidFloatCases := []g.String{"abc", "123abc", "12.34.56", "", " "}

	for _, input := range invalidFloatCases {
		result := input.ToFloat()
		if result.IsOk() {
			t.Errorf("ToFloat test failed for %s. Expected error, got result: %v", input, result.Ok())
		}
	}
}

func TestStringLower(t *testing.T) {
	// Test cases for lowercase conversion
	lowerCases := []struct {
		input    g.String
		expected g.String
	}{
		{"HELLO", "hello"},
		{"Hello World", "hello world"},
		{"123", "123"},
		{"", ""},
		{"AbCdEfG", "abcdefg"},
	}

	for _, testCase := range lowerCases {
		result := testCase.input.Lower()
		if !result.Eq(testCase.expected) {
			t.Errorf("Lower test failed for %s. Expected: %s, Got: %s", testCase.input, testCase.expected, result)
		}
	}
}

func TestStringTrim(t *testing.T) {
	// Test cases for Trim
	trimCases := []struct {
		input    g.String
		cutset   g.String
		expected g.String
	}{
		{"   Hello, World!   ", " ", "Hello, World!"},
		{"Hello, World!", ",! ", "Hello, World"},
		{"  Golang  ", " Go", "lang"},
		{"", "", ""},
	}

	for _, testCase := range trimCases {
		result := testCase.input.Trim(testCase.cutset)
		if !result.Eq(testCase.expected) {
			t.Errorf(
				"Trim test failed for %s with cutset %s. Expected: %s, Got: %s",
				testCase.input,
				testCase.cutset,
				testCase.expected,
				result,
			)
		}
	}
}

func TestStringTrimLeft(t *testing.T) {
	// Test cases for TrimLeft
	trimLeftCases := []struct {
		input    g.String
		cutset   g.String
		expected g.String
	}{
		{"   Hello, World!   ", " ", "Hello, World!   "},
		{"Hello, World!", ",! ", "Hello, World!"},
		{"  Golang  ", " Go", "lang  "},
		{"", "", ""},
	}

	for _, testCase := range trimLeftCases {
		result := testCase.input.TrimLeft(testCase.cutset)
		if !result.Eq(testCase.expected) {
			t.Errorf(
				"TrimLeft test failed for %s with cutset %s. Expected: %s, Got: %s",
				testCase.input,
				testCase.cutset,
				testCase.expected,
				result,
			)
		}
	}
}

func TestStringTrimPrefix(t *testing.T) {
	// Test cases for TrimPrefix
	trimPrefixCases := []struct {
		input    g.String
		prefix   g.String
		expected g.String
	}{
		{"Hello, World!", "Hello, ", "World!"},
		{"prefix-prefix-suffix", "prefix-", "prefix-suffix"},
		{"no prefix", "prefix-", "no prefix"},
		{"", "prefix-", ""},
	}

	for _, testCase := range trimPrefixCases {
		result := testCase.input.TrimPrefix(testCase.prefix)
		if !result.Eq(testCase.expected) {
			t.Errorf(
				"TrimPrefix test failed for %s with prefix %s. Expected: %s, Got: %s",
				testCase.input,
				testCase.prefix,
				testCase.expected,
				result,
			)
		}
	}
}

func TestStringTrimSuffix(t *testing.T) {
	// Test cases for TrimSuffix
	trimSuffixCases := []struct {
		input    g.String
		suffix   g.String
		expected g.String
	}{
		{"Hello, World!", ", World!", "Hello"},
		{"prefix-prefix-suffix", "-suffix", "prefix-prefix"},
		{"no suffix", "-suffix", "no suffix"},
		{"", "-suffix", ""},
	}

	for _, testCase := range trimSuffixCases {
		result := testCase.input.TrimSuffix(testCase.suffix)
		if !result.Eq(testCase.expected) {
			t.Errorf(
				"TrimSuffix test failed for %s with suffix %s. Expected: %s, Got: %s",
				testCase.input,
				testCase.suffix,
				testCase.expected,
				result,
			)
		}
	}
}

func TestStringReplace(t *testing.T) {
	// Test cases for Replace
	replaceCases := []struct {
		input    g.String
		oldS     g.String
		newS     g.String
		n        int
		expected g.String
	}{
		{"Hello, World!", "Hello", "Hi", 1, "Hi, World!"},
		{"Hello, Hello, Hello!", "Hello", "Hi", -1, "Hi, Hi, Hi!"},
		{"prefix-prefix-suffix", "prefix", "pre", 1, "pre-prefix-suffix"},
		{"no match", "match", "replacement", 1, "no replacement"},
		{"", "", "", 0, ""},
	}

	for _, testCase := range replaceCases {
		result := testCase.input.Replace(testCase.oldS, testCase.newS, testCase.n)
		if !result.Eq(testCase.expected) {
			t.Errorf(
				"Replace test failed for %s with oldS %s and newS %s. Expected: %s, Got: %s",
				testCase.input,
				testCase.oldS,
				testCase.newS,
				testCase.expected,
				result,
			)
		}
	}
}

func TestStringContainsAnyChars(t *testing.T) {
	// Test cases for ContainsAnyChars
	containsAnyCharsCases := []struct {
		input    g.String
		chars    g.String
		expected bool
	}{
		{"Hello, World!", "aeiou", true}, // Contains vowels
		{"1234567890", "aeiou", false},   // Does not contain vowels
		{"Hello, World!", "abc", false},  // Contains a, b, or c
		{"Hello, World!", "123", false},  // Does not contain 1, 2, or 3
		{"", "aeiou", false},             // Empty string
	}

	for _, testCase := range containsAnyCharsCases {
		result := testCase.input.ContainsAnyChars(testCase.chars)
		if result != testCase.expected {
			t.Errorf(
				"ContainsAnyChars test failed for %s with chars %s. Expected: %t, Got: %t",
				testCase.input,
				testCase.chars,
				testCase.expected,
				result,
			)
		}
	}
}

func TestStringSplitLines(t *testing.T) {
	// Test case 1: String with multiple lines.
	str1 := g.NewString("hello\nworld\nhow\nare\nyou\n")
	expected1 := g.Slice[g.String]{"hello", "world", "how", "are", "you"}
	result1 := str1.SplitLines()
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed: Expected %v, got %v", expected1, result1)
	}

	// Test case 2: String with single line.
	str2 := g.NewString("hello")
	expected2 := g.Slice[g.String]{"hello"}
	result2 := str2.SplitLines()
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed: Expected %v, got %v", expected2, result2)
	}
}

func TestStringSplitN(t *testing.T) {
	// Test case 1: String with multiple segments, n > 0.
	str1 := g.NewString("hello,world,how,are,you")
	sep1 := g.NewString(",")
	n1 := 3
	expected1 := g.Slice[g.String]{"hello", "world", "how,are,you"}
	result1 := str1.SplitN(sep1, n1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed: Expected %v, got %v", expected1, result1)
	}

	// Test case 2: String with multiple segments, n < 0.
	str2 := g.NewString("hello,world,how,are,you")
	sep2 := g.NewString(",")
	n2 := -1
	expected2 := g.Slice[g.String]{"hello", "world", "how", "are", "you"}
	result2 := str2.SplitN(sep2, n2)
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed: Expected %v, got %v", expected2, result2)
	}

	// Test case 3: String with single segment, n > 0.
	str3 := g.NewString("hello")
	sep3 := g.NewString(",")
	n3 := 1
	expected3 := g.Slice[g.String]{"hello"}
	result3 := str3.SplitN(sep3, n3)
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed: Expected %v, got %v", expected3, result3)
	}
}

func TestStringFields(t *testing.T) {
	// Test case 1: String with multiple words separated by whitespace.
	str1 := g.NewString("hello world how are you")
	expected1 := g.Slice[g.String]{"hello", "world", "how", "are", "you"}
	result1 := str1.Fields()
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed: Expected %v, got %v", expected1, result1)
	}

	// Test case 2: String with single word.
	str2 := g.NewString("hello")
	expected2 := g.Slice[g.String]{"hello"}
	result2 := str2.Fields()
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed: Expected %v, got %v", expected2, result2)
	}

	// Test case 3: Empty string.
	str3 := g.NewString("")
	expected3 := g.Slice[g.String]{}
	result3 := str3.Fields()
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed: Expected %v, got %v", expected3, result3)
	}

	// Test case 4: String with leading and trailing whitespace.
	str4 := g.NewString("   hello   world   ")
	expected4 := g.Slice[g.String]{"hello", "world"}
	result4 := str4.Fields()
	if !reflect.DeepEqual(result4, expected4) {
		t.Errorf("Test case 4 failed: Expected %v, got %v", expected4, result4)
	}
}

func TestStringCount(t *testing.T) {
	// Test case 1: Count occurrences of substring in a string with multiple occurrences.
	str1 := g.NewString("hello world hello hello")
	substr1 := g.NewString("hello")
	expected1 := 3
	result1 := str1.Count(substr1)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %d, got %d", expected1, result1)
	}

	// Test case 2: Count occurrences of substring in a string with no occurrences.
	str2 := g.NewString("abcdefg")
	substr2 := g.NewString("xyz")
	expected2 := 0
	result2 := str2.Count(substr2)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %d, got %d", expected2, result2)
	}

	// Test case 3: Count occurrences of substring in an empty string.
	str3 := g.NewString("")
	substr3 := g.NewString("hello")
	expected3 := 0
	result3 := str3.Count(substr3)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: Expected %d, got %d", expected3, result3)
	}
}

func TestStringEqFold(t *testing.T) {
	// Test case 1: Strings are equal case-insensitively.
	str1 := g.NewString("Hello")
	str2 := g.NewString("hello")
	expected1 := true
	result1 := str1.EqFold(str2)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %t, got %t", expected1, result1)
	}

	// Test case 2: Strings are not equal case-insensitively.
	str3 := g.NewString("world")
	expected2 := false
	result2 := str1.EqFold(str3)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %t, got %t", expected2, result2)
	}

	// Test case 3: Empty strings.
	str4 := g.NewString("")
	str5 := g.NewString("")
	expected3 := true
	result3 := str4.EqFold(str5)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: Expected %t, got %t", expected3, result3)
	}
}

func TestStringLastIndex(t *testing.T) {
	// Test case 1: Substring is present in the string.
	str1 := g.NewString("hello world hello")
	substr1 := g.NewString("hello")
	expected1 := 12
	result1 := str1.LastIndex(substr1)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %d, got %d", expected1, result1)
	}

	// Test case 2: Substring is not present in the string.
	substr2 := g.NewString("foo")
	expected2 := -1
	result2 := str1.LastIndex(substr2)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %d, got %d", expected2, result2)
	}

	// Test case 3: Empty string.
	str3 := g.NewString("")
	substr3 := g.NewString("hello")
	expected3 := -1
	result3 := str3.LastIndex(substr3)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: Expected %d, got %d", expected3, result3)
	}
}

func TestStringIndexRune(t *testing.T) {
	// Test case 1: Rune is present in the string.
	str1 := g.NewString("hello")
	rune1 := 'e'
	expected1 := 1
	result1 := str1.IndexRune(rune1)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %d, got %d", expected1, result1)
	}

	// Test case 2: Rune is not present in the string.
	rune2 := 'x'
	expected2 := -1
	result2 := str1.IndexRune(rune2)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %d, got %d", expected2, result2)
	}

	// Test case 3: Empty string.
	str3 := g.NewString("")
	rune3 := 'h'
	expected3 := -1
	result3 := str3.IndexRune(rune3)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: Expected %d, got %d", expected3, result3)
	}
}

func TestStringNotEmpty(t *testing.T) {
	// Test case 1: String is not empty.
	str1 := g.NewString("hello")
	expected1 := true
	result1 := str1.NotEmpty()
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %t, got %t", expected1, result1)
	}

	// Test case 2: String is empty.
	str2 := g.NewString("")
	expected2 := false
	result2 := str2.NotEmpty()
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %t, got %t", expected2, result2)
	}
}

func TestRepeat(t *testing.T) {
	// Test case 1: Repeat count is positive.
	str := g.NewString("abc")
	count := 3
	expected1 := g.NewString("abcabcabc")
	result1 := str.Repeat(count)
	if !result1.Eq(expected1) {
		t.Errorf("Test case 1 failed: Expected %s, got %s", expected1, result1)
	}

	// Test case 2: Repeat count is zero.
	count = 0
	expected2 := g.NewString("")
	result2 := str.Repeat(count)
	if !result2.Eq(expected2) {
		t.Errorf("Test case 2 failed: Expected %s, got %s", expected2, result2)
	}
}

func TestStringLeftJustify(t *testing.T) {
	// Test case 1: Original string length is less than the specified length.
	str1 := g.NewString("Hello")
	length1 := 10
	pad1 := g.NewString(".")
	expected1 := g.NewString("Hello.....")
	result1 := str1.LeftJustify(length1, pad1)
	if !result1.Eq(expected1) {
		t.Errorf("Test case 1 failed: Expected %s, got %s", expected1, result1)
	}

	// Test case 2: Original string length is equal to the specified length.
	str2 := g.NewString("Hello")
	length2 := 5
	pad2 := g.NewString(".")
	expected2 := g.NewString("Hello")
	result2 := str2.LeftJustify(length2, pad2)
	if !result2.Eq(expected2) {
		t.Errorf("Test case 2 failed: Expected %s, got %s", expected2, result2)
	}

	// Test case 3: Original string length is greater than the specified length.
	str3 := g.NewString("Hello")
	length3 := 3
	pad3 := g.NewString(".")
	expected3 := g.NewString("Hello")
	result3 := str3.LeftJustify(length3, pad3)
	if !result3.Eq(expected3) {
		t.Errorf("Test case 3 failed: Expected %s, got %s", expected3, result3)
	}

	// Test case 4: Empty padding string.
	str4 := g.NewString("Hello")
	length4 := 10
	pad4 := g.NewString("")
	expected4 := g.NewString("Hello")
	result4 := str4.LeftJustify(length4, pad4)
	if !result4.Eq(expected4) {
		t.Errorf("Test case 4 failed: Expected %s, got %s", expected4, result4)
	}
}

func TestStringRightJustify(t *testing.T) {
	// Test case 1: Original string length is less than the specified length.
	str1 := g.NewString("Hello")
	length1 := 10
	pad1 := g.NewString(".")
	expected1 := g.NewString(".....Hello")
	result1 := str1.RightJustify(length1, pad1)
	if !result1.Eq(expected1) {
		t.Errorf("Test case 1 failed: Expected %s, got %s", expected1, result1)
	}

	// Test case 2: Original string length is equal to the specified length.
	str2 := g.NewString("Hello")
	length2 := 5
	pad2 := g.NewString(".")
	expected2 := g.NewString("Hello")
	result2 := str2.RightJustify(length2, pad2)
	if !result2.Eq(expected2) {
		t.Errorf("Test case 2 failed: Expected %s, got %s", expected2, result2)
	}

	// Test case 3: Original string length is greater than the specified length.
	str3 := g.NewString("Hello")
	length3 := 3
	pad3 := g.NewString(".")
	expected3 := g.NewString("Hello")
	result3 := str3.RightJustify(length3, pad3)
	if !result3.Eq(expected3) {
		t.Errorf("Test case 3 failed: Expected %s, got %s", expected3, result3)
	}

	// Test case 4: Empty padding string.
	str4 := g.NewString("Hello")
	length4 := 10
	pad4 := g.NewString("")
	expected4 := g.NewString("Hello")
	result4 := str4.RightJustify(length4, pad4)
	if !result4.Eq(expected4) {
		t.Errorf("Test case 4 failed: Expected %s, got %s", expected4, result4)
	}
}

func TestStringCenter(t *testing.T) {
	// Test case 1: Original string length is less than the specified length.
	str1 := g.NewString("Hello")
	length1 := 10
	pad1 := g.NewString(".")
	expected1 := g.NewString("..Hello...")
	result1 := str1.Center(length1, pad1)
	if !result1.Eq(expected1) {
		t.Errorf("Test case 1 failed: Expected %s, got %s", expected1, result1)
	}

	// Test case 2: Original string length is equal to the specified length.
	str2 := g.NewString("Hello")
	length2 := 5
	pad2 := g.NewString(".")
	expected2 := g.NewString("Hello")
	result2 := str2.Center(length2, pad2)
	if !result2.Eq(expected2) {
		t.Errorf("Test case 2 failed: Expected %s, got %s", expected2, result2)
	}

	// Test case 3: Original string length is greater than the specified length.
	str3 := g.NewString("Hello")
	length3 := 3
	pad3 := g.NewString(".")
	expected3 := g.NewString("Hello")
	result3 := str3.Center(length3, pad3)
	if !result3.Eq(expected3) {
		t.Errorf("Test case 3 failed: Expected %s, got %s", expected3, result3)
	}

	// Test case 4: Empty padding string.
	str4 := g.NewString("Hello")
	length4 := 10
	pad4 := g.NewString("")
	expected4 := g.NewString("Hello")
	result4 := str4.Center(length4, pad4)
	if !result4.Eq(expected4) {
		t.Errorf("Test case 4 failed: Expected %s, got %s", expected4, result4)
	}
}

func TestStringEndsWith(t *testing.T) {
	// Test case 1: String ends with one of the provided suffixes.
	str1 := g.NewString("example.com")
	suffixes1 := g.Slice[g.String]{g.NewString(".com"), g.NewString(".net")}
	expected1 := true
	result1 := str1.EndsWith(suffixes1...)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %t, got %t", expected1, result1)
	}

	// Test case 2: String ends with multiple provided suffixes.
	str2 := g.NewString("example.net")
	suffixes2 := g.Slice[g.String]{g.NewString(".com"), g.NewString(".net")}
	expected2 := true
	result2 := str2.EndsWith(suffixes2...)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %t, got %t", expected2, result2)
	}

	// Test case 3: String does not end with any of the provided suffixes.
	str3 := g.NewString("example.org")
	suffixes3 := g.Slice[g.String]{g.NewString(".com"), g.NewString(".net")}
	expected3 := false
	result3 := str3.EndsWith(suffixes3...)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: Expected %t, got %t", expected3, result3)
	}
}
