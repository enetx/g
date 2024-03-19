package g_test

import (
	"bytes"
	"io"
	"reflect"
	"regexp"
	"testing"
	"unicode"

	"github.com/enetx/g"
)

func TestNewBytes(t *testing.T) {
	// Test case with string input
	strInput := "hello"
	bytesFromStr := g.NewBytes(strInput)
	expectedBytesFromStr := g.Bytes("hello")
	if !bytes.Equal(bytesFromStr, expectedBytesFromStr) {
		t.Errorf("Conversion from string failed. Expected: %s, Got: %s", expectedBytesFromStr, bytesFromStr)
	}

	// Test case with []byte input
	byteSliceInput := []byte{104, 101, 108, 108, 111}
	bytesFromSlice := g.NewBytes(byteSliceInput)
	expectedBytesFromSlice := g.Bytes("hello")
	if !bytes.Equal(bytesFromSlice, expectedBytesFromSlice) {
		t.Errorf("Conversion from byte slice failed. Expected: %s, Got: %s", expectedBytesFromSlice, bytesFromSlice)
	}
}

func TestBytesReplace(t *testing.T) {
	// Test case where old byte sequence exists and is replaced
	bs1 := g.Bytes("hello world")
	oldB1 := g.Bytes("world")
	newB1 := g.Bytes("gopher")
	replaced1 := bs1.Replace(oldB1, newB1, -1)
	expected1 := g.Bytes("hello gopher")
	if !bytes.Equal(replaced1, expected1) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected1, replaced1)
	}

	// Test case with multiple occurrences of old byte sequence
	bs2 := g.Bytes("hello world hello world")
	oldB2 := g.Bytes("world")
	newB2 := g.Bytes("gopher")
	replaced2 := bs2.Replace(oldB2, newB2, -1)
	expected2 := g.Bytes("hello gopher hello gopher")
	if !bytes.Equal(replaced2, expected2) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected2, replaced2)
	}

	// Test case with limited replacements
	bs3 := g.Bytes("hello world hello world")
	oldB3 := g.Bytes("world")
	newB3 := g.Bytes("gopher")
	replaced3 := bs3.Replace(oldB3, newB3, 1)
	expected3 := g.Bytes("hello gopher hello world")
	if !bytes.Equal(replaced3, expected3) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected3, replaced3)
	}

	// Test case where old byte sequence doesn't exist
	bs4 := g.Bytes("hello world")
	oldB4 := g.Bytes("gopher")
	newB4 := g.Bytes("earth")
	replaced4 := bs4.Replace(oldB4, newB4, -1)
	if !bytes.Equal(replaced4, bs4) {
		t.Errorf("Expected no change when old byte sequence doesn't exist. Got: %s", replaced4)
	}
}

func TestReplaceAll(t *testing.T) {
	// Test case where old byte sequence exists and is replaced
	bs1 := g.Bytes("hello world")
	oldB1 := g.Bytes("world")
	newB1 := g.Bytes("gopher")
	replaced1 := bs1.ReplaceAll(oldB1, newB1)
	expected1 := g.Bytes("hello gopher")
	if !bytes.Equal(replaced1, expected1) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected1, replaced1)
	}

	// Test case with multiple occurrences of old byte sequence
	bs2 := g.Bytes("hello world hello world")
	oldB2 := g.Bytes("world")
	newB2 := g.Bytes("gopher")
	replaced2 := bs2.ReplaceAll(oldB2, newB2)
	expected2 := g.Bytes("hello gopher hello gopher")
	if !bytes.Equal(replaced2, expected2) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected2, replaced2)
	}

	// Test case where old byte sequence doesn't exist
	bs3 := g.Bytes("hello world")
	oldB3 := g.Bytes("gopher")
	newB3 := g.Bytes("earth")
	replaced3 := bs3.ReplaceAll(oldB3, newB3)
	if !bytes.Equal(replaced3, bs3) {
		t.Errorf("Expected no change when old byte sequence doesn't exist. Got: %s", replaced3)
	}
}

func TestBytesReplaceRegexp(t *testing.T) {
	// Test case where pattern matches and is replaced
	bs1 := g.Bytes("hello world hello world")
	pattern1 := regexp.MustCompile("world")
	newB1 := g.Bytes("gopher")
	replaced1 := bs1.ReplaceRegexp(pattern1, newB1)
	expected1 := g.Bytes("hello gopher hello gopher")
	if !bytes.Equal(replaced1, expected1) {
		t.Errorf("Replacement failed. Expected: %s, Got: %s", expected1, replaced1)
	}

	// Test case where pattern matches and is replaced with capture group
	bs2 := g.Bytes("apple apple apple")
	pattern2 := regexp.MustCompile(`(\w+)`)
	newB2 := g.Bytes("${1}s")
	replaced2 := bs2.ReplaceRegexp(pattern2, newB2)
	expected2 := g.Bytes("apples apples apples")
	if !bytes.Equal(replaced2, expected2) {
		t.Errorf("Replacement with capture group failed. Expected: %s, Got: %s", expected2, replaced2)
	}

	// Test case where pattern doesn't match
	bs3 := g.Bytes("hello world")
	pattern3 := regexp.MustCompile("gopher")
	newB3 := g.Bytes("earth")
	replaced3 := bs3.ReplaceRegexp(pattern3, newB3)
	if !bytes.Equal(replaced3, bs3) {
		t.Errorf("Expected no change when pattern doesn't match. Got: %s", replaced3)
	}
}

func TestBytesFindRegexp(t *testing.T) {
	// Test case where pattern matches and is found
	bs1 := g.Bytes("hello world")
	pattern1 := regexp.MustCompile("world")
	found1 := bs1.FindRegexp(pattern1)
	expected1 := g.Bytes("world")
	if found1.IsNone() {
		t.Errorf("Expected to find matching pattern, but found none")
	} else if !bytes.Equal(found1.Unwrap(), expected1) {
		t.Errorf("Found pattern does not match expected result. Expected: %s, Got: %s", expected1, found1.Unwrap())
	}

	// Test case where pattern doesn't match
	bs2 := g.Bytes("hello world")
	pattern2 := regexp.MustCompile("gopher")
	found2 := bs2.FindRegexp(pattern2)
	if found2.IsSome() {
		t.Errorf("Expected not to find matching pattern, but found one")
	}
}

func TestBytesTrim(t *testing.T) {
	// Test case where cutset matches characters at the beginning and end
	bs1 := g.Bytes("!!hello world!!")
	cutset1 := g.String("! ")
	trimmed1 := bs1.Trim(cutset1)
	expected1 := g.Bytes("hello world")
	if !bytes.Equal(trimmed1, expected1) {
		t.Errorf("Trimming failed. Expected: %s, Got: %s", expected1, trimmed1)
	}

	// Test case where cutset matches characters only at the beginning
	bs2 := g.Bytes("!!!hello world")
	cutset2 := g.String("!")
	trimmed2 := bs2.Trim(cutset2)
	expected2 := g.Bytes("hello world")
	if !bytes.Equal(trimmed2, expected2) {
		t.Errorf("Trimming failed. Expected: %s, Got: %s", expected2, trimmed2)
	}

	// Test case where cutset matches characters only at the end
	bs3 := g.Bytes("hello world!!!")
	cutset3 := g.String("!")
	trimmed3 := bs3.Trim(cutset3)
	expected3 := g.Bytes("hello world")
	if !bytes.Equal(trimmed3, expected3) {
		t.Errorf("Trimming failed. Expected: %s, Got: %s", expected3, trimmed3)
	}

	// Test case where cutset doesn't match any characters
	bs4 := g.Bytes("hello world")
	cutset4 := g.String("-")
	trimmed4 := bs4.Trim(cutset4)
	if !bytes.Equal(trimmed4, bs4) {
		t.Errorf("Expected no change when cutset doesn't match any characters. Got: %s", trimmed4)
	}
}

func TestBytesTrimLeft(t *testing.T) {
	// Test case where cutset matches characters at the beginning
	bs1 := g.Bytes("!!hello world!!")
	cutset1 := g.String("! ")
	trimmed1 := bs1.TrimLeft(cutset1)
	expected1 := g.Bytes("hello world!!")
	if !bytes.Equal(trimmed1, expected1) {
		t.Errorf("Trimming left failed. Expected: %s, Got: %s", expected1, trimmed1)
	}

	// Test case where cutset doesn't match any characters at the beginning
	bs2 := g.Bytes("hello world")
	cutset2 := g.String("-")
	trimmed2 := bs2.TrimLeft(cutset2)
	if !bytes.Equal(trimmed2, bs2) {
		t.Errorf("Expected no change when cutset doesn't match any characters. Got: %s", trimmed2)
	}
}

func TestBytesTrimRight(t *testing.T) {
	// Test case where cutset matches characters at the end
	bs1 := g.Bytes("!!hello world!!")
	cutset1 := g.String("! ")
	trimmed1 := bs1.TrimRight(cutset1)
	expected1 := g.Bytes("!!hello world")
	if !bytes.Equal(trimmed1, expected1) {
		t.Errorf("Trimming right failed. Expected: %s, Got: %s", expected1, trimmed1)
	}

	// Test case where cutset doesn't match any characters at the end
	bs2 := g.Bytes("hello world")
	cutset2 := g.String("-")
	trimmed2 := bs2.TrimRight(cutset2)
	if !bytes.Equal(trimmed2, bs2) {
		t.Errorf("Expected no change when cutset doesn't match any characters. Got: %s", trimmed2)
	}
}

func TestBytesTrimPrefix(t *testing.T) {
	// Test case where cutset matches the prefix
	bs1 := g.Bytes("prefix_hello world")
	cutset1 := g.Bytes("prefix_")
	trimmed1 := bs1.TrimPrefix(cutset1)
	expected1 := g.Bytes("hello world")
	if !bytes.Equal(trimmed1, expected1) {
		t.Errorf("Trimming prefix failed. Expected: %s, Got: %s", expected1, trimmed1)
	}

	// Test case where cutset doesn't match the prefix
	bs2 := g.Bytes("hello world")
	cutset2 := g.Bytes("nonexistent_")
	trimmed2 := bs2.TrimPrefix(cutset2)
	if !bytes.Equal(trimmed2, bs2) {
		t.Errorf("Expected no change when cutset doesn't match the prefix. Got: %s", trimmed2)
	}
}

func TestBytesTrimSuffix(t *testing.T) {
	// Test case where cutset matches the suffix
	bs1 := g.Bytes("hello world_suffix")
	cutset1 := g.Bytes("_suffix")
	trimmed1 := bs1.TrimSuffix(cutset1)
	expected1 := g.Bytes("hello world")
	if !bytes.Equal(trimmed1, expected1) {
		t.Errorf("Trimming suffix failed. Expected: %s, Got: %s", expected1, trimmed1)
	}

	// Test case where cutset doesn't match the suffix
	bs2 := g.Bytes("hello world")
	cutset2 := g.Bytes("_nonexistent")
	trimmed2 := bs2.TrimSuffix(cutset2)
	if !bytes.Equal(trimmed2, bs2) {
		t.Errorf("Expected no change when cutset doesn't match the suffix. Got: %s", trimmed2)
	}
}

func TestBytesSplit(t *testing.T) {
	// Test case where separator exists
	bs1 := g.Bytes("hello world gopher")
	separator1 := g.NewBytes(" ")
	split1 := bs1.Split(separator1)
	expected1 := g.SliceOf(g.NewBytes("hello"), g.NewBytes("world"), g.NewBytes("gopher"))
	if !reflect.DeepEqual(split1, expected1) {
		t.Errorf("Split failed. Expected: %v, Got: %v", expected1, split1)
	}

	// Test case where separator doesn't exist
	bs2 := g.Bytes("helloworldgopher")
	separator2 := g.NewBytes(" ")
	split2 := bs2.Split(separator2)
	expected2 := g.Slice[g.Bytes]{g.NewBytes("helloworldgopher")}
	if !reflect.DeepEqual(split2, expected2) {
		t.Errorf("Split failed. Expected: %v, Got: %v", expected2, split2)
	}
}

func TestBytesAdd(t *testing.T) {
	// Test case where bytes are added
	bs1 := g.Bytes("hello")
	obs1 := g.NewBytes(" world")
	added1 := bs1.Add(obs1)
	expected1 := g.Bytes("hello world")
	if !bytes.Equal(added1, expected1) {
		t.Errorf("Add failed. Expected: %s, Got: %s", expected1, added1)
	}
}

func TestBytesAddPrefix(t *testing.T) {
	// Test case where bytes are added as a prefix
	bs1 := g.Bytes("world")
	obs1 := g.NewBytes("hello ")
	prefixed1 := bs1.AddPrefix(obs1)
	expected1 := g.Bytes("hello world")
	if !bytes.Equal(prefixed1, expected1) {
		t.Errorf("AddPrefix failed. Expected: %s, Got: %s", expected1, prefixed1)
	}
}

func TestBytesStd(t *testing.T) {
	// Test case where Bytes is converted to a byte slice
	bs1 := g.Bytes("hello world")
	std1 := bs1.Std()
	expected1 := []byte("hello world")
	if !bytes.Equal(std1, expected1) {
		t.Errorf("Std failed. Expected: %v, Got: %v", expected1, std1)
	}
}

func TestBytesClone(t *testing.T) {
	// Test case where Bytes is cloned
	bs1 := g.Bytes("hello world")
	cloned1 := bs1.Clone()
	if !bytes.Equal(cloned1, bs1) {
		t.Errorf("Clone failed. Expected: %s, Got: %s", bs1, cloned1)
	}
}

func TestBytesContainsAnyChars(t *testing.T) {
	// Test case where Bytes contains any characters from the input String
	bs1 := g.Bytes("hello")
	chars1 := g.String("aeiou")
	contains1 := bs1.ContainsAnyChars(chars1)
	if !contains1 {
		t.Errorf("ContainsAnyChars failed. Expected: true, Got: %t", contains1)
	}

	// Test case where Bytes doesn't contain any characters from the input String
	bs2 := g.Bytes("hello")
	chars2 := g.String("xyz")
	contains2 := bs2.ContainsAnyChars(chars2)
	if contains2 {
		t.Errorf("ContainsAnyChars failed. Expected: false, Got: %t", contains2)
	}
}

func TestBytesContainsRune(t *testing.T) {
	// Test case where Bytes contains the specified rune
	bs1 := g.Bytes("hello")
	rune1 := 'e'
	contains1 := bs1.ContainsRune(rune1)
	if !contains1 {
		t.Errorf("ContainsRune failed. Expected: true, Got: %t", contains1)
	}

	// Test case where Bytes doesn't contain the specified rune
	bs2 := g.Bytes("hello")
	rune2 := 'x'
	contains2 := bs2.ContainsRune(rune2)
	if contains2 {
		t.Errorf("ContainsRune failed. Expected: false, Got: %t", contains2)
	}
}

func TestBytesCount(t *testing.T) {
	// Test case where Bytes contains multiple occurrences of the specified Bytes
	bs1 := g.Bytes("hello hello hello")
	obs1 := g.Bytes("hello")
	count1 := bs1.Count(obs1)
	expected1 := 3
	if count1 != expected1 {
		t.Errorf("Count failed. Expected: %d, Got: %d", expected1, count1)
	}

	// Test case where Bytes doesn't contain the specified Bytes
	bs2 := g.Bytes("hello")
	obs2 := g.Bytes("world")
	count2 := bs2.Count(obs2)
	expected2 := 0
	if count2 != expected2 {
		t.Errorf("Count failed. Expected: %d, Got: %d", expected2, count2)
	}
}

func TestBytesCompare(t *testing.T) {
	testCases := []struct {
		bs1      g.Bytes
		bs2      g.Bytes
		expected g.Int
	}{
		{[]byte("apple"), []byte("banana"), -1},
		{[]byte("banana"), []byte("apple"), 1},
		{[]byte("banana"), []byte("banana"), 0},
		{[]byte("apple"), []byte("Apple"), 1},
		{[]byte(""), []byte(""), 0},
	}

	for _, tc := range testCases {
		result := tc.bs1.Compare(tc.bs2)
		if !result.Eq(tc.expected) {
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
	bs1 := g.Bytes("Hello World")
	obs1 := g.Bytes("hello world")
	eqFold1 := bs1.EqFold(obs1)
	if !eqFold1 {
		t.Errorf("EqFold failed. Expected: true, Got: %t", eqFold1)
	}

	// Test case where the byte slices are not equal regardless of case
	bs2 := g.Bytes("Hello World")
	obs2 := g.Bytes("gopher")
	eqFold2 := bs2.EqFold(obs2)
	if eqFold2 {
		t.Errorf("EqFold failed. Expected: false, Got: %t", eqFold2)
	}
}

func TestBytesEq(t *testing.T) {
	testCases := []struct {
		bs1      g.Bytes
		bs2      g.Bytes
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
		bs1      g.Bytes
		bs2      g.Bytes
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
		bs1      g.Bytes
		bs2      g.Bytes
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
		bs1      g.Bytes
		bs2      g.Bytes
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
		input    g.Bytes
		expected g.Bytes
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
		bs       g.Bytes
		expected []byte
	}{
		{"Empty Bytes", g.Bytes{}, []byte{}},
		{"Single byte Bytes", g.Bytes{0x41}, []byte{0x41}},
		{
			"Multiple bytes Bytes",
			g.Bytes{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			[]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
		},
		{
			"Bytes with various values",
			g.Bytes{0x00, 0xff, 0x80, 0x7f},
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
		bs       g.Bytes
		bss      []g.Bytes
		expected bool
	}{
		{
			bs:       g.Bytes("Hello, world!"),
			bss:      []g.Bytes{g.Bytes("world"), g.Bytes("Go")},
			expected: true,
		},
		{
			bs:       g.Bytes("Welcome to the HumanGo-1!"),
			bss:      []g.Bytes{g.Bytes("Go-3"), g.Bytes("Go-4")},
			expected: false,
		},
		{
			bs:       g.Bytes("Have a great day!"),
			bss:      []g.Bytes{g.Bytes(""), g.Bytes(" ")},
			expected: true,
		},
		{
			bs:       g.Bytes(""),
			bss:      []g.Bytes{g.Bytes("Hello"), g.Bytes("world")},
			expected: false,
		},
		{
			bs:       g.Bytes(""),
			bss:      []g.Bytes{},
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
		bs       g.Bytes
		bss      []g.Bytes
		expected bool
	}{
		{
			bs:       g.Bytes("Hello, world!"),
			bss:      []g.Bytes{g.Bytes("Hello"), g.Bytes("world")},
			expected: true,
		},
		{
			bs:       g.Bytes("Welcome to the HumanGo-1!"),
			bss:      []g.Bytes{g.Bytes("Go-3"), g.Bytes("Go-4")},
			expected: false,
		},
		{
			bs:       g.Bytes("Have a great day!"),
			bss:      []g.Bytes{g.Bytes("Have"), g.Bytes("a")},
			expected: true,
		},
		{
			bs:       g.Bytes(""),
			bss:      []g.Bytes{g.Bytes("Hello"), g.Bytes("world")},
			expected: false,
		},
		{
			bs:       g.Bytes("Hello, world!"),
			bss:      []g.Bytes{},
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
	bs := g.Bytes("hello world")
	obs := g.Bytes("world")
	idx := bs.Index(obs)
	expected := 6
	if idx != expected {
		t.Errorf("Index failed. Expected: %d, Got: %d", expected, idx)
	}

	// Test case where obs is not present in bs
	bs = g.Bytes("hello world")
	obs = g.Bytes("gopher")
	idx = bs.Index(obs)
	expected = -1
	if idx != expected {
		t.Errorf("Index failed. Expected: %d, Got: %d", expected, idx)
	}
}

func TestBytesIndexRegexp(t *testing.T) {
	// Test case where a match is found
	bs := g.Bytes("apple banana")
	pattern := regexp.MustCompile(`banana`)
	idx := bs.IndexRegexp(pattern)
	expected := g.Some(g.Slice[g.Int]{g.NewInt(6), g.NewInt(12)})
	if idx.IsNone() || !reflect.DeepEqual(idx.Some(), expected.Some()) {
		t.Errorf("IndexRegexp failed. Expected: %v, Got: %v", expected, idx)
	}

	// Test case where no match is found
	bs = g.Bytes("apple banana")
	pattern = regexp.MustCompile(`orange`)
	idx = bs.IndexRegexp(pattern)
	expected = g.None[g.Slice[g.Int]]()
	if idx.IsSome() || !reflect.DeepEqual(idx.IsNone(), expected.IsNone()) {
		t.Errorf("IndexRegexp failed. Expected: %v, Got: %v", expected, idx)
	}
}

func TestBytesRepeat(t *testing.T) {
	// Test case where the Bytes are repeated 3 times
	bs := g.Bytes("hello")
	repeated := bs.Repeat(3)
	expected := g.Bytes("hellohellohello")
	if !bytes.Equal(repeated, expected) {
		t.Errorf("Repeat failed. Expected: %s, Got: %s", expected, repeated)
	}

	// Test case where the Bytes are repeated 0 times
	bs = g.Bytes("hello")
	repeated = bs.Repeat(0)
	expected = g.Bytes("")
	if !bytes.Equal(repeated, expected) {
		t.Errorf("Repeat failed. Expected: %s, Got: %s", expected, repeated)
	}
}

func TestToRunes(t *testing.T) {
	// Test case where the Bytes are converted to runes
	bs := g.Bytes("hello")
	runes := bs.ToRunes()
	expected := []rune{'h', 'e', 'l', 'l', 'o'}
	if !reflect.DeepEqual(runes, expected) {
		t.Errorf("ToRunes failed. Expected: %v, Got: %v", expected, runes)
	}
}

func TestBytesLower(t *testing.T) {
	// Test case where the Bytes are converted to lowercase
	bs := g.Bytes("Hello World")
	lower := bs.Lower()
	expected := g.Bytes("hello world")
	if !reflect.DeepEqual(lower, expected) {
		t.Errorf("Lower failed. Expected: %s, Got: %s", expected, lower)
	}
}

func TestBytesUpper(t *testing.T) {
	// Test case where the Bytes are converted to uppercase
	bs := g.Bytes("hello world")
	upper := bs.Upper()
	expected := g.Bytes("HELLO WORLD")
	if !reflect.DeepEqual(upper, expected) {
		t.Errorf("Upper failed. Expected: %s, Got: %s", expected, upper)
	}
}

func TestBytesTrimSpace(t *testing.T) {
	// Test case where white space characters are trimmed from the beginning and end
	bs := g.Bytes("  hello world  ")
	trimmed := bs.TrimSpace()
	expected := g.Bytes("hello world")
	if !bytes.Equal(trimmed, expected) {
		t.Errorf("TrimSpace failed. Expected: %s, Got: %s", expected, trimmed)
	}

	// Test case where there are no white space characters
	bs = g.Bytes("hello world")
	trimmed = bs.TrimSpace()
	expected = g.Bytes("hello world")
	if !bytes.Equal(trimmed, expected) {
		t.Errorf("TrimSpace failed. Expected: %s, Got: %s", expected, trimmed)
	}

	// Test case where the Bytes is empty
	bs = g.Bytes("")
	trimmed = bs.TrimSpace()
	expected = g.Bytes("")
	if !bytes.Equal(trimmed, expected) {
		t.Errorf("TrimSpace failed. Expected: %s, Got: %s", expected, trimmed)
	}
}

func TestBytesTitle(t *testing.T) {
	// Test case where the Bytes are converted to title case
	bs := g.Bytes("hello world")
	title := bs.Title()
	expected := g.Bytes("Hello World")
	if !reflect.DeepEqual(title, expected) {
		t.Errorf("Title failed. Expected: %s, Got: %s", expected, title)
	}

	// Test case where the Bytes are already in title case
	bs = g.Bytes("Hello World")
	title = bs.Title()
	if !reflect.DeepEqual(title, bs) {
		t.Errorf("Title failed. Expected: %s, Got: %s", bs, title)
	}

	// Test case where the Bytes are empty
	bs = g.Bytes("")
	title = bs.Title()
	if !reflect.DeepEqual(title, bs) {
		t.Errorf("Title failed. Expected: %s, Got: %s", bs, title)
	}
}

func TestBytesNotEmpty(t *testing.T) {
	// Test case where the Bytes is not empty
	bs := g.Bytes("hello")
	if !bs.NotEmpty() {
		t.Errorf("NotEmpty failed. Expected: true, Got: false")
	}

	// Test case where the Bytes is empty
	bs = g.Bytes("")
	if bs.NotEmpty() {
		t.Errorf("NotEmpty failed. Expected: false, Got: true")
	}
}

func TestBytesMap(t *testing.T) {
	// Test case where the function converts each rune to uppercase
	bs := g.Bytes("hello")
	uppercase := bs.Map(func(r rune) rune {
		return unicode.ToUpper(r)
	})

	expected := g.Bytes("HELLO")
	if !bytes.Equal(uppercase, expected) {
		t.Errorf("Map failed. Expected: %s, Got: %s", expected, uppercase)
	}

	// Test case where the function removes spaces
	bs = g.Bytes("hello world")
	noSpaces := bs.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1 // Remove rune
		}
		return r
	})

	expected = g.Bytes("helloworld")
	if !bytes.Equal(noSpaces, expected) {
		t.Errorf("Map failed. Expected: %s, Got: %s", expected, noSpaces)
	}
}

func TestBytesLenRunes(t *testing.T) {
	// Test case where the Bytes contain ASCII characters
	bs := g.Bytes("hello world")
	lenRunes := bs.LenRunes()
	expected := 11
	if lenRunes != expected {
		t.Errorf("LenRunes failed. Expected: %d, Got: %d", expected, lenRunes)
	}

	// Test case where the Bytes contain Unicode characters
	bs = g.Bytes("你好，世界")
	lenRunes = bs.LenRunes()
	expected = 5
	if lenRunes != expected {
		t.Errorf("LenRunes failed. Expected: %d, Got: %d", expected, lenRunes)
	}

	// Test case where the Bytes are empty
	bs = g.Bytes("")
	lenRunes = bs.LenRunes()
	expected = 0
	if lenRunes != expected {
		t.Errorf("LenRunes failed. Expected: %d, Got: %d", expected, lenRunes)
	}
}

func TestBytesLastIndex(t *testing.T) {
	// Test case where obs is present in bs
	bs := g.Bytes("hello world")
	obs := g.Bytes("o")
	lastIndex := bs.LastIndex(obs)
	expected := 7
	if lastIndex != expected {
		t.Errorf("LastIndex failed. Expected: %d, Got: %d", expected, lastIndex)
	}

	// Test case where obs is not present in bs
	bs = g.Bytes("hello world")
	obs = g.Bytes("z")
	lastIndex = bs.LastIndex(obs)
	expected = -1
	if lastIndex != expected {
		t.Errorf("LastIndex failed. Expected: %d, Got: %d", expected, lastIndex)
	}
}

func TestBytesIndexByte(t *testing.T) {
	// Test case where b is present in bs
	bs := g.Bytes("hello world")
	b := byte('o')
	indexByte := bs.IndexByte(b)
	expected := 4
	if indexByte != expected {
		t.Errorf("IndexByte failed. Expected: %d, Got: %d", expected, indexByte)
	}

	// Test case where b is not present in bs
	bs = g.Bytes("hello world")
	b = byte('z')
	indexByte = bs.IndexByte(b)
	expected = -1
	if indexByte != expected {
		t.Errorf("IndexByte failed. Expected: %d, Got: %d", expected, indexByte)
	}
}

func TestBytesLastIndexByte(t *testing.T) {
	// Test case where b is present in bs
	bs := g.Bytes("hello world")
	b := byte('o')
	lastIndexByte := bs.LastIndexByte(b)
	expected := 7
	if lastIndexByte != expected {
		t.Errorf("LastIndexByte failed. Expected: %d, Got: %d", expected, lastIndexByte)
	}

	// Test case where b is not present in bs
	bs = g.Bytes("hello world")
	b = byte('z')
	lastIndexByte = bs.LastIndexByte(b)
	expected = -1
	if lastIndexByte != expected {
		t.Errorf("LastIndexByte failed. Expected: %d, Got: %d", expected, lastIndexByte)
	}
}

func TestBytesIndexRune(t *testing.T) {
	// Test case where r is present in bs
	bs := g.Bytes("hello world")
	r := 'o'
	indexRune := bs.IndexRune(r)
	expected := 4
	if indexRune != expected {
		t.Errorf("IndexRune failed. Expected: %d, Got: %d", expected, indexRune)
	}

	// Test case where r is not present in bs
	bs = g.Bytes("hello world")
	r = 'z'
	indexRune = bs.IndexRune(r)
	expected = -1
	if indexRune != expected {
		t.Errorf("IndexRune failed. Expected: %d, Got: %d", expected, indexRune)
	}
}

func TestBytesFindAllSubmatchRegexpN(t *testing.T) {
	// Test case where matches are found
	bs := g.Bytes("hello world")
	pattern := regexp.MustCompile(`\b\w+\b`)
	matches := bs.FindAllSubmatchRegexpN(pattern, -1)
	if matches.IsSome() {
		expected := g.Slice[g.Slice[g.Bytes]]{
			{g.Bytes("hello")},
			{g.Bytes("world")},
		}
		if !matches.Some().Eq(expected) {
			t.Errorf("FindAllSubmatchRegexpN failed. Expected: %s, Got: %s", expected, matches.Some())
		}
	} else {
		t.Errorf("FindAllSubmatchRegexpN failed. Expected matches, Got None")
	}

	// Test case where no matches are found
	bs = g.Bytes("")
	pattern = regexp.MustCompile(`\b\w+\b`)
	matches = bs.FindAllSubmatchRegexpN(pattern, -1)
	if matches.IsSome() {
		t.Errorf("FindAllSubmatchRegexpN failed. Expected None, Got matches")
	}
}

func TestBytesFindAllRegexp(t *testing.T) {
	// Test case where matches are found
	bs := g.Bytes("hello world")
	pattern := regexp.MustCompile(`\b\w+\b`)
	matches := bs.FindAllRegexp(pattern)
	if matches.IsSome() {
		expected := g.Slice[g.Bytes]{
			g.Bytes("hello"),
			g.Bytes("world"),
		}
		if !matches.Some().Eq(expected) {
			t.Errorf("FindAllRegexp failed. Expected: %s, Got: %s", expected, matches.Some())
		}
	} else {
		t.Errorf("FindAllRegexp failed. Expected matches, Got None")
	}

	// Test case where no matches are found
	bs = g.Bytes("")
	pattern = regexp.MustCompile(`\b\w+\b`)
	matches = bs.FindAllRegexp(pattern)
	if matches.IsSome() {
		t.Errorf("FindAllRegexp failed. Expected None, Got matches")
	}
}

func TestBytesFindSubmatchRegexp(t *testing.T) {
	// Test case where a match is found
	bs := g.Bytes("hello world")
	pattern := regexp.MustCompile(`\b\w+\b`)
	match := bs.FindSubmatchRegexp(pattern)
	if match.IsSome() {
		expected := g.SliceOf(g.Bytes("hello"))
		if !match.Some().Eq(expected) {
			t.Errorf("FindSubmatchRegexp failed. Expected: %s, Got: %s", expected, match.Some())
		}
	} else {
		t.Errorf("FindSubmatchRegexp failed. Expected match, Got None")
	}

	// Test case where no match is found
	bs = g.Bytes("")
	pattern = regexp.MustCompile(`\b\w+\b`)
	match = bs.FindSubmatchRegexp(pattern)
	if match.IsSome() {
		t.Errorf("FindSubmatchRegexp failed. Expected None, Got match")
	}
}

func TestBytesFindAllSubmatchRegexp(t *testing.T) {
	// Test case where matches are found
	bs := g.Bytes("hello world")
	pattern := regexp.MustCompile(`\b\w+\b`)
	matches := bs.FindAllSubmatchRegexp(pattern)
	if matches.IsSome() {
		expected := g.Slice[g.Slice[g.Bytes]]{
			{g.Bytes("hello")},
			{g.Bytes("world")},
		}
		if !matches.Some().Eq(expected) {
			t.Errorf("FindAllSubmatchRegexp failed. Expected: %s, Got: %s", expected, matches.Some())
		}
	} else {
		t.Errorf("FindAllSubmatchRegexp failed. Expected matches, Got None")
	}

	// Test case where no matches are found
	bs = g.Bytes("")
	pattern = regexp.MustCompile(`\b\w+\b`)
	matches = bs.FindAllSubmatchRegexp(pattern)
	if matches.IsSome() {
		t.Errorf("FindAllSubmatchRegexp failed. Expected None, Got matches")
	}
}

func TestBytesHashingFunctions(t *testing.T) {
	// Test case for MD5 hashing
	input := g.Bytes("hello world")
	expectedMD5 := g.Bytes("5eb63bbbe01eeed093cb22bb8f5acdc3")
	md5Hash := input.Hash().MD5()
	if md5Hash.Ne(expectedMD5) {
		t.Errorf("MD5 hashing failed. Expected: %s, Got: %s", expectedMD5, md5Hash)
	}

	// Test case for SHA1 hashing
	expectedSHA1 := g.Bytes("2aae6c35c94fcfb415dbe95f408b9ce91ee846ed")
	sha1Hash := input.Hash().SHA1()
	if sha1Hash.Ne(expectedSHA1) {
		t.Errorf("SHA1 hashing failed. Expected: %s, Got: %s", expectedSHA1, sha1Hash)
	}

	// Test case for SHA256 hashing
	expectedSHA256 := g.Bytes("b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9")
	sha256Hash := input.Hash().SHA256()
	if sha256Hash.Ne(expectedSHA256) {
		t.Errorf("SHA256 hashing failed. Expected: %s, Got: %s", expectedSHA256, sha256Hash)
	}

	// Test case for SHA512 hashing
	expectedSHA512 := g.Bytes(
		"309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f",
	)
	sha512Hash := input.Hash().SHA512()
	if sha512Hash.Ne(expectedSHA512) {
		t.Errorf("SHA512 hashing failed. Expected: %s, Got: %s", expectedSHA512, sha512Hash)
	}
}
