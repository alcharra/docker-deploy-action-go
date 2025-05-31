package validator

import "strings"

func isBindMount(name string) bool {
	return strings.HasPrefix(name, "./") ||
		strings.HasPrefix(name, "/") ||
		strings.HasPrefix(name, "../") ||
		strings.HasPrefix(name, "~")
}
