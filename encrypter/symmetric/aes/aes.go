package aes_lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"log"
)

func (ae *AESEncrypt) encrypt(data []byte) ([]byte, error) {
	//创建新的AES密钥块
	block, err := aes.NewCipher(ae.key)
	if err != nil {
		log.Fatal(err)
	}
	blockSize := block.BlockSize() // 获取秘钥块的长度
	iv := make([]byte, aes.BlockSize)
	if len(ae.iv) != 0 {
		copy(iv, ae.iv)
	}
	var adBytes = getAdditionalWithKey(ae.key)
	var out []byte

	switch ae.encryptType {
	case ECB:
		data = pkcs7Padding(data, blockSize) // 补全码
		blockMode := NewECBEncrypter(block)  // 加密模式
		out = make([]byte, len(data))        // 创建数组
		blockMode.CryptBlocks(out, data)     // 加密

		return out, nil
	case CBC:
		data = pkcs7Padding(data, blockSize)                           // 补全码
		blockMode := cipher.NewCBCEncrypter(block, ae.key[:blockSize]) // 加密模式
		out = make([]byte, len(data))                                  // 创建数组
		blockMode.CryptBlocks(out, data)                               // 加密

		return out, nil
	case CFB:
		out = make([]byte, len(data))
		cfb := cipher.NewCFBEncrypter(block, iv)
		cfb.XORKeyStream(out, data)

		return out, nil
	case OFB:
		out = make([]byte, len(data))
		cfb := cipher.NewOFB(block, iv)
		cfb.XORKeyStream(out, data)

		return out, nil
	case GCM:
		nonce := ae.key[:NONCE_LENGTH]
		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		out = aesGCM.Seal(nil, nonce, data, adBytes)

		return out, nil
	case GMCWithTagSize:
		//GCM模式初始化
		nonce := ae.key[:NONCE_LENGTH]
		aesGCM, err := cipher.NewGCMWithTagSize(block, TAG_LENGTH_BYTES)
		if err != nil {
			return nil, err
		}
		out = aesGCM.Seal(nil, nonce, data, adBytes)

		return out, nil
	case GMCWithNonceSize:
		//初始化给定密钥块的GCM模式，可提供加密和完整性检查
		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}

		//实现多次加密同一明文产生不同密文
		nonce := make([]byte, aesGCM.NonceSize())
		out = aesGCM.Seal(nonce, nonce, []byte(data), nil)

		return out, nil
	case CTR:
		ofb := cipher.NewCTR(block, iv)

		plaintext := data[0:len(data)]
		out = make([]byte, len(plaintext))
		ofb.XORKeyStream(out, []byte(plaintext))

		return out, nil
	default:
		out = make([]byte, len(data))
		block.Encrypt(out, data)
		return out, nil
	}

	return out, nil
}

func (ae *AESEncrypt) decrypt(data []byte) ([]byte, error) {
	//AES密码初始化
	block, err := aes.NewCipher(ae.key)
	if err != nil {
		log.Fatal(err)
	}
	blockSize := block.BlockSize() // 获取秘钥块的长度
	iv := make([]byte, aes.BlockSize)
	if len(ae.iv) != 0 {
		copy(iv, ae.iv)
	}
	var adBytes = getAdditionalWithKey(ae.key)
	var out []byte

	switch ae.encryptType {
	case ECB:
		blockMode := NewECBDecrypter(block) // 加密模式
		out = make([]byte, len(data))       // 创建数组
		blockMode.CryptBlocks(out, data)    // 解密
		out = pkcs7UnPadding(out)           // 补全码

		return out, nil
	case CBC:
		blockMode := cipher.NewCBCDecrypter(block, ae.key[:blockSize]) // 加密模式
		out = make([]byte, len(data))                                  // 创建数组
		blockMode.CryptBlocks(out, data)                               // 解密
		out = pkcs7UnPadding(out)                                      // 去除补全码

		return out, nil
	case CFB:
		cfbdec := cipher.NewCFBDecrypter(block, iv)
		out = make([]byte, len(data))
		cfbdec.XORKeyStream(out, data)

		return out, nil
	case OFB:
		ofb := cipher.NewOFB(block, iv)

		plaintext := data[0:len(data)]
		out = make([]byte, len(plaintext))
		ofb.XORKeyStream(out, []byte(plaintext))

		return out, nil
	case GCM:
		nonce := ae.key[:NONCE_LENGTH]
		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		out, err = aesGCM.Open(nil, nonce, data, adBytes)
		if err != nil {
			return nil, err
		}
		return out, nil
	case GMCWithTagSize:
		nonce := ae.key[:NONCE_LENGTH]
		aesGCM, err := cipher.NewGCMWithTagSize(block, TAG_LENGTH_BYTES)
		if err != nil {
			return nil, err
		}
		out, err = aesGCM.Open(nil, nonce, data, adBytes)
		if err != nil {
			return nil, err
		}
		return out, nil
	case GMCWithNonceSize:
		//GCM模式初始化
		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}

		//提取nonce
		nonceSize := aesGCM.NonceSize()
		if len(data) < nonceSize {
			return nil, fmt.Errorf("ciphertext too short")
		}
		//数据拆分为nonce和密文
		nonce, ciphertext := data[:nonceSize], data[nonceSize:]

		//密文解密
		out, err := aesGCM.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return nil, err
		}
		return out, nil
	case CTR:
		ofb := cipher.NewCTR(block, iv)

		plaintext := data[0:len(data)]
		out = make([]byte, len(plaintext))
		ofb.XORKeyStream(out, []byte(plaintext))

		return out, nil
	default:
		out = make([]byte, len(data))
		block.Decrypt(out, data)
		return out, nil
	}
}

func getAdditionalWithKey(key []byte) []byte {
	var adBytes [4]byte
	w := &bytes.Buffer{}
	length := uint32(len(key) + 1)

	binary.BigEndian.PutUint32(adBytes[:], length)
	if _, err := w.Write(adBytes[:]); err != nil {
		return nil
	}

	return adBytes[:]
}
