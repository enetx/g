package rx_test

import (
	"github.com/enetx/g/rx"
	"regexp"
	"testing"

	. "github.com/enetx/g"
)

func TestBytes_Regexp_Find(t *testing.T) {
	testBytes := Bytes("Hello 123 World")
	pattern := regexp.MustCompile(`\d+`)

	result := rx.From(pattern).FindBytes(testBytes)

	if result.IsNone() {
		t.Error("Expected to find digits, but got None")
	}

	expected := Bytes("123")
	if !result.Unwrap().Eq(expected) {
		t.Errorf("Find result mismatch: got %s, want %s", result.Unwrap(), expected)
	}
}

func TestBytes_Regexp_Find_NoMatch(t *testing.T) {
	testBytes := Bytes("Hello World")
	pattern := regexp.MustCompile(`\d+`)

	result := rx.From(pattern).FindBytes(testBytes)

	if result.IsSome() {
		t.Error("Expected None for no match, but got Some")
	}
}

func TestBytes_Regexp_Match(t *testing.T) {
	testBytes := Bytes("test@example.com")
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	result := rx.From(emailPattern).IsMatchBytes(testBytes)

	if !result {
		t.Error("Expected email pattern to match")
	}
}

func TestBytes_Regexp_Match_False(t *testing.T) {
	testBytes := Bytes("not-an-email")
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	result := rx.From(emailPattern).IsMatchBytes(testBytes)

	if result {
		t.Error("Expected email pattern not to match")
	}
}

func TestBytes_Regexp_MatchAny(t *testing.T) {
	testBytes := Bytes("The number is 42")

	patterns := []rx.Regex{
		rx.From(regexp.MustCompile(`[a-z]+`)),
		rx.From(regexp.MustCompile(`\d+`)),
		rx.From(regexp.MustCompile(`[A-Z]+`)),
	}

	result := rx.MatchAnyBytes(testBytes, patterns...)

	if !result {
		t.Error("Expected at least one pattern to match")
	}
}

func TestBytes_Regexp_MatchAny_None(t *testing.T) {
	testBytes := Bytes("123")

	patterns := []rx.Regex{
		rx.From(regexp.MustCompile(`[a-z]+`)),
		rx.From(regexp.MustCompile(`[A-Z]+`)),
	}

	result := rx.MatchAnyBytes(testBytes, patterns...)

	if result {
		t.Error("Expected no patterns to match")
	}
}

func TestBytes_Regexp_FindAll(t *testing.T) {
	testBytes := Bytes("abc 123 def 456 ghi")
	pattern := regexp.MustCompile(`\d+`)

	matches := rx.From(pattern).FindIterBytes(testBytes).Collect().Slice()

	if matches.IsEmpty() {
		t.Error("Expected matches, got none")
	}

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}

	expected1 := Bytes("123")
	expected2 := Bytes("456")

	if !matches[0].Eq(expected1) {
		t.Errorf("First match mismatch: got %s, want %s", matches[0], expected1)
	}

	if !matches[1].Eq(expected2) {
		t.Errorf("Second match mismatch: got %s, want %s", matches[1], expected2)
	}
}

func TestBytes_Regexp_FindAll_NoMatches(t *testing.T) {
	testBytes := Bytes("Hello World")
	pattern := regexp.MustCompile(`\d+`)

	matches := rx.From(pattern).FindIterBytes(testBytes).Collect().Slice()

	if !matches.IsEmpty() {
		t.Error("Expected no matches, got some")
	}
}

func TestBytes_Regexp_Replace(t *testing.T) {
	testBytes := Bytes("Hello 123 World 456")
	pattern := regexp.MustCompile(`\d+`)
	replacement := Bytes("XXX")

	// ReplaceBytes replaces the FIRST match only.
	result := rx.From(pattern).ReplaceBytes(testBytes, replacement)
	expected := Bytes("Hello XXX World 456")

	if !result.Eq(expected) {
		t.Errorf("Replace result mismatch: got %s, want %s", result, expected)
	}

	// ReplaceAllBytes replaces every match.
	all := rx.From(pattern).ReplaceAllBytes(testBytes, replacement)
	expectedAll := Bytes("Hello XXX World XXX")

	if !all.Eq(expectedAll) {
		t.Errorf("ReplaceAll result mismatch: got %s, want %s", all, expectedAll)
	}
}

func TestBytes_Regexp_EmptyBytes(t *testing.T) {
	emptyBytes := Bytes("")
	pattern := regexp.MustCompile(`.*`)

	// Find should return None for meaningful patterns on empty bytes
	digitPattern := regexp.MustCompile(`\d+`)
	result := rx.From(digitPattern).FindBytes(emptyBytes)
	if result.IsSome() {
		t.Error("Expected None when searching for digits in empty bytes")
	}

	// Match should work with patterns that match empty strings
	result2 := rx.From(pattern).IsMatchBytes(emptyBytes)
	if !result2 {
		t.Error("Pattern .* should match empty bytes")
	}
}

func TestBytes_Regexp_ComplexPatterns(t *testing.T) {
	testBytes := Bytes("Visit https://example.com or http://test.org")
	urlPattern := regexp.MustCompile(`https?://[^\s]+`)

	results := rx.From(urlPattern).FindIterBytes(testBytes).Collect().Slice()

	if results.IsEmpty() {
		t.Error("Expected matches, got none")
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(results))
	}

	expected1 := Bytes("https://example.com")
	expected2 := Bytes("http://test.org")

	if !results[0].Eq(expected1) {
		t.Errorf("First URL mismatch: got %s, want %s", results[0], expected1)
	}

	if !results[1].Eq(expected2) {
		t.Errorf("Second URL mismatch: got %s, want %s", results[1], expected2)
	}
}

func TestBytes_Regexp_Find_EmptyMatch(t *testing.T) {
	// A pattern that matches the empty string at position 0 must return Some(""),
	// not None (which is reserved for genuine no-match).
	testBytes := Bytes("abc")
	pattern := regexp.MustCompile(`x*`)

	result := rx.From(pattern).FindBytes(testBytes)
	if result.IsNone() {
		t.Fatal("Expected Some for empty match, got None")
	}
	if !result.Unwrap().Eq(Bytes("")) {
		t.Errorf("Expected empty match, got %q", result.Unwrap())
	}
}

func TestBytes_Regexp_ReplaceBy(t *testing.T) {
	testBytes := Bytes("a1b2c3")
	pattern := regexp.MustCompile(`\d`)

	result := rx.From(pattern).ReplaceAllBytesBy(testBytes, func(m Bytes) Bytes {
		// Clone first: the match aliases the source backing array (see Bytes.Append docs).
		return m.Clone().Append(Bytes("!"))
	})

	expected := Bytes("a1!b2!c3!")
	if !result.Eq(expected) {
		t.Errorf("ReplaceBy = %q, want %q", result, expected)
	}
}

func TestBytes_Regexp_Split(t *testing.T) {
	testBytes := Bytes("a1b2c3d")
	pattern := regexp.MustCompile(`\d`)

	result := rx.From(pattern).SplitBytes(testBytes).Collect().Slice()
	expected := Slice[Bytes]{Bytes("a"), Bytes("b"), Bytes("c"), Bytes("d")}

	if result.Len() != expected.Len() {
		t.Fatalf("Split length = %d, want %d", result.Len(), expected.Len())
	}
	for i := range expected {
		if !result[i].Eq(expected[i]) {
			t.Errorf("Split[%d] = %q, want %q", i, result[i], expected[i])
		}
	}
}

func TestBytes_Regexp_SplitN(t *testing.T) {
	testBytes := Bytes("a1b2c3d")
	pattern := regexp.MustCompile(`\d`)

	parts := rx.From(pattern).SplitBytesN(testBytes, 2).Collect().Slice()
	expected := Slice[Bytes]{Bytes("a"), Bytes("b2c3d")}
	if parts.Len() != expected.Len() {
		t.Fatalf("SplitN length = %d, want %d", parts.Len(), expected.Len())
	}
	for i := range expected {
		if !parts[i].Eq(expected[i]) {
			t.Errorf("SplitN[%d] = %q, want %q", i, parts[i], expected[i])
		}
	}

	// n == 0 yields no parts
	if !rx.From(pattern).SplitBytesN(testBytes, 0).Collect().Slice().IsEmpty() {
		t.Error("SplitN(_, 0) should yield no parts")
	}
}

func TestBytes_Regexp_ReplaceBytesBy_First(t *testing.T) {
	re := rx.From(regexp.MustCompile(`\d+`))

	// Only the FIRST match is transformed.
	got := re.ReplaceBytesBy(Bytes("a1b22c333"), func(m Bytes) Bytes {
		return m.Clone().Append(Bytes("!"))
	})
	if !got.Eq(Bytes("a1!b22c333")) {
		t.Errorf("ReplaceBytesBy = %q, want %q", got, "a1!b22c333")
	}

	// No match returns the input untouched.
	src := Bytes("abc")
	if got := re.ReplaceBytesBy(src, func(Bytes) Bytes { return Bytes("X") }); !got.Eq(src) {
		t.Errorf("ReplaceBytesBy without match = %q, want %q", got, src)
	}
}
