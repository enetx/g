package g

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math"
	"time"
)

// Scan implements the database/sql.Scanner interface.
//
// Behavior:
//   - If src is nil, the Option is set to None (SQL NULL).
//   - If T implements sql.Scanner, its Scan method is used.
//   - If src can be directly assigned to T, it is used as-is.
//   - Otherwise, common database type conversions are attempted
//     (e.g. int64 -> int, []byte -> string).
//
// Supported conversions (driver-dependent):
//   - INTEGER    -> int, g.Int
//   - REAL       -> float32, float64, g.Float
//   - TEXT       -> string, g.String
//   - BLOB       -> []byte, g.Bytes
//   - BOOLEAN    -> bool
//   - TIMESTAMP  -> time.Time
//
// This method does not use reflection.
func (o *Option[T]) Scan(src any) error {
	if src == nil {
		*o = None[T]()
		return nil
	}

	var v T

	if scanner, ok := any(&v).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return err
		}
		*o = Some(v)
		return nil
	}

	if val, ok := src.(T); ok {
		*o = Some(val)
		return nil
	}

	if converted, ok := convertToT[T](src); ok {
		*o = Some(converted)
		return nil
	}

	return fmt.Errorf("Option.Scan: cannot scan %T into %T", src, v)
}

// Value implements the database/sql/driver.Valuer interface.
//
// Behavior:
//   - If the Option is None, it returns nil (SQL NULL).
//   - If T implements driver.Valuer, its Value method is used.
//   - If the value is a valid driver.Value, it is returned directly.
//   - Otherwise, common conversions are applied (e.g. int -> int64).
func (o Option[T]) Value() (driver.Value, error) {
	if o.IsNone() {
		return nil, nil
	}

	if valuer, ok := any(o.v).(driver.Valuer); ok {
		return valuer.Value()
	}

	val := any(o.v)

	switch val.(type) {
	case int64, float64, bool, []byte, string, time.Time:
		return val, nil
	}

	if converted, ok := convertToDriverValue(val); ok {
		return converted, nil
	}

	return nil, fmt.Errorf("Option.Value: unsupported type %T", o.v)
}

// convertToT converts src to type T for common database driver conversions.
func convertToT[T any](src any) (T, bool) {
	var zero T

	switch any(zero).(type) {
	case int:
		if i64, ok := src.(int64); ok {
			return any(int(i64)).(T), true
		}
	case float32:
		if v, ok := src.(float64); ok {
			return any(float32(v)).(T), true
		}
	case string:
		switch v := src.(type) {
		case string:
			return any(v).(T), true
		case []byte:
			return any(string(v)).(T), true
		}
	case bool:
		if v, ok := src.(bool); ok {
			return any(v).(T), true
		}
	case time.Time:
		if v, ok := src.(time.Time); ok {
			return any(v).(T), true
		}
	case String:
		switch v := src.(type) {
		case string:
			return any(String(v)).(T), true
		case []byte:
			return any(String(v)).(T), true
		}
	case Int:
		if v, ok := src.(int64); ok {
			return any(Int(v)).(T), true
		}
	case Float:
		if v, ok := src.(float64); ok {
			return any(Float(v)).(T), true
		}
	case Bytes:
		if v, ok := src.([]byte); ok {
			return any(Bytes(v)).(T), true
		}
	}

	return zero, false
}

// convertToDriverValue converts Go values into valid driver.Value types.
func convertToDriverValue(val any) (driver.Value, bool) {
	switch v := val.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		if v <= math.MaxInt64 {
			return int64(v), true
		}
		return nil, false
	case float32:
		return float64(v), true
	case Int:
		return v.Int64(), true
	case Float:
		return v.Std(), true
	case String:
		return v.Std(), true
	case Bytes:
		return v.Std(), true
	}

	return nil, false
}
