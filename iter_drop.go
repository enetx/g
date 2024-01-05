package g

type dropIter[T any] struct {
	baseIter[T]
	iter      iterator[T]
	count     uint
	dropped   bool
	exhausted bool
}

func drop[T any](iter iterator[T], count uint) *dropIter[T] {
	iterator := &dropIter[T]{iter: iter, count: count}
	iterator.baseIter = baseIter[T]{iterator}

	return iterator
}

func (iter *dropIter[T]) Next() Option[T] {
	if iter.exhausted {
		return None[T]()
	}

	if !iter.dropped {
		iter.dropped = true

		for i := uint(0); i < iter.count; i++ {
			if iter.delegateNext().IsNone() {
				return None[T]()
			}
		}
	}

	return iter.delegateNext()
}

func (iter *dropIter[T]) delegateNext() Option[T] {
	next := iter.iter.Next()
	if next.IsNone() {
		iter.exhausted = true
	}

	return next
}
