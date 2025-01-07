package ipdb

import (
	"sync"
)

var (
	dbNameCache = make(map[string]DBImp)
	dbTypeCache = make(map[QueryType]DBImp)
	queryCache  = sync.Map{}
)

var (
	NameDBMap = make(NameMap)
	TypeDBMap = make(TypeMap)
)

type QueryType uint

type DBImp interface {
	Find(query string, params ...string) (result interface{}, err error)
	Name() string
}

const (
	TypeV4 = iota
	TypeV6
)

//
//var (
//	_ DB = &qqwry.QQwry{}
//	_ DB = &geoip.GeoIP{}
//	_ DB = &ip2region.Ip2Region{}
//	_ DB = &ip2location.IP2Location{}
//)
