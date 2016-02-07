package utils

func StringInSlice(valueToCheck string, array []string) (res bool) {
	res = false
	for _, v := range array {
		if v == valueToCheck {
			res = true
		}
	}
	return
}
