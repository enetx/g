package g_test

import (
	"testing"

	"github.com/enetx/g"
)

func TestStringCompressionAndDecompression(t *testing.T) {
	// Test data
	inputData := g.NewString("hello world")

	// Test Zlib compression and decompression
	zlibCompressed := inputData.Comp().Zlib()
	zlibDecompressed := zlibCompressed.Decomp().Zlib()
	if zlibDecompressed.IsErr() || zlibDecompressed.Unwrap().Ne(inputData) {
		t.Errorf(
			"Zlib compression and decompression failed. Input: %s, Decompressed: %s",
			inputData,
			zlibDecompressed.Unwrap(),
		)
	}

	// Test Gzip compression and decompression
	gzipCompressed := inputData.Comp().Gzip()
	gzipDecompressed := gzipCompressed.Decomp().Gzip()
	if gzipDecompressed.IsErr() || gzipDecompressed.Unwrap().Ne(inputData) {
		t.Errorf(
			"Gzip compression and decompression failed. Input: %s, Decompressed: %s",
			inputData,
			gzipDecompressed.Unwrap(),
		)
	}

	// Test Flate compression and decompression
	flateCompressed := inputData.Comp().Flate()
	flateDecompressed := flateCompressed.Decomp().Flate()
	if flateDecompressed.IsErr() || flateDecompressed.Unwrap().Ne(inputData) {
		t.Errorf(
			"Flate compression and decompression failed. Input: %s, Decompressed: %s",
			inputData,
			flateDecompressed.Unwrap(),
		)
	}
}
