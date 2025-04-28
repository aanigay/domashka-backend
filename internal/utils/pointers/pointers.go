package pointers

func To[T any](v T) *T {
	return &v
}

func From[T any](ptr *T) T {
	if ptr == nil {
		return *(new(T))
	}

	return *ptr
}
