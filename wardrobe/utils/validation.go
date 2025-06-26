package utils

func CheckIfEmpty[T any](data []T) any {
	if len(data) > 0 {
		return data
	}
	return nil
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
