package g_test

import (
	"errors"
	"testing"

	. "github.com/enetx/g"
)

// Lazy SeqResult transformers must be consumer-driven on Err (Rust iterator
// semantics): an Err element is yielded downstream like any other item, and
// the source keeps iterating as long as the consumer keeps accepting values.
// Fail-fast remains the consumer's choice (break / TryCollect / Fold / ...).

var errCont = errors.New("boom")

// mixed yields Ok(1), Err, Ok(2) and reports whether the element after the
// Err was pulled from the source.
func mixed(reachedAfterErr *bool) SeqResult[int] {
	return func(yield func(Result[int]) bool) {
		if !yield(Ok(1)) {
			return
		}
		if !yield(Err[int](errCont)) {
			return
		}
		*reachedAfterErr = true
		yield(Ok(2))
	}
}

// collectAll drains a sequence with an always-continue consumer and returns
// (ok values, err count).
func collectAll[V any](seq SeqResult[V]) ([]V, int) {
	var oks []V
	errs := 0

	seq.Range(func(r Result[V]) bool {
		if r.IsErr() {
			errs++
		} else {
			oks = append(oks, r.Ok())
		}
		return true
	})

	return oks, errs
}

func TestSeqResultMapContinuesPastErr(t *testing.T) {
	reached := false
	oks, errs := collectAll(mixed(&reached).Map(func(v int) int { return v * 10 }))

	if !reached {
		t.Fatal("Map: source not pulled past Err — transformer is not consumer-driven")
	}
	if errs != 1 || len(oks) != 2 || oks[0] != 10 || oks[1] != 20 {
		t.Fatalf("Map: got oks=%v errs=%d, want [10 20] / 1", oks, errs)
	}
}

func TestSeqResultFlatMapContinuesPastErr(t *testing.T) {
	reached := false
	src := SeqResult[Slice[int]](func(yield func(Result[Slice[int]]) bool) {
		if !yield(Ok(Slice[int]{1, 2})) {
			return
		}
		if !yield(Err[Slice[int]](errCont)) {
			return
		}
		reached = true
		yield(Ok(Slice[int]{3}))
	})

	oks, errs := collectAll(src.FlatMap(Slice[int].Iter))

	if !reached {
		t.Fatal("FlatMap: source not pulled past Err — transformer is not consumer-driven")
	}
	if errs != 1 || len(oks) != 3 || oks[0] != 1 || oks[1] != 2 || oks[2] != 3 {
		t.Fatalf("FlatMap: got oks=%v errs=%d, want [1 2 3] / 1", oks, errs)
	}
}

func TestSeqResultFilterExcludeContinuePastErr(t *testing.T) {
	reached := false
	oks, errs := collectAll(mixed(&reached).Filter(func(v int) bool { return v > 0 }))
	if !reached || errs != 1 || len(oks) != 2 {
		t.Fatalf("Filter: reached=%v oks=%v errs=%d, want true / [1 2] / 1", reached, oks, errs)
	}

	reached = false
	oks, errs = collectAll(mixed(&reached).Exclude(func(v int) bool { return v < 0 }))
	if !reached || errs != 1 || len(oks) != 2 {
		t.Fatalf("Exclude: reached=%v oks=%v errs=%d, want true / [1 2] / 1", reached, oks, errs)
	}
}

func TestSeqResultDedupUniqueContinuePastErr(t *testing.T) {
	reached := false
	oks, errs := collectAll(mixed(&reached).Dedup())
	if !reached || errs != 1 || len(oks) != 2 {
		t.Fatalf("Dedup: reached=%v oks=%v errs=%d, want true / [1 2] / 1", reached, oks, errs)
	}

	reached = false
	oks, errs = collectAll(mixed(&reached).Unique())
	if !reached || errs != 1 || len(oks) != 2 {
		t.Fatalf("Unique: reached=%v oks=%v errs=%d, want true / [1 2] / 1", reached, oks, errs)
	}
}

func TestSeqResultSkipStepByTakeContinuePastErr(t *testing.T) {
	reached := false
	oks, errs := collectAll(mixed(&reached).Skip(1))
	if !reached || errs != 1 || len(oks) != 1 || oks[0] != 2 {
		t.Fatalf("Skip: reached=%v oks=%v errs=%d, want true / [2] / 1", reached, oks, errs)
	}

	reached = false
	oks, errs = collectAll(mixed(&reached).StepBy(1))
	if !reached || errs != 1 || len(oks) != 2 {
		t.Fatalf("StepBy: reached=%v oks=%v errs=%d, want true / [1 2] / 1", reached, oks, errs)
	}

	reached = false
	oks, errs = collectAll(mixed(&reached).Take(2))
	if !reached || errs != 1 || len(oks) != 2 {
		t.Fatalf("Take: reached=%v oks=%v errs=%d, want true / [1 2] / 1 (Err not counted toward n)", reached, oks, errs)
	}
}

func TestSeqResultChainContinuesPastErr(t *testing.T) {
	reached := false
	second := SeqResult[int](func(yield func(Result[int]) bool) { yield(Ok(3)) })

	oks, errs := collectAll(mixed(&reached).Chain(second))
	if !reached || errs != 1 || len(oks) != 3 || oks[2] != 3 {
		t.Fatalf("Chain: reached=%v oks=%v errs=%d, want true / [1 2 3] / 1", reached, oks, errs)
	}
}

func TestSeqResultIntersperseInspectScanContinuePastErr(t *testing.T) {
	reached := false
	oks, errs := collectAll(mixed(&reached).Intersperse(0))
	if !reached || errs != 1 || len(oks) != 3 {
		t.Fatalf("Intersperse: reached=%v oks=%v errs=%d, want true / [1 0 2] / 1", reached, oks, errs)
	}

	reached = false
	seen := 0
	oks, errs = collectAll(mixed(&reached).Inspect(func(int) { seen++ }))
	if !reached || errs != 1 || len(oks) != 2 || seen != 2 {
		t.Fatalf("Inspect: reached=%v oks=%v errs=%d seen=%d, want true / [1 2] / 1 / 2", reached, oks, errs, seen)
	}

	reached = false
	oks, errs = collectAll(mixed(&reached).Scan(0, func(acc, v int) int { return acc + v }))
	// init(0), acc after 1 (=1), Err passed through, acc after 2 (=3)
	if !reached || errs != 1 || len(oks) != 3 || oks[2] != 3 {
		t.Fatalf("Scan: reached=%v oks=%v errs=%d, want true / [0 1 3] / 1", reached, oks, errs)
	}
}

// Consumer-side fail-fast must still work: breaking on the Err element stops
// the source exactly there.
func TestSeqResultConsumerBreakStillFailFast(t *testing.T) {
	reached := false
	var got []Result[int]

	mixed(&reached).Map(func(v int) int { return v }).Range(func(r Result[int]) bool {
		got = append(got, r)
		return r.IsOk() // break on Err
	})

	if reached {
		t.Fatal("consumer break: source was pulled past Err despite consumer stopping")
	}
	if len(got) != 2 || !got[1].IsErr() {
		t.Fatalf("consumer break: got %v, want [Ok(1) Err]", got)
	}
}

// Terminal fail-fast contracts are unchanged: TryCollect short-circuits at
// the first Err even though transformers now continue.
func TestSeqResultTryCollectStillFailFast(t *testing.T) {
	reached := false
	r := mixed(&reached).Map(func(v int) int { return v }).TryCollect()

	if r.IsOk() {
		t.Fatalf("TryCollect: want Err, got %v", r.Ok())
	}
	if reached {
		t.Fatal("TryCollect: source was pulled past Err — terminal must short-circuit")
	}
}
