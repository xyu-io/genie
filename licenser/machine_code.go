package licenser

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/shirou/gopsutil/host"
)

func MachineCode() (string, error) {
	sysUUID := ""
	hostInfo, err := host.Info()
	if err == nil {
		sysUUID = hostInfo.HostID
	}
	if sysUUID == "" {
		return "", errors.New("get machine uuid info failed")
	}
	return md5Code(sysUUID), nil
}

func md5Code(code string) string {
	hash := md5.Sum([]byte(code))
	return hex.EncodeToString(hash[:])
}
