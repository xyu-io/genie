package sm4_lib

import (
	"encoding/json"
	"github.com/tjfoc/gmsm/sm4"
	"testing"
	"time"
)

var pubPEMData = []byte(`
WDhUa0MxZFFtZ2oyRzlaRw==
`)

var dMsg = map[string]any{}

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
	key := make([]byte, sm4.BlockSize) // 分组长度16字节长度，128位
	key = []byte(`AKS8uB2yXSMVGjKE`)

	t.Logf("key = %v\n", key)
	err := sm4.WriteKeyToPemFile("key.pem", key, nil)
	if err != nil {
		t.Fatalf("WriteKeyToPem error")
	}

	key, err = sm4.ReadKeyFromPemFile("key.pem", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("key = %v\n", key)
}

func TestSM4Encryption(t *testing.T) {
	// 加密
	sm4Cli, _ := New(
		WithType(CBC),
		WithKeyFromPem(pubPEMData),
	)
	msgBytes, _ := json.Marshal(msg)

	encryptText, err := sm4Cli.Encrypt(msgBytes)
	if err != nil {
		panic(err)
		return
	}
	t.Logf("密文： %s", string(encryptText))
	// 解密
	decryptText, err := sm4Cli.Decrypt(encryptText)
	if err != nil {
		panic(err)
		return
	}
	_ = json.Unmarshal(decryptText, &dMsg)
	jsonData, err := json.MarshalIndent(dMsg, "", "  ")
	t.Logf("解密： %s", string(jsonData))

	time.Sleep(time.Second)
}

func TestSM4EncryptionWithBase64(t *testing.T) {
	// 加密
	sm4Cli, _ := New(
		WithType(CBC),
		WithKeyFromPem(pubPEMData),
	)
	msgBytes, _ := json.Marshal(msg)

	encryptText, err := sm4Cli.EncryptWithBase64(msgBytes)
	if err != nil {
		panic(err)
		return
	}
	t.Logf("密文： %s", string(encryptText))
	// 解密
	decryptText, err := sm4Cli.DecryptWithBase64(string(encryptText))
	if err != nil {
		panic(err)
		return
	}
	_ = json.Unmarshal(decryptText, &dMsg)
	jsonData, err := json.MarshalIndent(dMsg, "", "  ")
	t.Logf("解密： %s", string(jsonData))

	time.Sleep(time.Second)
}

func TestSM4Login(t *testing.T) {
	key := []byte(`AKS8uB2yXSMVGjKE`)

	pem, _ := sm4.WriteKeyToPem(key, nil)
	sm4Cli, _ := New(
		WithType(ECB),
		WithKeyFromPem(pem),
	)
	encryptText := []byte("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")

	decryptText, err := sm4Cli.DecryptWithBase64(string(encryptText))
	if err != nil {
		panic(err)
		return
	}

	t.Logf("解密： %s", string(decryptText))

	time.Sleep(time.Second)
}
