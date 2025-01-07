package hash

import (
	"encoding/hex"
	"encoding/json"
	"github.com/tjfoc/gmsm/sm3"
)

// Sm3Hash encodes sm3 hash by bytes data
func Sm3Hash(data []byte) string {
	h := sm3.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// Sm3String encodes md5 hash by string data
func Sm3String(data string) string {
	return Sm3Hash([]byte(data))
}

func Sm3Data(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return Sm3Hash(b)
}

func Sm3Sum(data []byte) []byte {
	return sm3.Sm3Sum(data)
}
