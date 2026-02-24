package g

import (
	"encoding/base64"
	"encoding/hex"
)

type (
	// A struct that wraps Bytes for encoding.
	bencode struct{ bytes Bytes }

	// A struct that wraps Bytes for decoding.
	bdecode struct{ bytes Bytes }
)

// Encode returns a bencode struct wrapping the given Bytes.
func (bs Bytes) Encode() bencode { return bencode{bs} }

// Decode returns a bdecode struct wrapping the given Bytes.
func (bs Bytes) Decode() bdecode { return bdecode{bs} }

// Base64 encodes the wrapped Bytes using standard Base64 (with padding).
func (e bencode) Base64() String {
	return String(base64.StdEncoding.EncodeToString(e.bytes))
}

// Base64Raw encodes the wrapped Bytes using standard Base64 without padding.
func (e bencode) Base64Raw() String {
	return String(base64.RawStdEncoding.EncodeToString(e.bytes))
}

// Base64URL encodes the wrapped Bytes using URL-safe Base64 (with padding).
func (e bencode) Base64URL() String {
	return String(base64.URLEncoding.EncodeToString(e.bytes))
}

// Base64RawURL encodes the wrapped Bytes using URL-safe Base64 without padding.
func (e bencode) Base64RawURL() String {
	return String(base64.RawURLEncoding.EncodeToString(e.bytes))
}

// Base64 decodes the wrapped Bytes as standard Base64 (with padding) and returns Result[Bytes].
func (d bdecode) Base64() Result[Bytes] {
	out := make(Bytes, base64.StdEncoding.DecodedLen(len(d.bytes)))
	n, err := base64.StdEncoding.Decode(out, d.bytes)
	if err != nil {
		return Err[Bytes](err)
	}

	return Ok(out[:n])
}

// Base64Raw decodes the wrapped Bytes as standard Base64 without padding and returns Result[Bytes].
func (d bdecode) Base64Raw() Result[Bytes] {
	out := make(Bytes, base64.RawStdEncoding.DecodedLen(len(d.bytes)))
	n, err := base64.RawStdEncoding.Decode(out, d.bytes)
	if err != nil {
		return Err[Bytes](err)
	}

	return Ok(out[:n])
}

// Base64URL decodes the wrapped Bytes as URL-safe Base64 (with padding) and returns Result[Bytes].
func (d bdecode) Base64URL() Result[Bytes] {
	out := make(Bytes, base64.URLEncoding.DecodedLen(len(d.bytes)))
	n, err := base64.URLEncoding.Decode(out, d.bytes)
	if err != nil {
		return Err[Bytes](err)
	}

	return Ok(out[:n])
}

// Base64RawURL decodes the wrapped Bytes as URL-safe Base64 without padding and returns Result[Bytes].
func (d bdecode) Base64RawURL() Result[Bytes] {
	out := make(Bytes, base64.RawURLEncoding.DecodedLen(len(d.bytes)))
	n, err := base64.RawURLEncoding.Decode(out, d.bytes)
	if err != nil {
		return Err[Bytes](err)
	}

	return Ok(out[:n])
}

// Hex hex-encodes the wrapped Bytes and returns the result as Bytes.
func (e bencode) Hex() Bytes {
	out := make(Bytes, hex.EncodedLen(len(e.bytes)))
	hex.Encode(out, e.bytes)
	return out
}

// Hex hex-decodes the wrapped Bytes and returns the decoded result as Result[Bytes].
func (d bdecode) Hex() Result[Bytes] {
	out := make(Bytes, hex.DecodedLen(len(d.bytes)))
	n, err := hex.Decode(out, d.bytes)
	if err != nil {
		return Err[Bytes](err)
	}

	return Ok(out[:n])
}

// XOR encodes the wrapped Bytes using XOR cipher with the given key.
func (e bencode) XOR(key Bytes) Bytes {
	if len(key) == 0 {
		return e.bytes.Clone()
	}

	out := make(Bytes, len(e.bytes))
	for i, b := range e.bytes {
		out[i] = b ^ key[i%len(key)]
	}

	return out
}

// XOR decodes the wrapped Bytes using XOR cipher with the given key.
func (d bdecode) XOR(key Bytes) Bytes { return d.bytes.Encode().XOR(key) }

// Binary converts the wrapped Bytes to its binary representation as String.
func (e bencode) Binary() String {
	var b Builder
	b.Grow(e.bytes.Len() * 8)

	for _, c := range e.bytes {
		for bit := 7; bit >= 0; bit-- {
			b.WriteByte('0' + (c>>uint(bit))&1)
		}
	}

	return b.String()
}

// Binary converts the wrapped binary Bytes back to raw Bytes as Result[Bytes].
func (d bdecode) Binary() Result[Bytes] {
	if len(d.bytes)%8 != 0 {
		return Err[Bytes](ErrInvalidBinaryLength)
	}

	out := make(Bytes, 0, len(d.bytes)/8)

	for i := 0; i+8 <= len(d.bytes); i += 8 {
		var b byte
		for j := range 8 {
			c := d.bytes[i+j]
			if c != '0' && c != '1' {
				return Err[Bytes](ErrInvalidBinaryDigit)
			}
			b = b<<1 | (c - '0')
		}

		out = append(out, b)
	}

	return Ok(out)
}
