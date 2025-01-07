package aes_lib

import (
	"crypto/aes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	ECB = iota + 1 // 不安全
	CBC
	CFB
	OFB
	GCM
	GMCWithNonceSize
	GMCWithTagSize
	CTR
)

const (
	NONCE_LENGTH     = 12
	TAG_LENGTH_BYTES = 16
)

//var AdBytes = []byte("")

type AESEncrypt struct {
	key         []byte
	iv          []byte
	encryptType int
}

type Option func(*AESEncrypt)

func WithIV(iv []byte) Option {
	return func(aesenc *AESEncrypt) {
		aesenc.iv = iv
	}
}

func WithType(enType int) Option {
	return func(aesenc *AESEncrypt) {
		aesenc.encryptType = enType
	}
}

func WithKeyFromFile(file string) Option {
	return func(aesenc *AESEncrypt) {
		key, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("read key file err: %s", err.Error())
			return
		}
		aesenc.key = key
	}
}

func WithKeyFromBytes(pem []byte) Option {
	return func(aesenc *AESEncrypt) {
		aesenc.key = pem
	}
}

func New(opts ...Option) (*AESEncrypt, error) {

	ipt := &AESEncrypt{}

	for _, opt := range opts {
		opt(ipt)
	}

	return ipt, nil
}

// keySize is 8,12,16
func generateKeyAndIV(keySize int) ([]byte, []byte, error) {
	if keySize%aes.BlockSize != 0 && keySize > 16 {
		return nil, nil, errors.New("key size is not a multiple of the block size")
	}
	key := make([]byte, 32) // AES-128(16*8) AES-192(16*12) AES-256 (16*16)
	iv := make([]byte, aes.BlockSize)

	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, nil, err
	}

	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, nil, err
	}

	return key, iv, nil
}
