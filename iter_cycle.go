package g

type cycleIter[T any] struct {
	baseIter[T]
	iter  iterator[T]
	items []T
	index int
}

func cycle[T any](iter iterator[T]) *cycleIter[T] {
	iterator := &cycleIter[T]{iter: iter, items: make([]T, 0)}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *cycleIter[T]) Next() Option[T] {
	if iter.iter != nil {
		if next := iter.iter.Next(); next.IsSome() {
			iter.items = append(iter.items, next.Some())
			return next
		}

		iter.iter = nil
	}

	if len(iter.items) == 0 {
		return None[T]()
	}

	if iter.index == len(iter.items) {
		iter.index = 0
	}

	next := iter.items[iter.index]
	iter.index++

	return Some(next)
}
