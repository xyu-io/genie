package aes_lib

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"github.com/xyu-io/genie/outer"
	"testing"
)

var (
	key = []byte("ABCDEFGHIJKLMNOP") // 加密的密钥
)

var msg = map[string]any{
	"data": struct {
		User     string
		Password string
		TalentID int
		AuthKey  string
		String0  string
		String1  string
		String2  string
		String3  string
		String4  string
		String5  string
		String6  string
		String7  string
		Int0     int64
		Int1     int64
		Int2     int64
		Int3     int64
		Int4     int64
		Int5     int64
		Int6     int64
		Int7     int64
	}{
		User:     "TEST",
		Password: "XXXXWEFEWWEGER--",
		TalentID: 1234,
		AuthKey:  "XXXXXXXXXXXXX",
		String0:  "0",
		String1:  "1",
		String2:  "2",
		String3:  "3",
		String4:  "4",
		String5:  "5",
		String6:  "6",
		String7:  "7",
		Int0:     1234,
		Int1:     1235,
		Int2:     1235,
		Int3:     1235,
		Int4:     1235,
		Int5:     1235,
		Int6:     1235,
		Int7:     1235,
	},
	"code": 200,
	"msg":  "success",
}

func TestMakeKey(t *testing.T) {
	key, iv, err := generateKeyAndIV(aes.BlockSize)
	if err != nil {
		return
	}
	t.Log("aes_key", base64.StdEncoding.EncodeToString(key))
	t.Log("aes_key", base64.StdEncoding.EncodeToString(iv)) // 非ECB使用
}

func TestAesECB(t *testing.T) {
	aesCli, _ := New(
		WithType(ECB),
		WithKeyFromBytes([]byte(key)),
	)

	data, _ := json.Marshal(msg)
	encText, err := aesCli.encrypt([]byte(data))
	if err != nil {
		t.Error(err)
	}
	t.Logf("密文(b64)：%s", outer.EncodeWithBase64(encText))

	decText, err := aesCli.decrypt(encText)
	if err != nil {
		t.Error(err)
	}
	t.Logf("明文(utf8):%s", decText)
}

func TestAesCBC(t *testing.T) {
	aesCli, _ := New(
		WithType(CBC),
		WithKeyFromBytes([]byte(key)),
	)

	data, _ := json.Marshal(msg)
	encText, err := aesCli.encrypt([]byte(data))
	if err != nil {
		t.Error(err)
	}
	t.Logf("密文(b64)：%s", outer.EncodeWithBase64(encText))

	decText, err := aesCli.decrypt(encText)
	if err != nil {
		t.Error(err)
	}
	t.Logf("明文(utf8):%s", decText)
}

func TestAesCFB(t *testing.T) {
	aesCli, _ := New(
		WithType(CFB),
		WithKeyFromBytes([]byte(key)),
	)

	data, _ := json.Marshal(msg)
	encText, err := aesCli.encrypt([]byte(data))
	if err != nil {
		t.Error(err)
	}
	t.Logf("密文(b64)：%s", outer.EncodeWithBase64(encText))

	decText, err := aesCli.decrypt(encText)
	if err != nil {
		t.Error(err)
	}
	t.Logf("明文(utf8):%s", decText)
}

func TestAesOFB(t *testing.T) {
	aesCli, _ := New(
		WithType(OFB),
		WithKeyFromBytes([]byte(key)),
	)

	data, _ := json.Marshal(msg)
	encText, err := aesCli.encrypt([]byte(data))
	if err != nil {
		t.Error(err)
	}
	t.Logf("密文(b64)：%s", outer.EncodeWithBase64(encText))

	decText, err := aesCli.decrypt(encText)
	if err != nil {
		t.Error(err)
	}
	t.Logf("明文(utf8):%s", decText)
}

func TestAesGMC(t *testing.T) {
	aesCli, _ := New(
		WithType(GCM),
		WithKeyFromBytes([]byte(key)),
	)

	data, _ := json.Marshal(msg)
	encText, err := aesCli.encrypt([]byte(data))
	if err != nil {
		t.Error(err)
	}
	t.Logf("密文(b64)：%s", outer.EncodeWithBase64(encText))

	decText, err := aesCli.decrypt(encText)
	if err != nil {
		t.Error(err)
	}
	t.Logf("明文(utf8):%s", decText)
}

func TestAesGMCWithTagSize(t *testing.T) {
	aesCli, _ := New(
		WithType(GMCWithTagSize),
		WithKeyFromBytes([]byte(key)),
	)

	data, _ := json.Marshal(msg)
	encText, err := aesCli.encrypt([]byte(data))
	if err != nil {
		t.Error(err)
	}
	t.Logf("密文(b64)：%s", outer.EncodeWithBase64(encText))

	decText, err := aesCli.decrypt(encText)
	if err != nil {
		t.Error(err)
	}
	t.Logf("明文(utf8):%s", decText)
}

func TestAesGMCWithNonceSize(t *testing.T) {
	aesCli, _ := New(
		WithType(GMCWithNonceSize),
		WithKeyFromBytes([]byte(key)),
	)

	data, _ := json.Marshal(msg)
	encText, err := aesCli.encrypt([]byte(data))
	if err != nil {
		t.Error(err)
	}
	t.Logf("密文(b64)：%s", outer.EncodeWithBase64(encText))

	decText, err := aesCli.decrypt(encText)
	if err != nil {
		t.Error(err)
	}
	t.Logf("明文(utf8):%s", decText)
}

func TestAesCTR(t *testing.T) {
	aesCli, _ := New(
		WithType(CTR),
		WithKeyFromBytes([]byte(key)),
	)

	data, _ := json.Marshal(msg)
	encText, err := aesCli.encrypt([]byte(data))
	if err != nil {
		t.Error(err)
	}
	t.Logf("密文(b64)：%s", outer.EncodeWithBase64(encText))

	decText, err := aesCli.decrypt(encText)
	if err != nil {
		t.Error(err)
	}
	t.Logf("明文(utf8):%s", decText)
}
