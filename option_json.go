package g

import "encoding/json"

// MarshalJSON implements the json.Marshaler interface for Option[T].
// Some(value) is marshaled as the JSON representation of value.
// None is marshaled as null.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return []byte("null"), nil
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
