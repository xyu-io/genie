package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"testing"
)

func TestGenRsaKey(t *testing.T) {
	// 生成 RSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Logf("生成 RSA 密钥对失败: %v", err)
	}
	// 将私钥转换为 ASN.1 PKCS#1 DER 编码
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	// 将 DER 编码的私钥转换为 PEM 格式
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	})

	// 将公钥提取为 *rsa.PublicKey 类型
	publicKey := &privateKey.PublicKey
	// 将公钥转换为 ASN.1 PKIX DER 编码
	pubDER := x509.MarshalPKCS1PublicKey(publicKey)
	// 将 DER 编码的公钥转换为 PEM 格式
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDER,
	})
	// 打印私钥和公钥
	fmt.Println("私钥:")
	fmt.Println(string(privPEM))
	err = os.WriteFile("private.pem", privPEM, 0666)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("公钥:")
	fmt.Println(string(pubPEM))
	err = os.WriteFile("public.pem", pubPEM, 0666)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEncrypt(t *testing.T) {
	// 待加密的数据
	data := map[string]string{
		"data": "hello",
	}
	message, _ := json.Marshal(data)

	publicKey, err := publicKeyFromFile("public.pem")
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := privateKeyFromFile("private.pem")
	if err != nil {
		t.Fatal(err)
	}

	// 使用公钥和 OAEP 填充方案加密数据
	label := []byte("OAEP Encrypted")
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encrypting message: %s\n", err)
		return
	}
	fmt.Printf("Ciphertext: %x\n", ciphertext)

	// 使用私钥和 OAEP 填充方案解密数据
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decrypting message: %s\n", err)
		return
	}
	fmt.Printf("Plaintext: %s\n", plaintext)
}
