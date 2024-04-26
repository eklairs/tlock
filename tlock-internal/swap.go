package tlockinternal

func Swap[T any](slice []T, index1, index2 int) []T {
	slice[index1], slice[index2] = slice[index2], slice[index1]

	return slice
}
