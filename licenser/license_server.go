package licenser

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	instance *LicServer
)

type LicServer struct {
	license *License

	key         []byte
	licensePath string
	app         string
}

func NewAuthServer(app, lic string, key []byte) *LicServer {
	instance = &LicServer{
		license:     &License{},
		licensePath: lic,
		key:         key,
		app:         app,
	}

	return instance
}

func GetLicServer() *LicServer {
	return instance
}

// VerifyLicense verify license used by middleware
func VerifyLicense() (string, error) {
	if instance == nil {
		return "", errors.New("license not init")
	}
	if err := instance.CheckLicense(); err != nil {
		log.Warn("product license check failed, data receiver is stop")
		mCode, codeErr := MachineCode()
		if codeErr != nil {
			log.Error(codeErr)
		} else {
			log.Infof("please update the product license file, machine code:  \n%s\n", mCode)
		}
		return mCode, err
	}
	return "", nil
}

// License get license
func (au *LicServer) License() License {

	return *au.license
}

// CheckLicense check license is invalid or expire
func (au *LicServer) CheckLicense() error {
	// 1.license exists check
	if au.license == nil || au.license.Expires == 0 {
		err := au.LoadLicense()
		if err != nil {
			return err
		}
	}

	// 2.product name is need not verify (it is verified by app start)
	// 3.license expire check
	if !au.checkExpire(*au.license) {
		return errors.New("license is expired or license is invalid")
	}

	return nil
}

// LoadLicense load license form file
func (au *LicServer) LoadLicense(handle ...handleFunc) error {
	// 1.load license form file
	licenseTxt, err := readFromFile(au.licensePath)
	if err != nil {
		return err
	}

	return au.licenseDecodeAndCheck(licenseTxt, au.key, handle...)
}

// decode license and check license is invalid
func (au *LicServer) licenseDecodeAndCheck(license string, key []byte, handle ...handleFunc) error {
	// 2.decrypt license code
	ciphertext, err := Decrypt(license, []byte(key))
	if err != nil {
		return err
	}

	// 3.decode license info
	var lic License
	err = json.Unmarshal([]byte(ciphertext), &lic)
	if err != nil {
		return err
	}

	// 4.check product name
	if !au.checkProduct(lic) {
		return errors.New("license does not match this product")
	}

	// 5.check machine uuid
	if !au.checkMachineCode(lic.MachineCode) {
		return errors.New("machine code is invalid or get system machine code failed")
	}

	// 6.license expire check
	if !au.checkExpire(lic) {
		return errors.New("license is expired or license is invalid")
	}

	// 7.middleware
	if len(handle) > 0 {
		for _, fn := range handle {
			fn(license)
		}
	}
	// 8.store to cache
	au.license = &lic
	return nil
}

// check machine uuid is invalid
func (au *LicServer) checkMachineCode(code string) bool {
	mCode, err := MachineCode()
	if err != nil || mCode != code {
		return false
	}

	return true
}

// check product name
func (au *LicServer) checkProduct(lic License) bool {
	if lic.AppName != au.app {
		return false
	}

	return true
}

// license expire check
func (au *LicServer) checkExpire(lic License) bool {
	newTime := time.Now().Unix()
	if lic.Expires < newTime || lic.Started > newTime {
		return false
	}
	return true
}
