package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestGobEncodingDecoding(t *testing.T) {
	file, err := NewFile("testfile.gob").CreateTemp().Result()
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}

	defer file.Remove()

	// Encode data to the file
	dataToEncode := SliceOf(1, 2, 3, 4, 5)
	result := file.Encode().Gob(dataToEncode)
	if result.IsErr() {
		t.Fatalf("Gob encoding failed: %v", result.Err())
	}

	// Decode data from the file
	var decodedData Slice[int]
	result = file.Decode().Gob(&decodedData)
	if result.IsErr() {
		t.Fatalf("Gob decoding failed: %v", result.Err())
	}

	if dataToEncode.Ne(decodedData) {
		t.Errorf("Decoded data does not match the original data.")
	}
}

func TestJSONDecodeDuplicateKeyError(t *testing.T) {
	file, err := NewFile("dupkey.json").CreateTemp().Result()
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}

	defer file.Remove()

	// encoding/json/v2 rejects duplicate object member names.
	if r := file.Write(`{"a":1,"a":2}`); r.IsErr() {
		t.Fatalf("Failed to write test data: %v", r.Err())
	}

	var data Map[String, Int]
	result := file.Decode().JSON(&data)
	if !result.IsErr() {
		t.Error("expected Err for duplicate object keys, got Ok")
	}
}

func TestJSONEncodeStructV2Defaults(t *testing.T) {
	file, err := NewFile("structdefaults.json").CreateTemp().Result()
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}

	defer file.Remove()

	// encoding/json/v2 marshals nil slices as [] and nil maps as {},
	// and appends no trailing newline.
	type payload struct {
		S []int          `json:"s"`
		M map[string]int `json:"m"`
	}

	if r := file.Encode().JSON(payload{}); r.IsErr() {
		t.Fatalf("JSON encoding failed: %v", r.Err())
	}

	content := file.Read()
	if content.IsErr() {
		t.Fatalf("Failed to read file: %v", content.Err())
	}

	expected := String(`{"s":[],"m":{}}`)
	if content.Ok().Ne(expected) {
		t.Errorf("content = %q, want %q", content.Ok(), expected)
	}
}

func TestJSONDecodeCaseSensitive(t *testing.T) {
	file, err := NewFile("casesensitive.json").CreateTemp().Result()
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}

	defer file.Remove()

	if r := file.Write(`{"name":"go"}`); r.IsErr() {
		t.Fatalf("Failed to write test data: %v", r.Err())
	}

	// encoding/json/v2 matches struct fields case-sensitively:
	// "name" must not populate the field tagged "Name".
	var data struct {
		Name string `json:"Name"`
	}

	result := file.Decode().JSON(&data)
	if result.IsErr() {
		t.Fatalf("JSON decoding failed: %v", result.Err())
	}

	if data.Name != "" {
		t.Errorf("expected case-sensitive matching to skip \"name\", got %q", data.Name)
	}
}

func TestJSONEncodingDecoding(t *testing.T) {
	file, err := NewFile("testfile.json").CreateTemp().Result()
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}

	defer file.Remove()

	// Encode data to the file
	dataToEncode := SliceOf(1, 2, 3, 4, 5)
	result := file.Encode().JSON(dataToEncode)
	if result.IsErr() {
		t.Fatalf("JSON encoding failed: %v", result.Err())
	}

	// Decode data from the file
	var decodedData Slice[int]
	result = file.Decode().JSON(&decodedData)
	if result.IsErr() {
		t.Fatalf("JSON decoding failed: %v", result.Err())
	}

	if dataToEncode.Ne(decodedData) {
		t.Errorf("Decoded data does not match the original data.")
	}
}
