package utils

func Plural(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
