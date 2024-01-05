package g

type mapIter[T, U any] struct {
	baseIter[U]
	iter      iterator[T]
	fun       func(T) U
	exhausted bool
}

func mapiter[T, U any](iter iterator[T], f func(T) U) *mapIter[T, U] {
	iterator := &mapIter[T, U]{iter: iter, fun: f}
	iterator.baseIter = baseIter[U]{iterator}

	return iterator
}

func transform[T any](iter iterator[T], op func(T) T) *mapIter[T, T] {
	return mapiter[T, T](iter, op)
}

func (iter *mapIter[T, U]) Next() Option[U] {
	if iter.exhausted {
		return None[U]()
	}

	next := iter.iter.Next()
	if next.IsNone() {
		iter.exhausted = true
		return None[U]()
	}

	return Some(iter.fun(next.Some()))
}
