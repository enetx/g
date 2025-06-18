package g_test

import (
	"io"
	"math/big"
	"reflect"
	"testing"
	"unicode"

	. "github.com/enetx/g"
)

func TestSubString(t *testing.T) {
	tests := []struct {
		name     string
		input    String
		start    Int
		end      Int
		step     []Int
		expected string
	}{
		{
			name:     "Basic substring",
			input:    "Hello, world!",
			start:    0,
			end:      5,
			expected: "Hello",
		},
		{
			name:     "With negative start",
			input:    "Hello, world!",
			start:    -6,
			end:      12,
			expected: "world",
		},
		{
			name:     "With negative end",
			input:    "Hello, world!",
			start:    7,
			end:      -1,
			expected: "world",
		},
		{
			name:     "Start exceeds end",
			input:    "Hello, world!",
			start:    5,
			end:      1,
			expected: "",
		},
		{
			name:     "Step parameter",
			input:    "abcdef",
			start:    0,
			end:      6,
			step:     []Int{2},
			expected: "ace",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := String(tc.input)
			result := s.SubString(tc.start, tc.end, tc.step...)
			if result.Std() != tc.expected {
				t.Errorf("Test %s failed: expected %s, got %s", tc.name, tc.expected, result.Std())
			}
		})
	}
}

func TestStringTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    String
		max      Int
		expected String
	}{
		// Basic truncation
		{
			name:     "Basic truncation",
			input:    String("Hello, World!"),
			max:      5,
			expected: String("Hello..."),
		},
		// No truncation (length less than max)
		{
			name:     "No truncation (shorter than max)",
			input:    String("Short"),
			max:      10,
			expected: String("Short"),
		},
		// Exact length (no truncation)
		{
			name:     "Exact length",
			input:    String("Perfect"),
			max:      7,
			expected: String("Perfect"),
		},
		// Truncation of Unicode characters
		{
			name:     "Truncation with Unicode",
			input:    String("😊😊😊😊😊"),
			max:      3,
			expected: String("😊😊😊..."),
		},
		// Truncation with mixed characters
		{
			name:     "Truncation with mixed characters",
			input:    String("Hello😊World"),
			max:      6,
			expected: String("Hello😊..."),
		},
		// Empty input
		{
			name:     "Empty input",
			input:    String(""),
			max:      5,
			expected: String(""),
		},
		// Zero max length
		{
			name:     "Zero max length",
			input:    String("Zero length"),
			max:      0,
			expected: String("..."),
		},
		// Negative max length (invalid case)
		{
			name:     "Negative max length",
			input:    String("Negative case"),
			max:      -1,
			expected: String("Negative case"), // No truncation, invalid max
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Truncate(tt.max)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestStringBuilder(t *testing.T) {
	// Create a String
	str := String("hello")

	// Call the Builder method
	builder := str.Builder()

	// Check if the Builder has been properly initialized
	expected := "hello"
	if result := builder.String().Std(); result != expected {
		t.Errorf("Builder() = %s; want %s", result, expected)
	}
}

func TestStringMinMax(t *testing.T) {
	// Test cases for Min method
	minTestCases := []struct {
		inputs   []String
		expected String
	}{
		{[]String{"apple", "banana", "orange"}, "apple"},
		{[]String{"cat", "dog", "elephant"}, "cat"},
		{[]String{"123", "456", "789"}, "123"},
	}

	for _, testCase := range minTestCases {
		result := String(testCase.inputs[0]).Min(testCase.inputs[1:]...)
		if result != testCase.expected {
			t.Errorf("Min test failed. Expected: %s, Got: %s", testCase.expected, result)
		}
	}

	// Test cases for Max method
	maxTestCases := []struct {
		inputs   []String
		expected String
	}{
		{[]String{"apple", "banana", "orange"}, "orange"},
		{[]String{"cat", "dog", "elephant"}, "elephant"},
		{[]String{"123", "456", "789"}, "789"},
	}

	for _, testCase := range maxTestCases {
		result := String(testCase.inputs[0]).Max(testCase.inputs[1:]...)
		if result != testCase.expected {
			t.Errorf("Max test failed. Expected: %s, Got: %s", testCase.expected, result)
		}
	}
}

func TestChars(t *testing.T) {
	// Testing on a regular string with ASCII characters
	asciiStr := String("Hello")
	asciiExpected := Slice[String]{"H", "e", "l", "l", "o"}
	asciiResult := asciiStr.Chars().Collect()
	if !reflect.DeepEqual(asciiExpected, asciiResult) {
		t.Errorf("Expected %v, but got %v", asciiExpected, asciiResult)
	}

	// Testing on a string with Unicode characters (Russian)
	unicodeStr := String("Привет")
	unicodeExpected := Slice[String]{"П", "р", "и", "в", "е", "т"}
	unicodeResult := unicodeStr.Chars().Collect()
	if !reflect.DeepEqual(unicodeExpected, unicodeResult) {
		t.Errorf("Expected %v, but got %v", unicodeExpected, unicodeResult)
	}

	// Testing on a string with Unicode characters (Chinese)
	chineseStr := String("你好")
	chineseExpected := Slice[String]{"你", "好"}
	chineseResult := chineseStr.Chars().Collect()
	if !reflect.DeepEqual(chineseExpected, chineseResult) {
		t.Errorf("Expected %v, but got %v", chineseExpected, chineseResult)
	}

	// Additional test with a mix of ASCII and Unicode characters
	mixedStr := String("Hello 你好")
	mixedExpected := Slice[String]{"H", "e", "l", "l", "o", " ", "你", "好"}
	mixedResult := mixedStr.Chars().Collect()
	if !reflect.DeepEqual(mixedExpected, mixedResult) {
		t.Errorf("Expected %v, but got %v", mixedExpected, mixedResult)
	}

	// Testing on a string with special characters and symbols
	specialStr := String("Hello, 你好! How are you today? こんにちは")
	specialExpected := Slice[String]{
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
	specialResult := specialStr.Chars().Collect()
	if !reflect.DeepEqual(specialExpected, specialResult) {
		t.Errorf("Expected %v, but got %v", specialExpected, specialResult)
	}

	// Testing on a string with emojis
	emojiStr := String("Hello, 😊🌍🚀")
	emojiExpected := Slice[String]{"H", "e", "l", "l", "o", ",", " ", "😊", "🌍", "🚀"}
	emojiResult := emojiStr.Chars().Collect()
	if !reflect.DeepEqual(emojiExpected, emojiResult) {
		t.Errorf("Expected %v, but got %v", emojiExpected, emojiResult)
	}

	// Testing on an empty string
	emptyStr := String("")
	emptyExpected := Slice[String]{}
	emptyResult := emptyStr.Chars().Collect()
	if !reflect.DeepEqual(emptyExpected, emptyResult) {
		t.Errorf("Expected %v, but got %v", emptyExpected, emptyResult)
	}
}

func TestStringIsDigit(t *testing.T) {
	tests := []struct {
		name string
		str  String
		want bool
	}{
		{"empty", String(""), false},
		{"one", String("1"), true},
		{"nine", String("99999"), true},
		{"non-digit", String("1111a"), false},
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
		str  String
		want Int
	}{
		{
			name: "empty",
			str:  String(""),
			want: 0,
		},
		{
			name: "one digit",
			str:  String("1"),
			want: 1,
		},
		{
			name: "two digits",
			str:  String("12"),
			want: 12,
		},
		{
			name: "one letter",
			str:  String("a"),
			want: 0,
		},
		{
			name: "one digit and one letter",
			str:  String("1a"),
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
		str  String
		want String
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
		str  String
		s    String
		want String
	}{
		{
			name: "empty",
			str:  String(""),
			s:    String(""),
			want: String(""),
		},
		{
			name: "empty_hs",
			str:  String(""),
			s:    String("test"),
			want: String("test"),
		},
		{
			name: "empty_s",
			str:  String("test"),
			s:    String(""),
			want: String("test"),
		},
		{
			name: "not_empty",
			str:  String("test"),
			s:    String("test"),
			want: String("testtest"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.str.Append(tt.s); got != tt.want {
				t.Errorf("String.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringAddPrefix(t *testing.T) {
	tests := []struct {
		name string
		str  String
		s    String
		want String
	}{
		{
			name: "empty",
			str:  String(""),
			s:    String(""),
			want: String(""),
		},
		{
			name: "empty_hs",
			str:  String(""),
			s:    String("test"),
			want: String("test"),
		},
		{
			name: "empty_s",
			str:  String("test"),
			s:    String(""),
			want: String("test"),
		},
		{
			name: "not_empty",
			str:  String("rest"),
			s:    String("test"),
			want: String("testrest"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.str.Prepend(tt.s); got != tt.want {
				t.Errorf("String.AddPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringRandom(t *testing.T) {
	for i := range 100 {
		random := String("").Random(Int(i))

		if random.Len().Std() != i {
			t.Errorf("Random string length %d is not equal to %d", random.Len(), i)
		}
	}
}

func TestStringChunks(t *testing.T) {
	str := String("")
	chunks := str.Chunks(3).Collect()

	if chunks.Len() != 0 {
		t.Errorf("Expected empty slice, but got %v", chunks)
	}

	str = String("hello")
	chunks = str.Chunks(10).Collect()

	if chunks.Len() != 1 {
		t.Errorf("Expected 1 chunk, but got %v", chunks.Len())
	}

	if chunks[0] != str {
		t.Errorf("Expected chunk to be %v, but got %v", str, chunks.Get(0))
	}

	str = String("hello")
	chunks = str.Chunks(2).Collect()

	if chunks.Len() != 3 {
		t.Errorf("Expected 3 chunks, but got %v", chunks.Len())
	}

	expectedChunks := Slice[String]{"he", "ll", "o"}

	for i, c := range chunks {
		if c != expectedChunks[i] {
			t.Errorf("Expected chunk %v to be %v, but got %v", i, expectedChunks[i], c)
		}
	}

	str = String("hello world")
	chunks = str.Chunks(3).Collect()

	if chunks.Len() != 4 {
		t.Errorf("Expected 4 chunks, but got %v", chunks.Len())
	}

	expectedChunks = Slice[String]{"hel", "lo ", "wor", "ld"}

	for i, c := range chunks {
		if c != expectedChunks[i] {
			t.Errorf("Expected chunk %v to be %v, but got %v", i, expectedChunks[i], c)
		}
	}

	str = String("hello")
	chunks = str.Chunks(5).Collect()

	if chunks.Len() != 1 {
		t.Errorf("Expected 1 chunk, but got %v", chunks.Len())
	}

	if chunks.Get(0).Some() != str {
		t.Errorf("Expected chunk to be %v, but got %v", str, chunks.Get(0))
	}

	str = String("hello")
	chunks = str.Chunks(-1).Collect()

	if chunks.Len() != 0 {
		t.Errorf("Expected empty slice, but got %v", chunks)
	}
}

func TestStringCut(t *testing.T) {
	tests := []struct {
		name   string
		input  String
		start  String
		end    String
		output String
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

func TestStringCompare(t *testing.T) {
	testCases := []struct {
		str1     String
		str2     String
		expected int
	}{
		{"apple", "banana", -1},
		{"banana", "apple", 1},
		{"banana", "banana", 0},
		{"apple", "Apple", 1},
		{"", "", 0},
	}

	for _, tc := range testCases {
		result := tc.str1.Cmp(tc.str2)
		if int(result) != tc.expected {
			t.Errorf("Compare(%q, %q): expected %d, got %d", tc.str1, tc.str2, tc.expected, result)
		}
	}
}

func TestStringEq(t *testing.T) {
	testCases := []struct {
		str1     String
		str2     String
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
		str1     String
		str2     String
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
		str1     String
		str2     String
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
		str1     String
		str2     String
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

func TestStringLte(t *testing.T) {
	testCases := []struct {
		str1     String
		str2     String
		expected bool
	}{
		{"apple", "banana", true},
		{"banana", "apple", false},
		{"Apple", "apple", true},
		{"banana", "banana", true}, // Equal strings should return true
		{"", "", true},             // Empty strings should return true
	}

	for _, tc := range testCases {
		result := tc.str1.Lte(tc.str2)
		if result != tc.expected {
			t.Errorf("Lte(%q, %q): expected %t, got %t", tc.str1, tc.str2, tc.expected, result)
		}
	}
}

func TestStringGte(t *testing.T) {
	testCases := []struct {
		str1     String
		str2     String
		expected bool
	}{
		{"apple", "banana", false},
		{"banana", "apple", true},
		{"Apple", "apple", false},
		{"banana", "banana", true}, // Equal strings should return true
		{"", "", true},             // Empty strings should return true
	}

	for _, tc := range testCases {
		result := tc.str1.Gte(tc.str2)
		if result != tc.expected {
			t.Errorf("Gte(%q, %q): expected %t, got %t", tc.str1, tc.str2, tc.expected, result)
		}
	}
}

func TestStringReverse(t *testing.T) {
	testCases := []struct {
		in      String
		wantOut String
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
		input    String
		expected String
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
		str1     String
		str2     String
		expected Float
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
					"String(\"%s\").SimilarText(\"%s\") = %.2f%% but want %.2f%%\n",
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
		str      String
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
		input   String
		substrs Slice[String]
		want    bool
	}{
		{
			name:    "ContainsAny_OneSubstringMatch",
			input:   "This is an example",
			substrs: []String{"This", "missing"},
			want:    true,
		},
		{
			name:    "ContainsAny_NoSubstringMatch",
			input:   "This is an example",
			substrs: []String{"notfound", "missing"},
			want:    false,
		},
		{
			name:    "ContainsAny_EmptySubstrings",
			input:   "This is an example",
			substrs: []String{},
			want:    false,
		},
		{
			name:    "ContainsAny_EmptyInput",
			input:   "",
			substrs: []String{"notfound", "missing"},
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
		input   String
		substrs Slice[String]
		want    bool
	}{
		{
			name:    "ContainsAll_AllSubstringsMatch",
			input:   "This is an example",
			substrs: []String{"This", "example"},
			want:    true,
		},
		{
			name:    "ContainsAll_NotAllSubstringsMatch",
			input:   "This is an example",
			substrs: []String{"This", "missing"},
			want:    false,
		},
		{
			name:    "ContainsAll_EmptySubstrings",
			input:   "This is an example",
			substrs: []String{},
			want:    true,
		},
		{
			name:    "ContainsAll_EmptyInput",
			input:   "",
			substrs: []String{"notfound", "missing"},
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
		str      String
		oldS     String
		newS     String
		n        Int
		expected String
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
		input    String
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
	input1 := String("Hello, world! This is a test.")
	replaced1 := input1.ReplaceMulti("Hello", "Greetings", "world", "universe", "test", "example")
	expected1 := String("Greetings, universe! This is a example.")
	if replaced1 != expected1 {
		t.Errorf("Test case 1 failed: Expected '%s', got '%s'", expected1, replaced1)
	}

	// Test case 2: Replacements with special characters
	input2 := String("The price is $100.00, not $200.00!")
	replaced2 := input2.ReplaceMulti("$100.00", "$50.00", "$200.00", "$80.00")
	expected2 := String("The price is $50.00, not $80.00!")
	if replaced2 != expected2 {
		t.Errorf("Test case 2 failed: Expected '%s', got '%s'", expected2, replaced2)
	}

	// Test case 3: No replacements
	input3 := String("No replacements here.")
	replaced3 := input3.ReplaceMulti("Hello", "Greetings", "world", "universe")
	if replaced3 != input3 {
		t.Errorf("Test case 3 failed: Expected '%s', got '%s'", input3, replaced3)
	}

	// Test case 4: Empty string
	input4 := String("")
	replaced4 := input4.ReplaceMulti("Hello", "Greetings")
	if replaced4 != input4 {
		t.Errorf("Test case 4 failed: Expected '%s', got '%s'", input4, replaced4)
	}
}

func TestStringToFloat(t *testing.T) {
	// Test cases for valid float strings
	validFloatCases := []struct {
		input    String
		expected Float
	}{
		{"3.14", Float(3.14)},
		{"-123.456", Float(-123.456)},
		{"0.0", Float(0.0)},
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
	invalidFloatCases := []String{"abc", "123abc", "12.34.56", "", " "}

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
		input    String
		expected String
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

func TestStringStripPrefix(t *testing.T) {
	stripPrefixCases := []struct {
		input    String
		prefix   String
		expected String
	}{
		{"Hello, World!", "Hello, ", "World!"},
		{"prefix-prefix-suffix", "prefix-", "prefix-suffix"},
		{"no prefix", "prefix-", "no prefix"},
		{"", "prefix-", ""},
	}

	for _, testCase := range stripPrefixCases {
		result := testCase.input.StripPrefix(testCase.prefix)
		if !result.Eq(testCase.expected) {
			t.Errorf(
				"StripPrefix test failed for %s with prefix %s. Expected: %s, Got: %s",
				testCase.input,
				testCase.prefix,
				testCase.expected,
				result,
			)
		}
	}
}

func TestStringStripSuffix(t *testing.T) {
	stripSuffixCases := []struct {
		input    String
		suffix   String
		expected String
	}{
		{"Hello, World!", ", World!", "Hello"},
		{"prefix-prefix-suffix", "-suffix", "prefix-prefix"},
		{"no suffix", "-suffix", "no suffix"},
		{"", "-suffix", ""},
	}

	for _, testCase := range stripSuffixCases {
		result := testCase.input.StripSuffix(testCase.suffix)
		if !result.Eq(testCase.expected) {
			t.Errorf(
				"StripSuffix test failed for %s with suffix %s. Expected: %s, Got: %s",
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
		input    String
		oldS     String
		newS     String
		n        Int
		expected String
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
		input    String
		chars    String
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
	str1 := String("hello\nworld\nhow\nare\nyou\n")
	expected1 := Slice[String]{"hello", "world", "how", "are", "you"}
	result1 := str1.Lines().Collect()
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed: Expected %v, got %v", expected1, result1)
	}

	// Test case 2: String with single line.
	str2 := String("hello")
	expected2 := Slice[String]{"hello"}
	result2 := str2.Lines().Collect()
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed: Expected %v, got %v", expected2, result2)
	}
}

func TestStringSplitN(t *testing.T) {
	// Test case 1: String with multiple segments, n > 0.
	str1 := String("hello,world,how,are,you")
	sep1 := String(",")
	n1 := Int(3)
	expected1 := Slice[String]{"hello", "world", "how,are,you"}
	result1 := str1.SplitN(sep1, n1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed: Expected %v, got %v", expected1, result1)
	}

	// Test case 2: String with multiple segments, n < 0.
	str2 := String("hello,world,how,are,you")
	sep2 := String(",")
	n2 := Int(-1)
	expected2 := Slice[String]{"hello", "world", "how", "are", "you"}
	result2 := str2.SplitN(sep2, n2)
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed: Expected %v, got %v", expected2, result2)
	}

	// Test case 3: String with single segment, n > 0.
	str3 := String("hello")
	sep3 := String(",")
	n3 := Int(1)
	expected3 := Slice[String]{"hello"}
	result3 := str3.SplitN(sep3, n3)
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed: Expected %v, got %v", expected3, result3)
	}
}

func TestStringFields(t *testing.T) {
	// Test case 1: String with multiple words separated by whitespace.
	str1 := String("hello world how are you")
	expected1 := Slice[String]{"hello", "world", "how", "are", "you"}
	result1 := str1.Fields().Collect()
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed: Expected %v, got %v", expected1, result1)
	}

	// Test case 2: String with single word.
	str2 := String("hello")
	expected2 := Slice[String]{"hello"}
	result2 := str2.Fields().Collect()
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed: Expected %v, got %v", expected2, result2)
	}

	// Test case 3: Empty strin
	str3 := String("")
	expected3 := Slice[String]{}
	result3 := str3.Fields().Collect()
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed: Expected %v, got %v", expected3, result3)
	}

	// Test case 4: String with leading and trailing whitespace.
	str4 := String("   hello   world   ")
	expected4 := Slice[String]{"hello", "world"}
	result4 := str4.Fields().Collect()
	if !reflect.DeepEqual(result4, expected4) {
		t.Errorf("Test case 4 failed: Expected %v, got %v", expected4, result4)
	}
}

func TestStringCount(t *testing.T) {
	// Test case 1: Count occurrences of substring in a string with multiple occurrences.
	str1 := String("hello world hello hello")
	substr1 := String("hello")
	expected1 := Int(3)
	result1 := str1.Count(substr1)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %d, got %d", expected1, result1)
	}

	// Test case 2: Count occurrences of substring in a string with no occurrences.
	str2 := String("abcdefg")
	substr2 := String("xyz")
	expected2 := Int(0)
	result2 := str2.Count(substr2)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %d, got %d", expected2, result2)
	}

	// Test case 3: Count occurrences of substring in an empty strin
	str3 := String("")
	substr3 := String("hello")
	expected3 := Int(0)
	result3 := str3.Count(substr3)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: Expected %d, got %d", expected3, result3)
	}
}

func TestStringEqFold(t *testing.T) {
	// Test case 1: Strings are equal case-insensitively.
	str1 := String("Hello")
	str2 := String("hello")
	expected1 := true
	result1 := str1.EqFold(str2)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %t, got %t", expected1, result1)
	}

	// Test case 2: Strings are not equal case-insensitively.
	str3 := String("world")
	expected2 := false
	result2 := str1.EqFold(str3)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %t, got %t", expected2, result2)
	}

	// Test case 3: Empty strings.
	str4 := String("")
	str5 := String("")
	expected3 := true
	result3 := str4.EqFold(str5)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: Expected %t, got %t", expected3, result3)
	}
}

func TestStringLastIndex(t *testing.T) {
	// Test case 1: Substring is present in the strin
	str1 := String("hello world hello")
	substr1 := String("hello")
	expected1 := Int(12)
	result1 := str1.LastIndex(substr1)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %d, got %d", expected1, result1)
	}

	// Test case 2: Substring is not present in the strin
	substr2 := String("foo")
	expected2 := Int(-1)
	result2 := str1.LastIndex(substr2)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %d, got %d", expected2, result2)
	}

	// Test case 3: Empty strin
	str3 := String("")
	substr3 := String("hello")
	expected3 := Int(-1)
	result3 := str3.LastIndex(substr3)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: Expected %d, got %d", expected3, result3)
	}
}

func TestStringIndexRune(t *testing.T) {
	// Test case 1: Rune is present in the strin
	str1 := String("hello")
	rune1 := 'e'
	expected1 := Int(1)
	result1 := str1.IndexRune(rune1)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %d, got %d", expected1, result1)
	}

	// Test case 2: Rune is not present in the strin
	rune2 := 'x'
	expected2 := Int(-1)
	result2 := str1.IndexRune(rune2)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %d, got %d", expected2, result2)
	}

	// Test case 3: Empty strin
	str3 := String("")
	rune3 := 'h'
	expected3 := Int(-1)
	result3 := str3.IndexRune(rune3)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: Expected %d, got %d", expected3, result3)
	}
}

func TestStringNotEmpty(t *testing.T) {
	// Test case 1: String is not empty.
	str1 := String("hello")
	expected1 := true
	result1 := str1.NotEmpty()
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %t, got %t", expected1, result1)
	}

	// Test case 2: String is empty.
	str2 := String("")
	expected2 := false
	result2 := str2.NotEmpty()
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %t, got %t", expected2, result2)
	}
}

func TestRepeat(t *testing.T) {
	// Test case 1: Repeat count is positive.
	str := String("abc")
	count := Int(3)
	expected1 := String("abcabcabc")
	result1 := str.Repeat(count)
	if !result1.Eq(expected1) {
		t.Errorf("Test case 1 failed: Expected %s, got %s", expected1, result1)
	}

	// Test case 2: Repeat count is zero.
	count = 0
	expected2 := String("")
	result2 := str.Repeat(count)
	if !result2.Eq(expected2) {
		t.Errorf("Test case 2 failed: Expected %s, got %s", expected2, result2)
	}
}

func TestStringLeftJustify(t *testing.T) {
	// Test case 1: Original string length is less than the specified length.
	str1 := String("Hello")
	length1 := Int(10)
	pad1 := String(".")
	expected1 := String("Hello.....")
	result1 := str1.LeftJustify(length1, pad1)
	if !result1.Eq(expected1) {
		t.Errorf("Test case 1 failed: Expected %s, got %s", expected1, result1)
	}

	// Test case 2: Original string length is equal to the specified length.
	str2 := String("Hello")
	length2 := Int(5)
	pad2 := String(".")
	expected2 := String("Hello")
	result2 := str2.LeftJustify(length2, pad2)
	if !result2.Eq(expected2) {
		t.Errorf("Test case 2 failed: Expected %s, got %s", expected2, result2)
	}

	// Test case 3: Original string length is greater than the specified length.
	str3 := String("Hello")
	length3 := Int(3)
	pad3 := String(".")
	expected3 := String("Hello")
	result3 := str3.LeftJustify(length3, pad3)
	if !result3.Eq(expected3) {
		t.Errorf("Test case 3 failed: Expected %s, got %s", expected3, result3)
	}

	// Test case 4: Empty padding strin
	str4 := String("Hello")
	length4 := Int(10)
	pad4 := String("")
	expected4 := String("Hello")
	result4 := str4.LeftJustify(length4, pad4)
	if !result4.Eq(expected4) {
		t.Errorf("Test case 4 failed: Expected %s, got %s", expected4, result4)
	}
}

func TestStringRightJustify(t *testing.T) {
	// Test case 1: Original string length is less than the specified length.
	str1 := String("Hello")
	length1 := Int(10)
	pad1 := String(".")
	expected1 := String(".....Hello")
	result1 := str1.RightJustify(length1, pad1)
	if !result1.Eq(expected1) {
		t.Errorf("Test case 1 failed: Expected %s, got %s", expected1, result1)
	}

	// Test case 2: Original string length is equal to the specified length.
	str2 := String("Hello")
	length2 := Int(5)
	pad2 := String(".")
	expected2 := String("Hello")
	result2 := str2.RightJustify(length2, pad2)
	if !result2.Eq(expected2) {
		t.Errorf("Test case 2 failed: Expected %s, got %s", expected2, result2)
	}

	// Test case 3: Original string length is greater than the specified length.
	str3 := String("Hello")
	length3 := Int(3)
	pad3 := String(".")
	expected3 := String("Hello")
	result3 := str3.RightJustify(length3, pad3)
	if !result3.Eq(expected3) {
		t.Errorf("Test case 3 failed: Expected %s, got %s", expected3, result3)
	}

	// Test case 4: Empty padding strin
	str4 := String("Hello")
	length4 := Int(10)
	pad4 := String("")
	expected4 := String("Hello")
	result4 := str4.RightJustify(length4, pad4)
	if !result4.Eq(expected4) {
		t.Errorf("Test case 4 failed: Expected %s, got %s", expected4, result4)
	}
}

func TestStringCenter(t *testing.T) {
	// Test case 1: Original string length is less than the specified length.
	str1 := String("Hello")
	length1 := Int(10)
	pad1 := String(".")
	expected1 := String("..Hello...")
	result1 := str1.Center(length1, pad1)
	if !result1.Eq(expected1) {
		t.Errorf("Test case 1 failed: Expected %s, got %s", expected1, result1)
	}

	// Test case 2: Original string length is equal to the specified length.
	str2 := String("Hello")
	length2 := Int(5)
	pad2 := String(".")
	expected2 := String("Hello")
	result2 := str2.Center(length2, pad2)
	if !result2.Eq(expected2) {
		t.Errorf("Test case 2 failed: Expected %s, got %s", expected2, result2)
	}

	// Test case 3: Original string length is greater than the specified length.
	str3 := String("Hello")
	length3 := Int(3)
	pad3 := String(".")
	expected3 := String("Hello")
	result3 := str3.Center(length3, pad3)
	if !result3.Eq(expected3) {
		t.Errorf("Test case 3 failed: Expected %s, got %s", expected3, result3)
	}

	// Test case 4: Empty padding strin
	str4 := String("Hello")
	length4 := Int(10)
	pad4 := String("")
	expected4 := String("Hello")
	result4 := str4.Center(length4, pad4)
	if !result4.Eq(expected4) {
		t.Errorf("Test case 4 failed: Expected %s, got %s", expected4, result4)
	}
}

func TestStringEndsWithAny(t *testing.T) {
	// Test case 1: String ends with one of the provided suffixes.
	str1 := String("example.com")
	suffixes1 := Slice[String]{String(".com"), String(".net")}
	expected1 := true
	result1 := str1.EndsWithAny(suffixes1...)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: Expected %t, got %t", expected1, result1)
	}

	// Test case 2: String ends with multiple provided suffixes.
	str2 := String("example.net")
	suffixes2 := Slice[String]{String(".com"), String(".net")}
	expected2 := true
	result2 := str2.EndsWithAny(suffixes2...)
	if result2 != expected2 {
		t.Errorf("Test case 2 failed: Expected %t, got %t", expected2, result2)
	}

	// Test case 3: String does not end with any of the provided suffixes.
	str3 := String("example.org")
	suffixes3 := Slice[String]{String(".com"), String(".net")}
	expected3 := false
	result3 := str3.EndsWithAny(suffixes3...)
	if result3 != expected3 {
		t.Errorf("Test case 3 failed: Expected %t, got %t", expected3, result3)
	}
}

func TestStringStartsWithAny(t *testing.T) {
	// Test cases
	testCases := []struct {
		str      String
		prefixes []String
		expected bool
	}{
		{"http://example.com", []String{"http://", "https://"}, true},
		{"https://example.com", []String{"http://", "https://"}, true},
		{"ftp://example.com", []String{"http://", "https://"}, false},
		{"", []String{""}, true}, // Empty string should match empty prefix
		{"", []String{"non-empty"}, false},
	}

	// Test each case
	for _, tc := range testCases {
		// Wrap the input string
		s := String(tc.str)

		// Call the StartsWith method
		result := s.StartsWithAny(tc.prefixes...)

		// Assert the result
		if result != tc.expected {
			t.Errorf("StartsWith() returned %v; expected %v", result, tc.expected)
		}
	}
}

func TestStringSplitAfter(t *testing.T) {
	testCases := []struct {
		input     String
		separator String
		expected  Slice[String]
	}{
		{"hello,world,how,are,you", ",", Slice[String]{"hello,", "world,", "how,", "are,", "you"}},
		{"apple banana cherry", " ", Slice[String]{"apple ", "banana ", "cherry"}},
		{"a-b-c-d-e", "-", Slice[String]{"a-", "b-", "c-", "d-", "e"}},
		{"abcd", "a", Slice[String]{"a", "bcd"}},
		{"thisistest", "is", Slice[String]{"this", "is", "test"}},
		{"☺☻☹", "", Slice[String]{"☺", "☻", "☹"}},
		{"☺☻☹", "☹", Slice[String]{"☺☻☹", ""}},
		{"123", "", Slice[String]{"1", "2", "3"}},
	}

	for _, tc := range testCases {
		actual := tc.input.SplitAfter(tc.separator).Collect()

		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf(
				"Unexpected result for input: %s, separator: %s\nExpected: %v\nGot: %v",
				tc.input,
				tc.separator,
				tc.expected,
				actual,
			)
		}
	}
}

func TestStringFieldsBy(t *testing.T) {
	testCases := []struct {
		input    String
		fn       func(r rune) bool
		expected Slice[String]
	}{
		{"hello world", unicode.IsSpace, Slice[String]{"hello", "world"}},
		{"1,2,3,4,5", func(r rune) bool { return r == ',' }, Slice[String]{"1", "2", "3", "4", "5"}},
		{"camelCcase", unicode.IsUpper, Slice[String]{"camel", "case"}},
	}

	for _, tc := range testCases {
		actual := tc.input.FieldsBy(tc.fn).Collect()

		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("Unexpected result for input: %s\nExpected: %v\nGot: %v", tc.input, tc.expected, actual)
		}
	}
}

func TestStringRemove(t *testing.T) {
	tests := []struct {
		original String
		matches  Slice[String]
		expected String
	}{
		{
			original: "Hello, world! This is a test.",
			matches:  Slice[String]{"Hello", "test"},
			expected: ", world! This is a .",
		},
		{
			original: "This is a test string. This is a test string.",
			matches:  Slice[String]{"test", "string"},
			expected: "This is a  . This is a  .",
		},
		{
			original: "I love ice cream. Ice cream is delicious.",
			matches:  Slice[String]{"Ice", "cream"},
			expected: "I love ice .   is delicious.",
		},
	}

	for _, test := range tests {
		original := test.original
		modified := original.Remove(test.matches...)
		if modified != test.expected {
			t.Errorf("Remove(%q, %q) = %q, expected %q", test.original, test.matches, modified.Std(), test.expected)
		}
	}
}

func TestStringTrimSet(t *testing.T) {
	tests := []struct {
		input    String
		cutset   String
		expected String
	}{
		{"##Hello, world!##", "#", "Hello, world!"},
		{"**Magic**String**", "*", "Magic**String"},
		{"  trim spaces  ", " ", "trim spaces"},
		{"!@##@@!!SpecialChars!!@@##@!", "@#!", "SpecialChars"},
		{"NoChangeNeeded", "z", "NoChangeNeeded"},
	}

	for _, test := range tests {
		result := test.input.TrimSet(test.cutset)
		if result != test.expected {
			t.Errorf("TrimSet(%q, %q) = %q; want %q", test.input, test.cutset, result, test.expected)
		}
	}
}

func TestStringTrimStartSet(t *testing.T) {
	tests := []struct {
		input    String
		cutset   String
		expected String
	}{
		{"##Hello, world!##", "#", "Hello, world!##"},
		{"**Magic**String**", "*", "Magic**String**"},
		{"  trim spaces  ", " ", "trim spaces  "},
		{"!@##@@!!SpecialChars!!@@##@!", "@#!", "SpecialChars!!@@##@!"},
		{"NoChangeNeeded", "z", "NoChangeNeeded"},
	}

	for _, test := range tests {
		result := test.input.TrimStartSet(test.cutset)
		if result != test.expected {
			t.Errorf("TrimStartSet(%q, %q) = %q; want %q", test.input, test.cutset, result, test.expected)
		}
	}
}

func TestStringTrimEndSet(t *testing.T) {
	tests := []struct {
		input    String
		cutset   String
		expected String
	}{
		{"##Hello, world!##", "#", "##Hello, world!"},
		{"**Magic**String**", "*", "**Magic**String"},
		{"  trim spaces  ", " ", "  trim spaces"},
		{"!@##@@!!SpecialChars!!@@##@!", "@#!", "!@##@@!!SpecialChars"},
		{"NoChangeNeeded", "z", "NoChangeNeeded"},
	}

	for _, test := range tests {
		result := test.input.TrimEndSet(test.cutset)
		if result != test.expected {
			t.Errorf("TrimEndSet(%q, %q) = %q; want %q", test.input, test.cutset, result, test.expected)
		}
	}
}

func TestStringTrim(t *testing.T) {
	tests := []struct {
		input    String
		expected String
	}{
		{"  Hello, world!  ", "Hello, world!"},
		{"\t\tTabs\t\t", "Tabs"},
		{"\nNewLine\n", "NewLine"},
		{"NoTrimNeeded", "NoTrimNeeded"},
		{"   Multiple   Spaces   ", "Multiple   Spaces"},
	}

	for _, test := range tests {
		result := test.input.Trim()
		if result != test.expected {
			t.Errorf("Trim(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestStringTrimStart(t *testing.T) {
	tests := []struct {
		input    String
		expected String
	}{
		{"  Hello, world!  ", "Hello, world!  "},
		{"\t\tTabs\t\t", "Tabs\t\t"},
		{"\nNewLine\n", "NewLine\n"},
		{"NoTrimNeeded", "NoTrimNeeded"},
		{"   Multiple   Spaces   ", "Multiple   Spaces   "},
	}

	for _, test := range tests {
		result := test.input.TrimStart()
		if result != test.expected {
			t.Errorf("TrimStart(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestStringTrimEnd(t *testing.T) {
	tests := []struct {
		input    String
		expected String
	}{
		{"  Hello, world!  ", "  Hello, world!"},
		{"\t\tTabs\t\t", "\t\tTabs"},
		{"\nNewLine\n", "\nNewLine"},
		{"NoTrimNeeded", "NoTrimNeeded"},
		{"   Multiple   Spaces   ", "   Multiple   Spaces"},
	}

	for _, test := range tests {
		result := test.input.TrimEnd()
		if result != test.expected {
			t.Errorf("TrimEnd(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestStringToBigInt(t *testing.T) {
	tests := []struct {
		name     string
		input    String
		expected *big.Int
	}{
		{
			name:     "Decimal",
			input:    "12345",
			expected: big.NewInt(12345),
		},
		{
			name:     "Hexadecimal",
			input:    "0x1abc",
			expected: big.NewInt(6844),
		},
		{
			name:     "Octal",
			input:    "071",
			expected: big.NewInt(57),
		},
		{
			name:     "Invalid format",
			input:    "abc123",
			expected: nil,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToBigInt()
			if result.UnwrapOrDefault().Cmp(tc.expected) != 0 {
				t.Errorf("Failed %s: expected %v, got %v", tc.name, tc.expected, result)
			}
		})
	}
}

func TestStringTransform(t *testing.T) {
	original := String("hello world")
	expected := String("HELLO WORLD")
	result := original.Transform(String.Upper)

	if result != expected {
		t.Errorf("Transform failed: expected %q, got %q", expected, result)
	}
}
