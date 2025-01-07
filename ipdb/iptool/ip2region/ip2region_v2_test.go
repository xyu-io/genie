package ip2region

import (
	"fmt"
	"testing"
)

var ipAddr = "34.117.118.44"

func TestIp2RegionV1(t *testing.T) {
	ip2Region, err := NewIp2RegionV1("./ip2region.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	details, err := ip2Region.Find(ipAddr) //34.117.118.44 117.187.194.117 23.45.66.23
	if err != nil {
		return
	}
	fmt.Println(details)
}

func TestIp2RegionV2(t *testing.T) {
	ip2Region, err := NewIp2RegionV2("./ip2region.xdb")
	if err != nil {
		fmt.Printf("open file err:%v\n", err)
		return
	}

	str, err := ip2Region.Find(ipAddr)
	if err != nil {
		return
	}
	fmt.Println(str)
}
