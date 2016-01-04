package common

func String(s string) *string {
	str := new(string)
	*str = s
	return str
}
