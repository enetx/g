package g_test

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/enetx/g"
)

// go test -bench=. -benchmem -count=4

func BenchmarkSprintf(b *testing.B) {
	name := "World"

	b.Run("StringConcat", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = "Hello, " + name + "!"
		}
	})

	b.Run("StringBuilder", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var sb strings.Builder
			sb.WriteString("Hello, ")
			sb.WriteString(name)
			sb.WriteString("!")
			_ = sb.String()
		}
	})

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("Hello, %s!", name)
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Sprintf("Hello, {}!", name)
		}
	})
}

func BenchmarkSprintfPositional(b *testing.B) {
	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%s comes before %s", "Hello", "World")
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Sprintf("{1} comes before {0}", "World", "Hello")
		}
	})
}

func BenchmarkSprintfNamedAccess(b *testing.B) {
	data := Named{"email": "alice@example.com"}

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("Email: %s", data["email"])
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Sprintf("Email: {email}", data)
		}
	})
}

func BenchmarkSprintfFormatSpecifiers(b *testing.B) {
	num := 255

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("Hex: %x, Binary: %b", num, num)
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = Sprintf("Hex: {.$hex}, Binary: {$.bin}", num, num)
		}
	})
}
