package g_test

import (
	"testing"

	"github.com/enetx/g"
)

func TestBuilder(t *testing.T) {
	builder := g.NewBuilder()

	// Test Write
	builder.Write(g.NewString("hello"))
	expected := "hello"
	if result := builder.String().Std(); result != expected {
		t.Errorf("Write() = %s; want %s", result, expected)
	}

	// Test WriteBytes
	builder.WriteBytes([]byte(" world"))
	expected = "hello world"
	if result := builder.String().Std(); result != expected {
		t.Errorf("WriteBytes() = %s; want %s", result, expected)
	}

	// Test WriteByte
	builder.WriteByte('!')
	expected = "hello world!"
	if result := builder.String().Std(); result != expected {
		t.Errorf("WriteByte() = %s; want %s", result, expected)
	}

	// Test WriteRune
	builder.WriteRune('👋')
	expected = "hello world!👋"
	if result := builder.String().Std(); result != expected {
		t.Errorf("WriteRune() = %s; want %s", result, expected)
	}

	// Test Len
	if result := builder.Len().Std(); result != len(expected) {
		t.Errorf("Len() = %d; want %d", result, len(expected))
	}

	// Test Reset
	builder.Reset()
	if result := builder.Len(); result != 0 {
		t.Errorf("After Reset, Len() = %d; want 0", result)
	}
}

func TestBuilderGrow(t *testing.T) {
	// Create a new Builder
	builder := g.NewBuilder()

	// Grow the builder
	builder.Grow(16)

	// Check if the capacity has been increased
	expected := 16
	if result := builder.Cap().Std(); result != expected {
		t.Errorf("Grow(16) = %d; want %d", result, expected)
	}
}

func TestBuilderCap(t *testing.T) {
	// Create a new Builder
	builder := g.NewBuilder()

	// Check the initial capacity
	expected := 0
	if result := builder.Cap().Std(); result != expected {
		t.Errorf("Initial Cap() = %d; want %d", result, expected)
	}
}
