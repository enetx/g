package g

type takeIter[T any] struct {
	baseIter[T]
	iter  iterator[T]
	limit uint
}

func take[T any](iter iterator[T], limit uint) *takeIter[T] {
	iterator := &takeIter[T]{iter: iter, limit: limit}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *takeIter[T]) Next() Option[T] {
	if iter.limit == 0 {
		return None[T]()
	}

	next := iter.iter.Next()
	if next.IsNone() {
		iter.limit = 0
	} else {
		iter.limit--
	}

	return next
}
