package g

type chainIter[T any] struct {
	baseIter[T]
	iterators     []iterator[T]
	iteratorIndex int
}

func chain[T any](iterators ...iterator[T]) *chainIter[T] {
	iter := &chainIter[T]{iterators: iterators}
	iter.baseIter = baseIter[T]{iter}
	return iter
}

func (iter *chainIter[T]) Next() Option[T] {
	for {
		if iter.iteratorIndex == len(iter.iterators) {
			return None[T]()
		}

		if next := iter.iterators[iter.iteratorIndex].Next(); next.IsSome() {
			return next
		}

		iter.iteratorIndex++
	}
}
