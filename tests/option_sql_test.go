package g_test

import (
	"database/sql/driver"
	"errors"
	"math"
	"testing"
	"time"

	. "github.com/enetx/g"
)

type scanValue string

func (s *scanValue) Scan(src any) error {
	if str, ok := src.(string); ok {
		*s = scanValue(str)
		return nil
	}
	return errors.New("scan error")
}

type badScanner string

func (b *badScanner) Scan(src any) error {
	return errors.New("forced scan error")
}

type valuerValue string

func (v valuerValue) Value() (driver.Value, error) {
	return string(v) + "_val", nil
}

func TestOptionScanNil(t *testing.T) {
	var o Option[int]
	if err := o.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if !o.IsNone() {
		t.Fatal("expected None for nil")
	}
}

func TestOptionScanScannerSuccess(t *testing.T) {
	var o Option[scanValue]
	if err := o.Scan("hello"); err != nil {
		t.Fatalf("Scan error: %v", err)
	}
	if o.Some() != "hello" {
		t.Fatalf("expected 'hello', got %v", o.Some())
	}
}

func TestOptionScanScannerError(t *testing.T) {
	var o Option[badScanner]
	err := o.Scan("x")
	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestOptionScanDirectType(t *testing.T) {
	var o Option[string]
	if err := o.Scan("direct"); err != nil {
		t.Fatal(err)
	}
	if o.Some() != "direct" {
		t.Fatalf("expected 'direct', got %v", o.Some())
	}
}

func TestOptionScanNumericConversions(t *testing.T) {
	var i Option[int]
	if err := i.Scan(int64(42)); err != nil {
		t.Fatal(err)
	}
	if i.Some() != 42 {
		t.Fatalf("expected 42, got %v", i.Some())
	}

	var f Option[float32]
	if err := f.Scan(float64(3.14)); err != nil {
		t.Fatal(err)
	}
	if f.Some() != 3.14 {
		t.Fatalf("expected 3.14, got %v", f.Some())
	}
}

func TestOptionScanStringAndBytes(t *testing.T) {
	var s Option[string]
	if err := s.Scan([]byte("abc")); err != nil {
		t.Fatal(err)
	}
	if s.Some() != "abc" {
		t.Fatalf("expected 'abc', got %v", s.Some())
	}
}

func TestOptionScanBoolAndTime(t *testing.T) {
	var b Option[bool]
	if err := b.Scan(true); err != nil {
		t.Fatal(err)
	}
	if b.Some() != true {
		t.Fatalf("expected true, got %v", b.Some())
	}

	now := time.Now()
	var tm Option[time.Time]
	if err := tm.Scan(now); err != nil {
		t.Fatal(err)
	}
	if !tm.Some().Equal(now) {
		t.Fatalf("expected %v, got %v", now, tm.Some())
	}
}

func TestOptionScanGTypes(t *testing.T) {
	var s Option[String]
	if err := s.Scan([]byte("abc")); err != nil {
		t.Fatal(err)
	}
	if s.Some() != "abc" {
		t.Fatalf("expected 'abc', got %v", s.Some())
	}

	var i Option[Int]
	if err := i.Scan(int64(10)); err != nil {
		t.Fatal(err)
	}
	if i.Some() != 10 {
		t.Fatalf("expected 10, got %v", i.Some())
	}

	var f Option[Float]
	if err := f.Scan(float64(1.5)); err != nil {
		t.Fatal(err)
	}
	if f.Some() != 1.5 {
		t.Fatalf("expected 1.5, got %v", f.Some())
	}

	var b Option[Bytes]
	if err := b.Scan([]byte{1, 2}); err != nil {
		t.Fatal(err)
	}
	if len(b.Some()) != 2 {
		t.Fatalf("unexpected bytes %v", b.Some())
	}
}

func TestOptionScanIncompatible(t *testing.T) {
	var o Option[int]
	err := o.Scan("oops")
	if err == nil {
		t.Fatal("expected error for incompatible type")
	}
}

func TestOptionValueNone(t *testing.T) {
	n := None[int]()
	v, err := n.Value()
	if err != nil {
		t.Fatal(err)
	}
	if v != nil {
		t.Fatalf("expected nil, got %v", v)
	}
}

func TestOptionValueValuer(t *testing.T) {
	o := Some(valuerValue("x"))
	v, err := o.Value()
	if err != nil {
		t.Fatal(err)
	}
	if v != "x_val" {
		t.Fatalf("expected 'x_val', got %v", v)
	}
}

func TestOptionValueNativeTypes(t *testing.T) {
	o := Some(123)
	v, _ := o.Value()
	if v != int64(123) {
		t.Fatalf("expected 123, got %v", v)
	}

	o2 := Some(3.14)
	v2, _ := o2.Value()
	if v2 != 3.14 {
		t.Fatalf("expected 3.14, got %v", v2)
	}

	o3 := Some("abc")
	v3, _ := o3.Value()
	if v3 != "abc" {
		t.Fatalf("expected 'abc', got %v", v3)
	}

	now := time.Now()
	o4 := Some(now)
	v4, _ := o4.Value()
	if !v4.(time.Time).Equal(now) {
		t.Fatalf("expected %v, got %v", now, v4)
	}
}

func TestOptionValueFloat32(t *testing.T) {
	o := Some(float32(1.25))
	v, err := o.Value()
	if err != nil {
		t.Fatal(err)
	}
	if v != float64(1.25) {
		t.Fatalf("expected 1.25, got %v", v)
	}
}

func TestOptionValueUintVariants(t *testing.T) {
	tests := []any{
		uint(1),
		uint8(2),
		uint16(3),
		uint32(4),
	}

	for _, tt := range tests {
		o := Some(tt)
		v, err := o.Value()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := v.(int64); !ok {
			t.Fatalf("expected int64, got %T", v)
		}
	}
}

func TestOptionValueUintOverflow(t *testing.T) {
	o := Some(uint64(math.MaxInt64 + 1))
	v, err := o.Value()
	if err != nil {
		t.Logf("expected overflow handling, got err=%v, v=%v", err, v)
	} else {
		t.Fatalf("expected overflow error, got value %v", v)
	}
}

func TestOptionValueGTypes(t *testing.T) {
	o := Some(Int(42))
	v, _ := o.Value()
	if v != int64(42) {
		t.Fatalf("expected int64(42), got %v", v)
	}

	o2 := Some(Float(2.71))
	v2, _ := o2.Value()
	if v2 != 2.71 {
		t.Fatalf("expected 2.71, got %v", v2)
	}

	o3 := Some(String("hello"))
	v3, _ := o3.Value()
	if v3 != "hello" {
		t.Fatalf("expected 'hello', got %v", v3)
	}

	o4 := Some(Bytes([]byte("bytes")))
	v4, _ := o4.Value()
	if string(v4.([]byte)) != "bytes" {
		t.Fatalf("expected 'bytes', got %v", string(v4.([]byte)))
	}
}

func TestOptionValueUnsupported(t *testing.T) {
	type X struct{}
	o := Some(X{})
	_, err := o.Value()
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestOptionValueUnsupportedSlice(t *testing.T) {
	o := Some([]int{1, 2, 3})
	_, err := o.Value()
	if err == nil {
		t.Fatal("expected error for []int")
	}
}
