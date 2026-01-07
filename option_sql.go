package g

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

// Scan implements the database/sql.Scanner interface.
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
func (o Option[T]) Value() (driver.Value, error) {
	if o.IsNone() {
		return nil, nil
	}

	if valuer, ok := any(o.v).(driver.Valuer); ok {
		return valuer.Value()
	}

	val := any(o.v)
	switch val.(type) {
	case int64, float64, bool, []byte, string, time.Time, nil:
		return val, nil
	}

	if converted, ok := convertToDriverValue(val); ok {
		return converted, nil
	}

	return nil, fmt.Errorf("Option.Value: unsupported type %T", o.v)
}

// convertToT converts src to type T for common database type conversions
func convertToT[T any](src any) (T, bool) {
	var zero T

	switch any(zero).(type) {
	case int:
		if i64, ok := src.(int64); ok {
			return any(int(i64)).(T), true
		}
	case string:
		if b, ok := src.([]byte); ok {
			return any(string(b)).(T), true
		}
	case float32:
		if f64, ok := src.(float64); ok {
			return any(float32(f64)).(T), true
		}
	case String:
		if s, ok := src.(string); ok {
			return any(String(s)).(T), true
		}
		if b, ok := src.([]byte); ok {
			return any(String(b)).(T), true
		}
	case Int:
		if i64, ok := src.(int64); ok {
			return any(Int(i64)).(T), true
		}
	case Float:
		if f64, ok := src.(float64); ok {
			return any(Float(f64)).(T), true
		}
	case Bytes:
		if b, ok := src.([]byte); ok {
			return any(Bytes(b)).(T), true
		}
	}

	return zero, false
}

// convertToDriverValue converts val to a valid driver.Value type
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
		return int64(v), true
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
