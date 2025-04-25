package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestStringUnsafe(t *testing.T) {
	original := Bytes("hello unsafe")
	str := original.StringUnsafe()

	if str != "hello unsafe" {
		t.Errorf("StringUnsafe() = %q, want %q", str, "hello unsafe")
	}

	bs := str.BytesUnsafe()
	if !bs.Eq(original) {
		t.Errorf("BytesUnsafe() after StringUnsafe() mismatch: got %q, want %q", bs, original)
	}
}

func TestBytesUnsafe(t *testing.T) {
	original := String("unsafe back")
	bs := original.BytesUnsafe()

	if bs.String() != "unsafe back" {
		t.Errorf("BytesUnsafe() = %q, want %q", bs.String(), "unsafe back")
	}

	str := bs.StringUnsafe()
	if str != original {
		t.Errorf("StringUnsafe() after BytesUnsafe() mismatch: got %q, want %q", str, original)
	}
}

func BenchmarkUnsafeString(b *testing.B) {
	original := Bytes("this is a very fast conversion without copy")

	b.ReportAllocs()

	for b.Loop() {
		_ = original.StringUnsafe()
	}
}

func BenchmarkUnsafeStringStd(b *testing.B) {
	original := Bytes("this is a very fast conversion without copy")

	b.ReportAllocs()

	for b.Loop() {
		_ = String(string(original))
	}
}

func BenchmarkUnsafeBytes(b *testing.B) {
	original := String("this is a very fast conversion without copy")

	b.ReportAllocs()

	for b.Loop() {
		_ = original.BytesUnsafe()
	}
}

func BenchmarkUnsafeBytesStd(b *testing.B) {
	original := String("this is a very fast conversion without copy")

	b.ReportAllocs()

	for b.Loop() {
		_ = Bytes([]byte(original.Std()))
	}
}
