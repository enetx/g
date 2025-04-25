package g_test

import (
	"fmt"
	"testing"

	. "github.com/enetx/g"
)

// go test -bench=. -benchmem -count=4

func BenchmarkSprintf(b *testing.B) {
	name := "World"

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = fmt.Sprintf("Hello, %s!", name)
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = Format("Hello, {}!", name)
		}
	})
}

func BenchmarkSprintfPositional(b *testing.B) {
	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = fmt.Sprintf("%[2]s comes before %[1]s", "World", "Hello")
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = Format("{2} comes before {1}", "World", "Hello")
		}
	})
}

func BenchmarkSprintfNamedAccess(b *testing.B) {
	data := Named{"email": "alice@example.com"}

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = fmt.Sprintf("Email: %s", data["email"])
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = Format("Email: {email}", data)
		}
	})
}

func BenchmarkSprintfFormatSpecifiers(b *testing.B) {
	num := Int(255)

	b.Run("fmt.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = fmt.Sprintf("Hex: %x, Binary: %b", num, num)
		}
	})

	b.Run("g.Sprintf", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			_ = Format("Hex: {1.Hex}, Binary: {1.Binary}", num)
		}
	})
}
