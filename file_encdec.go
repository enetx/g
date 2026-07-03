package g

import (
	"encoding/gob"

	json "encoding/json/v2"
)

type (
	// fencode represents a wrapper for file encoding.
	fencode struct{ f *File }

	// fdecode represents a wrapper for file decoding.
	fdecode struct{ f *File }
)

// Encode returns an fencode struct wrapping the given file for encoding.
func (f *File) Encode() fencode { return fencode{f} }

// Decode returns an fdecode struct wrapping the given file for decoding.
func (f *File) Decode() fdecode { return fdecode{f} }

// Gob encodes the provided data using the encoding/gob package and writes it to the file.
// It returns a Result[*File] indicating the success or failure of the encoding operation.
//
// If the encoding operation is successful, the created file is closed automatically.
//
// Usage:
//
//	data := g.SliceOf(1, 2, 3, 4)
//	result := g.NewFile("somefile.gob").Encode().Gob(data)
//
// Parameters:
//   - data: The data to be encoded and written to the file.
//
// Returns:
//   - Result[*File]: A Result containing a *File if the operation is successful; otherwise, an error Result.
func (fe fencode) Gob(data any) Result[*File] {
	r := fe.f.Create()
	if r.IsErr() {
		return r
	}

	defer r.v.Close()

	if err := gob.NewEncoder(r.v.Std()).Encode(data); err != nil {
		return Err[*File](err)
	}

	return r
}

// Gob decodes data from the file using the encoding/gob package and populates the provided data structure.
// It returns a Result[*File] indicating the success or failure of the decoding operation.
//
// If the decoding operation is successful, the file is closed automatically.
//
// Usage:
//
//	var data g.Slice[int]
//	result := g.NewFile("somefile.gob").Decode().Gob(&data)
//
// Parameters:
//   - data: A pointer to the data structure where the decoded data will be stored.
//
// Returns:
//   - Result[*File]: A Result containing a *File if the operation is successful; otherwise, an error Result.
func (fd fdecode) Gob(data any) Result[*File] {
	r := fd.f.Open()
	if r.IsErr() {
		return r
	}

	defer r.v.Close()

	if err := gob.NewDecoder(r.v.Std()).Decode(data); err != nil {
		return Err[*File](err)
	}

	return r
}

// JSON encodes the provided data using the encoding/json/v2 package and writes it to the file.
// It returns a Result[*File] indicating the success or failure of the encoding operation.
//
// If the encoding operation is successful, the created file is closed automatically.
//
// Breaking changes (v2 semantics) compared to the previous encoding/json v1 implementation:
//   - nil slices are marshaled as [] and nil maps as {} (v1 emitted null for both);
//   - strings containing invalid UTF-8 yield Err (v1 replaced invalid sequences with U+FFFD);
//   - no trailing newline is written after the JSON value (v1's Encoder.Encode appended '\n');
//   - '<', '>', '&' and U+2028/U+2029 are emitted raw (v1's Encoder HTML-escaped them).
//
// Usage:
//
//	data := g.SliceOf(1, 2, 3, 4)
//	result := g.NewFile("somefile.json").Encode().JSON(data)
//
// Parameters:
//   - data: The data to be encoded and written to the file.
//
// Returns:
//   - Result[*File]: A Result containing a *File if the operation is successful; otherwise, an error Result.
func (fe fencode) JSON(data any) Result[*File] {
	r := fe.f.Create()
	if r.IsErr() {
		return r
	}

	defer r.v.Close()

	if err := json.MarshalWrite(r.v.Std(), data); err != nil {
		return Err[*File](err)
	}

	return r
}

// JSON decodes data from the file using the encoding/json/v2 package and populates the provided data structure.
// It returns a Result[*File] indicating the success or failure of the decoding operation.
//
// If the decoding operation is successful, the file is closed automatically.
//
// Breaking changes (v2 semantics) compared to the previous encoding/json v1 implementation:
//   - duplicate object member names yield Err (v1 silently kept the last value);
//   - struct field name matching is case-sensitive (v1 fell back to case-insensitive matching);
//   - JSON strings containing invalid UTF-8 yield Err (v1 decoded them with U+FFFD replacements);
//   - the file must contain exactly one JSON value: non-whitespace data after the
//     top-level value yields Err (v1's Decoder.Decode read one value and ignored the rest).
//
// Usage:
//
//	var data g.Slice[int]
//	result := g.NewFile("somefile.json").Decode().JSON(&data)
//
// Parameters:
//   - data: A pointer to the data structure where the decoded data will be stored.
//
// Returns:
//   - Result[*File]: A Result containing a *File if the operation is successful; otherwise, an error Result.
func (fd fdecode) JSON(data any) Result[*File] {
	r := fd.f.Open()
	if r.IsErr() {
		return r
	}

	defer r.v.Close()

	if err := json.UnmarshalRead(r.v.Std(), data); err != nil {
		return Err[*File](err)
	}

	return r
}
