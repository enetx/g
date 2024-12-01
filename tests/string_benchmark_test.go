package g_test

import (
	"strings"
	"testing"

	. "github.com/enetx/g"
)

func genText() String {
	text := String("").Builder()
	for range 10000 {
		text.Write(String("").Random(1000)).WriteByte('\n')
	}

	return text.String()
}

func BenchmarkSplit(b *testing.B) {
	text := genText()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		text.Split("\n").Collect()
	}
}

func BenchmarkSplitStd(b *testing.B) {
	text := genText().Std()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		strings.Split(text, "\n")
	}
}

func BenchmarkFields(b *testing.B) {
	text := genText()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		text.Fields().Collect()
	}
}

func BenchmarkFieldsStd(b *testing.B) {
	text := genText().Std()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		strings.Fields(text)
	}
}
