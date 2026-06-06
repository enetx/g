package g

type (
	// A struct that wraps Bytes for compression.
	bcompress struct{ bytes Bytes }

	// A struct that wraps Bytes for decompression.
	bdecompress struct{ bytes Bytes }
)

// Compress returns a bcompress struct wrapping the given Bytes.
func (bs Bytes) Compress() bcompress { return bcompress{bs} }

// Decompress returns a bdecompress struct wrapping the given Bytes.
func (bs Bytes) Decompress() bdecompress { return bdecompress{bs} }

// Zstd compresses the wrapped Bytes using the zstd compression algorithm and
// returns the compressed data as Bytes.
func (c bcompress) Zstd() Bytes { return c.bytes.String().Compress().Zstd().Bytes() }

// Zstd decompresses the wrapped Bytes using the zstd compression algorithm and
// returns the decompressed data as a Result[Bytes].
func (d bdecompress) Zstd() Result[Bytes] {
	return resultStringToBytes(d.bytes.String().Decompress().Zstd())
}

// Brotli compresses the wrapped Bytes using the Brotli compression algorithm and
// returns the compressed data as Bytes.
func (c bcompress) Brotli() Bytes { return c.bytes.String().Compress().Brotli().Bytes() }

// Brotli decompresses the wrapped Bytes using the Brotli compression algorithm and
// returns the decompressed data as a Result[Bytes].
func (d bdecompress) Brotli() Result[Bytes] {
	return resultStringToBytes(d.bytes.String().Decompress().Brotli())
}

// Zlib compresses the wrapped Bytes using the zlib compression algorithm and
// returns the compressed data as Bytes.
func (c bcompress) Zlib() Bytes { return c.bytes.String().Compress().Zlib().Bytes() }

// Zlib decompresses the wrapped Bytes using the zlib compression algorithm and
// returns the decompressed data as a Result[Bytes].
func (d bdecompress) Zlib() Result[Bytes] {
	return resultStringToBytes(d.bytes.String().Decompress().Zlib())
}

// Gzip compresses the wrapped Bytes using the gzip compression format and
// returns the compressed data as Bytes.
func (c bcompress) Gzip() Bytes { return c.bytes.String().Compress().Gzip().Bytes() }

// Gzip decompresses the wrapped Bytes using the gzip compression format and
// returns the decompressed data as a Result[Bytes].
func (d bdecompress) Gzip() Result[Bytes] {
	return resultStringToBytes(d.bytes.String().Decompress().Gzip())
}

// Flate compresses the wrapped Bytes using the flate (zlib) compression algorithm
// and returns the compressed data as Bytes.
// It accepts an optional compression level. If no level is provided, it defaults to 7.
func (c bcompress) Flate(level ...int) Bytes {
	return c.bytes.String().Compress().Flate(level...).Bytes()
}

// Flate decompresses the wrapped Bytes using the flate (zlib) compression algorithm
// and returns the decompressed data as a Result[Bytes].
func (d bdecompress) Flate() Result[Bytes] {
	return resultStringToBytes(d.bytes.String().Decompress().Flate())
}

// resultStringToBytes converts a Result[String] into a Result[Bytes], preserving any error.
func resultStringToBytes(r Result[String]) Result[Bytes] {
	if r.IsErr() {
		return Err[Bytes](r.Err())
	}

	return Ok(r.Ok().Bytes())
}
