package ip2region

import (
	"genie/ipdb/iptool/ip_base"
	"log"
	"sync"
)

type MyRegion struct {
	RegionDB *Ip2RegionV1
	Lock     sync.RWMutex
	closed   bool
}

func NewIp2RegionV1(path string) (*MyRegion, error) {
	instance := &MyRegion{}
	instance.RegionDB, err = New(path)
	defer instance.RegionDB.Close()
	if err != nil {
		log.Fatalln("func-InitIPdb-error:", err)
		return nil, err
	}
	return instance, nil
}

func (region *MyRegion) Name() string {
	return "ip2location"
}

func (region *MyRegion) IsChinaIP(ip string) bool {
	return region.CheckChinaIP(ip)
}

func (region *MyRegion) CloseIPdb(ip string) {
	region.RegionDB.Close()
}

func (region *MyRegion) CheckChinaIP(ip string) bool {
	ipInfo, err := region.RegionDB.MemorySearch(ip)
	if err != nil {
		return false
	}
	if ipInfo.Country != "中国" {
		return false
	}
	return true
}

func (region *MyRegion) Find(query string, params ...string) (interface{}, error) {
	var result ip_base.Result
	ipInfo, err := region.RegionDB.MemorySearch(query)
	if err != nil {
		return result, err
	}
	result.Country = ipInfo.Country
	result.State = ipInfo.Region
	result.Region = ipInfo.Province
	result.City = ipInfo.City
	result.ISP = ipInfo.ISP

	return result, nil
}

func (region *MyRegion) GetCityInfo(ip string) (ip_base.Result, error) {
	ipRes, err := region.Find(ip)
	if res, ok := ipRes.(ip_base.Result); ok && err == nil {
		return res, err
	}

	return ip_base.Result{}, err
}
