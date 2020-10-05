package d7024e

func leftPad(str string, pad rune, length int) string {
	for i := len(str); i < length; i++ {
		str = string(pad) + str
	}
	return str
}

// MinInt takes an arbitrary number of integers as in put and returns the smallest one
func MinInt(vars ...int) int {
	min := vars[0]
	for i := 0; i < len(vars); i++ {
		if min > vars[i] {
			min = vars[i]
		}
	}
	return min
}
