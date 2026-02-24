package g

import (
	"encoding/json"
	"html"
	"net/url"
	"strconv"
)

type (
	// A struct that wraps an String for encoding.
	encode struct{ str String }

	// A struct that wraps an String for decoding.
	decode struct{ str String }
)

// Encode returns an encode struct wrapping the given String.
func (s String) Encode() encode { return encode{s} }

// Decode returns a decode struct wrapping the given String.
func (s String) Decode() decode { return decode{s} }

// Base64 encodes the wrapped String using standard Base64 (with padding).
func (e encode) Base64() String { return e.str.BytesUnsafe().Encode().Base64() }

// Base64Raw encodes the wrapped String using standard Base64 without padding.
func (e encode) Base64Raw() String { return e.str.BytesUnsafe().Encode().Base64Raw() }

// Base64URL encodes the wrapped String using URL-safe Base64 (with padding).
func (e encode) Base64URL() String { return e.str.BytesUnsafe().Encode().Base64URL() }

// Base64RawURL encodes the wrapped String using URL-safe Base64 without padding.
func (e encode) Base64RawURL() String { return e.str.BytesUnsafe().Encode().Base64RawURL() }

// Base64 decodes the wrapped String using standard Base64 (with padding).
func (d decode) Base64() Result[String] {
	return TransformResult(d.str.BytesUnsafe().Decode().Base64(), Bytes.String)
}

// Base64Raw decodes the wrapped String using standard Base64 without padding.
func (d decode) Base64Raw() Result[String] {
	return TransformResult(d.str.BytesUnsafe().Decode().Base64Raw(), Bytes.String)
}

// Base64URL decodes the wrapped String using URL-safe Base64 (with padding).
func (d decode) Base64URL() Result[String] {
	return TransformResult(d.str.BytesUnsafe().Decode().Base64URL(), Bytes.String)
}

// Base64RawURL decodes the wrapped String using URL-safe Base64 without padding.
func (d decode) Base64RawURL() Result[String] {
	return TransformResult(d.str.BytesUnsafe().Decode().Base64RawURL(), Bytes.String)
}

// Hex hex-encodes the wrapped String and returns the encoded result as an String.
func (e encode) Hex() String { return e.str.BytesUnsafe().Encode().Hex().StringUnsafe() }

// Hex hex-decodes the wrapped String and returns the decoded result as Result[String].
func (d decode) Hex() Result[String] {
	return TransformResult(d.str.BytesUnsafe().Decode().Hex(), Bytes.String)
}

// XOR encodes the wrapped String using XOR cipher with the given key.
func (e encode) XOR(key String) String {
	return String(e.str.BytesUnsafe().Encode().XOR(key.BytesUnsafe()))
}

// XOR decodes the wrapped String using XOR cipher with the given key.
func (d decode) XOR(key String) String { return d.str.Encode().XOR(key) }

// Binary converts the wrapped String to its binary representation.
func (e encode) Binary() String { return e.str.BytesUnsafe().Encode().Binary() }

// Binary converts the wrapped binary String back to its original String.
func (d decode) Binary() Result[String] {
	return TransformResult(d.str.BytesUnsafe().Decode().Binary(), Bytes.String)
}

// JSON encodes the provided string as JSON and returns the result as Result[String].
func (e encode) JSON() Result[String] {
	jsonData, err := json.Marshal(e.str)
	if err != nil {
		return Err[String](err)
	}

	return Ok(String(jsonData))
}

// JSON decodes the provided JSON string and returns the result as Result[String].
func (d decode) JSON() Result[String] {
	var data String
	err := json.Unmarshal(d.str.BytesUnsafe(), &data)
	if err != nil {
		return Err[String](err)
	}

	return Ok(data)
}

// URL encodes the input string, escaping reserved characters as per RFC 2396.
// If safe characters are provided, they will not be encoded.
func (e encode) URL(safe ...String) String {
	reserved := String(";/?:@&=+$,")
	if len(safe) != 0 {
		reserved = safe[0]
	}

	var b Builder
	b.Grow(e.str.Len())

	for _, r := range e.str {
		if reserved.ContainsRune(r) {
			b.WriteRune(r)
			continue
		}

		_, _ = b.WriteString(String(url.QueryEscape(string(r))))
	}

	return b.String()
}

// URL URL-decodes the wrapped String and returns the decoded result as Result[String].
func (d decode) URL() Result[String] {
	result, err := url.QueryUnescape(d.str.Std())
	if err != nil {
		return Err[String](err)
	}

	return Ok(String(result))
}

// HTML HTML-encodes the wrapped String.
func (e encode) HTML() String { return String(html.EscapeString(e.str.Std())) }

// HTML HTML-decodes the wrapped String.
func (d decode) HTML() String { return String(html.UnescapeString(d.str.Std())) }

// Rot13 encodes the wrapped String using ROT13 cipher.
func (e encode) Rot13() String {
	rot := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+13)%26
		default:
			return r
		}
	}

	return e.str.Map(rot)
}

// Rot13 decodes the wrapped String using ROT13 cipher.
func (d decode) Rot13() String { return d.str.Encode().Rot13() }

// Octal returns the octal representation of the encoded string.
func (e encode) Octal() String {
	var b Builder
	var tmp [7]byte

	first := true

	for _, char := range e.str {
		if !first {
			b.WriteByte(' ')
		}

		_, _ = b.Write(strconv.AppendInt(tmp[:0], int64(char), 8))
		first = false
	}

	return b.String()
}

// Octal decodes the octal representation back to String.
func (d decode) Octal() Result[String] {
	var b Builder

	for v := range d.str.Split(" ") {
		n, err := strconv.ParseUint(v.Std(), 8, 32)
		if err != nil {
			return Err[String](err)
		}

		b.WriteRune(rune(n))
	}

	return Ok(b.String())
}
