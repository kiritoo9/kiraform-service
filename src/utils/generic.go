package utils

func Find[T any](slice []T, match func(T) bool) *T {
	for i := range slice {
		if match(slice[i]) {
			return &slice[i]
		}
	}
	return nil
}
