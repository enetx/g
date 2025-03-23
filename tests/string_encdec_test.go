package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

// Tests
func TestEncodeJSON(t *testing.T) {
	input := String("test")
	result := input.Encode().JSON()
	if result.IsErr() {
		t.Fatalf("Expected no error, got %v", result.Err())
	}

	expected := String("\"test\"")
	if result.Ok() != expected {
		t.Fatalf("Expected %s, got %s", expected, result.Ok())
	}
}

func TestDecodeJSON(t *testing.T) {
	input := String("\"test\"")
	result := input.Decode().JSON()
	if result.IsErr() {
		t.Fatalf("Expected no error, got %v", result.Err())
	}

	expected := "test"
	if result.Ok() != String(expected) {
		t.Fatalf("Expected %s, got %s", expected, result.Ok())
	}
}

func TestDecodeJSONError(t *testing.T) {
	input := String("invalid json")
	result := input.Decode().JSON()
	if !result.IsErr() {
		t.Fatal("Expected an error, but got none")
	}
}

func TestStringBase64Encode(t *testing.T) {
	tests := []struct {
		name string
		e    String
		want String
	}{
		{"empty", "", ""},
		{"hello", "hello", "aGVsbG8="},
		{"hello world", "hello world", "aGVsbG8gd29ybGQ="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Encode().Base64(); got != tt.want {
				t.Errorf("enc.Base64Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBase64Decode(t *testing.T) {
	tests := []struct {
		name string
		d    String
		want String
	}{
		{"base64 decode", "aGVsbG8gd29ybGQ=", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Decode().Base64().Unwrap(); got != tt.want {
				t.Errorf("dec.Base64Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringRot13Encoding(t *testing.T) {
	// Test case 1: Encoding uppercase letters
	inputData1 := String("HELLO")
	expectedEncoded1 := String("URYYB")
	result1 := inputData1.Encode().Rot13()
	if !result1.Eq(expectedEncoded1) {
		t.Errorf("Test case 1 failed: Expected encoded string is %s, got %s", expectedEncoded1, result1)
	}

	// Test case 2: Encoding lowercase letters
	inputData2 := String("hello")
	expectedEncoded2 := String("uryyb")
	result2 := inputData2.Encode().Rot13()
	if !result2.Eq(expectedEncoded2) {
		t.Errorf("Test case 2 failed: Expected encoded string is %s, got %s", expectedEncoded2, result2)
	}

	// Test case 3: Encoding mixed case letters
	inputData3 := String("Hello, World!")
	expectedEncoded3 := String("Uryyb, Jbeyq!")
	result3 := inputData3.Encode().Rot13()
	if !result3.Eq(expectedEncoded3) {
		t.Errorf("Test case 3 failed: Expected encoded string is %s, got %s", expectedEncoded3, result3)
	}

	// Test case 4: Encoding non-alphabetic characters
	inputData4 := String("12345 !@#$")
	expectedEncoded4 := String("12345 !@#$")
	result4 := inputData4.Encode().Rot13()
	if !result4.Eq(expectedEncoded4) {
		t.Errorf("Test case 4 failed: Expected encoded string is %s, got %s", expectedEncoded4, result4)
	}
}

func TestStringRot13Decoding(t *testing.T) {
	// Test cases for Rot13
	testCases := []struct {
		input    String
		expected string
	}{
		{"Uryyb", "Hello"},
		{"jbeyq", "world"},
		{"Grfg123", "Test123"},
		{"nopqrstuvwxyzabcdefghijklm", "abcdefghijklmnopqrstuvwxyz"},
		{"NOPQRSTUVWXYZABCDEFGHIJKLM", "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
	}

	for _, tc := range testCases {
		// Apply the Rot13 transformation
		result := tc.input.Decode().Rot13()
		// Check if the result matches the expected output
		if result.Std() != tc.expected {
			t.Errorf("Rot13(%s) returned %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestStringXOR(t *testing.T) {
	for range 100 {
		input := String("").Random(Int(30).RandomRange(100))
		key := String("").Random(10)
		obfuscated := input.Encode().XOR(key)
		deobfuscated := obfuscated.Decode().XOR(key)

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
		got := String(tt.input).Encode().XOR(String(tt.key))
		if got != String(tt.want) {
			t.Errorf("XOR(%q, %q) = %q; want %q", tt.input, tt.key, got, tt.want)
		}
	}
}

func TestGzFlateDecode(t *testing.T) {
	testCases := []struct {
		name     string
		input    String
		expected String
	}{
		{"Valid compressed data", "8kjNycnXUXCvcstJLElVBAAAAP//AQAA//8=", "Hello, GzFlate!"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.Decode().Base64().Unwrap().Decompress().Flate().Unwrap()
			if result.Ne(tc.expected) {
				t.Errorf("GzFlateDecode, expected: %s, got: %s", tc.expected, result)
			}
		})
	}
}

func TestGzFlateEncode(t *testing.T) {
	testCases := []struct {
		name     string
		input    String
		expected String
	}{
		{"Empty input", "", "AAAA//8BAAD//w=="},
		{"Valid input", "Hello, GzFlate!", "8kjNycnXUXCvcstJLElVBAAAAP//AQAA//8="},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.Compress().Flate().Encode().Base64()
			if result.Ne(tc.expected) {
				t.Errorf("GzFlateEncode, expected: %s, got: %s", tc.expected, result)
			}
		})
	}
}

func TestStringBinaryEncodingAndDecoding(t *testing.T) {
	// Test case 1: Encoding a string to binary
	inputStr1 := String("Hello")
	expectedBinary1 := String("0100100001100101011011000110110001101111")
	result1 := inputStr1.Encode().Binary()
	if !result1.Eq(expectedBinary1) {
		t.Errorf("Test case 1 failed: Expected binary string is %s, got %s", expectedBinary1, result1)
	}

	// Test case 2: Decoding a binary string back to the original string
	inputBinary2 := String("0100100001100101011011000110110001101111")
	expectedStr2 := String("Hello")
	result2 := inputBinary2.Decode().Binary()
	if result2.IsErr() {
		t.Errorf("Test case 2 failed: Error occurred during decoding: %v", result2.Err())
	} else if result2.Ok().Ne(expectedStr2) {
		t.Errorf("Test case 2 failed: Expected decoded string is %s, got %s", expectedStr2, result2.Ok())
	}
}

func TestStringHTMLEncodingAndDecoding(t *testing.T) {
	// Test case 1: Encoding HTML
	inputData1 := String("<p>Hello, <b>World</b>!</p>")
	expectedEncoded1 := String("&lt;p&gt;Hello, &lt;b&gt;World&lt;/b&gt;!&lt;/p&gt;")
	result1 := inputData1.Encode().HTML()
	if !result1.Eq(expectedEncoded1) {
		t.Errorf("Test case 1 failed: Expected encoded HTML is %s, got %s", expectedEncoded1, result1)
	}

	// Test case 2: Decoding HTML
	htmlData2 := String("&lt;a href=&quot;https://example.com&quot;&gt;Link&lt;/a&gt;")
	expectedDecoded2 := String("<a href=\"https://example.com\">Link</a>")
	result2 := htmlData2.Decode().HTML()
	if !result2.Eq(expectedDecoded2) {
		t.Errorf("Test case 2 failed: Expected decoded HTML is %s, got %s", expectedDecoded2, result2)
	}
}

func TestStringDecBase64_Success(t *testing.T) {
	// Input string encoded in Base64
	encodedStr := "SGVsbG8gV29ybGQh"

	// Create a dec instance wrapping the encoded string
	dec := String(encodedStr).Decode()

	// Decode the string using Base64
	decodedResult := dec.Base64()

	// Check if the result is successful
	if decodedResult.IsErr() {
		t.Errorf(
			"TestDec_Base64_Success: Expected decoding to be successful, but got an error: %v",
			decodedResult.Err(),
		)
	}

	// Check if the decoded string is correct
	expectedDecodedStr := "Hello World!"
	if decodedResult.Ok().Std() != expectedDecodedStr {
		t.Errorf(
			"TestDec_Base64_Success: Expected decoded string %s, got %s",
			expectedDecodedStr,
			decodedResult.Ok().Std(),
		)
	}
}

func TestDecStringBase64Failure(t *testing.T) {
	// Invalid Base64 encoded string (contains invalid character)
	invalidEncodedStr := "SGVsbG8gV29ybGQh==="

	// Create a dec instance wrapping the invalid encoded string
	dec := String(invalidEncodedStr).Decode()

	// Decode the string using Base64
	decodedResult := dec.Base64()

	// Check if the result is an error
	if !decodedResult.IsErr() {
		t.Errorf(
			"TestDec_Base64_Failure: Expected decoding to fail, but got a successful result: %s",
			decodedResult.Ok().Std(),
		)
	}
}

func TestStringHexEncode(t *testing.T) {
	// Test cases for Hex
	testCases := []struct {
		input    String
		expected String
	}{
		{"Hello", "48656c6c6f"},
		{"world", "776f726c64"},
		{"Test123", "54657374313233"},
		{"", ""},
	}

	for _, tc := range testCases {
		// Encode the string to hex using your package method
		result := tc.input.Encode().Hex()
		// Check if the result matches the expected output
		if result != tc.expected {
			t.Errorf("Hex(%s) returned %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestStringHexDec(t *testing.T) {
	// Test cases for Hex decoding
	testCases := []struct {
		input    String
		expected String
		err      error
	}{
		{"48656c6c6f20576f726c64", "Hello World", nil},
		{"74657374", "test", nil},
		{"", "", nil}, // Empty input should return empty output
	}

	for _, tc := range testCases {
		result := tc.input.Decode().Hex().Unwrap()
		if result != tc.expected {
			t.Errorf("Hex(%s) returned %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestStringOctalEnc(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    String
		expected String
	}{
		{"hello", "150 145 154 154 157"},
		{"world", "167 157 162 154 144"},
		{"", ""},
		{"123", "61 62 63"},
	}

	// Test each case
	for _, tc := range testCases {
		result := tc.input.Encode().Octal()
		if result != tc.expected {
			t.Errorf("Octal encoding is incorrect %s, exceted %s", result, tc.expected)
		}
	}
}

func TestStringOctalDec(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    String
		expected String
	}{
		{"150 145 154 154 157", "hello"},
		{"167 157 162 154 144", "world"},
		{"61 62 63", "123"},
	}

	// Test each case
	for _, tc := range testCases {
		result := tc.input.Decode().Octal()

		// Assert the result
		if result.IsErr() {
			t.Errorf("Error occurred during Octal decoding: %s", result.Err())
		}

		if result.Ok() != tc.expected {
			t.Errorf("Octal decoding is incorrect %s, exceted %s", result.Ok(), tc.expected)
		}
	}
}

func TestURLEncode(t *testing.T) {
	testCases := []struct {
		input    String
		safe     String
		expected String
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
		encoded := tc.input.Encode().URL(tc.safe)
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
		input    String
		expected String
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
		actual := test.input.Decode().URL().Unwrap()
		if actual != test.expected {
			t.Errorf("UnEscape(%s): expected %s, but got %s", test.input, test.expected, actual)
		}
	}
}
