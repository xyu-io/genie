package sm2_lib

import "strings"

func withNoSign(pubKey string) string {
	str := strings.Split(pubKey, "\n")
	if len(str) >= 3 {
		return str[1] + str[2]
	}
	return pubKey
}
