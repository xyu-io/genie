package licenser

import (
	"genie/looper"
	"testing"
	"time"
)

func TestAppLicense(t *testing.T) {
	t.Log("-----------user-----------")
	machineCode, err := MachineCode()
	if err != nil {
		return
	}
	t.Log("get machine code:", machineCode)

	t.Log("-----------agent-----------")
	expireTime := time.Now().Add(time.Duration(time.Second * 10)).Unix() //
	config := Option{
		App:         analyser,
		Org:         "ysgj.com",
		User:        "xuanyu.li",
		Expires:     expireTime,
		LicensePath: licPath,
		MachineCode: machineCode,
		Key:         key,
	}
	licAgent := NewAuthAgent(config)
	license, err := licAgent.MakeLicense()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("license:", license)
	t.Log("expire time: 10s")

	t.Log("-----------server-----------")
	licServer := NewAuthServer("agent", licPath, key)
	// 模拟重启程序时候授权接收和存储
	dbHandle := func(license string) {
		t.Log("license of store:", license)
	}

	err = licServer.LoadLicense(dbHandle)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("decode license:%+v\n", licServer.License())
	t.Log("app is running now, check internal is 2s")

	// 模拟程序秩序运行中的授权检测
	looper.TimeLoopThen(time.Second*2, true, func(ti time.Time) {
		// 检测过期
		err = GetLicServer().CheckLicense()
		if err != nil {
			t.Log(err)
			mCode, err := MachineCode()
			if err != nil {
				t.Log(err)
			}
			t.Logf("please update the authorization license file, machine code:  \n%s\n", mCode)
			return
		}
		t.Log("app check license status:", true)
	})

	time.Sleep(time.Second * 15)
}
