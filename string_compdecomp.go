package g

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
)

type (
	// A struct that wraps a String for compression.
	comp struct{ str String }

	// A struct that wraps a String for decompression.
	decomp struct{ str String }
)

// Comp returns a comp struct wrapping the given String.
func (s String) Comp() comp { return comp{s} }

// Decomp returns a decomp struct wrapping the given String.
func (s String) Decomp() decomp { return decomp{s} }

// Zlib compresses the wrapped String using the zlib compression algorithm and
// returns the compressed data as a String.
func (c comp) Zlib() String {
	// gzcompress() php
	buffer := new(bytes.Buffer)
	writer := zlib.NewWriter(buffer)

	_, _ = writer.Write(c.str.ToBytes())
	_ = writer.Flush()
	_ = writer.Close()

	return String(buffer.String())
}

// Zlib decompresses the wrapped String using the zlib compression algorithm and
// returns the decompressed data as a Result[String].
func (d decomp) Zlib() Result[String] {
	// gzuncompress() php
	reader, err := zlib.NewReader(d.str.Reader())
	if err != nil {
		return Err[String](err)
	}

	defer reader.Close()

	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, reader); err != nil {
		return Err[String](err)
	}

	return Ok(String(buffer.String()))
}

// Gzip compresses the wrapped String using the gzip compression format and
// returns the compressed data as a String.
func (c comp) Gzip() String {
	// gzencode() php
	buffer := new(bytes.Buffer)
	writer := gzip.NewWriter(buffer)

	_, _ = writer.Write(c.str.ToBytes())
	_ = writer.Flush()
	_ = writer.Close()

	return String(buffer.String())
}

// Gzip decompresses the wrapped String using the gzip compression format and
// returns the decompressed data as a Result[String].
func (d decomp) Gzip() Result[String] {
	// gzdecode() php
	reader, err := gzip.NewReader(d.str.Reader())
	if err != nil {
		return Err[String](err)
	}

	defer reader.Close()

	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, reader); err != nil {
		return Err[String](err)
	}

	return Ok(String(buffer.String()))
}

// Flate compresses the wrapped String using the flate (zlib) compression algorithm
// and returns the compressed data as a String.
// It accepts an optional compression level. If no level is provided, it defaults to 7.
func (c comp) Flate(level ...int) String {
	// gzdeflate() php
	buffer := new(bytes.Buffer)

	l := 7
	if len(level) != 0 {
		l = level[0]
	}

	writer, _ := flate.NewWriter(buffer, l)

	_, _ = writer.Write(c.str.ToBytes())
	_ = writer.Flush()
	_ = writer.Close()

	return String(buffer.String())
}

// Flate decompresses the wrapped String using the flate (zlib) compression algorithm
// and returns the decompressed data as a Result[String].
func (d decomp) Flate() Result[String] {
	// gzinflate() php
	reader := flate.NewReader(d.str.Reader())
	defer reader.Close()

	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, reader); err != nil {
		return Err[String](err)
	}

	return Ok(String(buffer.String()))
}
