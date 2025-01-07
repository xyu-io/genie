package ip2location

import (
	"errors"
	"genie/ipdb/iptool/ip_base"
	"log"
	"net"
	"os"

	"github.com/ip2location/ip2location-go/v9"
)

type IP2Location struct {
	db *ip2location.DB
}

// NewIP2Location from database file
func NewIP2Location(filePath string) (*IP2Location, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，请自行下载 IP2Location 库，并保存在", filePath)
		return nil, err
	} else {
		db, err := ip2location.OpenDB(filePath)

		if err != nil {
			log.Fatal(err)
		}
		return &IP2Location{db: db}, nil
	}
}

func (db IP2Location) Find(query string, params ...string) (result interface{}, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("Query should be valid IP")
	}
	record, err := db.db.Get_all(ip.String())

	if err != nil {
		return
	}

	result = ip_base.Result{
		Country: record.Country_long,
		Region:  record.Region,
		City:    record.City,
		//ISP:     record.Isp, isp等信息需要付费版本
	}
	return
}

func (db IP2Location) Name() string {
	return "ip2location"
}
