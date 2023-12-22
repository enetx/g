package g

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"net/url"
	"strconv"
	"strings"

	"gitlab.com/x0xO/g/pkg/iter"
)

type (
	// A struct that wraps an String for encoding.
	enc struct{ str String }

	// A struct that wraps an String for decoding.
	dec struct{ str String }
)

// Enc returns an enc struct wrapping the given String.
func (s String) Enc() enc { return enc{s} }

// Dec returns a dec struct wrapping the given String.
func (s String) Dec() dec { return dec{s} }

// Base64 encodes the wrapped String using Base64 and returns the encoded result as an String.
func (e enc) Base64() String { return String(base64.StdEncoding.EncodeToString(e.str.ToBytes())) }

// Base64 decodes the wrapped String using Base64 and returns the decoded result as an String.
func (d dec) Base64() Result[String] {
	decoded, err := base64.StdEncoding.DecodeString(d.str.Std())
	if err != nil {
		return Err[String](err)
	}

	return Ok(String(decoded))
}

// JSON encodes the provided data as JSON and returns the result as an String.
func (enc) JSON(data any) Result[String] {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return Err[String](err)
	}

	return Ok(String(jsonData))
}

// JSON decodes the wrapped String using JSON and unmarshals it into the provided data object.
func (d dec) JSON(data any) Result[String] {
	err := json.Unmarshal(d.str.ToBytes(), data)
	if err != nil {
		return Err[String](err)
	}

	return Ok(d.str)
}

// XML encodes the provided data as XML and returns the result as an String.
// The optional prefix and indent String values can be provided for XML indentation.
func (enc) XML(data any, s ...String) Result[String] {
	var (
		prefix String
		indent String
	)

	if len(s) > 1 {
		prefix = s[0]
		indent = s[1]
	}

	xmlData, err := xml.MarshalIndent(data, prefix.Std(), indent.Std())
	if err != nil {
		return Err[String](err)
	}

	return Ok(String(xmlData))
}

// XML decodes the wrapped String using XML and unmarshals it into the provided data object.
func (d dec) XML(data any) Result[String] {
	err := xml.Unmarshal(d.str.ToBytes(), data)
	if err != nil {
		return Err[String](err)
	}

	return Ok(d.str)
}

// URL URL-encodes the wrapped String and returns the encoded result as an String.
func (e enc) URL(safe ...String) String {
	// https://www.ietf.org/rfc/rfc2396.txt#2.2
	reserved := String(";/?:@&=+$,")
	if len(safe) != 0 {
		reserved = safe[0]
	}

	var enc strings.Builder

	e.str.ToRunes().ForEach(func(r rune) {
		if reserved.ContainsRune(r) {
			enc.WriteRune(r)
			return
		}

		enc.WriteString(url.QueryEscape(string(r)))
	})

	return String(enc.String())
}

// URL URL-decodes the wrapped String and returns the decoded result as an String.
func (d dec) URL() Result[String] {
	result, err := url.QueryUnescape(d.str.Std())
	if err != nil {
		return Err[String](err)
	}

	return Ok(String(result))
}

// HTML HTML-encodes the wrapped String and returns the encoded result as an String.
func (e enc) HTML() String { return String(html.EscapeString(e.str.Std())) }

// HTML HTML-decodes the wrapped String and returns the decoded result as an String.
func (d dec) HTML() String { return String(html.UnescapeString(d.str.Std())) }

// Rot13 encodes the wrapped String using ROT13 cipher and returns the encoded result as an
// String.
func (e enc) Rot13() String {
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

// Rot13 decodes the wrapped String using ROT13 cipher and returns the decoded result as an
// String.
func (d dec) Rot13() String { return d.str.Enc().Rot13() }

// XOR encodes the wrapped String using XOR cipher with the given key and returns the encoded
// result as an String.
func (e enc) XOR(key String) String {
	if key.Empty() {
		return e.str
	}

	encrypted := e.str.ToBytes()

	for i := range iter.N(e.str.Len()) {
		encrypted[i] ^= key[i%key.Len()]
	}

	return String(encrypted)
}

// XOR decodes the wrapped String using XOR cipher with the given key and returns the decoded
// result as an String.
func (d dec) XOR(key String) String { return d.str.Enc().XOR(key) }

// Hex hex-encodes the wrapped String and returns the encoded result as an String.
func (e enc) Hex() String {
	var result strings.Builder
	for i := range iter.N(e.str.Len()) {
		_, _ = fmt.Fprint(&result, Int(e.str[i]).ToHex())
	}

	return String(result.String())
}

// Hex hex-decodes the wrapped String and returns the decoded result as an String.
func (d dec) Hex() Result[String] {
	result, err := hex.DecodeString(d.str.Std())
	if err != nil {
		return Err[String](err)
	}

	return Ok(String(result))
}

// Octal returns the octal representation of the encoded string.
func (e enc) Octal() String {
	result := NewSlice[String](e.str.LenRunes())
	for i, char := range e.str.ToRunes() {
		result.Set(i, Int(char).ToOctal())
	}

	return result.Join(" ")
}

// Octal returns the octal representation of the decimal-encoded string.
func (d dec) Octal() Result[String] {
	var r Result[String]

	var result strings.Builder

	d.str.Split(" ").ForEach(func(oct String) {
		n, err := strconv.ParseUint(oct.Std(), 8, 32)
		if err != nil {
			r = Err[String](err)
			return
		}
		fmt.Fprint(&result, string(rune(n)))
	})

	if r.IsOk() {
		r = Ok(String(result.String()))
	}

	return r
}

// Binary converts the wrapped String to its binary representation as an String.
func (e enc) Binary() String {
	var result strings.Builder
	for i := range iter.N(e.str.Len()) {
		_, _ = fmt.Fprint(&result, Int(e.str[i]).ToBinary())
	}

	return String(result.String())
}

// Binary converts the wrapped binary String back to its original String representation.
func (d dec) Binary() Result[String] {
	var result Bytes

	for i := 0; i+8 <= d.str.Len(); i += 8 {
		b, err := strconv.ParseUint(d.str[i:i+8].Std(), 2, 8)
		if err != nil {
			return Err[String](err)
		}

		result = append(result, byte(b))
	}

	return Ok(result.ToString())
}
