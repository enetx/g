package g

import "strings"

// Builder represents a string builder.
type Builder struct{ builder *strings.Builder }

// NewBuilder creates a new instance of Builder.
func NewBuilder() *Builder { return &Builder{new(strings.Builder)} }

// Write appends a string to the current state of the builder.
func (b *Builder) Write(str String) *Builder {
	b.builder.WriteString(str.Std())
	return b
}

// WriteBytes appends a byte slice to the current state of the builder.
func (b *Builder) WriteBytes(bs Bytes) *Builder {
	b.builder.Write(bs)
	return b
}

// WriteByte appends a byte to the current state of the builder.
func (b *Builder) WriteByte(c byte) *Builder {
	b.builder.WriteByte(c)
	return b
}

// WriteRune appends a rune to the current state of the builder.
func (b *Builder) WriteRune(r rune) *Builder {
	b.builder.WriteRune(r)
	return b
}

// Grow increases the capacity of the builder by n bytes.
func (b *Builder) Grow(n Int) { b.builder.Grow(n.Std()) }

// Cap returns the current capacity of the builder.
func (b *Builder) Cap() Int { return Int(b.builder.Cap()) }

// Len returns the current length of the string in the builder.
func (b *Builder) Len() Int { return Int(b.builder.Len()) }

// Reset clears the content of the Builder, resetting it to an empty state.
func (b *Builder) Reset() { b.builder.Reset() }

// String returns the content of the builder as a string.
func (b *Builder) String() String { return String(b.builder.String()) }
