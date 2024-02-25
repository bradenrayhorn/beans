package fake

func filter[T any](s []T, filter func(it T) bool) []T {
	res := make([]T, 0)
	for _, v := range s {
		if filter(v) {
			res = append(res, v)
		}
	}
	return res
}

func mapVals[T any, K any](s []T, remap func(it T) K) []K {
	res := make([]K, 0)
	for _, v := range s {
		res = append(res, remap(v))
	}
	return res
}

func reduce[T any, K any](s []T, initial K, r func(it T, acc K) K) K {
	res := initial
	for _, v := range s {
		res = r(v, res)
	}
	return res
}

func find[T any](s []T, cmp func(it T) bool) *T {
	for _, v := range s {
		if cmp(v) {
			return &v
		}
	}
	return nil
}
