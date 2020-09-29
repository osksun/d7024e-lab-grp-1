package d7024e

func leftPad(str string, pad rune, lenght int) string {
	for i := len(str); i < lenght; i++ {
		str = string(pad) + str
	}
	return str
}
