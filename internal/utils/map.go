package utils

// MapValues get the values of map and return the slice
// Source - https://stackoverflow.com/a/71635953
// Posted by blackgreen, modified by community. See post 'Timeline' for change history
// Retrieved 2026-09-10, License - CC BY-SA 4.0
func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
