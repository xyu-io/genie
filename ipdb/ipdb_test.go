package ipdb

import (
	"fmt"
	"testing"
)

var IPs = []string{"3.5.145.255", "34.117.118.44", "117.187.194.117", "23.45.66.23", "178.162.217.107", "5.79.71.225", "193.36.119.50", "2408:873d:22:1a01:0:ff:b087:eecc"}

func TestIpDB(t *testing.T) {
	Load("../com/iptool/db_data")

	for _, ip := range IPs {
		fmt.Println("====================")
		country, err := FindCountry("", ip)
		if err != nil {
			return
		}
		fmt.Println(ip, "---->>", country)
	}

}
