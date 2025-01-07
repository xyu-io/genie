package sm4_lib

import (
	"encoding/base64"
	"errors"
	"github.com/tjfoc/gmsm/sm4"
)

func (sm4enc *SM4Encrypt) Encrypt(data []byte) ([]byte, error) {
	if len(sm4enc.key) == 0 {
		return nil, errors.New("SM4Encrypt key is empty")
	}

	var err error
	var ecbMsg []byte
	switch sm4enc.encryptType {
	case ECB:
		ecbMsg, err = sm4.Sm4Ecb(sm4enc.key, data, true)
		if err != nil {
			return nil, err
		}
	case CBC:
		ecbMsg, err = sm4.Sm4Cbc(sm4enc.key, data, true)
		if err != nil {
			return nil, err
		}
	case CFB:
		ecbMsg, err = sm4.Sm4CFB(sm4enc.key, data, true)
		if err != nil {
			return nil, err
		}
	case OFB:
		ecbMsg, err = sm4.Sm4OFB(sm4enc.key, data, true)
		if err != nil {
			return nil, err
		}
	default:
		ecbMsg, err = sm4.Sm4Cbc(sm4enc.key, data, true)
		if err != nil {
			return nil, err
		}
	}
	return ecbMsg, nil
}

func (sm4enc *SM4Encrypt) EncryptWithBase64(data []byte) ([]byte, error) {

	ecbMsg, err := sm4enc.Encrypt(data)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(ecbMsg)))
	base64.StdEncoding.Encode(dst, ecbMsg)
	return dst, nil
}

func (sm4enc *SM4Encrypt) Decrypt(data []byte) ([]byte, error) {
	var err error
	var ecbDec []byte
	switch sm4enc.encryptType {
	case ECB:
		ecbDec, err = sm4.Sm4Ecb(sm4enc.key, data, false)
		if err != nil {
			return nil, err
		}
	case CBC:
		ecbDec, err = sm4.Sm4Cbc(sm4enc.key, data, false)
		if err != nil {
			return nil, err
		}
	case CFB:
		ecbDec, err = sm4.Sm4CFB(sm4enc.key, data, false)
		if err != nil {
			return nil, err
		}
	case OFB:
		ecbDec, err = sm4.Sm4OFB(sm4enc.key, data, false)
		if err != nil {
			return nil, err
		}
	default:
		ecbDec, err = sm4.Sm4Cbc(sm4enc.key, data, false)
		if err != nil {
			return nil, err
		}
	}
	return ecbDec, nil
}

func (sm4enc *SM4Encrypt) DecryptWithBase64(base64Data string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(base64Data))
	if err != nil {
		return nil, err
	}

	return sm4enc.Decrypt(data)
}

func (sm4enc *SM4Encrypt) EncryptWithGCM(a, data []byte) (d []byte, t []byte, err error) {
	d, t, err = sm4.Sm4GCM(sm4enc.key, sm4.IV, data, a, true)
	if err != nil {
		return nil, nil, err
	}

	return
}

func (sm4enc *SM4Encrypt) DecryptWithGCM(a, data []byte) (d []byte, t []byte, err error) {
	d, t, err = sm4.Sm4GCM(sm4enc.key, sm4.IV, data, a, false)
	if err != nil {
		return nil, nil, err
	}

	return
}
