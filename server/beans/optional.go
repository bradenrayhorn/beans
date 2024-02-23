package beans

type Optional[T any] struct {
	set bool
	val T
}

func (n Optional[T]) Empty() bool {
	return !n.set
}

func (n Optional[T]) Value() (T, bool) {
	return n.val, n.set
}

func OptionalWrap[T any](val T) Optional[T] {
	return Optional[T]{set: true, val: val}
}
