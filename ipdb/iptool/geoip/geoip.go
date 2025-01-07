package geoip

import (
	"errors"
	"github.com/xyu-io/genie/ipdb/iptool/ip_base"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/viper"
)

// GeoIP2
type GeoIP struct {
	db *geoip2.Reader
}

// new geoip from database file
func NewGeoIP(filePath string) (*GeoIP, error) {
	// 判断文件是否存在
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，请自行下载 Geoip2 City库，并保存在", filePath)
		return nil, err
	} else {
		db, err := geoip2.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		return &GeoIP{db: db}, nil
	}
}

func (g GeoIP) Find(query string, params ...string) (result interface{}, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("Query should be valid IP")
	}
	record, err := g.db.City(ip)
	if err != nil {
		return
	}

	lang := viper.GetString("selected.lang")
	if lang == "" {
		lang = "zh-CN"
	}

	result = ip_base.Result{
		Country: getMapLang(record.Country.Names, lang),
		City:    getMapLang(record.City.Names, lang),
	}
	return
}

func (db GeoIP) Name() string {
	return "geoip"
}

const DefaultLang = "en"

func getMapLang(data map[string]string, lang string) string {
	res, found := data[lang]
	if found {
		return res
	}
	return data[DefaultLang]
}
