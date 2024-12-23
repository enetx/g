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
	collection := NewSlice[V](0)

	for v := range seq {
		if v.IsErr() {
			return Err[Slice[V]](v.Err())
		}
		collection = append(collection, v.Ok())
	}

	return Ok(collection)
}

func (seq SeqResult[V]) Count() Int {
	var counter Int
	seq(func(Result[V]) bool {
		counter++
		return true
	})

	return counter
}

func (seq SeqResult[V]) Map(transform func(V) V) SeqResult[V] { return transformResult(seq, transform) }

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

func (seq SeqResult[V]) Filter(fn func(V) bool) SeqResult[V] { return filterResult(seq, fn) }

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

func (seq SeqResult[V]) Exclude(fn func(V) bool) SeqResult[V] {
	return filterResult(seq, func(v V) bool { return !fn(v) })
}

func (seq SeqResult[V]) Dedup() SeqResult[V] {
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

func (seq SeqResult[V]) Unique() SeqResult[V] {
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

func (seq SeqResult[V]) Skip(n uint) SeqResult[V] {
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

func (seq SeqResult[V]) StepBy(n uint) SeqResult[V] {
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

func (seq SeqResult[V]) Take(n uint) SeqResult[V] {
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

func (seq SeqResult[V]) Chain(seqs ...SeqResult[V]) SeqResult[V] {
	return func(yield func(Result[V]) bool) {
		for _, seq := range append([]SeqResult[V]{seq}, seqs...) {
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

func (seq SeqResult[V]) Intersperse(sep V) SeqResult[V] {
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

func (seq SeqResult[V]) Inspect(fn func(v V)) SeqResult[V] {
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
