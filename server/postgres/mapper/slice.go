package mapper

func MapSlice[T any, K any](dbModels []T, mapper func(T) (K, error)) ([]K, error) {
	var models []K
	for _, m := range dbModels {
		mapped, err := mapper(m)

		if err != nil {
			return nil, err
		}

		models = append(models, mapped)
	}

	return models, nil
}

func MapSliceNoErr[T any, K any](objs []T, mapper func(T) K) []K {
	var models []K
	for _, m := range objs {
		mapped := mapper(m)

		models = append(models, mapped)
	}

	return models
}
