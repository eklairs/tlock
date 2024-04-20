package tlockvault

// Removes an index from a slice
func remove[T any](slice []T, s int) []T {
    return append(slice[:s], slice[s+1:]...)
}

