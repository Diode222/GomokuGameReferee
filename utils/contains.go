package utils

func StringContains(str string, strs []string) bool {
	for _, s := range strs {
		if str == s {
			return true
		}
	}
	return false
}
