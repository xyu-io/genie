package geoip

import (
	"fmt"
	"testing"
)

func TestGeoip(t *testing.T) {
	geoIP, err := NewGeoIP("./dbip.mmdb")
	if err != nil {
		fmt.Println(err)
		return
	}
	find, err := geoIP.Find("34.117.118.44")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(find)
}
