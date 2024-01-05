package g

type liftIter[T any] struct {
	baseIter[T]
	items Slice[T]
	index int
}

func lift[T any](items []T) *liftIter[T] {
	iterator := &liftIter[T]{items: items}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *liftIter[T]) Next() Option[T] {
	if iter.index >= iter.items.Len() {
		return None[T]()
	}

	iter.index++

	return Some(iter.items[iter.index-1])
}
