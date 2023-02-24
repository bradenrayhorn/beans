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
