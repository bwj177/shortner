package base62

import (
	"math"
	"strings"
)

// 62进制 0-9|A-Z|a-z
const base62String = `defgh345ABCDERSTUV09abc6GHklmnijr12opq78stuFIJKLxyzWXYZMNOPQvw` //乱序

func GetBase62(n uint64) string {

	res := []byte{}
	if n == 0 {
		return string(base62String[0])
	}
	for n > 0 {
		yushu := n % 62
		res = append(res, base62String[yushu])

		div := n / 62
		n = div
	}
	String62 := reverseByteToString(res)
	return String62
}

func StrToint64(s string) uint64 {
	res := 0
	l := len(s)
	for i := 0; i < len(s); i++ {
		res += strings.Index(base62String, string(s[l-1-i])) * int(math.Pow(62, float64(i)))
	}
	return uint64(res)
}

func reverseByteToString(s []byte) string {
	l := len(s)
	for i := 0; i < l/2; i++ {
		s[i], s[l-i-1] = s[l-1-i], s[i]
	}
	return string(s)
}
