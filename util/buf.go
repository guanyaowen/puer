package util

func FilterNull(s []byte) []byte {
	c := byte(0)
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			s = s[:i]
			break
		}
	}
	return s
}
