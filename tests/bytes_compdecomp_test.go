package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestBytesCompressionAndDecompression(t *testing.T) {
	inputData := Bytes("hello world")

	// Zstd
	zstdCompressed := inputData.Compress().Zstd()
	zstdDecompressed := zstdCompressed.Decompress().Zstd()
	if zstdDecompressed.IsErr() || zstdDecompressed.Unwrap().Ne(inputData) {
		t.Errorf("Zstd round-trip failed. Input: %s, Decompressed: %s", inputData, zstdDecompressed.Unwrap())
	}

	// Brotli
	brotliCompressed := inputData.Compress().Brotli()
	brotliDecompressed := brotliCompressed.Decompress().Brotli()
	if brotliDecompressed.IsErr() || brotliDecompressed.Unwrap().Ne(inputData) {
		t.Errorf("Brotli round-trip failed. Input: %s, Decompressed: %s", inputData, brotliDecompressed.Unwrap())
	}

	// Zlib
	zlibCompressed := inputData.Compress().Zlib()
	zlibDecompressed := zlibCompressed.Decompress().Zlib()
	if zlibDecompressed.IsErr() || zlibDecompressed.Unwrap().Ne(inputData) {
		t.Errorf("Zlib round-trip failed. Input: %s, Decompressed: %s", inputData, zlibDecompressed.Unwrap())
	}

	// Gzip
	gzipCompressed := inputData.Compress().Gzip()
	gzipDecompressed := gzipCompressed.Decompress().Gzip()
	if gzipDecompressed.IsErr() || gzipDecompressed.Unwrap().Ne(inputData) {
		t.Errorf("Gzip round-trip failed. Input: %s, Decompressed: %s", inputData, gzipDecompressed.Unwrap())
	}

	// Flate (default level)
	flateCompressed := inputData.Compress().Flate()
	flateDecompressed := flateCompressed.Decompress().Flate()
	if flateDecompressed.IsErr() || flateDecompressed.Unwrap().Ne(inputData) {
		t.Errorf("Flate round-trip failed. Input: %s, Decompressed: %s", inputData, flateDecompressed.Unwrap())
	}

	// Flate (explicit level)
	flateLvl := inputData.Compress().Flate(9)
	flateLvlDecompressed := flateLvl.Decompress().Flate()
	if flateLvlDecompressed.IsErr() || flateLvlDecompressed.Unwrap().Ne(inputData) {
		t.Errorf("Flate(9) round-trip failed. Input: %s, Decompressed: %s", inputData, flateLvlDecompressed.Unwrap())
	}
}

func TestBytesDecompressErrorPropagates(t *testing.T) {
	// Garbage input should surface an Err, not panic.
	garbage := Bytes("not a valid gzip stream")
	if r := garbage.Decompress().Gzip(); r.IsOk() {
		t.Error("Expected Err decompressing garbage gzip data")
	}
}
