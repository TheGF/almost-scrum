package functional

// ApplyStrings applies a function to a list of strings.
func ApplyStrings(xs []string, apply func(string) string) []string {
	result := make([]string, 0, len(xs))
	for _, x := range xs {
		result = append(result, apply(x))
	}
	return result
}

// FilterString filters from a slice all items that do not satisfy the given filter
func FilterString(xs []string, filter func(string) bool) []string {
	result := make([]string, 0, len(xs))
	for _, x := range xs {
		if filter(x) {
			result = append(result, x)
		}
	}
	return result
}
