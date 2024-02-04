package response

type Data[T any] struct {
	Data T `json:"data"`
}

func Wrap[T any](t T) Data[T] {
	return Data[T]{Data: t}
}

func WrapList[T any](t []T) Data[[]T] {
	return Data[[]T]{Data: t}
}
