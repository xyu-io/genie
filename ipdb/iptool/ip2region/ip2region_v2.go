package ip2region

import (
	"errors"
	"fmt"
	"genie/ipdb/iptool/ip_base"
	"io"
	"log"
	"os"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

type Ip2Region struct {
	seacher *xdb.Searcher
}

func NewIp2RegionV2(filePath string) (*Ip2Region, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新 ip2region 库")
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	searcher, err := xdb.NewWithBuffer(data)
	if err != nil {
		fmt.Printf("无法解析 ip2region xdb 数据库: %s\n", err)
		return nil, err
	}
	return &Ip2Region{
		seacher: searcher,
	}, nil
}

func (db Ip2Region) Find(query string, params ...string) (interface{}, error) {
	var result ip_base.Result
	if db.seacher != nil {
		res, err := db.seacher.SearchByStr(query)
		if err != nil {
			return "", err
		} else {
			ipArr := strings.Split(strings.ReplaceAll(res, "|0", "|"), "|")
			result.Country = ipArr[0]
			result.State = ipArr[1]
			result.Region = ipArr[2]
			result.City = ipArr[3]
			result.ISP = ipArr[4]

			return result, nil
		}
	}

	return "", errors.New("ip2region 未初始化")
}

func (db Ip2Region) Name() string {
	return "ip2region"
}
