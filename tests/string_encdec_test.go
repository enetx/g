package g_test

import (
	"reflect"
	"testing"

	"github.com/enetx/g"
)

func TestStringJSON(t *testing.T) {
	// Test case 1: Encoding a struct
	type Person struct {
		Name string
		Age  int
	}

	person := Person{Name: "John", Age: 30}
	expectedResult1 := g.String(`{"Name":"John","Age":30}`)
	result1 := g.String("").Enc().JSON(person).Ok()
	if !result1.Eq(expectedResult1) {
		t.Errorf("Test case 1 failed: Expected result is %v, got %v", expectedResult1, result1)
	}

	// Test case 2: Encoding a map
	data2 := map[string]any{"Name": "Alice", "Age": 25}
	expectedResult2 := g.String(`{"Age":25,"Name":"Alice"}`)
	result2 := g.String("").Enc().JSON(data2).Ok()
	if !result2.Eq(expectedResult2) {
		t.Errorf("Test case 2 failed: Expected result is %v, got %v", expectedResult2, result2)
	}

	// Test case 3: Encoding an array
	data3 := []int{1, 2, 3}
	expectedResult3 := g.String("[1,2,3]")
	result3 := g.String("").Enc().JSON(data3).Ok()
	if !result3.Eq(expectedResult3) {
		t.Errorf("Test case 3 failed: Expected result is %v, got %v", expectedResult3, result3)
	}

	// Test case 4: Encoding a nil value
	expectedResult4 := g.String("null")
	result4 := g.String("").Enc().JSON(nil).Ok()
	if !result4.Eq(expectedResult4) {
		t.Errorf("Test case 4 failed: Expected result is %v, got %v", expectedResult4, result4)
	}
}

func TestStringJSONDecode(t *testing.T) {
	// Test case 1: Decoding a valid JSON string into a struct
	type Person struct {
		Name string
		Age  int
	}

	inputJSON1 := `{"Name":"John","Age":30}`
	var person1 Person
	expectedResult1 := g.String(inputJSON1)
	result1 := g.String(inputJSON1).Dec().JSON(&person1).Ok()
	if !result1.Eq(expectedResult1) {
		t.Errorf("Test case 1 failed: Expected result is %v, got %v", expectedResult1, result1)
	}
	expectedPerson1 := Person{Name: "John", Age: 30}
	if !reflect.DeepEqual(person1, expectedPerson1) {
		t.Errorf("Test case 1 failed: Decoded struct is not equal to expected struct")
	}

	// Test case 2: Decoding a valid JSON string into a map
	inputJSON2 := `{"Name":"Alice","Age":25}`
	var data2 map[string]any
	expectedResult2 := g.String(inputJSON2)
	result2 := g.String(inputJSON2).Dec().JSON(&data2).Ok()
	if !result2.Eq(expectedResult2) {
		t.Errorf("Test case 2 failed: Expected result is %v, got %v", expectedResult2, result2)
	}
	expectedData2 := map[string]any{"Name": "Alice", "Age": float64(25)}
	if !reflect.DeepEqual(data2, expectedData2) {
		t.Errorf("Test case 2 failed: Decoded map is not equal to expected map")
	}
}

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

func TestStringRot13Encoding(t *testing.T) {
	// Test case 1: Encoding uppercase letters
	inputData1 := g.NewString("HELLO")
	expectedEncoded1 := g.NewString("URYYB")
	result1 := inputData1.Enc().Rot13()
	if !result1.Eq(expectedEncoded1) {
		t.Errorf("Test case 1 failed: Expected encoded string is %s, got %s", expectedEncoded1, result1)
	}

	// Test case 2: Encoding lowercase letters
	inputData2 := g.NewString("hello")
	expectedEncoded2 := g.NewString("uryyb")
	result2 := inputData2.Enc().Rot13()
	if !result2.Eq(expectedEncoded2) {
		t.Errorf("Test case 2 failed: Expected encoded string is %s, got %s", expectedEncoded2, result2)
	}

	// Test case 3: Encoding mixed case letters
	inputData3 := g.NewString("Hello, World!")
	expectedEncoded3 := g.NewString("Uryyb, Jbeyq!")
	result3 := inputData3.Enc().Rot13()
	if !result3.Eq(expectedEncoded3) {
		t.Errorf("Test case 3 failed: Expected encoded string is %s, got %s", expectedEncoded3, result3)
	}

	// Test case 4: Encoding non-alphabetic characters
	inputData4 := g.NewString("12345 !@#$")
	expectedEncoded4 := g.NewString("12345 !@#$")
	result4 := inputData4.Enc().Rot13()
	if !result4.Eq(expectedEncoded4) {
		t.Errorf("Test case 4 failed: Expected encoded string is %s, got %s", expectedEncoded4, result4)
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

func TestStringBinaryEncodingAndDecoding(t *testing.T) {
	// Test case 1: Encoding a string to binary
	inputStr1 := g.NewString("Hello")
	expectedBinary1 := g.NewString("0100100001100101011011000110110001101111")
	result1 := inputStr1.Enc().Binary()
	if !result1.Eq(expectedBinary1) {
		t.Errorf("Test case 1 failed: Expected binary string is %s, got %s", expectedBinary1, result1)
	}

	// Test case 2: Decoding a binary string back to the original string
	inputBinary2 := g.NewString("0100100001100101011011000110110001101111")
	expectedStr2 := g.NewString("Hello")
	result2 := inputBinary2.Dec().Binary()
	if result2.IsErr() {
		t.Errorf("Test case 2 failed: Error occurred during decoding: %v", result2.Err())
	} else if result2.Ok().Ne(expectedStr2) {
		t.Errorf("Test case 2 failed: Expected decoded string is %s, got %s", expectedStr2, result2.Ok())
	}
}

func TestXMLEncodingAndDecoding(t *testing.T) {
	// Define a struct to represent data for testing XML encoding and decoding
	type Person struct {
		Name  string `xml:"name"`
		Age   int    `xml:"age"`
		City  string `xml:"city"`
		Email string `xml:"email"`
	}

	// Test case 1: Encoding data to XML
	inputData1 := Person{Name: "John", Age: 30, City: "New York", Email: "john@example.com"}
	expectedXML1 := g.NewString(
		"<Person><name>John</name><age>30</age><city>New York</city><email>john@example.com</email></Person>",
	)
	result1 := g.NewString("").Enc().XML(inputData1)
	if !result1.Ok().Eq(expectedXML1) {
		t.Errorf("Test case 1 failed: Expected XML is %s, got %s", expectedXML1, result1.Ok())
	}

	// Test case 2: Decoding XML back to data
	xmlData2 := g.NewString(
		"<Person><name>Alice</name><age>25</age><city>London</city><email>alice@example.com</email></Person>",
	)
	var decodedData2 Person
	expectedData2 := Person{Name: "Alice", Age: 25, City: "London", Email: "alice@example.com"}
	result2 := xmlData2.Dec().XML(&decodedData2)
	if result2.IsErr() {
		t.Errorf("Test case 2 failed: Error occurred during XML decoding: %v", result2.Err())
	} else if decodedData2 != expectedData2 {
		t.Errorf("Test case 2 failed: Expected decoded data is %+v, got %+v", expectedData2, decodedData2)
	}
}

func TestStringHTMLEncodingAndDecoding(t *testing.T) {
	// Test case 1: Encoding HTML
	inputData1 := g.NewString("<p>Hello, <b>World</b>!</p>")
	expectedEncoded1 := g.NewString("&lt;p&gt;Hello, &lt;b&gt;World&lt;/b&gt;!&lt;/p&gt;")
	result1 := inputData1.Enc().HTML()
	if !result1.Eq(expectedEncoded1) {
		t.Errorf("Test case 1 failed: Expected encoded HTML is %s, got %s", expectedEncoded1, result1)
	}

	// Test case 2: Decoding HTML
	htmlData2 := g.NewString("&lt;a href=&quot;https://example.com&quot;&gt;Link&lt;/a&gt;")
	expectedDecoded2 := g.NewString("<a href=\"https://example.com\">Link</a>")
	result2 := htmlData2.Dec().HTML()
	if !result2.Eq(expectedDecoded2) {
		t.Errorf("Test case 2 failed: Expected decoded HTML is %s, got %s", expectedDecoded2, result2)
	}
}
