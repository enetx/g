package g_test

import (
	"database/sql/driver"
	"errors"
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

type valuerValue string

func (v valuerValue) Value() (driver.Value, error) {
	return string(v) + "_val", nil
}

func TestOptionScanAllBranches(t *testing.T) {
	var o Option[int]
	if err := o.Scan(nil); err != nil {
		t.Fatal(err)
	}
	if !o.IsNone() {
		t.Fatal("Expected None for nil")
	}

	var s Option[scanValue]
	if err := s.Scan("hello"); err != nil {
		t.Fatalf("Scan error: %v", err)
	}

	if s.Some() != "hello" {
		t.Fatalf("Expected 'hello', got %v", s.Some())
	}

	var si Option[int]
	if err := si.Scan(42); err != nil {
		t.Fatalf("Scan error: %v", err)
	}

	if si.Some() != 42 {
		t.Fatalf("Expected 42, got %v", si.Some())
	}

	var ci Option[int]
	if err := ci.Scan(int64(99)); err != nil {
		t.Fatalf("Scan int64 → int failed: %v", err)
	}

	if ci.Some() != 99 {
		t.Fatalf("Expected 99, got %v", ci.Some())
	}

	var cf Option[float32]
	if err := cf.Scan(float64(3.14)); err != nil {
		t.Fatalf("Scan float64 → float32 failed: %v", err)
	}

	if cf.Some() != 3.14 {
		t.Fatalf("Expected 3.14, got %v", cf.Some())
	}

	var cs Option[string]
	if err := cs.Scan([]byte("abc")); err != nil {
		t.Fatalf("Scan []byte → string failed: %v", err)
	}

	if cs.Some() != "abc" {
		t.Fatalf("Expected 'abc', got %v", cs.Some())
	}

	var bad Option[int]
	err := bad.Scan("oops")
	if err == nil {
		t.Fatal("Expected error for incompatible type")
	}
}

func TestOptionValueAllBranches(t *testing.T) {
	n := None[int]()
	v, err := n.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if v != nil {
		t.Fatalf("Expected nil, got %v", v)
	}

	o := Some(valuerValue("x"))
	val, err := o.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if val != "x_val" {
		t.Fatalf("Expected 'x_val', got %v", val)
	}

	o2 := Some(123)
	v2, _ := o2.Value()
	if v2 != int64(123) {
		t.Fatalf("Expected 123, got %v", v2)
	}

	o3 := Some(3.14)
	v3, _ := o3.Value()
	if v3 != 3.14 {
		t.Fatalf("Expected 3.14, got %v", v3)
	}

	o4 := Some("abc")
	v4, _ := o4.Value()
	if v4 != "abc" {
		t.Fatalf("Expected 'abc', got %v", v4)
	}

	now := time.Now()
	o5 := Some(now)
	v5, _ := o5.Value()
	if v5.(time.Time) != now {
		t.Fatalf("Expected %v, got %v", now, v5)
	}

	o6 := Some(Int(42))
	v6, _ := o6.Value()
	if v6 != int64(42) {
		t.Fatalf("Expected int64(42), got %v", v6)
	}

	o7 := Some(Float(2.71))
	v7, _ := o7.Value()
	if v7 != 2.71 {
		t.Fatalf("Expected 2.71, got %v", v7)
	}

	o8 := Some(String("hello"))
	v8, _ := o8.Value()
	if v8 != "hello" {
		t.Fatalf("Expected 'hello', got %v", v8)
	}

	o9 := Some(Bytes([]byte("bytes")))
	v9, _ := o9.Value()
	if string(v9.([]byte)) != "bytes" {
		t.Fatalf("Expected 'bytes', got %v", string(v9.([]byte)))
	}

	type X struct{}
	o10 := Some(X{})
	_, err = o10.Value()
	if err == nil {
		t.Fatal("Expected error for unsupported type")
	}
}
