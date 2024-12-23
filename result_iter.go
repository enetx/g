package g

import (
	"iter"

	"github.com/enetx/g/f"
)

func (seq SeqResult[V]) Pull() (func() (Result[V], bool), func()) {
	return iter.Pull(iter.Seq[Result[V]](seq))
}

func (seq SeqResult[V]) All(fn func(v V) bool) Result[bool] {
	for v := range seq {
		if v.IsErr() {
			return Err[bool](v.Err())
		}

		if !fn(v.Ok()) {
			return Ok(false)
		}
	}

	return Ok(true)
}

func (seq SeqResult[V]) Any(fn func(v V) bool) Result[bool] {
	for v := range seq {
		if v.IsErr() {
			return Err[bool](v.Err())
		}

		if fn(v.Ok()) {
			return Ok(true)
		}
	}

	return Ok(false)
}

func (seq SeqResult[V]) Collect() Result[Slice[V]] {
	collection := make([]V, 0)
	var err error

	seq(func(v Result[V]) bool {
		if v.IsErr() {
			err = v.Err()
			return false
		}
		collection = append(collection, v.Ok())
		return true
	})

	if err != nil {
		return Err[Slice[V]](err)
	}

	return Ok[Slice[V]](collection)
}

func (seq SeqResult[V]) Map(transform func(V) V) SeqResult[V] { return transformResult(seq, transform) }

func (seq SeqResult[V]) Filter(fn func(V) bool) SeqResult[V] { return filterResult(seq, fn) }

func (seq SeqResult[V]) Exclude(fn func(V) bool) SeqResult[V] { return excludeResult(seq, fn) }

func (seq SeqResult[V]) Dedup() SeqResult[V] { return dedupResult(seq) }

func (seq SeqResult[V]) Unique() SeqResult[V] { return uniqueResult(seq) }

func (seq SeqResult[V]) ForEach(fn func(v Result[V])) {
	seq(func(v Result[V]) bool {
		fn(v)
		return true
	})
}

func (seq SeqResult[V]) Range(fn func(v Result[V]) bool) {
	seq(func(v Result[V]) bool {
		return fn(v)
	})
}

func (seq SeqResult[V]) Skip(n uint) SeqResult[V] { return skipResult(seq, n) }

func (seq SeqResult[V]) StepBy(n uint) SeqResult[V] { return stepbyResult(seq, n) }

func (seq SeqResult[V]) Take(n uint) SeqResult[V] { return takeResult(seq, n) }

func (seq SeqResult[V]) Chain(seqs ...SeqResult[V]) SeqResult[V] {
	return chainResult(append([]SeqResult[V]{seq}, seqs...)...)
}

func (seq SeqResult[V]) Intersperse(sep V) SeqResult[V] { return intersperseResult(seq, sep) }

func (seq SeqResult[V]) Inspect(fn func(v V)) SeqResult[V] { return inspectResult(seq, fn) }

func inspectResult[V any](seq SeqResult[V], fn func(V)) SeqResult[V] {
	return func(yield func(Result[V]) bool) {
		seq(func(v Result[V]) bool {
			if v.IsErr() {
				yield(v)
				return false
			}
			fn(v.Ok())
			return yield(v)
		})
	}
}

func transformResult[V, U any](seq SeqResult[V], fn func(V) U) SeqResult[U] {
	return func(yield func(Result[U]) bool) {
		seq(func(v Result[V]) bool {
			if v.IsErr() {
				yield(Err[U](v.Err()))
				return false
			}
			return yield(Ok(fn(v.Ok())))
		})
	}
}

func filterResult[V any](seq SeqResult[V], fn func(V) bool) SeqResult[V] {
	return func(yield func(Result[V]) bool) {
		seq(func(v Result[V]) bool {
			if v.IsErr() {
				yield(v)
				return false
			}
			if fn(v.Ok()) {
				return yield(v)
			}
			return true
		})
	}
}

func excludeResult[V any](seq SeqResult[V], fn func(V) bool) SeqResult[V] {
	return filterResult(seq, func(v V) bool { return !fn(v) })
}

func dedupResult[V any](seq SeqResult[V]) SeqResult[V] {
	var current V
	comparable := f.IsComparable(current)

	return func(yield func(Result[V]) bool) {
		seq(func(v Result[V]) bool {
			if v.IsErr() {
				yield(v)
				return false
			}

			if comparable {
				if f.Eq[any](current)(v.Ok()) {
					return true
				}
			} else {
				if f.Eqd(current)(v.Ok()) {
					return true
				}
			}

			current = v.Ok()
			return yield(v)
		})
	}
}

func uniqueResult[V any](seq SeqResult[V]) SeqResult[V] {
	return func(yield func(Result[V]) bool) {
		seen := NewSet[any]()
		seq(func(v Result[V]) bool {
			if v.IsErr() {
				yield(v)
				return false
			}
			if !seen.Contains(v.Ok()) {
				seen.Add(v.Ok())
				return yield(v)
			}
			return true
		})
	}
}

func takeResult[V any](seq SeqResult[V], n uint) SeqResult[V] {
	return func(yield func(Result[V]) bool) {
		seq(func(v Result[V]) bool {
			if v.IsErr() {
				yield(v)
				return false
			}
			if n == 0 {
				return false
			}
			n--
			return yield(v)
		})
	}
}

func stepbyResult[V any](seq SeqResult[V], n uint) SeqResult[V] {
	return func(yield func(Result[V]) bool) {
		i := uint(0)
		seq(func(v Result[V]) bool {
			if v.IsErr() {
				yield(v)
				return false
			}
			i++
			if (i-1)%n == 0 {
				return yield(v)
			}
			return true
		})
	}
}

func skipResult[V any](seq SeqResult[V], n uint) SeqResult[V] {
	return func(yield func(Result[V]) bool) {
		seq(func(v Result[V]) bool {
			if v.IsErr() {
				yield(v)
				return false
			}
			if n > 0 {
				n--
				return true
			}
			return yield(v)
		})
	}
}

func chainResult[V any](seqs ...SeqResult[V]) SeqResult[V] {
	return func(yield func(Result[V]) bool) {
		for _, seq := range seqs {
			seq(func(v Result[V]) bool {
				if v.IsErr() {
					yield(v)
					return false
				}
				return yield(v)
			})
		}
	}
}

func intersperseResult[V any](seq SeqResult[V], sep V) SeqResult[V] {
	return func(yield func(Result[V]) bool) {
		first := true

		seq(func(v Result[V]) bool {
			if v.IsErr() {
				yield(v)
				return false
			}

			if !first && !yield(Ok(sep)) {
				return false
			}

			first = false
			return yield(v)
		})
	}
}
