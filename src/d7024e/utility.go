package d7024e

func leftPad(str string, pad rune, lenght int) string {
	for i := len(str); i < lenght; i++ {
		str = string(pad) + str
	}
	return str
}

// MinInt takes an arbitrary number of integers as in put and returns the smallest one
func MinInt(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}
