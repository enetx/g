package g

import "encoding/json"

// jsonNull is the JSON literal returned when marshaling a None Option.
// It is a package-level value to avoid allocating a new []byte on every None marshal.
//
// NOTE: callers (encoding/json) must not mutate the returned slice; the standard
// library treats Marshaler output as read-only, so sharing this backing array is safe.
var jsonNull = []byte("null")

// MarshalJSON implements the json.Marshaler interface for Option[T].
// Some(value) is marshaled as the JSON representation of value.
// None is marshaled as null.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return jsonNull, nil
	}

	return json.Marshal(o.v)
}

// UnmarshalJSON implements the json.Unmarshaler interface for Option[T].
// JSON null is unmarshaled as None.
// Any other valid JSON value is unmarshaled as Some(value).
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = None[T]()
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*o = Some(v)
	return nil
}
