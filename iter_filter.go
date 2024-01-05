package g

type filterIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	fun       func(T) bool
	exhausted bool
}

func filter[T any](iter iterator[T], fun func(T) bool) *filterIter[T] {
	iterator := &filterIter[T]{iter: iter, fun: fun}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *filterIter[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	for next := iter.iter.Next(); next.IsSome(); next = iter.iter.Next() {
		if iter.fun(next.Some()) {
			return next
		}
	}

	iter.exhausted = true

	return None[T]()
}

func exclude[T any](iter iterator[T], fun func(T) bool) *filterIter[T] {
	inverse := func(t T) bool { return !fun(t) }
	iterator := &filterIter[T]{iter: iter, fun: inverse}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}
