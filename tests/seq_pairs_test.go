package g_test

import (
	"reflect"
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

// pairsSeq returns a SeqPairs[int, string] with pairs (1,a) (2,b) (3,c) (4,d).
func pairsSeq() SeqPairs[int, string] {
	return SliceOf(1, 2, 3, 4).Iter().Zip(SliceOf("a", "b", "c", "d").Iter())
}

// emptyPairsSeq returns an empty SeqPairs[int, string].
func emptyPairsSeq() SeqPairs[int, string] {
	return SliceOf[int]().Iter().Zip(SliceOf[string]().Iter())
}

func assertPairs(t *testing.T, got, want []Pair[int, string]) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected %d pairs, got %d: %v", len(want), len(got), got)
	}

	if len(want) != 0 && !reflect.DeepEqual(got, want) {
		t.Errorf("expected pairs %v, got %v", want, got)
	}
}

func TestSeqPairsKeys(t *testing.T) {
	keys := pairsSeq().Keys().Collect()
	if !keys.Eq(SliceOf(1, 2, 3, 4)) {
		t.Errorf("expected keys [1 2 3 4], got %v", keys)
	}

	if empty := emptyPairsSeq().Keys().Collect(); !empty.IsEmpty() {
		t.Errorf("expected no keys, got %v", empty)
	}
}

func TestSeqPairsValues(t *testing.T) {
	values := pairsSeq().Values().Collect()
	if !values.Eq(SliceOf("a", "b", "c", "d")) {
		t.Errorf("expected values [a b c d], got %v", values)
	}

	if empty := emptyPairsSeq().Values().Collect(); !empty.IsEmpty() {
		t.Errorf("expected no values, got %v", empty)
	}
}

func TestSeqPairsUnzip(t *testing.T) {
	keys, values := pairsSeq().Unzip()

	if !keys.Eq(SliceOf(1, 2, 3, 4)) {
		t.Errorf("expected keys [1 2 3 4], got %v", keys)
	}

	if !values.Eq(SliceOf("a", "b", "c", "d")) {
		t.Errorf("expected values [a b c d], got %v", values)
	}

	emptyKeys, emptyValues := emptyPairsSeq().Unzip()
	if !emptyKeys.IsEmpty() || !emptyValues.IsEmpty() {
		t.Errorf("expected empty unzip, got %v / %v", emptyKeys, emptyValues)
	}
}

func TestSeqPairsMap(t *testing.T) {
	doubled := pairsSeq().Map(func(k int, _ string) int { return k * 2 }).Collect()
	if !doubled.Eq(SliceOf(2, 4, 6, 8)) {
		t.Errorf("expected [2 4 6 8], got %v", doubled)
	}

	// Cross-type mapping: pairs to strings.
	joined := pairsSeq().Map(func(_ int, v string) string { return v + "!" }).Collect()
	if !joined.Eq(SliceOf("a!", "b!", "c!", "d!")) {
		t.Errorf("expected [a! b! c! d!], got %v", joined)
	}
}

func TestSeqPairsFilter(t *testing.T) {
	even := pairsSeq().Filter(func(k int, _ string) bool { return k%2 == 0 }).Collect()
	assertPairs(t, even, []Pair[int, string]{
		{Key: 2, Value: "b"},
		{Key: 4, Value: "d"},
	})

	none := pairsSeq().Filter(func(int, string) bool { return false }).Collect()
	assertPairs(t, none, nil)
}

func TestSeqPairsExclude(t *testing.T) {
	odd := pairsSeq().Exclude(func(k int, _ string) bool { return k%2 == 0 }).Collect()
	assertPairs(t, odd, []Pair[int, string]{
		{Key: 1, Value: "a"},
		{Key: 3, Value: "c"},
	})

	all := pairsSeq().Exclude(func(int, string) bool { return false }).Collect()
	if len(all) != 4 {
		t.Errorf("expected all 4 pairs, got %v", all)
	}
}

func TestSeqPairsTake(t *testing.T) {
	head := pairsSeq().Take(2).Collect()
	assertPairs(t, head, []Pair[int, string]{
		{Key: 1, Value: "a"},
		{Key: 2, Value: "b"},
	})

	assertPairs(t, pairsSeq().Take(0).Collect(), nil)

	over := pairsSeq().Take(10).Collect()
	if len(over) != 4 {
		t.Errorf("expected all 4 pairs, got %v", over)
	}
}

func TestSeqPairsSkip(t *testing.T) {
	tail := pairsSeq().Skip(2).Collect()
	assertPairs(t, tail, []Pair[int, string]{
		{Key: 3, Value: "c"},
		{Key: 4, Value: "d"},
	})

	all := pairsSeq().Skip(0).Collect()
	if len(all) != 4 {
		t.Errorf("expected all 4 pairs, got %v", all)
	}

	assertPairs(t, pairsSeq().Skip(10).Collect(), nil)
}

func TestSeqPairsTakeWhile(t *testing.T) {
	head := pairsSeq().TakeWhile(func(k int, _ string) bool { return k < 3 }).Collect()
	assertPairs(t, head, []Pair[int, string]{
		{Key: 1, Value: "a"},
		{Key: 2, Value: "b"},
	})

	// Elements after the first failing one are not yielded, even if they match.
	trailing := SliceOf(1, 2, 1).Iter().Zip(SliceOf("a", "b", "c").Iter()).
		TakeWhile(func(k int, _ string) bool { return k < 2 }).
		Collect()
	assertPairs(t, trailing, []Pair[int, string]{{Key: 1, Value: "a"}})

	assertPairs(t, pairsSeq().TakeWhile(func(int, string) bool { return false }).Collect(), nil)
	assertPairs(t, emptyPairsSeq().TakeWhile(func(int, string) bool { return true }).Collect(), nil)
}

func TestSeqPairsSkipWhile(t *testing.T) {
	tail := pairsSeq().SkipWhile(func(k int, _ string) bool { return k < 3 }).Collect()
	assertPairs(t, tail, []Pair[int, string]{
		{Key: 3, Value: "c"},
		{Key: 4, Value: "d"},
	})

	// Once the predicate fails, later matching elements are still yielded.
	trailing := SliceOf(1, 2, 1).Iter().Zip(SliceOf("a", "b", "c").Iter()).
		SkipWhile(func(k int, _ string) bool { return k < 2 }).
		Collect()
	assertPairs(t, trailing, []Pair[int, string]{
		{Key: 2, Value: "b"},
		{Key: 1, Value: "c"},
	})

	all := pairsSeq().SkipWhile(func(int, string) bool { return false }).Collect()
	if len(all) != 4 {
		t.Errorf("expected all 4 pairs, got %v", all)
	}

	assertPairs(t, pairsSeq().SkipWhile(func(int, string) bool { return true }).Collect(), nil)
}

func TestSeqPairsFind(t *testing.T) {
	found := pairsSeq().Find(func(k int, _ string) bool { return k == 3 })
	if found.IsNone() || found.Some() != (Pair[int, string]{Key: 3, Value: "c"}) {
		t.Errorf("expected Some({3 c}), got %v", found)
	}

	missing := pairsSeq().Find(func(k int, _ string) bool { return k == 42 })
	if missing.IsSome() {
		t.Errorf("expected None, got %v", missing)
	}
}

func TestSeqPairsForEach(t *testing.T) {
	var (
		sum    int
		concat string
	)

	pairsSeq().ForEach(func(k int, v string) {
		sum += k
		concat += v
	})

	if sum != 10 || concat != "abcd" {
		t.Errorf("expected sum 10 and concat abcd, got %d / %q", sum, concat)
	}
}

func TestSeqPairsCount(t *testing.T) {
	if count := pairsSeq().Count(); count != 4 {
		t.Errorf("expected count 4, got %d", count)
	}

	if count := emptyPairsSeq().Count(); count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
}

func TestSeqPairsInspect(t *testing.T) {
	var inspected int

	pairs := pairsSeq().
		Inspect(func(int, string) { inspected++ }).
		Collect()

	if inspected != 4 {
		t.Errorf("expected 4 inspected pairs, got %d", inspected)
	}

	if len(pairs) != 4 {
		t.Errorf("expected pairs to pass through unchanged, got %v", pairs)
	}
}

func TestSeqPairsFold(t *testing.T) {
	sum := pairsSeq().Fold(0, func(acc, k int, _ string) int { return acc + k })
	if sum != 10 {
		t.Errorf("expected sum 10, got %d", sum)
	}

	// The accumulator type may differ from the key and value types.
	concat := pairsSeq().Fold(String(""), func(acc String, _ int, v string) String { return acc + String(v) })
	if concat != "abcd" {
		t.Errorf("expected abcd, got %q", concat)
	}

	if init := emptyPairsSeq().Fold(42, func(acc, _ int, _ string) int { return acc + 1 }); init != 42 {
		t.Errorf("expected init 42 for empty sequence, got %d", init)
	}
}

func TestSeqPairsAll(t *testing.T) {
	if !pairsSeq().All(func(k int, _ string) bool { return k > 0 }) {
		t.Error("expected All to be true for positive keys")
	}

	if pairsSeq().All(func(k int, _ string) bool { return k > 1 }) {
		t.Error("expected All to be false when a key fails the predicate")
	}

	if !emptyPairsSeq().All(func(int, string) bool { return false }) {
		t.Error("expected All to be true for an empty sequence")
	}

	// All stops at the first failing pair.
	var checked int
	pairsSeq().All(func(int, string) bool { checked++; return false })
	if checked != 1 {
		t.Errorf("expected All to short-circuit after 1 pair, checked %d", checked)
	}
}

func TestSeqPairsAny(t *testing.T) {
	if !pairsSeq().Any(func(k int, v string) bool { return k == 3 && v == "c" }) {
		t.Error("expected Any to be true for pair (3, c)")
	}

	if pairsSeq().Any(func(k int, _ string) bool { return k > 10 }) {
		t.Error("expected Any to be false when no pair matches")
	}

	if emptyPairsSeq().Any(func(int, string) bool { return true }) {
		t.Error("expected Any to be false for an empty sequence")
	}

	// Any stops at the first matching pair.
	var checked int
	pairsSeq().Any(func(int, string) bool { checked++; return true })
	if checked != 1 {
		t.Errorf("expected Any to short-circuit after 1 pair, checked %d", checked)
	}
}

func TestSeqPairsCollect(t *testing.T) {
	assertPairs(t, pairsSeq().Collect(), []Pair[int, string]{
		{Key: 1, Value: "a"},
		{Key: 2, Value: "b"},
		{Key: 3, Value: "c"},
		{Key: 4, Value: "d"},
	})

	assertPairs(t, emptyPairsSeq().Collect(), nil)
}

func TestSeqPairsFilterByKey(t *testing.T) {
	keys, values := pairsSeq().FilterByKey(f.Gt(2)).Unzip()

	if !keys.Eq(SliceOf(3, 4)) {
		t.Errorf("FilterByKey(f.Gt(2)): expected keys [3 4], got %v", keys)
	}
	if !values.Eq(SliceOf("c", "d")) {
		t.Errorf("FilterByKey(f.Gt(2)): expected values [c d], got %v", values)
	}

	if n := emptyPairsSeq().FilterByKey(f.Gt(0)).Count(); n != 0 {
		t.Errorf("FilterByKey on empty seq: expected 0 pairs, got %d", n)
	}
}

func TestSeqPairsFilterByValue(t *testing.T) {
	keys, values := pairsSeq().FilterByValue(f.Eq("b")).Unzip()

	if !keys.Eq(SliceOf(2)) {
		t.Errorf("FilterByValue(f.Eq(\"b\")): expected keys [2], got %v", keys)
	}
	if !values.Eq(SliceOf("b")) {
		t.Errorf("FilterByValue(f.Eq(\"b\")): expected values [b], got %v", values)
	}
}
