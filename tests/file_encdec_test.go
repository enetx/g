package g_test

import (
	"testing"

	"github.com/enetx/g"
)

func TestGobEncodingDecoding(t *testing.T) {
	file, err := g.NewFile("testfile.gob").CreateTemp().Result()
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}

	defer file.Remove()

	// Encode data to the file
	dataToEncode := g.SliceOf(1, 2, 3, 4, 5)
	result := file.Enc().Gob(dataToEncode)
	if result.IsErr() {
		t.Fatalf("Gob encoding failed: %v", result.Err())
	}

	// Decode data from the file
	var decodedData g.Slice[int]
	result = file.Dec().Gob(&decodedData)
	if result.IsErr() {
		t.Fatalf("Gob decoding failed: %v", result.Err())
	}

	if dataToEncode.Ne(decodedData) {
		t.Errorf("Decoded data does not match the original data.")
	}
}

func TestJSONEncodingDecoding(t *testing.T) {
	file, err := g.NewFile("testfile.json").CreateTemp().Result()
	if err != nil {
		t.Fatalf("Failed to create a temporary file: %v", err)
	}

	defer file.Remove()

	// Encode data to the file
	dataToEncode := g.SliceOf(1, 2, 3, 4, 5)
	result := file.Enc().JSON(dataToEncode)
	if result.IsErr() {
		t.Fatalf("JSON encoding failed: %v", result.Err())
	}

	// Decode data from the file
	var decodedData g.Slice[int]
	result = file.Dec().JSON(&decodedData)
	if result.IsErr() {
		t.Fatalf("JSON decoding failed: %v", result.Err())
	}

	if dataToEncode.Ne(decodedData) {
		t.Errorf("Decoded data does not match the original data.")
	}
}
