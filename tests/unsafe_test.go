package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestStringUnsafeMethods(t *testing.T) {
	t.Run("Lower", func(t *testing.T) {
		in := String("HELLO World!")
		got := in.Lower()
		want := String("hello world!")
		if got != want {
			t.Errorf("Lower() = %q, want %q", got, want)
		}
	})

	t.Run("Upper", func(t *testing.T) {
		in := String("hello World!")
		got := in.Upper()
		want := String("HELLO WORLD!")
		if got != want {
			t.Errorf("Upper() = %q, want %q", got, want)
		}
	})

	t.Run("Title", func(t *testing.T) {
		in := String("hello world")
		got := in.Title()
		want := String("Hello World")
		if got != want {
			t.Errorf("Title() = %q, want %q", got, want)
		}
	})

	t.Run("Trim", func(t *testing.T) {
		in := String("  padded  ")
		got := in.Trim()
		want := String("padded")
		if got != want {
			t.Errorf("Trim() = %q, want %q", got, want)
		}
	})

	t.Run("TrimStart", func(t *testing.T) {
		in := String("  start")
		got := in.TrimStart()
		want := String("start")
		if got != want {
			t.Errorf("TrimStart() = %q, want %q", got, want)
		}
	})

	t.Run("TrimEnd", func(t *testing.T) {
		in := String("end  ")
		got := in.TrimEnd()
		want := String("end")
		if got != want {
			t.Errorf("TrimEnd() = %q, want %q", got, want)
		}
	})

	t.Run("TrimSet", func(t *testing.T) {
		in := String("@@@trimmed@@@")
		got := in.TrimSet("@")
		want := String("trimmed")
		if got != want {
			t.Errorf("TrimSet() = %q, want %q", got, want)
		}
	})

	t.Run("TrimStartSet", func(t *testing.T) {
		in := String("***start")
		got := in.TrimStartSet("*")
		want := String("start")
		if got != want {
			t.Errorf("TrimStartSet() = %q, want %q", got, want)
		}
	})

	t.Run("TrimEndSet", func(t *testing.T) {
		in := String("end###")
		got := in.TrimEndSet("#")
		want := String("end")
		if got != want {
			t.Errorf("TrimEndSet() = %q, want %q", got, want)
		}
	})

	t.Run("StripPrefix", func(t *testing.T) {
		in := String("prefixText")
		got := in.StripPrefix("prefix")
		want := String("Text")
		if got != want {
			t.Errorf("StripPrefix() = %q, want %q", got, want)
		}
	})

	t.Run("StripSuffix", func(t *testing.T) {
		in := String("TextSuffix")
		got := in.StripSuffix("Suffix")
		want := String("Text")
		if got != want {
			t.Errorf("StripSuffix() = %q, want %q", got, want)
		}
	})

	t.Run("Reverse", func(t *testing.T) {
		in := String("abc")
		got := in.Reverse()
		want := String("cba")
		if got != want {
			t.Errorf("Reverse() = %q, want %q", got, want)
		}
	})

	t.Run("Invalid UTF-8 Reverse", func(t *testing.T) {
		in := String(string([]byte{0xff, 0xfe, 0xfd}))
		_ = in.Reverse() // just ensure no panic
	})

	t.Run("Mixed Runes and ASCII Lower", func(t *testing.T) {
		in := String("ABCф")
		got := in.Lower()
		want := String("abcф")
		if got != want {
			t.Errorf("Lower() = %q, want %q", got, want)
		}
	})
}
