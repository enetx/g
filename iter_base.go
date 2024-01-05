package g

import "context"

type iterator[T any] interface{ Next() Option[T] }

func collect[T any](iter iterator[T]) Slice[T] {
	values := NewSlice[T]()

	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		values = values.Append(next.Some())
	}

	return values
}

func fold[T, U any](iter iterator[T], init U, fn func(U, T) U) U {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		init = fn(init, next.Some())
	}

	return init
}

func find[T any](iter iterator[T], fn func(v T) bool) Option[T] {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		if fn(next.Some()) {
			return next
		}
	}

	return None[T]()
}

func foreach[T any](iter iterator[T], fn func(T)) {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		fn(next.Some())
	}
}

func rangeb[T any](iter iterator[T], fn func(T) bool) {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		if !fn(next.Some()) {
			break
		}
	}
}

func all[T any](iter iterator[T], fn func(T) bool) bool {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		if !fn(next.Some()) {
			return false
		}
	}

	return true
}

func anyb[T any](iter iterator[T], fn func(T) bool) bool {
	for next := iter.Next(); next.IsSome(); next = iter.Next() {
		if fn(next.Some()) {
			return true
		}
	}

	return false
}

func tochannel[T any](iter iterator[T], ctxs ...context.Context) chan T {
	ch := make(chan T)

	ctx := context.Background()
	if len(ctxs) != 0 {
		ctx = ctxs[0]
	}

	go func() {
		defer close(ch)

		for next := iter.Next(); next.IsSome(); next = iter.Next() {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- next.Some()
			}
		}
	}()

	return ch
}

type baseIter[T any] struct{ iterator[T] }

func (iter *baseIter[T]) All(fn func(T) bool) bool               { return all[T](iter, fn) }
func (iter *baseIter[T]) Any(fn func(T) bool) bool               { return anyb[T](iter, fn) }
func (iter *baseIter[T]) Collect() Slice[T]                      { return collect[T](iter.iterator) }
func (iter *baseIter[T]) Cycle() *cycleIter[T]                   { return cycle[T](iter) }
func (iter *baseIter[T]) Drop(n uint) *dropIter[T]               { return drop[T](iter, n) }
func (iter *baseIter[T]) Enumerate() *enumerateIter[T]           { return enumerate[T](iter) }
func (iter *baseIter[T]) Exclude(fn func(T) bool) *filterIter[T] { return exclude[T](iter, fn) }
func (iter *baseIter[T]) Filter(fn func(T) bool) *filterIter[T]  { return filter[T](iter, fn) }
func (iter *baseIter[T]) Flatten() *flattenIter[T]               { return flatten[T](iter) }
func (iter *baseIter[T]) Fold(init T, fn func(T, T) T) T         { return fold[T, T](iter, init, fn) }
func (iter *baseIter[T]) ForEach(fn func(T))                     { foreach[T](iter, fn) }
func (iter *baseIter[T]) Map(fn func(T) T) *mapIter[T, T]        { return transform[T](iter, fn) }
func (iter *baseIter[T]) Range(fn func(T) bool)                  { rangeb[T](iter, fn) }
func (iter *baseIter[T]) Take(n uint) *takeIter[T]               { return take[T](iter, n) }
func (iter *baseIter[T]) ToChannel() chan T                      { return tochannel[T](iter) }

func (iter *baseIter[T]) Chain(iterators ...iterator[T]) *chainIter[T] {
	return chain[T](append([]iterator[T]{iter}, iterators...)...)
}
