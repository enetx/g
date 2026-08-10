package rand_test

import (
	"math"
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/rand"
)

func TestN(t *testing.T) {
	// Test normal case where i > 0 should return random number in range [0, i)
	ten := Int(10)
	result := rand.N(ten.Std())
	if result < 0 || result >= 10 {
		t.Errorf("Random(10) should return value in [0, 10), got %d", result)
	}

	one := Int(1)
	result = rand.N(one.Std())
	if result != 0 {
		t.Errorf("Random(1) should return 0, got %d", result)
	}
}

func TestN_Panic(t *testing.T) {
	// rand.N matches math/rand/v2 semantics: n <= 0 panics.
	defer func() {
		if recover() == nil {
			t.Errorf("rand.N(0) must panic like math/rand/v2.N")
		}
	}()

	rand.N(Int(0))
}

func TestRangeInclusive(t *testing.T) {
	for range 100 {
		min := rand.N(Int(100))
		max := rand.N(Int(100)).Add(min)

		r := rand.RangeInclusive(min, max)
		if r.Lt(min) || r.Gt(max) {
			t.Errorf("RandomRange(%d, %d) = %d, want in range [%d, %d]", min, max, r, min, max)
		}
	}
}

func TestRangeInclusive_DistributionSanity(t *testing.T) {
	const (
		lo    = Int(-3)
		hi    = Int(3)
		iters = 1_000_00 // 100k
	)
	counts := make(map[Int]int)

	for range iters {
		x := rand.RangeInclusive(lo, hi)
		counts[x]++
	}

	want := float64(iters) / float64(hi-lo+1)
	tol := want * 0.05
	for v := lo; v <= hi; v++ {
		got := float64(counts[v])
		if diff := got - want; diff < -tol || diff > tol {
			t.Fatalf("value %d: got %d, want ~%.0f (+/-%0.f)", v, counts[v], want, tol)
		}
	}
}

func TestRangeInclusive_FullRangeMoves(t *testing.T) {
	lo := Int(math.MinInt64)
	hi := Int(math.MaxInt64)
	a := rand.RangeInclusive(lo, hi)
	b := rand.RangeInclusive(lo, hi)
	if a == b {
		t.Logf("two draws equal (ok but unlikely); a=%d b=%d", a, b)
	}
}

func TestShuffleSlice(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	rand.Shuffle(sl)

	if sl.Len() != 10 {
		t.Error("Expected length of 10, got ", sl.Len())
	}
}

func TestSample(t *testing.T) {
	sl := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := rand.Sample(sl, 5)

	if result.Len() != 5 {
		t.Errorf("Expected result length to be 5, got %d", result.Len())
	}

	for _, item := range result {
		if !sl.Contains(item) {
			t.Errorf("Expected result to contain only items from the original slice, got %d", item)
		}
	}
}

func TestShuffleMapOrd(t *testing.T) {
	mo := NewMapOrd[int, int]()
	for i := 1; i <= 5; i++ {
		mo.Insert(i, i)
	}

	clone := mo.Clone()

	// A permutation can legitimately equal the original (1/n! chance); retry a
	// few times so the test only fails when Shuffle is genuinely broken.
	changed := false
	for range 5 {
		rand.Shuffle(mo)
		if !mo.Eq(clone) {
			changed = true
			break
		}
	}

	if !changed {
		t.Errorf("The order of elements has not changed after 5 shuffles")
	}
}

func TestSecureBytes(t *testing.T) {
	b := rand.SecureBytes(32)
	if b.Len() != 32 {
		t.Fatalf("SecureBytes(32) length = %d, want 32", b.Len())
	}

	// Zero and negative lengths yield nil.
	if rand.SecureBytes(0) != nil {
		t.Error("SecureBytes(0) should be nil")
	}

	if rand.SecureBytes(-5) != nil {
		t.Error("SecureBytes(-5) should be nil")
	}

	// Two draws must differ (collision probability is negligible for 32 bytes).
	if b.Eq(rand.SecureBytes(32)) {
		t.Error("two SecureBytes(32) draws returned identical output")
	}
}

func TestSecureString(t *testing.T) {
	s := rand.SecureString(32)
	if len(s.Runes()) != 32 {
		t.Fatalf("SecureString(32) rune length = %d, want 32", len(s.Runes()))
	}

	// The default set is alphanumeric.
	const charset = ASCIILetters + Digits
	for _, r := range s.Runes() {
		if !charset.ContainsRune(r) {
			t.Fatalf("SecureString produced %q outside the default charset", r)
		}
	}

	// An explicit set restricts the output.
	code := rand.SecureString(6, Digits)
	if len(code.Runes()) != 6 {
		t.Fatalf("SecureString(6, Digits) rune length = %d, want 6", len(code.Runes()))
	}

	for _, r := range code.Runes() {
		if !Digits.ContainsRune(r) {
			t.Fatalf("SecureString(6, Digits) produced non-digit %q", r)
		}
	}

	// Zero/negative lengths and an empty explicit set yield "".
	if rand.SecureString(0) != "" || rand.SecureString(-1) != "" {
		t.Error("SecureString with non-positive length should be empty")
	}

	if rand.SecureString(10, "") != "" {
		t.Error("SecureString with an empty explicit set should be empty")
	}

	// Every character of a small set appears over enough draws — a smoke test
	// that selection is not stuck or truncating the set.
	seen := NewSet[rune]()
	for range 200 {
		for _, r := range rand.SecureString(5, Digits).Runes() {
			seen.Insert(r)
		}
	}

	if seen.Len() != 10 {
		t.Errorf("expected all 10 digits over 1000 draws, saw %d", seen.Len())
	}
}
