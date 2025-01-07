package licenser

import (
	"encoding/json"
	"github.com/go-logr/logr"
	"time"
)

type LicAgent struct {
	log logr.Logger
	Option
}

type Option struct {
	App         string
	Org         string
	User        string
	Expires     int64
	LicensePath string
	MachineCode string
	Key         []byte
}

func NewAuthAgent(op Option) *LicAgent {
	return &LicAgent{
		Option: op,
	}
}

// MakeLicense make license and write to file.
// return license code
func (au *LicAgent) MakeLicense() (string, error) {
	lic := License{
		AppName:     au.App,
		OrgName:     au.Org,
		UserName:    au.Org,
		MachineCode: au.MachineCode,
		Started:     time.Now().Unix(),
		Expires:     au.Expires,
	}

	licBytes, _ := json.Marshal(lic)

	license, err := Encrypt([]byte(licBytes), []byte(au.Key))
	if err != nil {
		return "", err
	}

	err = writeToFile(licPath, license)
	if err != nil {
		return "", err
	}

	return license, nil
}
