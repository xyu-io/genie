package sm2_lib

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"log"
	"strings"
	"sync"
)

type Sm2Encryption struct {
	mu         sync.RWMutex
	PrivateKey *sm2.PrivateKey
	PublicKey  *sm2.PublicKey
}

var instance *Sm2Encryption

func NewSm2() *Sm2Encryption {
	if instance == nil {
		privateKey, _ := sm2.GenerateKey(rand.Reader) // 生成密钥对
		instance = &Sm2Encryption{
			PrivateKey: privateKey,
			PublicKey:  &sm2.PublicKey{},
		}
	}
	return instance
}

func (sm2x *Sm2Encryption) SetPublicKey(pubKeyStr string) (string, error) {
	sm2x.mu.Lock()
	defer sm2x.mu.Unlock()
	// 将公钥字节流解析为 big.Int 类型
	if !strings.Contains(pubKeyStr, "PUBLIC KEY") {
		pubKeyStr = "-----BEGIN PUBLIC KEY-----\r\n" +
			pubKeyStr + "\r\n" +
			"-----END PUBLIC KEY-----"
	}
	publicKeyBytes := []byte(pubKeyStr)
	publicKey, err := x509.ReadPublicKeyFromPem(publicKeyBytes)
	if err != nil {
		return "", err
	}
	sm2x.PublicKey = publicKey
	//sm2x.PublicKey = &sm2x.PrivateKey.PublicKey
	return sm2x.GetPublicKey(), nil
}

// GetPublicKey 获取公钥
func (sm2x *Sm2Encryption) GetPublicKey() string {
	if instance == nil {
		instance = NewSm2()
		return instance.GetPublicKey()
	}
	pub := sm2x.PrivateKey.PublicKey
	// 将公钥转换为字符串
	//pubKeyBytes := sm2.Compress(&pub)
	//pubKeyStr := hex.EncodeToString(pubKeyBytes)
	pubKeyToPem, err := x509.WritePublicKeyToPem(&pub)
	if err != nil {
		panic(err)
	}

	return withNoSign(string(pubKeyToPem))
}

// EncryptWithNoError sm2加密
func (sm2x *Sm2Encryption) EncryptWithNoError(msg string) string {
	sm2x.mu.Lock()
	defer sm2x.mu.Unlock()
	if &sm2x.PublicKey == nil || sm2x.PublicKey.Curve == nil {
		return msg
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return ""
	}
	cipherBytes, err := sm2.Encrypt(sm2x.PublicKey, data, nil, sm2.C1C3C2)
	if err != nil {
		log.Println(err)
		return msg
	}
	ciphertext := hex.EncodeToString(cipherBytes)
	return ciphertext
}

// Encrypt sm2加密
func (sm2x *Sm2Encryption) Encrypt(msg interface{}) (string, error) {
	sm2x.mu.Lock()
	defer sm2x.mu.Unlock()
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if sm2x.PublicKey.Y == nil {
		return "", errors.New("接口请求被拒绝，请先交换公钥")
	}
	cipherBytes, err := sm2.Encrypt(sm2x.PublicKey, data, nil, sm2.C1C3C2)
	if err != nil {
		log.Println("Encrypt : ", err)
		return "", err
	}
	ciphertext := hex.EncodeToString(cipherBytes)
	return ciphertext, nil
}

// Decrypt sm2解密
func (sm2x *Sm2Encryption) Decrypt(ciphertext string) ([]byte, error) {
	cipherByte, err := hex.DecodeString(ciphertext)
	if err != nil {
		fmt.Printf("hex.DecodeString: ", err.Error())
		return []byte(ciphertext), err
	}
	plaintext, err := sm2.Decrypt(sm2x.PrivateKey, cipherByte, sm2.C1C3C2)
	if err != nil {
		return []byte(ciphertext), err
	}

	return plaintext, nil
}
