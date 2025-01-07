package outer

import (
	"encoding/base64"
	"encoding/hex"
)

func EncodeWithBase64(data []byte) []byte {
	// 编码 Base64 字符串并返回
	out := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(out, data)
	return out
}

func DecodeWithBase64(data []byte) []byte {
	// 编码 Base64 字符串并返回
	out := make([]byte, len(data))
	_, err := base64.StdEncoding.Decode(out, data)
	if err != nil {
		return nil
	}

	return out
}

func EncodeWithHex(data []byte) string {
	return hex.EncodeToString(data)
}

func DecodeWithHex(data string) ([]byte, error) {
	out, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return out, nil
}
