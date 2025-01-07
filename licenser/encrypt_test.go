package licenser

import "testing"

func TestEncrypt(t *testing.T) {
	var dbPwd = "xaas-tdp-server@admin.123"
	var ckPwd = "Ck@remote#123"
	encrypt, err := Encrypt([]byte(dbPwd), pwdKey)
	if err != nil {
		return
	}
	t.Log("db", encrypt)

	encrypt, err = Encrypt([]byte(ckPwd), pwdKey)
	if err != nil {
		return
	}
	t.Log("ck", encrypt)
}
