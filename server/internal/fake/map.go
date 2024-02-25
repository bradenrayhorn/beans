package fake

func values[T any, K comparable](m map[K]T) []T {
	v := make([]T, 0, len(m))
	for _, value := range m {
		v = append(v, value)
	}
	return v
}
