package tlockinternal

// Removes an index from a slice
func Remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
