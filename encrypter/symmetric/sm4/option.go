package sm4_lib

import (
	"fmt"
	"github.com/tjfoc/gmsm/sm4"
	"strings"
)

const (
	ECB int = iota + 1
	CBC
	CFB
	OFB
	GMC
)

type SM4Encrypt struct {
	key         []byte
	encryptType int
}

type Option func(*SM4Encrypt)

func WithType(enType int) Option {
	return func(sm4enc *SM4Encrypt) {
		sm4enc.encryptType = enType
	}
}

func WithKeyFromFile(file string) Option {
	return func(sm4enc *SM4Encrypt) {
		key, err := sm4.ReadKeyFromPemFile(file, nil)
		if err != nil {
			fmt.Printf("read sm4 key file, %s", err.Error())
			return
		}
		sm4enc.key = key
	}
}

func WithKeyFromPem(pem []byte) Option {
	return func(sm4enc *SM4Encrypt) {
		key, err := sm4.ReadKeyFromPem(pem, nil)
		if err != nil {
			if !strings.Contains(string(pem), "SM4 KEY") {
				newPem := "-----BEGIN SM4 KEY-----\r\n" +
					string(pem) + "\r\n" +
					"-----END SM4 KEY-----"
				key, err = sm4.ReadKeyFromPem([]byte(newPem), nil)
				sm4enc.key = key
				return
			}

			return
		}
		sm4enc.key = key
	}
}

func New(opts ...Option) (*SM4Encrypt, error) {

	ipt := &SM4Encrypt{}

	for _, opt := range opts {
		opt(ipt)
	}
	return ipt, nil
}
