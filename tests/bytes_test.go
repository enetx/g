package g_test

import (
	"bytes"
	"io"
	"testing"

	"gitlab.com/x0xO/g"
)

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
