package ipdb

import (
	"errors"
	"fmt"
	"github.com/xyu-io/genie/ipdb/iptool/ip_base"
	"log"
	"net"
	"slices"
	"strings"

	log2 "github.com/sirupsen/logrus"
)

func Load(pathPrefix string) {
	NameDBMap.From(DefaultDBList(pathPrefix))
	TypeDBMap.From(DefaultDBList(pathPrefix))
	log2.Info("ip db resource init end")
}

func getDB(typ QueryType, dbName string) (db DBImp) {
	if db, found := dbTypeCache[typ]; found {
		return db
	}

	var err error
	switch typ {
	case TypeV4:
		selected := defaultV4DB
		if dbName != "" {
			selected = dbName
		}
		dbUtil := getDbByName(selected)
		if dbUtil != nil && slices.Contains(dbUtil.Types, TypeIPv4) {
			db = dbUtil.get()
		} else {
			err = errors.New(fmt.Sprintf("db %s not found", dbName))
		}
	case TypeV6:
		selected := defaultV6DB
		if dbName != "" {
			selected = dbName
		}

		dbUtil := getDbByName(selected)
		if dbUtil != nil && slices.Contains(dbUtil.Types, TypeIPv6) {
			db = dbUtil.get()
		} else {
			err = errors.New(fmt.Sprintf("db %s not found", dbName))
		}

	default:
		panic("Query type not supported!")
	}

	if err != nil || db == nil {
		log.Fatalln("Database init failed:", err)
	}

	dbTypeCache[typ] = db
	return
}

// FindCountry 查询IP国家归属信息
func FindCountry(dbName, query string) (string, error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return "unknown", errors.New("query should be valid IP")
	}

	var typ QueryType
	if strings.Contains(query, ".") && len(strings.Split(query, ".")) == 4 {
		typ = TypeV4
	} else {
		typ = TypeV6
	}

	if result, found := queryCache.Load(query); found {
		return result.(*ip_base.Result).Country, nil
	}

	db := getDB(typ, dbName)
	result, err := db.Find(query)
	if err != nil {
		fmt.Println("find error:", err)
		return "unknown", err
	}
	res := &ip_base.Result{}

	if resp, ok := result.(ip_base.Result); ok {
		res = &resp
	}

	queryCache.Store(query, res)

	return res.Country, nil
}

// FindAll 入口, 用于查询IP信息，自动识别IP类型，可指定数据库或置空系统自动选择
func FindAll(dbName, query string) (*ip_base.Result, error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("query should be valid IP")
	}

	var typ QueryType
	if strings.Contains(query, ".") && len(strings.Split(query, ".")) == 4 {
		typ = TypeV4
	} else {
		typ = TypeV6
	}

	if result, found := queryCache.Load(query); found {
		return result.(*ip_base.Result), nil
	}

	db := getDB(typ, dbName)
	result, err := db.Find(query)
	if err != nil {
		fmt.Println("find error:", err)
		return nil, err
	}
	res := &ip_base.Result{}

	if resp, ok := result.(ip_base.Result); ok {
		res = &resp
	}

	queryCache.Store(query, res)

	return res, nil
}
