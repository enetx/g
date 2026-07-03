package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

var parityCases = []struct {
	name string
	data Bytes
}{
	{"empty", Bytes("")},
	{"ascii", Bytes("Hello, World! 123 <tag> & 'quote' \"dq\"")},
	{"unicode", Bytes("Привет, мир! 你好 🌍")},
}

// binaryGarbage deliberately avoids '+' (0x2b), which URL-decode maps to a space.
var binaryGarbage = Bytes{0x00, 0xff, 0xfe, 0x01, 0x80, 0x7f, 0xc3}

func TestBytesJSONRoundtripParity(t *testing.T) {
	for _, tt := range parityCases {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.data.Encode().JSON().Unwrap()

			strEnc := tt.data.String().Encode().JSON().Unwrap()
			if !enc.Eq(strEnc.Bytes()) {
				t.Errorf("encode parity: Bytes = %q, String = %q", enc, strEnc)
			}

			dec := enc.Decode().JSON().Unwrap()
			if !dec.Eq(tt.data) {
				t.Errorf("roundtrip = %q, want %q", dec, tt.data)
			}

			strDec := strEnc.Decode().JSON().Unwrap()
			if !dec.Eq(strDec.Bytes()) {
				t.Errorf("decode parity: Bytes = %q, String = %q", dec, strDec)
			}
		})
	}
}

func TestBytesJSONBinaryGarbage(t *testing.T) {
	// encoding/json/v2 rejects invalid UTF-8 instead of lossily replacing it
	// with U+FFFD, so non-UTF-8 input yields Err for both Bytes and String.
	if !binaryGarbage.Encode().JSON().IsErr() {
		t.Error("expected Err for non-UTF-8 Bytes input")
	}

	if !binaryGarbage.String().Encode().JSON().IsErr() {
		t.Error("error parity: expected Err for non-UTF-8 String input")
	}
}

func TestBytesJSONDecodeError(t *testing.T) {
	if !Bytes("{invalid").Decode().JSON().IsErr() {
		t.Error("expected error for malformed JSON input")
	}

	if !Bytes("123").Decode().JSON().IsErr() {
		t.Error("expected error for non-string JSON input")
	}
}

func TestBytesURLRoundtripParity(t *testing.T) {
	for _, tt := range parityCases {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.data.Encode().URL()

			strEnc := tt.data.String().Encode().URL()
			if !enc.Eq(strEnc.Bytes()) {
				t.Errorf("encode parity: Bytes = %q, String = %q", enc, strEnc)
			}

			dec := enc.Decode().URL().Unwrap()
			if !dec.Eq(tt.data) {
				t.Errorf("roundtrip = %q, want %q", dec, tt.data)
			}

			strDec := strEnc.Decode().URL().Unwrap()
			if !dec.Eq(strDec.Bytes()) {
				t.Errorf("decode parity: Bytes = %q, String = %q", dec, strDec)
			}
		})
	}
}

func TestBytesURLBinaryGarbage(t *testing.T) {
	enc := binaryGarbage.Encode().URL()

	dec := enc.Decode().URL().Unwrap()
	if !dec.Eq(binaryGarbage) {
		t.Errorf("roundtrip = %v, want %v", dec, binaryGarbage)
	}
}

func TestBytesURLSafeParity(t *testing.T) {
	data := Bytes("a b/c d")

	enc := data.Encode().URL(Bytes(" "))
	if enc.Ne(Bytes("a b%2Fc d")) {
		t.Errorf("encode with safe = %q, want %q", enc, "a b%2Fc d")
	}

	strEnc := data.String().Encode().URL(String(" "))
	if !enc.Eq(strEnc.Bytes()) {
		t.Errorf("safe parity: Bytes = %q, String = %q", enc, strEnc)
	}
}

func TestBytesURLDecodeError(t *testing.T) {
	if !Bytes("%zz").Decode().URL().IsErr() {
		t.Error("expected error for invalid percent-encoding")
	}

	if !Bytes("%f").Decode().URL().IsErr() {
		t.Error("expected error for truncated percent-encoding")
	}
}

func TestBytesHTMLRoundtripParity(t *testing.T) {
	for _, tt := range parityCases {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.data.Encode().HTML()

			strEnc := tt.data.String().Encode().HTML()
			if !enc.Eq(strEnc.Bytes()) {
				t.Errorf("encode parity: Bytes = %q, String = %q", enc, strEnc)
			}

			dec := enc.Decode().HTML()
			if !dec.Eq(tt.data) {
				t.Errorf("roundtrip = %q, want %q", dec, tt.data)
			}

			strDec := strEnc.Decode().HTML()
			if !dec.Eq(strDec.Bytes()) {
				t.Errorf("decode parity: Bytes = %q, String = %q", dec, strDec)
			}
		})
	}
}

func TestBytesHTMLBinaryGarbage(t *testing.T) {
	enc := binaryGarbage.Encode().HTML()

	dec := enc.Decode().HTML()
	if !dec.Eq(binaryGarbage) {
		t.Errorf("roundtrip = %v, want %v", dec, binaryGarbage)
	}
}

func TestBytesRot13RoundtripParity(t *testing.T) {
	for _, tt := range parityCases {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.data.Encode().Rot13()

			strEnc := tt.data.String().Encode().Rot13()
			if !enc.Eq(strEnc.Bytes()) {
				t.Errorf("encode parity: Bytes = %q, String = %q", enc, strEnc)
			}

			dec := enc.Decode().Rot13()
			if !dec.Eq(tt.data) {
				t.Errorf("roundtrip = %q, want %q", dec, tt.data)
			}
		})
	}
}

func TestBytesRot13Known(t *testing.T) {
	enc := Bytes("Hello").Encode().Rot13()
	if enc.Ne(Bytes("Uryyb")) {
		t.Errorf("Rot13(Hello) = %q, want %q", enc, "Uryyb")
	}
}

func TestBytesRot13BinaryGarbage(t *testing.T) {
	enc := binaryGarbage.Encode().Rot13()

	dec := enc.Decode().Rot13()
	if !dec.Eq(binaryGarbage) {
		t.Errorf("roundtrip = %v, want %v", dec, binaryGarbage)
	}
}

func TestBytesOctalRoundtrip(t *testing.T) {
	cases := append(parityCases, struct {
		name string
		data Bytes
	}{"binary", binaryGarbage})

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			enc := tt.data.Encode().Octal()

			dec := enc.Decode().Octal().Unwrap()
			if !dec.Eq(tt.data) {
				t.Errorf("roundtrip = %v, want %v", dec, tt.data)
			}
		})
	}
}

func TestBytesOctalASCIIParity(t *testing.T) {
	// Octal parity holds only for ASCII: String encodes code points, Bytes encodes bytes.
	data := Bytes("Hello, World! 123")

	enc := data.Encode().Octal()

	strEnc := data.String().Encode().Octal()
	if !enc.Eq(strEnc.Bytes()) {
		t.Errorf("encode parity: Bytes = %q, String = %q", enc, strEnc)
	}

	dec := enc.Decode().Octal().Unwrap()

	strDec := strEnc.Decode().Octal().Unwrap()
	if !dec.Eq(strDec.Bytes()) {
		t.Errorf("decode parity: Bytes = %q, String = %q", dec, strDec)
	}
}

func TestBytesOctalUnicodeDivergence(t *testing.T) {
	// For multibyte runes the representations must differ by design:
	// 'ж' is one code point (octal 1066) but two bytes (320 266).
	data := Bytes("ж")

	enc := data.Encode().Octal()
	if enc.Ne(Bytes("320 266")) {
		t.Errorf("byte-wise octal = %q, want %q", enc, "320 266")
	}

	strEnc := data.String().Encode().Octal()
	if enc.Eq(strEnc.Bytes()) {
		t.Error("expected byte-wise Bytes octal to differ from rune-wise String octal")
	}
}

func TestBytesOctalDecodeError(t *testing.T) {
	cases := []struct {
		name string
		data Bytes
	}{
		{"non-octal digit", Bytes("8")},
		{"out of byte range", Bytes("777")},
		{"garbage token", Bytes("abc")},
		{"empty token", Bytes("101  102")},
		{"trailing space", Bytes("101 ")},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.data.Decode().Octal().IsErr() {
				t.Errorf("expected error for %q", tt.data)
			}
		})
	}
}
