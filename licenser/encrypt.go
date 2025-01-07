package licenser

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func Encrypt(data, key []byte) (string, error) {
	//创建新的AES密钥块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//初始化给定密钥块的GCM模式，可提供加密和完整性检查
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//实现多次加密同一明文产生不同密文
	nonce := make([]byte, aesGCM.NonceSize())
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(data), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(data string, key []byte) ([]byte, error) {
	// 解码 Base64 字符串
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	//AES密码初始化
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//GCM模式初始化
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//提取nonce
	nonceSize := aesGCM.NonceSize()
	if len(decoded) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	//数据拆分为nonce和密文
	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]

	//密文解密
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
