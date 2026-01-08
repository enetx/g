package g_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"reflect"
	"regexp"
	"testing"
	"unicode"
	"unicode/utf8"

	. "github.com/enetx/g"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TestBytes(t *testing.T) {
	// Test case with string input
	strInput := "hello"
	bytesFromStr := Bytes(strInput)
	expectedBytesFromStr := Bytes("hello")
	if !bytes.Equal(bytesFromStr, expectedBytesFromStr) {
		t.Errorf("Conversion from string failed. Expected: %s, Got: %s", expectedBytesFromStr, bytesFromStr)
	}

	// Test case with []byte input
	byteSliceInput := []byte{104, 101, 108, 108, 111}
	bytesFromSlice := Bytes(byteSliceInput)
	expectedBytesFromSlice := Bytes("hello")
	if !bytes.Equal(bytesFromSlice, expectedBytesFromSlice) {
		t.Errorf("Conversion from byte slice failed. Expected: %s, Got: %s", expectedBytesFromSlice, bytesFromSlice)
	}
}

func TestBytesReplace(t *testing.T) {
	// Test case where old byte sequence exists and is replaced
	bs1 := Bytes("hello world")
	oldB1 := Bytes("world")
	newB1 := Bytes("gopher")
	replaced1 := bs1.Replace(oldB1, newB1, -1)
	expected1 := Bytes("hello gopher")
	if !bytes.Equal(replaced1, expected1) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected1, replaced1)
	}

	// Test case with multiple occurrences of old byte sequence
	bs2 := Bytes("hello world hello world")
	oldB2 := Bytes("world")
	newB2 := Bytes("gopher")
	replaced2 := bs2.Replace(oldB2, newB2, -1)
	expected2 := Bytes("hello gopher hello gopher")
	if !bytes.Equal(replaced2, expected2) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected2, replaced2)
	}

	// Test case with limited replacements
	bs3 := Bytes("hello world hello world")
	oldB3 := Bytes("world")
	newB3 := Bytes("gopher")
	replaced3 := bs3.Replace(oldB3, newB3, 1)
	expected3 := Bytes("hello gopher hello world")
	if !bytes.Equal(replaced3, expected3) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected3, replaced3)
	}

	// Test case where old byte sequence doesn't exist
	bs4 := Bytes("hello world")
	oldB4 := Bytes("gopher")
	newB4 := Bytes("earth")
	replaced4 := bs4.Replace(oldB4, newB4, -1)
	if !bytes.Equal(replaced4, bs4) {
		t.Errorf("Expected no change when old byte sequence doesn't exist. Got: %s", replaced4)
	}
}

func TestReplaceAll(t *testing.T) {
	// Test case where old byte sequence exists and is replaced
	bs1 := Bytes("hello world")
	oldB1 := Bytes("world")
	newB1 := Bytes("gopher")
	replaced1 := bs1.ReplaceAll(oldB1, newB1)
	expected1 := Bytes("hello gopher")
	if !bytes.Equal(replaced1, expected1) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected1, replaced1)
	}

	// Test case with multiple occurrences of old byte sequence
	bs2 := Bytes("hello world hello world")
	oldB2 := Bytes("world")
	newB2 := Bytes("gopher")
	replaced2 := bs2.ReplaceAll(oldB2, newB2)
	expected2 := Bytes("hello gopher hello gopher")
	if !bytes.Equal(replaced2, expected2) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected2, replaced2)
	}

	// Test case where old byte sequence doesn't exist
	bs3 := Bytes("hello world")
	oldB3 := Bytes("gopher")
	newB3 := Bytes("earth")
	replaced3 := bs3.ReplaceAll(oldB3, newB3)
	if !bytes.Equal(replaced3, bs3) {
		t.Errorf("Expected no change when old byte sequence doesn't exist. Got: %s", replaced3)
	}
}

func TestBytesRxReplace(t *testing.T) {
	// Test case where pattern matches and is replaced
	bs1 := Bytes("hello world hello world")
	pattern1 := regexp.MustCompile("world")
	newB1 := Bytes("gopher")
	replaced1 := bs1.Regexp().Replace(pattern1, newB1)
	expected1 := Bytes("hello gopher hello gopher")
	if !bytes.Equal(replaced1, expected1) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected1, replaced1)
	}

	// Test case where pattern matches and is replaced with capture group
	bs2 := Bytes("apple apple apple")
	pattern2 := regexp.MustCompile(`(\w+)`)
	newB2 := Bytes("${1}s")
	replaced2 := bs2.Regexp().Replace(pattern2, newB2)
	expected2 := Bytes("apples apples apples")
	if !bytes.Equal(replaced2, expected2) {
		t.Errorf("Replacement with capture group failed. Expected: %s, Got: %s", expected2, replaced2)
	}

	// Test case where pattern doesn't match
	bs3 := Bytes("hello world")
	pattern3 := regexp.MustCompile("gopher")
	newB3 := Bytes("earth")
	replaced3 := bs3.Regexp().Replace(pattern3, newB3)
	if !bytes.Equal(replaced3, bs3) {
		t.Errorf("Expected no change when pattern doesn't match. Got: %s", replaced3)
	}
}

func TestBytesRxFind(t *testing.T) {
	// Test case where pattern matches and is found
	bs1 := Bytes("hello world")
	pattern1 := regexp.MustCompile("world")
	found1 := bs1.Regexp().Find(pattern1)
	expected1 := Bytes("world")
	if found1.IsNone() {
		t.Errorf("Expected to find matching pattern, but found none")
	} else if !bytes.Equal(found1.Unwrap(), expected1) {
		t.Errorf("Found pattern does not match expected result. Expected: %s, Got: %s", expected1, found1.Unwrap())
	}

	// Test case where pattern doesn't match
	bs2 := Bytes("hello world")
	pattern2 := regexp.MustCompile("gopher")
	found2 := bs2.Regexp().Find(pattern2)
	if found2.IsSome() {
		t.Errorf("Expected not to find matching pattern, but found one")
	}
}

func TestBytesRxMatch(t *testing.T) {
	pattern := regexp.MustCompile(`\d+`)
	bs := Bytes("123abc456")
	if !bs.Regexp().Match(pattern) {
		t.Errorf("Expected match for pattern %v, but got none", pattern)
	}

	bs = Bytes("abc")
	if bs.Regexp().Match(pattern) {
		t.Errorf("Expected no match for pattern %v, but got one", pattern)
	}
}

func TestBytesRxMatchAny(t *testing.T) {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\d+`),
		regexp.MustCompile(`[a-z]+`),
	}

	bs := Bytes("123abc456")
	if !bs.Regexp().MatchAny(patterns...) {
		t.Errorf("Expected match for one of the patterns, but got none")
	}

	bs = Bytes("123")
	if !bs.Regexp().MatchAny(patterns...) {
		t.Errorf("Expected match for one of the patterns, but got none")
	}

	bs = Bytes("!@#")
	if bs.Regexp().MatchAny(patterns...) {
		t.Errorf("Expected no match for any of the patterns, but got one")
	}
}

func TestBytesRxMatchAll(t *testing.T) {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\d+`),
		regexp.MustCompile(`[a-z]+`),
	}

	bs := Bytes("123abc")
	if !bs.Regexp().MatchAll(patterns...) {
		t.Errorf("Expected match for all patterns, but got none")
	}

	bs = Bytes("123")
	if bs.Regexp().MatchAll(patterns...) {
		t.Errorf("Expected no match for all patterns, but got one")
	}

	bs = Bytes("abc")
	if bs.Regexp().MatchAll(patterns...) {
		t.Errorf("Expected no match for all patterns, but got one")
	}
}

func TestBytesStripPrefix(t *testing.T) {
	// Test case where cutset matches the prefix
	bs1 := Bytes("prefix_hello world")
	cutset1 := Bytes("prefix_")
	trimmed1 := bs1.StripPrefix(cutset1)
	expected1 := Bytes("hello world")
	if !bytes.Equal(trimmed1, expected1) {
		t.Errorf("Trimming prefix failed. Expected: %s, Got: %s", expected1, trimmed1)
	}

	// Test case where cutset doesn't match the prefix
	bs2 := Bytes("hello world")
	cutset2 := Bytes("nonexistent_")
	trimmed2 := bs2.StripPrefix(cutset2)
	if !bytes.Equal(trimmed2, bs2) {
		t.Errorf("Expected no change when cutset doesn't match the prefix. Got: %s", trimmed2)
	}
}

func TestBytesStripSuffix(t *testing.T) {
	// Test case where cutset matches the suffix
	bs1 := Bytes("hello world_suffix")
	cutset1 := Bytes("_suffix")
	trimmed1 := bs1.StripSuffix(cutset1)
	expected1 := Bytes("hello world")
	if !bytes.Equal(trimmed1, expected1) {
		t.Errorf("Trimming suffix failed. Expected: %s, Got: %s", expected1, trimmed1)
	}

	// Test case where cutset doesn't match the suffix
	bs2 := Bytes("hello world")
	cutset2 := Bytes("_nonexistent")
	trimmed2 := bs2.StripSuffix(cutset2)
	if !bytes.Equal(trimmed2, bs2) {
		t.Errorf("Expected no change when cutset doesn't match the suffix. Got: %s", trimmed2)
	}
}

func TestBytesSplit(t *testing.T) {
	// Test case where separator exists
	bs1 := Bytes("hello world gopher")
	separator1 := Bytes(" ")
	split1 := bs1.Split(separator1).Collect()
	expected1 := SliceOf(Bytes("hello"), Bytes("world"), Bytes("gopher"))
	if !reflect.DeepEqual(split1, expected1) {
		t.Errorf("Split failed. Expected: %v, Got: %v", expected1, split1)
	}

	// Test case where separator doesn't exist
	bs2 := Bytes("helloworldgopher")
	separator2 := Bytes(" ")
	split2 := bs2.Split(separator2).Collect()
	expected2 := Slice[Bytes]{Bytes("helloworldgopher")}
	if !reflect.DeepEqual(split2, expected2) {
		t.Errorf("Split failed. Expected: %v, Got: %v", expected2, split2)
	}
}

func TestBytesAppend(t *testing.T) {
	// Test case where bytes are added
	bs1 := Bytes("hello")
	obs1 := Bytes(" world")
	added1 := bs1.Append(obs1)
	expected1 := Bytes("hello world")
	if !bytes.Equal(added1, expected1) {
		t.Errorf("Add failed. Expected: %s, Got: %s", expected1, added1)
	}
}

func TestBytesPrepend(t *testing.T) {
	// Test case where bytes are added as a prefix
	bs1 := Bytes("world")
	obs1 := Bytes("hello ")
	prefixed1 := bs1.Prepend(obs1)
	expected1 := Bytes("hello world")
	if !bytes.Equal(prefixed1, expected1) {
		t.Errorf("AddPrefix failed. Expected: %s, Got: %s", expected1, prefixed1)
	}

	// Test case where prefix is empty
	bs2 := Bytes("world")
	obs2 := Bytes("")
	prefixed2 := bs2.Prepend(obs2)
	expected2 := Bytes("world")
	if !bytes.Equal(prefixed2, expected2) {
		t.Errorf("Prepend with empty prefix failed. Expected: %s, Got: %s", expected2, prefixed2)
	}

	// Test case where original bytes are empty
	bs3 := Bytes("")
	obs3 := Bytes("hello")
	prefixed3 := bs3.Prepend(obs3)
	expected3 := Bytes("hello")
	if !bytes.Equal(prefixed3, expected3) {
		t.Errorf("Prepend to empty bytes failed. Expected: %s, Got: %s", expected3, prefixed3)
	}

	// Test case where both are empty
	bs4 := Bytes("")
	obs4 := Bytes("")
	prefixed4 := bs4.Prepend(obs4)
	expected4 := Bytes("")
	if !bytes.Equal(prefixed4, expected4) {
		t.Errorf("Prepend empty to empty failed. Expected: %s, Got: %s", expected4, prefixed4)
	}
}

func TestBytesStd(t *testing.T) {
	// Test case where Bytes is converted to a byte slice
	bs1 := Bytes("hello world")
	std1 := bs1.Std()
	expected1 := []byte("hello world")
	if !bytes.Equal(std1, expected1) {
		t.Errorf("Std failed. Expected: %v, Got: %v", expected1, std1)
	}
}

func TestBytesClone(t *testing.T) {
	// Test case where Bytes is cloned
	bs1 := Bytes("hello world")
	cloned1 := bs1.Clone()
	if !bytes.Equal(cloned1, bs1) {
		t.Errorf("Clone failed. Expected: %s, Got: %s", bs1, cloned1)
	}
}

func TestBytesContainsAnyChars(t *testing.T) {
	// Test case where Bytes contains any characters from the input String
	bs1 := Bytes("hello")
	chars1 := String("aeiou")
	contains1 := bs1.ContainsAnyChars(chars1)
	if !contains1 {
		t.Errorf("ContainsAnyChars failed. Expected: true, Got: %t", contains1)
	}

	// Test case where Bytes doesn't contain any characters from the input String
	bs2 := Bytes("hello")
	chars2 := String("xyz")
	contains2 := bs2.ContainsAnyChars(chars2)
	if contains2 {
		t.Errorf("ContainsAnyChars failed. Expected: false, Got: %t", contains2)
	}
}

func TestBytesContainsRune(t *testing.T) {
	// Test case where Bytes contains the specified rune
	bs1 := Bytes("hello")
	rune1 := 'e'
	contains1 := bs1.ContainsRune(rune1)
	if !contains1 {
		t.Errorf("ContainsRune failed. Expected: true, Got: %t", contains1)
	}

	// Test case where Bytes doesn't contain the specified rune
	bs2 := Bytes("hello")
	rune2 := 'x'
	contains2 := bs2.ContainsRune(rune2)
	if contains2 {
		t.Errorf("ContainsRune failed. Expected: false, Got: %t", contains2)
	}
}

func TestBytesCount(t *testing.T) {
	// Test case where Bytes contains multiple occurrences of the specified Bytes
	bs1 := Bytes("hello hello hello")
	obs1 := Bytes("hello")
	count1 := bs1.Count(obs1)
	expected1 := Int(3)
	if count1 != expected1 {
		t.Errorf("Count failed. Expected: %d, Got: %d", expected1, count1)
	}

	// Test case where Bytes doesn't contain the specified Bytes
	bs2 := Bytes("hello")
	obs2 := Bytes("world")
	count2 := bs2.Count(obs2)
	expected2 := Int(0)
	if count2 != expected2 {
		t.Errorf("Count failed. Expected: %d, Got: %d", expected2, count2)
	}
}

func TestBytesCompare(t *testing.T) {
	testCases := []struct {
		bs1      Bytes
		bs2      Bytes
		expected Int
	}{
		{[]byte("apple"), []byte("banana"), -1},
		{[]byte("banana"), []byte("apple"), 1},
		{[]byte("banana"), []byte("banana"), 0},
		{[]byte("apple"), []byte("Apple"), 1},
		{[]byte(""), []byte(""), 0},
	}

	for _, tc := range testCases {
		result := Int(tc.bs1.Cmp(tc.bs2))
		if result != tc.expected {
			t.Errorf(
				"Bytes.Compare(%q, %q): expected %d, got %d",
				tc.bs1,
				tc.bs2,
				tc.expected,
				result,
			)
		}
	}
}

func TestBytesEqFold(t *testing.T) {
	// Test case where the byte slices are equal regardless of case
	bs1 := Bytes("Hello World")
	obs1 := Bytes("hello world")
	eqFold1 := bs1.EqFold(obs1)
	if !eqFold1 {
		t.Errorf("EqFold failed. Expected: true, Got: %t", eqFold1)
	}

	// Test case where the byte slices are not equal regardless of case
	bs2 := Bytes("Hello World")
	obs2 := Bytes("gopher")
	eqFold2 := bs2.EqFold(obs2)
	if eqFold2 {
		t.Errorf("EqFold failed. Expected: false, Got: %t", eqFold2)
	}
}

func TestBytesContains(t *testing.T) {
	testCases := []struct {
		bs       Bytes
		obs      Bytes
		expected bool
	}{
		{[]byte("hello world"), []byte("world"), true},
		{[]byte("hello world"), []byte("hello"), true},
		{[]byte("hello world"), []byte("gopher"), false},
		{[]byte("hello"), []byte("hello world"), false},
		{[]byte(""), []byte(""), true},
		{[]byte("test"), []byte(""), true},
		{[]byte(""), []byte("test"), false},
		{[]byte("abcdef"), []byte("cde"), true},
		{[]byte("abcdef"), []byte("xyz"), false},
	}

	for _, tc := range testCases {
		result := tc.bs.Contains(tc.obs)
		if result != tc.expected {
			t.Errorf(
				"Bytes.Contains(%q, %q): expected %t, got %t",
				tc.bs,
				tc.obs,
				tc.expected,
				result,
			)
		}
	}
}

func TestBytesEq(t *testing.T) {
	testCases := []struct {
		bs1      Bytes
		bs2      Bytes
		expected bool
	}{
		{[]byte("apple"), []byte("banana"), false},
		{[]byte("banana"), []byte("banana"), true},
		{[]byte("Apple"), []byte("apple"), false},
		{[]byte(""), []byte(""), true},
	}

	for _, tc := range testCases {
		result := tc.bs1.Eq(tc.bs2)
		if result != tc.expected {
			t.Errorf(
				"Bytes.Eq(%q, %q): expected %t, got %t",
				tc.bs1,
				tc.bs2,
				tc.expected,
				result,
			)
		}
	}
}

func TestBytesNe(t *testing.T) {
	testCases := []struct {
		bs1      Bytes
		bs2      Bytes
		expected bool
	}{
		{[]byte("apple"), []byte("banana"), true},
		{[]byte("banana"), []byte("banana"), false},
		{[]byte("Apple"), []byte("apple"), true},
		{[]byte(""), []byte(""), false},
	}

	for _, tc := range testCases {
		result := tc.bs1.Ne(tc.bs2)
		if result != tc.expected {
			t.Errorf(
				"Bytes.Ne(%q, %q): expected %t, got %t",
				tc.bs1,
				tc.bs2,
				tc.expected,
				result,
			)
		}
	}
}

func TestBytesGt(t *testing.T) {
	testCases := []struct {
		bs1      Bytes
		bs2      Bytes
		expected bool
	}{
		{[]byte("apple"), []byte("banana"), false},
		{[]byte("banana"), []byte("apple"), true},
		{[]byte("Apple"), []byte("apple"), false},
		{[]byte("banana"), []byte("banana"), false},
		{[]byte(""), []byte(""), false},
	}

	for _, tc := range testCases {
		result := tc.bs1.Gt(tc.bs2)
		if result != tc.expected {
			t.Errorf(
				"Bytes.Gt(%q, %q): expected %t, got %t",
				tc.bs1,
				tc.bs2,
				tc.expected,
				result,
			)
		}
	}
}

func TestBytesLt(t *testing.T) {
	testCases := []struct {
		bs1      Bytes
		bs2      Bytes
		expected bool
	}{
		{[]byte("apple"), []byte("banana"), true},
		{[]byte("banana"), []byte("apple"), false},
		{[]byte("Apple"), []byte("apple"), true},
		{[]byte("banana"), []byte("banana"), false},
		{[]byte(""), []byte(""), false},
	}

	for _, tc := range testCases {
		result := tc.bs1.Lt(tc.bs2)
		if result != tc.expected {
			t.Errorf(
				"Bytes.Lt(%q, %q): expected %t, got %t",
				tc.bs1,
				tc.bs2,
				tc.expected,
				result,
			)
		}
	}
}

func TestBytesNormalizeNFC(t *testing.T) {
	testCases := []struct {
		input    Bytes
		expected Bytes
	}{
		{[]byte("Mëtàl Hëàd"), []byte("Mëtàl Hëàd")},
		{[]byte("Café"), []byte("Café")},
		{[]byte("Ĵūņě"), []byte("Ĵūņě")},
		{[]byte("A\u0308"), []byte("Ä")},
		{[]byte("o\u0308"), []byte("ö")},
		{[]byte("u\u0308"), []byte("ü")},
		{[]byte("O\u0308"), []byte("Ö")},
		{[]byte("U\u0308"), []byte("Ü")},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			output := tc.input.NormalizeNFC()
			if string(output) != string(tc.expected) {
				t.Errorf("Bytes.NormalizeNFC(%q) = %q; want %q", tc.input, output, tc.expected)
			}
		})
	}
}

func TestBytesReader(t *testing.T) {
	tests := []struct {
		name     string
		bs       Bytes
		expected []byte
	}{
		{"Empty Bytes", Bytes{}, []byte{}},
		{"Single byte Bytes", Bytes{0x41}, []byte{0x41}},
		{
			"Multiple bytes Bytes",
			Bytes{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			[]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
		},
		{
			"Bytes with various values",
			Bytes{0x00, 0xff, 0x80, 0x7f},
			[]byte{0x00, 0xff, 0x80, 0x7f},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := test.bs.Reader()
			resultBytes, err := io.ReadAll(reader)
			if err != nil {
				t.Fatalf("Error reading from *bytes.Reader: %v", err)
			}

			if !bytes.Equal(resultBytes, test.expected) {
				t.Errorf("Bytes.Reader() content = %v, expected %v", resultBytes, test.expected)
			}
		})
	}
}

func TestBytesContainsAny(t *testing.T) {
	testCases := []struct {
		bs       Bytes
		bss      []Bytes
		expected bool
	}{
		{
			bs:       Bytes("Hello, world!"),
			bss:      []Bytes{Bytes("world"), Bytes("Go")},
			expected: true,
		},
		{
			bs:       Bytes("Welcome to the HumanGo-1!"),
			bss:      []Bytes{Bytes("Go-3"), Bytes("Go-4")},
			expected: false,
		},
		{
			bs:       Bytes("Have a great day!"),
			bss:      []Bytes{Bytes(""), Bytes(" ")},
			expected: true,
		},
		{
			bs:       Bytes(""),
			bss:      []Bytes{Bytes("Hello"), Bytes("world")},
			expected: false,
		},
		{
			bs:       Bytes(""),
			bss:      []Bytes{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		result := tc.bs.ContainsAny(tc.bss...)
		if result != tc.expected {
			t.Errorf(
				"Bytes.ContainsAny(%v, %v) = %v; want %v",
				tc.bs,
				tc.bss,
				result,
				tc.expected,
			)
		}
	}
}

func TestBytesContainsAll(t *testing.T) {
	testCases := []struct {
		bs       Bytes
		bss      []Bytes
		expected bool
	}{
		{
			bs:       Bytes("Hello, world!"),
			bss:      []Bytes{Bytes("Hello"), Bytes("world")},
			expected: true,
		},
		{
			bs:       Bytes("Welcome to the HumanGo-1!"),
			bss:      []Bytes{Bytes("Go-3"), Bytes("Go-4")},
			expected: false,
		},
		{
			bs:       Bytes("Have a great day!"),
			bss:      []Bytes{Bytes("Have"), Bytes("a")},
			expected: true,
		},
		{
			bs:       Bytes(""),
			bss:      []Bytes{Bytes("Hello"), Bytes("world")},
			expected: false,
		},
		{
			bs:       Bytes("Hello, world!"),
			bss:      []Bytes{},
			expected: true,
		},
	}

	for _, tc := range testCases {
		result := tc.bs.ContainsAll(tc.bss...)
		if result != tc.expected {
			t.Errorf(
				"Bytes.ContainsAll(%v, %v) = %v; want %v",
				tc.bs,
				tc.bss,
				result,
				tc.expected,
			)
		}
	}
}

func TestBytesIndex(t *testing.T) {
	// Test case where obs is present in bs
	bs := Bytes("hello world")
	obs := Bytes("world")
	idx := bs.Index(obs)
	expected := Int(6)
	if idx != expected {
		t.Errorf("Index failed. Expected: %d, Got: %d", expected, idx)
	}

	// Test case where obs is not present in bs
	bs = Bytes("hello world")
	obs = Bytes("gopher")
	idx = bs.Index(obs)
	expected = Int(-1)
	if idx != expected {
		t.Errorf("Index failed. Expected: %d, Got: %d", expected, idx)
	}
}

func TestBytesRxIndex(t *testing.T) {
	// Test case where a match is found
	bs := Bytes("apple banana")
	pattern := regexp.MustCompile(`banana`)
	idx := bs.Regexp().Index(pattern)
	expected := Some(Slice[Int]{6, 12})
	if idx.IsNone() || !reflect.DeepEqual(idx.Some(), expected.Some()) {
		t.Errorf("IndexRegexp failed. Expected: %v, Got: %v", expected, idx)
	}

	// Test case where no match is found
	bs = Bytes("apple banana")
	pattern = regexp.MustCompile(`orange`)
	idx = bs.Regexp().Index(pattern)
	expected = None[Slice[Int]]()
	if idx.IsSome() || !reflect.DeepEqual(idx.IsNone(), expected.IsNone()) {
		t.Errorf("IndexRegexp failed. Expected: %v, Got: %v", expected, idx)
	}
}

func TestBytesRepeat(t *testing.T) {
	// Test case where the Bytes are repeated 3 times
	bs := Bytes("hello")
	repeated := bs.Repeat(3)
	expected := Bytes("hellohellohello")
	if !bytes.Equal(repeated, expected) {
		t.Errorf("Repeat failed. Expected: %s, Got: %s", expected, repeated)
	}

	// Test case where the Bytes are repeated 0 times
	bs = Bytes("hello")
	repeated = bs.Repeat(0)
	expected = Bytes("")
	if !bytes.Equal(repeated, expected) {
		t.Errorf("Repeat failed. Expected: %s, Got: %s", expected, repeated)
	}
}

func TestToRunes(t *testing.T) {
	// Test case where the Bytes are converted to runes
	bs := Bytes("hello")
	runes := bs.Runes()
	expected := []rune{'h', 'e', 'l', 'l', 'o'}
	if !reflect.DeepEqual(runes, expected) {
		t.Errorf("ToRunes failed. Expected: %v, Got: %v", expected, runes)
	}
}

func TestBytesLower(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		expected Bytes
	}{
		{"Empty", Bytes(""), Bytes("")},
		{"AlreadyLower", Bytes("hello world"), Bytes("hello world")},
		{"ASCII Mixed", Bytes("Hello WORLD"), Bytes("hello world")},
		{"DigitsAndLetters", Bytes("ABC123xyz"), Bytes("abc123xyz")},
		{"Punctuation", Bytes("Hello-World!"), Bytes("hello-world!")},
		{"Cyrillic", Bytes("ПрИвЕт Мир"), cases.Lower(language.English).Bytes([]byte("ПрИвЕт Мир"))},
		{"Chinese", Bytes("你好世界"), cases.Lower(language.English).Bytes([]byte("你好世界"))},
		{"Emoji", Bytes("Go🚀Lang"), cases.Lower(language.English).Bytes([]byte("Go🚀Lang"))},
		{"MixedSeparator", Bytes("foo_bar"), Bytes("foo_bar")},
		{"InvalidUTF8", Bytes([]byte{0xff, 0xfe, 0xfd}), Bytes([]byte{0xff, 0xfe, 0xfd})},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Lower()
			if !bytes.Equal([]byte(got), []byte(tc.expected)) {
				t.Errorf("%s: Lower(%q) = %q; want %q",
					tc.name, []byte(tc.input), []byte(got), []byte(tc.expected))
			}

			if utf8.Valid([]byte(tc.input)) && !utf8.Valid([]byte(got)) {
				t.Errorf("%s: result is invalid UTF-8: %x", tc.name, []byte(got))
			}
		})
	}
}

func TestBytesUpper(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		expected Bytes
	}{
		{"Empty", Bytes(""), Bytes("")},
		{"AlreadyUpper", Bytes("HELLO WORLD"), Bytes("HELLO WORLD")},
		{"ASCII Mixed", Bytes("Hello world"), Bytes("HELLO WORLD")},
		{"DigitsAndLetters", Bytes("abc123XYZ"), Bytes("ABC123XYZ")},
		{"Punctuation", Bytes("foo-bar!"), Bytes("FOO-BAR!")},
		{"Cyrillic", Bytes("ПрИвЕт Мир"), cases.Upper(language.English).Bytes([]byte("ПрИвЕт Мир"))},
		{"Chinese", Bytes("你好世界"), cases.Upper(language.English).Bytes([]byte("你好世界"))},
		{"Emoji", Bytes("Go🚀Lang"), cases.Upper(language.English).Bytes([]byte("Go🚀Lang"))},
		{"MixedSeparator", Bytes("foo_bar"), Bytes("FOO_BAR")},
		{"InvalidUTF8", Bytes([]byte{0xff, 0xfe, 0xfd}), Bytes([]byte{0xff, 0xfe, 0xfd})},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Upper()
			if !bytes.Equal([]byte(got), []byte(tc.expected)) {
				t.Errorf("%s: Upper(%q) = %q; want %q",
					tc.name, []byte(tc.input), []byte(got), []byte(tc.expected))
			}

			if utf8.Valid([]byte(tc.input)) && !utf8.Valid([]byte(got)) {
				t.Errorf("%s: result is invalid UTF-8: %x", tc.name, []byte(got))
			}
		})
	}
}

func TestBytesTrimSpace(t *testing.T) {
	// Test case where white space characters are trimmed from the beginning and end
	bs := Bytes("  hello world  ")
	trimmed := bs.Trim()
	expected := Bytes("hello world")
	if !bytes.Equal(trimmed, expected) {
		t.Errorf("TrimSpace failed. Expected: %s, Got: %s", expected, trimmed)
	}

	// Test case where there are no white space characters
	bs = Bytes("hello world")
	trimmed = bs.Trim()
	expected = Bytes("hello world")
	if !bytes.Equal(trimmed, expected) {
		t.Errorf("TrimSpace failed. Expected: %s, Got: %s", expected, trimmed)
	}

	// Test case where the Bytes is empty
	bs = Bytes("")
	trimmed = bs.Trim()
	expected = Bytes("")
	if !bytes.Equal(trimmed, expected) {
		t.Errorf("TrimSpace failed. Expected: %s, Got: %s", expected, trimmed)
	}
}

func TestBytesTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    Bytes
		expected Bytes
	}{
		{"Empty", Bytes(""), Bytes("")},
		{"SingleWordLower", Bytes("hello"), Bytes("Hello")},
		{"SingleWordUpper", Bytes("HELLO"), Bytes("Hello")},
		{"Sentence", Bytes("hello world"), Bytes("Hello World")},
		{"AlreadyTitle", Bytes("Hello World"), Bytes("Hello World")},
		{"LeadingTrailingSpaces", Bytes("  hello world  "), Bytes("  Hello World  ")},
		{"TabsNewline", Bytes("foo\tbar\nbaz"), Bytes("Foo\tBar\nBaz")},
		{"Punctuation", Bytes("hello-world"), Bytes("Hello-World")},
		{"NumbersAndWords", Bytes("123abc 456def"), Bytes("123Abc 456Def")},
		{"NumbersThenWord", Bytes("123abc abc"), Bytes("123Abc Abc")},
		{"Chinese", Bytes("你好 世界"), Bytes("你好 世界")},            // no casing in CJK
		{"Cyrillic", Bytes("привет мир"), Bytes("Привет Мир")}, // uses Unicode fallback
		{"MixedASCIIUnicode", Bytes("Go 🚀 Language"), cases.Title(language.English).Bytes([]byte("Go 🚀 Language"))},
		{"Emoji", Bytes("😊 emoji 😊"), cases.Title(language.English).Bytes([]byte("😊 emoji 😊"))},
		{"MixedSeparator", Bytes("rock-n-roll"), Bytes("Rock-N-Roll")},
		{"MultipleSpaces", Bytes("  multiple   spaces "), Bytes("  Multiple   Spaces ")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.Title()
			if !bytes.Equal([]byte(got), []byte(tc.expected)) {
				t.Errorf("%s: Title(%q) = %q; want %q",
					tc.name,
					[]byte(tc.input),
					[]byte(got),
					[]byte(tc.expected),
				)
			}

			if !utf8.Valid([]byte(got)) {
				t.Errorf("%s: result is invalid UTF-8: %x", tc.name, []byte(got))
			}
		})
	}
}

func TestBytesNotEmpty(t *testing.T) {
	// Test case where the Bytes is not empty
	bs := Bytes("hello")
	if !bs.NotEmpty() {
		t.Errorf("NotEmpty failed. Expected: true, Got: false")
	}

	// Test case where the Bytes is empty
	bs = Bytes("")
	if bs.NotEmpty() {
		t.Errorf("NotEmpty failed. Expected: false, Got: true")
	}
}

func TestBytesMap(t *testing.T) {
	// Test case where the function converts each rune to uppercase
	bs := Bytes("hello")
	uppercase := bs.Map(func(r rune) rune {
		return unicode.ToUpper(r)
	})

	expected := Bytes("HELLO")
	if !bytes.Equal(uppercase, expected) {
		t.Errorf("Map failed. Expected: %s, Got: %s", expected, uppercase)
	}

	// Test case where the function removes spaces
	bs = Bytes("hello world")
	noSpaces := bs.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1 // Remove rune
		}
		return r
	})

	expected = Bytes("helloworld")
	if !bytes.Equal(noSpaces, expected) {
		t.Errorf("Map failed. Expected: %s, Got: %s", expected, noSpaces)
	}
}

func TestBytesLenRunes(t *testing.T) {
	// Test case where the Bytes contain ASCII characters
	bs := Bytes("hello world")
	lenRunes := bs.LenRunes()
	expected := Int(11)
	if lenRunes != expected {
		t.Errorf("LenRunes failed. Expected: %d, Got: %d", expected, lenRunes)
	}

	// Test case where the Bytes contain Unicode characters
	bs = Bytes("你好，世界")
	lenRunes = bs.LenRunes()
	expected = Int(5)
	if lenRunes != expected {
		t.Errorf("LenRunes failed. Expected: %d, Got: %d", expected, lenRunes)
	}

	// Test case where the Bytes are empty
	bs = Bytes("")
	lenRunes = bs.LenRunes()
	expected = Int(0)
	if lenRunes != expected {
		t.Errorf("LenRunes failed. Expected: %d, Got: %d", expected, lenRunes)
	}
}

func TestBytesLastIndex(t *testing.T) {
	// Test case where obs is present in bs
	bs := Bytes("hello world")
	obs := Bytes("o")
	lastIndex := bs.LastIndex(obs)
	expected := Int(7)
	if lastIndex != expected {
		t.Errorf("LastIndex failed. Expected: %d, Got: %d", expected, lastIndex)
	}

	// Test case where obs is not present in bs
	bs = Bytes("hello world")
	obs = Bytes("z")
	lastIndex = bs.LastIndex(obs)
	expected = Int(-1)
	if lastIndex != expected {
		t.Errorf("LastIndex failed. Expected: %d, Got: %d", expected, lastIndex)
	}
}

func TestBytesIndexByte(t *testing.T) {
	// Test case where b is present in bs
	bs := Bytes("hello world")
	b := byte('o')
	indexByte := bs.IndexByte(b)
	expected := Int(4)
	if indexByte != expected {
		t.Errorf("IndexByte failed. Expected: %d, Got: %d", expected, indexByte)
	}

	// Test case where b is not present in bs
	bs = Bytes("hello world")
	b = byte('z')
	indexByte = bs.IndexByte(b)
	expected = -1
	if indexByte != expected {
		t.Errorf("IndexByte failed. Expected: %d, Got: %d", expected, indexByte)
	}
}

func TestBytesLastIndexByte(t *testing.T) {
	// Test case where b is present in bs
	bs := Bytes("hello world")
	b := byte('o')
	lastIndexByte := bs.LastIndexByte(b)
	expected := Int(7)
	if lastIndexByte != expected {
		t.Errorf("LastIndexByte failed. Expected: %d, Got: %d", expected, lastIndexByte)
	}

	// Test case where b is not present in bs
	bs = Bytes("hello world")
	b = byte('z')
	lastIndexByte = bs.LastIndexByte(b)
	expected = -1
	if lastIndexByte != expected {
		t.Errorf("LastIndexByte failed. Expected: %d, Got: %d", expected, lastIndexByte)
	}
}

func TestBytesIndexRune(t *testing.T) {
	// Test case where r is present in bs
	bs := Bytes("hello world")
	r := 'o'
	indexRune := bs.IndexRune(r)
	expected := Int(4)
	if indexRune != expected {
		t.Errorf("IndexRune failed. Expected: %d, Got: %d", expected, indexRune)
	}

	// Test case where r is not present in bs
	bs = Bytes("hello world")
	r = 'z'
	indexRune = bs.IndexRune(r)
	expected = -1
	if indexRune != expected {
		t.Errorf("IndexRune failed. Expected: %d, Got: %d", expected, indexRune)
	}
}

func TestBytesRxFindAllSubmatchN(t *testing.T) {
	// Test case where matches are found
	bs := Bytes("hello world")
	pattern := regexp.MustCompile(`\b\w+\b`)
	matches := bs.Regexp().FindAllSubmatchN(pattern, -1)
	if matches.IsSome() {
		expected := Slice[Slice[Bytes]]{
			{Bytes("hello")},
			{Bytes("world")},
		}
		if !matches.Some().Eq(expected) {
			t.Errorf("FindAllSubmatchRegexpN failed. Expected: %s, Got: %s", expected, matches.Some())
		}
	} else {
		t.Errorf("FindAllSubmatchRegexpN failed. Expected matches, Got None")
	}

	// Test case where no matches are found
	bs = Bytes("")
	pattern = regexp.MustCompile(`\b\w+\b`)
	matches = bs.Regexp().FindAllSubmatchN(pattern, -1)
	if matches.IsSome() {
		t.Errorf("FindAllSubmatchRegexpN failed. Expected None, Got matches")
	}
}

func TestBytesRxFindAll(t *testing.T) {
	// Test case where matches are found
	bs := Bytes("hello world")
	pattern := regexp.MustCompile(`\b\w+\b`)
	matches := bs.Regexp().FindAll(pattern)
	if matches.IsSome() {
		expected := Slice[Bytes]{
			Bytes("hello"),
			Bytes("world"),
		}
		if !matches.Some().Eq(expected) {
			t.Errorf("FindAllRegexp failed. Expected: %s, Got: %s", expected, matches.Some())
		}
	} else {
		t.Errorf("FindAllRegexp failed. Expected matches, Got None")
	}

	// Test case where no matches are found
	bs = Bytes("")
	pattern = regexp.MustCompile(`\b\w+\b`)
	matches = bs.Regexp().FindAll(pattern)
	if matches.IsSome() {
		t.Errorf("FindAllRegexp failed. Expected None, Got matches")
	}
}

func TestBytesRxFindSubmatch(t *testing.T) {
	// Test case where a match is found
	bs := Bytes("hello world")
	pattern := regexp.MustCompile(`\b\w+\b`)
	match := bs.Regexp().FindSubmatch(pattern)
	if match.IsSome() {
		expected := SliceOf(Bytes("hello"))
		if !match.Some().Eq(expected) {
			t.Errorf("FindSubmatchRegexp failed. Expected: %s, Got: %s", expected, match.Some())
		}
	} else {
		t.Errorf("FindSubmatchRegexp failed. Expected match, Got None")
	}

	// Test case where no match is found
	bs = Bytes("")
	pattern = regexp.MustCompile(`\b\w+\b`)
	match = bs.Regexp().FindSubmatch(pattern)
	if match.IsSome() {
		t.Errorf("FindSubmatchRegexp failed. Expected None, Got match")
	}
}

func TestBytesRxFindAllSubmatch(t *testing.T) {
	// Test case where matches are found
	bs := Bytes("hello world")
	pattern := regexp.MustCompile(`\b\w+\b`)
	matches := bs.Regexp().FindAllSubmatch(pattern)
	if matches.IsSome() {
		expected := Slice[Slice[Bytes]]{
			{Bytes("hello")},
			{Bytes("world")},
		}
		if !matches.Some().Eq(expected) {
			t.Errorf("FindAllSubmatchRegexp failed. Expected: %s, Got: %s", expected, matches.Some())
		}
	} else {
		t.Errorf("FindAllSubmatchRegexp failed. Expected matches, Got None")
	}

	// Test case where no matches are found
	bs = Bytes("")
	pattern = regexp.MustCompile(`\b\w+\b`)
	matches = bs.Regexp().FindAllSubmatch(pattern)
	if matches.IsSome() {
		t.Errorf("FindAllSubmatchRegexp failed. Expected None, Got matches")
	}
}

func TestBytesHashingFunctions(t *testing.T) {
	// Test case for MD5 hashing
	input := Bytes("hello world")
	expectedMD5 := Bytes("5eb63bbbe01eeed093cb22bb8f5acdc3")
	md5Hash := input.Hash().MD5()
	if md5Hash.Ne(expectedMD5) {
		t.Errorf("MD5 hashing failed. Expected: %s, Got: %s", expectedMD5, md5Hash)
	}

	// Test case for SHA1 hashing
	expectedSHA1 := Bytes("2aae6c35c94fcfb415dbe95f408b9ce91ee846ed")
	sha1Hash := input.Hash().SHA1()
	if sha1Hash.Ne(expectedSHA1) {
		t.Errorf("SHA1 hashing failed. Expected: %s, Got: %s", expectedSHA1, sha1Hash)
	}

	// Test case for SHA256 hashing
	expectedSHA256 := Bytes("b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9")
	sha256Hash := input.Hash().SHA256()
	if sha256Hash.Ne(expectedSHA256) {
		t.Errorf("SHA256 hashing failed. Expected: %s, Got: %s", expectedSHA256, sha256Hash)
	}

	// Test case for SHA512 hashing
	expectedSHA512 := Bytes(
		"309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f",
	)
	sha512Hash := input.Hash().SHA512()
	if sha512Hash.Ne(expectedSHA512) {
		t.Errorf("SHA512 hashing failed. Expected: %s, Got: %s", expectedSHA512, sha512Hash)
	}
}

func TestBytesSplitAfter(t *testing.T) {
	testCases := []struct {
		input     Bytes
		separator Bytes
		expected  Slice[Bytes]
	}{
		{
			Bytes("hello,world,how,are,you"),
			Bytes(","),
			Slice[Bytes]{Bytes("hello,"), Bytes("world,"), Bytes("how,"), Bytes("are,"), Bytes("you")},
		},
		{
			Bytes("apple banana cherry"),
			Bytes(" "),
			Slice[Bytes]{Bytes("apple "), Bytes("banana "), Bytes("cherry")},
		},

		{
			Bytes("a-b-c-d-e"),
			Bytes("-"),
			Slice[Bytes]{Bytes("a-"), Bytes("b-"), Bytes("c-"), Bytes("d-"), Bytes("e")},
		},
		{Bytes("abcd"), Bytes("a"), Slice[Bytes]{Bytes("a"), Bytes("bcd")}},
		{Bytes("thisistest"), Bytes("is"), Slice[Bytes]{Bytes("this"), Bytes("is"), Bytes("test")}},
		{Bytes("☺☻☹"), Bytes(""), Slice[Bytes]{Bytes("☺"), Bytes("☻"), Bytes("☹")}},
		{Bytes("☺☻☹"), Bytes("☹"), Slice[Bytes]{Bytes("☺☻☹"), Bytes("")}},
		{Bytes("123"), Bytes(""), Slice[Bytes]{Bytes("1"), Bytes("2"), Bytes("3")}},
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

func TestBytesFields(t *testing.T) {
	bs1 := Bytes("hello world how are you")
	expected1 := Slice[Bytes]{Bytes("hello"), Bytes("world"), Bytes("how"), Bytes("are"), Bytes("you")}
	result1 := bs1.Fields().Collect()
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed: Expected %v, got %v", expected1, result1)
	}

	bs2 := Bytes("hello")
	expected2 := Slice[Bytes]{Bytes("hello")}
	result2 := bs2.Fields().Collect()
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed: Expected %v, got %v", expected2, result2)
	}

	bs3 := Bytes("")
	expected3 := Slice[Bytes]{}
	result3 := bs3.Fields().Collect()
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed: Expected %v, got %v", expected3, result3)
	}

	bs4 := Bytes("   hello   world   ")
	expected4 := Slice[Bytes]{Bytes("hello"), Bytes("world")}
	result4 := bs4.Fields().Collect()
	if !reflect.DeepEqual(result4, expected4) {
		t.Errorf("Test case 4 failed: Expected %v, got %v", expected4, result4)
	}
}

func TestBytesFieldsBy(t *testing.T) {
	testCases := []struct {
		input    Bytes
		fn       func(r rune) bool
		expected Slice[Bytes]
	}{
		{Bytes("hello world"), unicode.IsSpace, Slice[Bytes]{Bytes("hello"), Bytes("world")}},
		{
			Bytes("1,2,3,4,5"),
			func(r rune) bool { return r == ',' },
			Slice[Bytes]{Bytes("1"), Bytes("2"), Bytes("3"), Bytes("4"), Bytes("5")},
		},
		{Bytes("camelCcase"), unicode.IsUpper, Slice[Bytes]{Bytes("camel"), Bytes("case")}},
	}

	for _, tc := range testCases {
		actual := tc.input.FieldsBy(tc.fn).Collect()

		if !reflect.DeepEqual(actual, tc.expected) {
			t.Errorf("Unexpected result for input: %s\nExpected: %v\nGot: %v", tc.input, tc.expected, actual)
		}
	}
}

func TestBytesTrimSet(t *testing.T) {
	tests := []struct {
		input    Bytes
		cutset   String
		expected Bytes
	}{
		{[]byte("##Hello, world!##"), "#", []byte("Hello, world!")},
		{[]byte("!!!Amazing!!!"), "!", []byte("Amazing")},
		{[]byte("Spaces    "), " ", []byte("Spaces")},
		{[]byte("--Dashes--"), "-", []byte("Dashes")},
		{[]byte("NoTrimNeeded"), "x", []byte("NoTrimNeeded")},
		{[]byte("123Numbers123"), "123", []byte("Numbers")},
	}

	for _, test := range tests {
		result := test.input.TrimSet(test.cutset)
		if !bytes.Equal(result, test.expected) {
			t.Errorf("TrimSet(%q, %v) = %q; want %q", test.input, test.cutset, result, test.expected)
		}
	}
}

func TestBytesTrimStartSet(t *testing.T) {
	tests := []struct {
		input    Bytes
		cutset   String
		expected Bytes
	}{
		{[]byte("##Hello, world!##"), "#", []byte("Hello, world!##")},
		{[]byte("!!!Amazing!!!"), "!", []byte("Amazing!!!")},
		{[]byte("Spaces    "), " ", []byte("Spaces    ")},
		{[]byte("--Dashes--"), "-", []byte("Dashes--")},
		{[]byte("NoTrimNeeded"), "x", []byte("NoTrimNeeded")},
		{[]byte("123Numbers123"), "123", []byte("Numbers123")},
	}

	for _, test := range tests {
		result := test.input.TrimStartSet(test.cutset)
		if !bytes.Equal(result, test.expected) {
			t.Errorf("TrimStartSet(%q, %v) = %q; want %q", test.input, test.cutset, result, test.expected)
		}
	}
}

func TestBytesTrimEndSet(t *testing.T) {
	tests := []struct {
		input    Bytes
		cutset   String
		expected Bytes
	}{
		{[]byte("##Hello, world!##"), "#", []byte("##Hello, world!")},
		{[]byte("!!!Amazing!!!"), "!", []byte("!!!Amazing")},
		{[]byte("Spaces    "), " ", []byte("Spaces")},
		{[]byte("--Dashes--"), "-", []byte("--Dashes")},
		{[]byte("NoTrimNeeded"), "x", []byte("NoTrimNeeded")},
		{[]byte("123Numbers123"), "123", []byte("123Numbers")},
	}

	for _, test := range tests {
		result := test.input.TrimEndSet(test.cutset)
		if !bytes.Equal(result, test.expected) {
			t.Errorf("TrimEndSet(%q, %v) = %q; want %q", test.input, test.cutset, result, test.expected)
		}
	}
}

func TestBytesTrim(t *testing.T) {
	tests := []struct {
		input    Bytes
		expected Bytes
	}{
		{[]byte("  Hello, world!  "), []byte("Hello, world!")},
		{[]byte("\t\tTabbed\t\t"), []byte("Tabbed")},
		{[]byte("\nNewLine\n"), []byte("NewLine")},
		{[]byte("NoTrimNeeded"), []byte("NoTrimNeeded")},
		{[]byte("   Multiple   Spaces   "), []byte("Multiple   Spaces")},
		{[]byte(" \t \n Mixed \r\n "), []byte("Mixed")},
	}

	for _, test := range tests {
		result := test.input.Trim()
		if !bytes.Equal(result, test.expected) {
			t.Errorf("Trim(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestBytesTrimStart(t *testing.T) {
	tests := []struct {
		input    Bytes
		expected Bytes
	}{
		{[]byte("  Hello, world!  "), []byte("Hello, world!  ")},
		{[]byte("\t\tTabbed\t\t"), []byte("Tabbed\t\t")},
		{[]byte("\nNewLine\n"), []byte("NewLine\n")},
		{[]byte("NoTrimNeeded"), []byte("NoTrimNeeded")},
		{[]byte("   Multiple   Spaces   "), []byte("Multiple   Spaces   ")},
		{[]byte(" \t \n Mixed \r\n "), []byte("Mixed \r\n ")},
	}

	for _, test := range tests {
		result := test.input.TrimStart()
		if !bytes.Equal(result, test.expected) {
			t.Errorf("TrimStart(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestBytesTrimEnd(t *testing.T) {
	tests := []struct {
		input    Bytes
		expected Bytes
	}{
		{[]byte("  Hello, world!  "), []byte("  Hello, world!")},
		{[]byte("\t\tTabbed\t\t"), []byte("\t\tTabbed")},
		{[]byte("\nNewLine\n"), []byte("\nNewLine")},
		{[]byte("NoTrimNeeded"), []byte("NoTrimNeeded")},
		{[]byte("   Multiple   Spaces   "), []byte("   Multiple   Spaces")},
		{[]byte(" \t \n Mixed \r\n "), []byte(" \t \n Mixed")},
	}

	for _, test := range tests {
		result := test.input.TrimEnd()
		if !bytes.Equal(result, test.expected) {
			t.Errorf("TrimEnd(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}

func TestBytesTransform(t *testing.T) {
	original := Bytes("hello world")
	expected := Bytes("HELLO WORLD")
	result := original.Transform(Bytes.Upper)

	if !bytes.Equal(result, expected) {
		t.Errorf("Transform failed: expected %q, got %q", expected, result)
	}
}

var reverseTests = []struct {
	name  string
	input Bytes
	want  Bytes
}{
	{"Empty", Bytes(""), Bytes("")},
	{"Single ASCII", Bytes("A"), Bytes("A")},
	{"ASCII", Bytes("ABCdef"), Bytes("fedCBA")},
	{"Single Unicode", Bytes("Ж"), Bytes("Ж")},
	{"Unicode", Bytes("Привет"), Bytes("тевирП")},
	{"Chinese", Bytes("你好世界"), Bytes("界世好你")},
	{"Hindi", Bytes("नमस्ते"), Bytes("ेत्समन")},
	{"Mixed ASCII+Unicode", Bytes("Go🚀Lang"), Bytes("gnaL🚀oG")},
	{"Family Emoji", Bytes("👨‍👩‍👧‍👦"), Bytes("👦‍👧‍👩‍👨")},
	{"Combining", Bytes("é"), Bytes("é")},
	{"Emoji Sequence", Bytes("🙂🙃🙂"), Bytes("🙂🙃🙂")},
	{"Raw Bytes", Bytes([]byte{0, 1, 2, 3}), Bytes([]byte{3, 2, 1, 0})},
	{"Variation Selector", Bytes("✈️"), Bytes("️✈")},
	{"Invalid UTF-8 Bytes", Bytes([]byte{0xff, 0xfe, 0xfd}), Bytes([]byte{0xfd, 0xfe, 0xff})},
}

func TestBytesReverse(t *testing.T) {
	for _, tt := range reverseTests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Reverse()
			if string(got) != string(tt.want) {
				t.Errorf("Reverse(%q) = %q; want %q", string(tt.input), string(got), string(tt.want))
			}
		})
	}
}

func TestNewBytes(t *testing.T) {
	// Test NewBytes with no arguments - should create empty bytes
	bs1 := NewBytes()
	if len(bs1) != 0 {
		t.Errorf("NewBytes() should create empty bytes, got length %d", len(bs1))
	}

	// Test NewBytes with length only
	bs2 := NewBytes(5)
	if len(bs2) != 5 {
		t.Errorf("NewBytes(5) should create bytes with length 5, got %d", len(bs2))
	}
	if cap(bs2) != 5 {
		t.Errorf("NewBytes(5) should create bytes with capacity 5, got %d", cap(bs2))
	}

	// Test NewBytes with length and capacity
	bs3 := NewBytes(3, 10)
	if len(bs3) != 3 {
		t.Errorf("NewBytes(3, 10) should create bytes with length 3, got %d", len(bs3))
	}
	if cap(bs3) != 10 {
		t.Errorf("NewBytes(3, 10) should create bytes with capacity 10, got %d", cap(bs3))
	}
}

func TestBytesStringUnsafe(t *testing.T) {
	bs := Bytes("hello world")
	str := bs.StringUnsafe()
	expected := String("hello world")
	if str != expected {
		t.Errorf("StringUnsafe failed. Expected: %s, Got: %s", expected, str)
	}
}

func TestBytesLen(t *testing.T) {
	bs1 := Bytes("")
	if bs1.Len() != 0 {
		t.Errorf("Len() for empty bytes should be 0, got %d", bs1.Len())
	}

	bs2 := Bytes("hello")
	if bs2.Len() != 5 {
		t.Errorf("Len() for 'hello' should be 5, got %d", bs2.Len())
	}
}

func TestBytesReset(t *testing.T) {
	bs := NewBytes(5, 10)
	// Set some data
	copy(bs, "hello")

	if bs.Len() != 5 {
		t.Errorf("Initial length should be 5, got %d", bs.Len())
	}

	bs.Reset()
	if bs.Len() != 0 {
		t.Errorf("Length after Reset() should be 0, got %d", bs.Len())
	}
	if cap(bs) != 10 {
		t.Errorf("Capacity after Reset() should be preserved (10), got %d", cap(bs))
	}
}

func TestBytesPrint(t *testing.T) {
	bs := Bytes("test print")
	result := bs.Print()
	if !bytes.Equal(result, bs) {
		t.Errorf("Print() should return original bytes unchanged")
	}
}

func TestBytesPrintln(t *testing.T) {
	bs := Bytes("test println")
	result := bs.Println()
	if !bytes.Equal(result, bs) {
		t.Errorf("Println() should return original bytes unchanged")
	}
}

func wantInt64BE(b Bytes) int64 {
	var buf [8]byte
	if len(b) > 8 {
		b = b[len(b)-8:]
	}

	copy(buf[8-len(b):], b)

	if len(b) > 0 && b[0]&0x80 != 0 && len(b) < 8 {
		for i := 0; i < 8-len(b); i++ {
			buf[i] = 0xFF
		}
	}

	return int64(binary.BigEndian.Uint64(buf[:]))
}

func wantInt64LE(b Bytes) int64 {
	var buf [8]byte
	if len(b) > 8 {
		b = b[:8]
	}
	copy(buf[:len(b)], b)

	if len(b) > 0 && b[len(b)-1]&0x80 != 0 && len(b) < 8 {
		for i := len(b); i < 8; i++ {
			buf[i] = 0xFF
		}
	}

	return int64(binary.LittleEndian.Uint64(buf[:]))
}

func TestBytesIntSigned_Orders(t *testing.T) {
	type tc struct {
		name string
		in   Bytes
	}

	cases := []tc{
		{name: "empty", in: Bytes{}},
		{name: "single positive", in: Bytes{5}},
		{name: "single -1 (0xFF)", in: Bytes{0xFF}},
		{name: "single -128 (0x80)", in: Bytes{0x80}},
		{name: "two bytes BE negative 0xFFFE (-2)", in: Bytes{0xFF, 0xFE}},
		{name: "two bytes LE negative 0xFFFE (-2)", in: Bytes{0xFE, 0xFF}},
		{name: "short positive 3 bytes", in: Bytes{1, 2, 3}},
		{name: "long -> truncate (mixed)", in: Bytes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{name: "long negative BE after trunc", in: Bytes{0x00, 0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8}},
		{name: "long negative LE first8", in: Bytes{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x80, 0x99}},
	}

	for _, c := range cases {
		t.Run("BE/"+c.name, func(t *testing.T) {
			want := wantInt64BE(c.in)
			got := int64(c.in.IntBE())
			if got != want {
				t.Fatalf("IntBE(%v): want %d, got %d", c.in, want, got)
			}
		})
		t.Run("LE/"+c.name, func(t *testing.T) {
			want := wantInt64LE(c.in)
			got := int64(c.in.IntLE())
			if got != want {
				t.Fatalf("IntLE(%v): want %d, got %d", c.in, want, got)
			}
		})
	}
}

func TestBytesFloat_Orders(t *testing.T) {
	type tc struct {
		name string
		in   Bytes
		want float64
	}
	cases := []tc{
		{name: "zero BE", in: Bytes{0, 0, 0, 0, 0, 0, 0, 0}, want: 0.0},
		{name: "negative zero BE", in: Bytes{0x80, 0, 0, 0, 0, 0, 0, 0}, want: math.Copysign(0, -1)},
		{name: "positive infinity BE", in: Bytes{0x7F, 0xF0, 0, 0, 0, 0, 0, 0}, want: math.Inf(1)},
		{name: "negative infinity BE", in: Bytes{0xFF, 0xF0, 0, 0, 0, 0, 0, 0}, want: math.Inf(-1)},
		{name: "NaN BE", in: Bytes{0x7F, 0xF8, 0, 0, 0, 0, 0, 1}, want: math.NaN()},
		{name: "simple positive BE", in: Bytes{0x40, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, want: 42.0},
		{name: "simple negative BE", in: Bytes{0xC0, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, want: -42.0},

		// Invalid length cases
		{name: "empty", in: Bytes{}, want: 0.0},
		{name: "too short", in: Bytes{1, 2, 3}, want: 0.0},
		{name: "too long", in: Bytes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, want: 0.0},
	}

	for _, c := range cases {
		t.Run("FloatBE/"+c.name, func(t *testing.T) {
			got := float64(c.in.FloatBE())

			// Special handling for NaN since NaN != NaN
			if math.IsNaN(c.want) {
				if !math.IsNaN(got) {
					t.Fatalf("FloatBE(%v): want NaN, got %g", c.in, got)
				}
				return
			}

			// Special handling for negative zero
			if c.want == 0 && math.Signbit(c.want) {
				if got != 0 || !math.Signbit(got) {
					t.Fatalf("FloatBE(%v): want -0, got %g (signbit: %t)", c.in, got, math.Signbit(got))
				}
				return
			}

			if got != c.want {
				t.Fatalf("FloatBE(%v): want %g, got %g", c.in, c.want, got)
			}
		})

		t.Run("FloatLE/"+c.name, func(t *testing.T) {
			// Convert BE bytes to LE for testing
			var leBytes Bytes
			if len(c.in) == 8 {
				leBytes = make(Bytes, 8)
				for i := 0; i < 8; i++ {
					leBytes[i] = c.in[7-i]
				}
			} else {
				leBytes = c.in // Invalid length cases
			}

			got := float64(leBytes.FloatLE())

			// Special handling for NaN since NaN != NaN
			if math.IsNaN(c.want) {
				if !math.IsNaN(got) {
					t.Fatalf("FloatLE(%v): want NaN, got %g", leBytes, got)
				}
				return
			}

			// Special handling for negative zero
			if c.want == 0 && math.Signbit(c.want) {
				if got != 0 || !math.Signbit(got) {
					t.Fatalf("FloatLE(%v): want -0, got %g (signbit: %t)", leBytes, got, math.Signbit(got))
				}
				return
			}

			if got != c.want {
				t.Fatalf("FloatLE(%v): want %g, got %g", leBytes, c.want, got)
			}
		})
	}
}

func TestBytesIsLower(t *testing.T) {
	tests := []struct {
		name string
		in   Bytes
		want bool
	}{
		{"Empty", Bytes(""), false},
		{"OnlyDigits", Bytes("12345"), false},
		{"OnlyPunct", Bytes("!?-+"), false},
		{"ASCII_lower", Bytes("hello"), true},
		{"ASCII_upper", Bytes("HELLO"), false},
		{"ASCII_mixed", Bytes("Hello"), false},
		{"LowerWithDigits", Bytes("abc123!"), true},
		{"UpperWithDigits", Bytes("ABC123!"), false},
		{"MixedWithPunct", Bytes("abc-DEF"), false},
		{"Cyrillic_lower", Bytes("привет мир"), true},
		{"Cyrillic_upper", Bytes("ПРИВЕТ"), false},
		{"Cyrillic_mixed", Bytes("Привет"), false},
		{"Greek_lower", Bytes("γειασου"), true},
		{"Greek_upper", Bytes("ΚΑΛΗΜΕΡΑ"), false},
		{"Greek_mixed", Bytes("Γεια"), false},
		{"Latin_German_eszett", Bytes("straße"), true},
		{"Latin_Turkish_lower", Bytes("ıi"), true},
		{"Latin_Turkish_upper", Bytes("İI"), false},
		{"CombiningLower", Bytes("e\u0301gal"), true},
		{"CombiningMixed", Bytes("E\u0301gal"), false},
		{"InvalidUTF8", Bytes([]byte{0xff, 0xfe, 0xfd}), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.IsLower()
			if got != tc.want {
				t.Errorf("IsLower(%q) = %v; want %v", []byte(tc.in), got, tc.want)
			}
		})
	}
}

func TestBytesIsUpper(t *testing.T) {
	tests := []struct {
		name string
		in   Bytes
		want bool
	}{
		{"Empty", Bytes(""), false},
		{"OnlyDigits", Bytes("12345"), false},
		{"OnlyPunct", Bytes("!?-+"), false},
		{"ASCII_upper", Bytes("HELLO"), true},
		{"ASCII_lower", Bytes("hello"), false},
		{"ASCII_mixed", Bytes("Hello"), false},
		{"UpperWithDigits", Bytes("ABC123!"), true},
		{"LowerWithDigits", Bytes("abc123!"), false},
		{"MixedWithPunct", Bytes("ABC-def"), false},
		{"Cyrillic_upper", Bytes("ПРИВЕТ"), true},
		{"Cyrillic_lower", Bytes("привет"), false},
		{"Cyrillic_mixed", Bytes("Привет"), false},
		{"Greek_upper", Bytes("ΚΑΛΗΜΕΡΑ"), true},
		{"Greek_lower", Bytes("γειασου"), false},
		{"Greek_mixed", Bytes("Γεια"), false},
		{"Latin_German_eszett", Bytes("STRAẞE"), true},
		{"Latin_Turkish_upper", Bytes("İI"), true},
		{"Latin_Turkish_lower", Bytes("ıi"), false},
		{"CombiningUpper", Bytes("E\u0301GAL"), true},
		{"CombiningMixed", Bytes("e\u0301GAL"), false},
		{"InvalidUTF8", Bytes([]byte{0xff, 0xfe, 0xfd}), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.IsUpper()
			if got != tc.want {
				t.Errorf("IsUpper(%q) = %v; want %v", []byte(tc.in), got, tc.want)
			}
		})
	}
}

func TestBytesScan(t *testing.T) {
	var b Bytes

	if err := b.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) error: %v", err)
	}
	if b != nil {
		t.Fatalf("Expected nil, got %v", b)
	}

	input := []byte{1, 2, 3}
	if err := b.Scan(input); err != nil {
		t.Fatalf("Scan([]byte) error: %v", err)
	}
	if string(b) != string(input) {
		t.Fatalf("Expected %v, got %v", input, b)
	}

	err := b.Scan("not bytes")
	if err == nil {
		t.Fatal("Expected error for unsupported type")
	}
}

func TestBytesValue(t *testing.T) {
	b := Bytes{1, 2, 3}
	val, err := b.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if bv, ok := val.([]byte); !ok || string(bv) != string(b) {
		t.Fatalf("Expected %v, got %v", b, val)
	}

	var empty Bytes
	val2, err := empty.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if val2 != nil && len(val2.([]byte)) != 0 {
		t.Fatalf("Expected nil or empty slice, got %v", val2)
	}
}
