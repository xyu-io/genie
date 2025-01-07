package iptool

import (
	"fmt"
	"github.com/xyu-io/genie/ipdb/iptool/geoip"
	"github.com/xyu-io/genie/ipdb/iptool/ip2location"
	"github.com/xyu-io/genie/ipdb/iptool/ip2region"
	"github.com/xyu-io/genie/ipdb/iptool/qqwry"
	"testing"
)

var IPs = []string{"34.117.118.44", "117.187.194.117", "23.45.66.23", "178.162.217.107", "5.79.71.225", "193.36.119.50", "2408:873d:22:1a01:0:ff:b087:eecc"}

func TestIpTool(t *testing.T) {
	for _, ip := range IPs {
		geoipIns, err := geoip.NewGeoIP("./geoip/dbip.mmdb")
		if err != nil {
			t.Error(err)
		}
		find, err := geoipIns.Find(ip)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println("geoip>> ", find)
		}

		ip2Loc, err := ip2location.NewIP2Location("./ip2location/IP2LOCATION-LITE-DB11.IPV6.BIN")
		find, err = ip2Loc.Find(ip)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println("ip2location>> ", find)
		}

		ip2RegionV1, err := ip2region.NewIp2RegionV1("./ip2region/ip2region.db")
		find, err = ip2RegionV1.Find(ip)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println("ip2RegionV1>> ", find)
		}

		ip2RegionV2, err := ip2region.NewIp2RegionV2("./ip2region/ip2region.xdb")
		find, err = ip2RegionV2.Find(ip)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println("ip2RegionV2>> ", find)
		}

		qqwry, err := qqwry.NewQQwry("./qqwry/qqwry_utf.dat")
		find, err = qqwry.Find(ip)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println("qqwry>> ", find)
		}
	}
}
