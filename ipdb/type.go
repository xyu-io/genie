package ipdb

import (
	"github.com/xyu-io/genie/ipdb/iptool/geoip"
	"github.com/xyu-io/genie/ipdb/iptool/ip2location"
	"github.com/xyu-io/genie/ipdb/iptool/ip2region"
	"github.com/xyu-io/genie/ipdb/iptool/qqwry"
	"log"
	"strings"
)

var (
	defaultV4DB = "ip2location" //"ip2location" //  "geoip" //"ip2region_v1"
	defaultV6DB = "ip2location"
)

type DB struct {
	Name      string
	NameAlias []string `yaml:"name-alias,omitempty" mapstructure:"name-alias"`
	Format    Format
	FilePath  string
	File      string

	Languages []string
	Types     []Type
}

func (d *DB) get() (db DBImp) {
	if db, found := dbNameCache[d.Name]; found {
		return db
	}
	filePath := d.File
	if d.FilePath != "" {
		filePath = strings.ReplaceAll(d.FilePath+"/"+d.File, "//", "/")
	}

	var err error
	switch d.Format {
	case FormatQQWry:
		db, err = qqwry.NewQQwry(filePath)
	case FormatMMDB:
		db, err = geoip.NewGeoIP(filePath)
	case FormatIP2RegionV1:
		db, err = ip2region.NewIp2RegionV1(filePath)
	case FormatIP2RegionV2:
		db, err = ip2region.NewIp2RegionV2(filePath)
	case FormatIP2Location:
		db, err = ip2location.NewIP2Location(filePath)
	default:
		panic("DB format not supported!")
	}

	if err != nil {
		log.Fatalln("Database init failed:", err)
	}

	dbNameCache[d.Name] = db
	return
}

type Format string

const (
	FormatMMDB        Format = "mmdb"
	FormatQQWry              = "qqwry"
	FormatIP2RegionV1        = "ip2region_v1"
	FormatIP2RegionV2        = "ip2region_v2"
	FormatIP2Location        = "ip2location"

	FormatZXIPv6Wry = "zxipv6wry"
	FormatIPIP      = "ipip"
)

var (
	LanguagesAll = []string{"ALL"}
	LanguagesZH  = []string{"zh-CN"}
	LanguagesEN  = []string{"en"}
)

type Type string

const (
	TypeIPv4 Type = "IPv4"
	TypeIPv6      = "IPv6"
	// TypeCDN       = "CDN"
)

var (
	TypesAll  = []Type{TypeIPv4, TypeIPv6}
	TypesIP   = []Type{TypeIPv4, TypeIPv6}
	TypesIPv4 = []Type{TypeIPv4}
	TypesIPv6 = []Type{TypeIPv6}
	//	TypesCDN  = []Type{TypeCDN}
)

type List []*DB
type NameMap map[string]*DB
type TypeMap map[Type][]*DB

func (m *NameMap) From(dbs List) {
	for _, db := range dbs {
		(*m)[db.Name] = db

		if alias := db.NameAlias; alias != nil {
			for _, aName := range alias {
				(*m)[aName] = db
			}
		}
	}
}

func (m *TypeMap) From(dbs List) {
	for _, db := range dbs {
		for _, typ := range db.Types {
			dbsInType := (*m)[typ]
			if dbsInType == nil {
				dbsInType = []*DB{db}
			} else {
				dbsInType = append(dbsInType, db)
			}
			(*m)[typ] = dbsInType
		}
	}
}

func getDbByName(name string) (db *DB) {
	if dbInfo, found := NameDBMap[name]; found {
		return dbInfo
	}

	defaultNameDBMap := NameMap{}
	defaultNameDBMap.From(DefaultDBList(""))
	if dbInfo, found := defaultNameDBMap[name]; found {
		return dbInfo
	}

	log.Fatalf("DB with name %s not found!\n", name)
	return
}
