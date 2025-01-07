package ip2location

import (
	"fmt"
	"testing"
)

func TestIp2location(t *testing.T) {
	location, err := NewIP2Location("./IP2LOCATION-LITE-DB11.IPV6.BIN")
	if err != nil {
		return
	}
	find, err := location.Find("117.187.194.117")
	if err != nil {
		return
	}
	fmt.Println(find)
}
