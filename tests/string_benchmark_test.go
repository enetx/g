package g_test

import (
	"strings"
	"testing"

	. "github.com/enetx/g"
)

func genText() String {
	text := String("").Builder()
	for range 10000 {
		text.WriteString(String("").Random(1000))
		text.WriteByte('\n')
	}

	return text.String()
}

func BenchmarkSplit(b *testing.B) {
	text := genText()

	for b.Loop() {
		text.Split("\n").Collect()
	}
}

func BenchmarkSplitStd(b *testing.B) {
	text := genText().Std()

	for b.Loop() {
		strings.Split(text, "\n")
	}
}

func BenchmarkFields(b *testing.B) {
	text := genText()

	for b.Loop() {
		text.Fields().Collect()
	}
}

func BenchmarkFieldsStd(b *testing.B) {
	text := genText().Std()

	for b.Loop() {
		strings.Fields(text)
	}
}
