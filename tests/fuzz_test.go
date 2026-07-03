package g_test

// Fuzz targets for panic-safety invariants. The seed corpora below double as
// regression unit cases: `go test ./tests -run Fuzz` executes the seeds only.
// Full fuzzing (run manually, one target at a time):
//
//	go test ./tests -run '^$' -fuzz FuzzFormat
//	go test ./tests -run '^$' -fuzz FuzzSubBytes
//	go test ./tests -run '^$' -fuzz FuzzOctalDecode
//	go test ./tests -run '^$' -fuzz FuzzResultUnmarshalJSON

import (
	"encoding/json"
	"math"
	"strings"
	"testing"

	. "github.com/enetx/g"
)

// FuzzFormat checks that g.Format never panics, regardless of the template
// shape, for both positional and named arguments.
func FuzzFormat(f *testing.F) {
	f.Add("{} + {} = {}", "sum", 42)
	f.Add("{1} {2} {1}", "hello", 7)
	f.Add("My name is {name} and I am {age} years old.", "Alice", 30)
	f.Add("Hello, {name.Trim.Upper}!", "  bob  ", 1)
	f.Add("{name?fallback} {missing?}", "x", 0)
	f.Add("no placeholders at all", "y", -1)
	f.Add("{", "unclosed", 2)
	f.Add("}{", "reversed", 3)
	f.Add(`\{escaped\} {}`, "esc", 4)
	f.Add("{9} {name.NoSuchMethod}", "outofrange", 5)
	f.Add("", "", 0)

	f.Fuzz(func(_ *testing.T, template, a string, b int) {
		_ = Format(template, a, b)
		_ = Format(template, Named{"name": String(a), "age": Int(b)})
	})
}

// FuzzSubBytes checks that Bytes.SubBytes clamps arbitrary start/end/step
// combinations instead of panicking.
func FuzzSubBytes(f *testing.F) {
	f.Add([]byte("hello world"), 0, 5, 1)
	f.Add([]byte("hello"), -3, -1, 1)
	f.Add([]byte("hello"), 5, 0, -1)
	f.Add([]byte("hello"), 100, 0, -1)
	f.Add([]byte("hello"), -100, 100, 2)
	f.Add([]byte(""), 1, 10, 3)
	f.Add([]byte("abc"), 0, 3, 0)
	f.Add([]byte("привет"), 1, 11, 2)
	f.Add([]byte("x"), math.MinInt, math.MaxInt, -2)

	f.Fuzz(func(_ *testing.T, data []byte, start, end, step int) {
		bs := Bytes(data)
		_ = bs.SubBytes(Int(start), Int(end), Int(step))
		_ = bs.SubBytes(Int(start), Int(end))
	})
}

// FuzzOctalDecode checks that String.Decode().Octal() never panics: invalid
// input must surface as an Err result, not a crash.
func FuzzOctalDecode(f *testing.F) {
	f.Add("")
	f.Add("110 145 154 154 157")
	f.Add("777")
	f.Add("141 142 143")
	f.Add("not octal")
	f.Add("999")
	f.Add("-141")
	f.Add("4177777")  // utf8.MaxRune in octal
	f.Add("4200000")  // just above utf8.MaxRune
	f.Add("154000")   // 0xD800, start of the surrogate range
	f.Add("141  142") // double separator -> empty token
	f.Add(strings.Repeat("7", 100))
	f.Add(" 141 ")

	f.Fuzz(func(_ *testing.T, s string) {
		_ = String(s).Decode().Octal()
	})
}

// FuzzResultUnmarshalJSON checks that unmarshaling arbitrary bytes into a
// Result[int] never panics: malformed documents must yield an error.
func FuzzResultUnmarshalJSON(f *testing.F) {
	f.Add([]byte(`{"ok":42}`))
	f.Add([]byte(`{"err":"boom"}`))
	f.Add([]byte(`{"ok":null}`))
	f.Add([]byte(`{"err":null}`))
	f.Add([]byte(`{"ok":1,"err":"x"}`))
	f.Add([]byte(`{"ok":1,"ok":2}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`null`))
	f.Add([]byte(``))
	f.Add([]byte(`[1,2,3]`))
	f.Add([]byte(`{"ok":"not an int"}`))
	f.Add([]byte(`{"ok":9223372036854775808}`))
	f.Add([]byte("{\"ok\":\xff\xfe}"))

	f.Fuzz(func(_ *testing.T, data []byte) {
		var r Result[int]
		_ = json.Unmarshal(data, &r)
	})
}
