package qqwry

import (
	"errors"
	"fmt"
	"genie/ipdb/iptool/ip_base"
	"log"
	"os"

	"github.com/kayon/iploc"
)

type QQwry struct {
	loc *iploc.Locator
}

func NewQQwry(filePath string) (*QQwry, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新 ip2region 库")
		return nil, err
	}

	loc, err := iploc.Open(filePath)
	if err != nil {
		fmt.Println("open qqwry.dat error", err.Error())
		return nil, err
	}

	return &QQwry{
		loc: loc,
	}, nil
}

func (q *QQwry) Find(query string, params ...string) (interface{}, error) {
	var result ip_base.Result
	detail := q.loc.Find(query)
	if detail == nil {
		return result, errors.New("ip not found")
	}

	result.Country = detail.Country
	result.City = detail.City
	result.ISP = detail.Region

	return result, nil
}

func (q *QQwry) Name() string {
	return "qqwry"
}
