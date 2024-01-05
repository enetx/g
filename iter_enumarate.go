package g

type enumerateIter[T any] struct {
	baseIter[T]
	counter   uint
	exhausted bool
}

func enumerate[T any](iterator iterator[T]) *enumerateIter[T] {
	return &enumerateIter[T]{baseIter: baseIter[T]{iterator}}
}

func (iter *enumerateIter[T]) Next() Option[Map[uint, T]] {
	if iter.exhausted {
		return None[Map[uint, T]]()
	}

	next := iter.baseIter.Next()
	if next.IsNone() {
		iter.exhausted = true
		return None[Map[uint, T]]()
	}

	enext := Map[uint, T]{iter.counter: next.Some()}
	iter.counter++

	return Some(enext)
}
