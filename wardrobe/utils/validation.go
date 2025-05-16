package utils

func CheckIfEmpty[T any](data []T) any {
	if len(data) > 0 {
		return data
	}
	return nil
}
