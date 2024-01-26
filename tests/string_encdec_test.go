package g_test

import (
	"testing"

	"gitlab.com/x0xO/g"
)

func TestStringBase64Encode(t *testing.T) {
	tests := []struct {
		name string
		e    g.String
		want g.String
	}{
		{"empty", "", ""},
		{"hello", "hello", "aGVsbG8="},
		{"hello world", "hello world", "aGVsbG8gd29ybGQ="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Enc().Base64(); got != tt.want {
				t.Errorf("enc.Base64Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBase64Decode(t *testing.T) {
	tests := []struct {
		name string
		d    g.String
		want g.String
	}{
		{"base64 decode", "aGVsbG8gd29ybGQ=", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Dec().Base64().Unwrap(); got != tt.want {
				t.Errorf("dec.Base64Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringRot13(t *testing.T) {
	input := g.String("hello world")
	expected := g.String("uryyb jbeyq")
	actual := input.Enc().Rot13()

	if actual != expected {
		t.Errorf("Rot13Encode(%q) = %q; expected %q", input, actual, expected)
	}

	input = g.String("uryyb jbeyq")
	expected = g.String("hello world")
	actual = input.Dec().Rot13()

	if actual != expected {
		t.Errorf("Rot13Decode(%q) = %q; expected %q", input, actual, expected)
	}
}

func TestStringXOR(t *testing.T) {
	for range 100 {
		input := g.NewString("").Random(g.NewInt(30).RandomRange(100).Std())
		key := g.NewString("").Random(10)
		obfuscated := input.Enc().XOR(key)
		deobfuscated := obfuscated.Dec().XOR(key)

		if input != deobfuscated {
			t.Errorf("expected %s, but got %s", input, deobfuscated)
		}
	}
}

func TestXOR(t *testing.T) {
	tests := []struct {
		input string
		key   string
		want  string
	}{
		{"01", "qsCDE", "AB"},
		{"123", "ABCDE", "ppp"},
		{"12345", "98765", "\x08\x0a\x04\x02\x00"},
		{"Hello", "wORLD", "?*> +"},
		// {"Hello,", "World!", "\x0f\x0a\x1e\x00\x0b\x0d"},
		// {"`c345", "QQ", "12345"},
		{"abcde", "01234", "QSQWQ"},
		{"lowercase", "9?'      ", "UPPERCASE"},
		{"test", "", "test"},
		{"test", "test", "\x00\x00\x00\x00"},
	}

	for _, tt := range tests {
		got := g.NewString(tt.input).Enc().XOR(g.String(tt.key))
		if got != g.String(tt.want) {
			t.Errorf("XOR(%q, %q) = %q; want %q", tt.input, tt.key, got, tt.want)
		}
	}
}

func TestGzFlateDecode(t *testing.T) {
	testCases := []struct {
		name     string
		input    g.String
		expected g.String
	}{
		{"Valid compressed data", "8kjNycnXUXCvcstJLElVBAAAAP//AQAA//8=", "Hello, GzFlate!"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.Dec().Base64().Unwrap().Decomp().Flate().Unwrap()
			if result.Ne(tc.expected) {
				t.Errorf("GzFlateDecode, expected: %s, got: %s", tc.expected, result)
			}
		})
	}
}

func TestGzFlateEncode(t *testing.T) {
	testCases := []struct {
		name     string
		input    g.String
		expected g.String
	}{
		{"Empty input", "", "AAAA//8BAAD//w=="},
		{"Valid input", "Hello, GzFlate!", "8kjNycnXUXCvcstJLElVBAAAAP//AQAA//8="},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.Comp().Flate().Enc().Base64()
			if result.Ne(tc.expected) {
				t.Errorf("GzFlateEncode, expected: %s, got: %s", tc.expected, result)
			}
		})
	}
}
