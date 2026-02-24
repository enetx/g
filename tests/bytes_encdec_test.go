package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestBytesBase64Encode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want String
	}{
		{"empty", Bytes(""), ""},
		{"hello", Bytes("hello"), "aGVsbG8="},
		{"hello world", Bytes("hello world"), "aGVsbG8gd29ybGQ="},
		{"binary", Bytes{0xff, 0x00, 0xab}, "/wCr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Encode().Base64(); got != tt.want {
				t.Errorf("Bytes.Encode().Base64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesBase64Decode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want Bytes
	}{
		{"hello world", Bytes("aGVsbG8gd29ybGQ="), Bytes("hello world")},
		{"empty", Bytes(""), Bytes("")},
		{"binary", Bytes("/wCr"), Bytes{0xff, 0x00, 0xab}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.Decode().Base64().Unwrap()
			if !got.Eq(tt.want) {
				t.Errorf("Bytes.Decode().Base64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesBase64DecodeError(t *testing.T) {
	result := Bytes("!!!invalid!!!").Decode().Base64()
	if !result.IsErr() {
		t.Error("expected error for invalid base64 input")
	}
}

func TestBytesBase64RawEncode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want String
	}{
		{"empty", Bytes(""), ""},
		{"hello", Bytes("hello"), "aGVsbG8"},
		{"hello world", Bytes("hello world"), "aGVsbG8gd29ybGQ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Encode().Base64Raw(); got != tt.want {
				t.Errorf("Bytes.Encode().Base64Raw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesBase64RawDecode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want Bytes
	}{
		{"hello world", Bytes("aGVsbG8gd29ybGQ"), Bytes("hello world")},
		{"empty", Bytes(""), Bytes("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.Decode().Base64Raw().Unwrap()
			if !got.Eq(tt.want) {
				t.Errorf("Bytes.Decode().Base64Raw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesBase64URLEncode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want String
	}{
		{"empty", Bytes(""), ""},
		{"hello", Bytes("hello"), "aGVsbG8="},
		{"url unsafe", Bytes{0xfb, 0xf0, 0x3f}, "-_A_"},
		{"binary", Bytes{0xff, 0x00, 0xab}, "_wCr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Encode().Base64URL(); got != tt.want {
				t.Errorf("Bytes.Encode().Base64URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesBase64URLDecode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want Bytes
	}{
		{"hello", Bytes("aGVsbG8="), Bytes("hello")},
		{"url unsafe", Bytes("-_A_"), Bytes{0xfb, 0xf0, 0x3f}},
		{"empty", Bytes(""), Bytes("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.Decode().Base64URL().Unwrap()
			if !got.Eq(tt.want) {
				t.Errorf("Bytes.Decode().Base64URL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesBase64RawURLEncode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want String
	}{
		{"empty", Bytes(""), ""},
		{"hello", Bytes("hello"), "aGVsbG8"},
		{"binary", Bytes{0xff, 0x00, 0xab}, "_wCr"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Encode().Base64RawURL(); got != tt.want {
				t.Errorf("Bytes.Encode().Base64RawURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesBase64RawURLDecode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want Bytes
	}{
		{"hello", Bytes("aGVsbG8"), Bytes("hello")},
		{"binary", Bytes("_wCr"), Bytes{0xff, 0x00, 0xab}},
		{"empty", Bytes(""), Bytes("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.Decode().Base64RawURL().Unwrap()
			if !got.Eq(tt.want) {
				t.Errorf("Bytes.Decode().Base64RawURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesBase64RoundTrip(t *testing.T) {
	input := Bytes("Hello, World! 🌍")

	if got := input.Encode().Base64().BytesUnsafe().Decode().Base64().Unwrap(); !got.Eq(input) {
		t.Error("Base64 round trip failed")
	}

	if got := input.Encode().Base64Raw().BytesUnsafe().Decode().Base64Raw().Unwrap(); !got.Eq(input) {
		t.Error("Base64Raw round trip failed")
	}

	if got := input.Encode().Base64URL().BytesUnsafe().Decode().Base64URL().Unwrap(); !got.Eq(input) {
		t.Error("Base64URL round trip failed")
	}

	if got := input.Encode().Base64RawURL().BytesUnsafe().Decode().Base64RawURL().Unwrap(); !got.Eq(input) {
		t.Error("Base64RawURL round trip failed")
	}
}

func TestBytesBase64URLvsStd(t *testing.T) {
	input := Bytes{0xfb, 0xf0, 0x3f}

	std := input.Encode().Base64()
	url := input.Encode().Base64URL()

	if !std.Contains("+") && !std.Contains("/") {
		t.Skip("test input doesn't produce +/ in std base64")
	}

	if url.Contains("+") || url.Contains("/") {
		t.Errorf("Base64URL should not contain + or /, got %s", url)
	}
}

func TestBytesBase64DecodeErrors(t *testing.T) {
	invalid := Bytes("!!!invalid!!!")

	if result := invalid.Decode().Base64(); !result.IsErr() {
		t.Error("Base64: expected error for invalid input")
	}

	if result := invalid.Decode().Base64Raw(); !result.IsErr() {
		t.Error("Base64Raw: expected error for invalid input")
	}

	if result := invalid.Decode().Base64URL(); !result.IsErr() {
		t.Error("Base64URL: expected error for invalid input")
	}

	if result := invalid.Decode().Base64RawURL(); !result.IsErr() {
		t.Error("Base64RawURL: expected error for invalid input")
	}
}

func TestBytesHexEncode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want Bytes
	}{
		{"hello", Bytes("hello"), Bytes("68656c6c6f")},
		{"empty", Bytes(""), Bytes("")},
		{"binary", Bytes{0x00, 0xff, 0x0a}, Bytes("00ff0a")},
		{"single byte", Bytes{0xab}, Bytes("ab")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Encode().Hex(); !got.Eq(tt.want) {
				t.Errorf("Bytes.Encode().Hex() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestBytesHexDecode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want Bytes
	}{
		{"hello", Bytes("68656c6c6f"), Bytes("hello")},
		{"empty", Bytes(""), Bytes("")},
		{"binary", Bytes("00ff0a"), Bytes{0x00, 0xff, 0x0a}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.Decode().Hex().Unwrap()
			if !got.Eq(tt.want) {
				t.Errorf("Bytes.Decode().Hex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesHexDecodeError(t *testing.T) {
	result := Bytes("zzzz").Decode().Hex()
	if !result.IsErr() {
		t.Error("expected error for invalid hex input")
	}
}

func TestBytesHexRoundTrip(t *testing.T) {
	input := Bytes{0x00, 0x01, 0xfe, 0xff}
	got := input.Encode().Hex().Decode().Hex().Unwrap()

	if !got.Eq(input) {
		t.Errorf("Hex round trip failed: got %v, want %v", got, input)
	}
}

func TestBytesXOREncode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		key  Bytes
		want Bytes
	}{
		{"simple", Bytes{0x01, 0x02, 0x03}, Bytes{0xff}, Bytes{0xfe, 0xfd, 0xfc}},
		{"empty key", Bytes{0x01, 0x02}, Bytes{}, Bytes{0x01, 0x02}},
		{"same key", Bytes("hello"), Bytes("hello"), Bytes{0, 0, 0, 0, 0}},
		{"key wrap", Bytes{0x01, 0x02, 0x03, 0x04}, Bytes{0xff, 0x00}, Bytes{0xfe, 0x02, 0xfc, 0x04}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.Encode().XOR(tt.key)
			if !got.Eq(tt.want) {
				t.Errorf("Bytes.Encode().XOR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesXORRoundTrip(t *testing.T) {
	input := Bytes("Hello, World!")
	key := Bytes("secret")

	encrypted := input.Encode().XOR(key)
	decrypted := encrypted.Decode().XOR(key)

	if !decrypted.Eq(input) {
		t.Errorf("XOR round trip failed: got %v, want %v", decrypted, input)
	}
}

func TestBytesXORNoMutation(t *testing.T) {
	input := Bytes{0x01, 0x02, 0x03}
	original := input.Clone()
	key := Bytes{0xff}

	input.Encode().XOR(key)

	if !input.Eq(original) {
		t.Error("XOR should not mutate original bytes")
	}
}

func TestBytesBinaryEncode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want String
	}{
		{"hello", Bytes("H"), "01001000"},
		{"zero", Bytes{0x00}, "00000000"},
		{"ff", Bytes{0xff}, "11111111"},
		{"multi", Bytes{0x00, 0xff}, "0000000011111111"},
		{"empty", Bytes(""), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.Encode().Binary(); got != tt.want {
				t.Errorf("Bytes.Encode().Binary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesBinaryDecode(t *testing.T) {
	tests := []struct {
		name string
		data Bytes
		want Bytes
	}{
		{"H", Bytes("01001000"), Bytes("H")},
		{"zero", Bytes("00000000"), Bytes{0x00}},
		{"ff", Bytes("11111111"), Bytes{0xff}},
		{"multi", Bytes("0000000011111111"), Bytes{0x00, 0xff}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.Decode().Binary().Unwrap()
			if !got.Eq(tt.want) {
				t.Errorf("Bytes.Decode().Binary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesBinaryDecodeErrors(t *testing.T) {
	// invalid length (not multiple of 8)
	if result := Bytes("01001").Decode().Binary(); !result.IsErr() {
		t.Error("expected error for invalid binary length")
	}

	// invalid digit
	if result := Bytes("0100100x").Decode().Binary(); !result.IsErr() {
		t.Error("expected error for invalid binary digit")
	}
}

func TestBytesBinaryRoundTrip(t *testing.T) {
	input := Bytes("Hello!")
	got := input.Encode().Binary().BytesUnsafe().Decode().Binary().Unwrap()

	if !got.Eq(input) {
		t.Errorf("Binary round trip failed: got %v, want %v", got, input)
	}
}
