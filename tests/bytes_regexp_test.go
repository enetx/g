package g_test

import (
	"regexp"
	"testing"

	"github.com/enetx/g"
)

func TestBytes_Regexp_Find(t *testing.T) {
	testBytes := g.Bytes("Hello 123 World")
	pattern := regexp.MustCompile(`\d+`)

	result := testBytes.Regexp().Find(pattern)

	if result.IsNone() {
		t.Error("Expected to find digits, but got None")
	}

	expected := g.Bytes("123")
	if !result.Unwrap().Eq(expected) {
		t.Errorf("Find result mismatch: got %s, want %s", result.Unwrap(), expected)
	}
}

func TestBytes_Regexp_Find_NoMatch(t *testing.T) {
	testBytes := g.Bytes("Hello World")
	pattern := regexp.MustCompile(`\d+`)

	result := testBytes.Regexp().Find(pattern)

	if result.IsSome() {
		t.Error("Expected None for no match, but got Some")
	}
}

func TestBytes_Regexp_Match(t *testing.T) {
	testBytes := g.Bytes("test@example.com")
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	result := testBytes.Regexp().Match(emailPattern)

	if !result {
		t.Error("Expected email pattern to match")
	}
}

func TestBytes_Regexp_Match_False(t *testing.T) {
	testBytes := g.Bytes("not-an-email")
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	result := testBytes.Regexp().Match(emailPattern)

	if result {
		t.Error("Expected email pattern not to match")
	}
}

func TestBytes_Regexp_MatchAny(t *testing.T) {
	testBytes := g.Bytes("The number is 42")

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`[a-z]+`),
		regexp.MustCompile(`\d+`),
		regexp.MustCompile(`[A-Z]+`),
	}

	result := testBytes.Regexp().MatchAny(patterns...)

	if !result {
		t.Error("Expected at least one pattern to match")
	}
}

func TestBytes_Regexp_MatchAny_None(t *testing.T) {
	testBytes := g.Bytes("123")

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`[a-z]+`),
		regexp.MustCompile(`[A-Z]+`),
	}

	result := testBytes.Regexp().MatchAny(patterns...)

	if result {
		t.Error("Expected no patterns to match")
	}
}

func TestBytes_Regexp_FindAll(t *testing.T) {
	testBytes := g.Bytes("abc 123 def 456 ghi")
	pattern := regexp.MustCompile(`\d+`)

	result := testBytes.Regexp().FindAll(pattern)

	if result.IsNone() {
		t.Error("Expected Some result, got None")
	}

	matches := result.Unwrap()
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}

	expected1 := g.Bytes("123")
	expected2 := g.Bytes("456")

	if !matches[0].Eq(expected1) {
		t.Errorf("First match mismatch: got %s, want %s", matches[0], expected1)
	}

	if !matches[1].Eq(expected2) {
		t.Errorf("Second match mismatch: got %s, want %s", matches[1], expected2)
	}
}

func TestBytes_Regexp_FindAll_NoMatches(t *testing.T) {
	testBytes := g.Bytes("Hello World")
	pattern := regexp.MustCompile(`\d+`)

	result := testBytes.Regexp().FindAll(pattern)

	if result.IsSome() {
		t.Error("Expected None for no matches, got Some")
	}
}

func TestBytes_Regexp_Replace(t *testing.T) {
	testBytes := g.Bytes("Hello 123 World 456")
	pattern := regexp.MustCompile(`\d+`)
	replacement := g.Bytes("XXX")

	result := testBytes.Regexp().Replace(pattern, replacement)
	expected := g.Bytes("Hello XXX World XXX")

	if !result.Eq(expected) {
		t.Errorf("Replace result mismatch: got %s, want %s", result, expected)
	}
}

func TestBytes_Regexp_EmptyBytes(t *testing.T) {
	emptyBytes := g.Bytes("")
	pattern := regexp.MustCompile(`.*`)

	// Find should return None for meaningful patterns on empty bytes
	digitPattern := regexp.MustCompile(`\d+`)
	result := emptyBytes.Regexp().Find(digitPattern)
	if result.IsSome() {
		t.Error("Expected None when searching for digits in empty bytes")
	}

	// Match should work with patterns that match empty strings
	result2 := emptyBytes.Regexp().Match(pattern)
	if !result2 {
		t.Error("Pattern .* should match empty bytes")
	}
}

func TestBytes_Regexp_ComplexPatterns(t *testing.T) {
	testBytes := g.Bytes("Visit https://example.com or http://test.org")
	urlPattern := regexp.MustCompile(`https?://[^\s]+`)

	result := testBytes.Regexp().FindAll(urlPattern)

	if result.IsNone() {
		t.Error("Expected Some result, got None")
	}

	results := result.Unwrap()
	if len(results) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(results))
	}

	expected1 := g.Bytes("https://example.com")
	expected2 := g.Bytes("http://test.org")

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
	testBytes := g.Bytes("abc")
	pattern := regexp.MustCompile(`x*`)

	result := testBytes.Regexp().Find(pattern)
	if result.IsNone() {
		t.Fatal("Expected Some for empty match, got None")
	}
	if !result.Unwrap().Eq(g.Bytes("")) {
		t.Errorf("Expected empty match, got %q", result.Unwrap())
	}
}

func TestBytes_Regexp_ReplaceBy(t *testing.T) {
	testBytes := g.Bytes("a1b2c3")
	pattern := regexp.MustCompile(`\d`)

	result := testBytes.Regexp().ReplaceBy(pattern, func(m g.Bytes) g.Bytes {
		// Clone first: the match aliases the source backing array (see Bytes.Append docs).
		return m.Clone().Append(g.Bytes("!"))
	})

	expected := g.Bytes("a1!b2!c3!")
	if !result.Eq(expected) {
		t.Errorf("ReplaceBy = %q, want %q", result, expected)
	}
}

func TestBytes_Regexp_Split(t *testing.T) {
	testBytes := g.Bytes("a1b2c3d")
	pattern := regexp.MustCompile(`\d`)

	result := testBytes.Regexp().Split(pattern)
	expected := g.Slice[g.Bytes]{g.Bytes("a"), g.Bytes("b"), g.Bytes("c"), g.Bytes("d")}

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
	testBytes := g.Bytes("a1b2c3d")
	pattern := regexp.MustCompile(`\d`)

	result := testBytes.Regexp().SplitN(pattern, 2)
	if result.IsNone() {
		t.Fatal("Expected Some, got None")
	}

	parts := result.Unwrap()
	expected := g.Slice[g.Bytes]{g.Bytes("a"), g.Bytes("b2c3d")}
	if parts.Len() != expected.Len() {
		t.Fatalf("SplitN length = %d, want %d", parts.Len(), expected.Len())
	}
	for i := range expected {
		if !parts[i].Eq(expected[i]) {
			t.Errorf("SplitN[%d] = %q, want %q", i, parts[i], expected[i])
		}
	}

	// n == 0 yields an empty slice -> None
	if testBytes.Regexp().SplitN(pattern, 0).IsSome() {
		t.Error("SplitN(_, 0) should be None")
	}
}
