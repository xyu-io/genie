package ip_base

import (
	"regexp"
	"strings"
)

type Result struct {
	Country string `json:"country"` // 国家
	State   string `json:"state"`   // 四大洲
	Region  string `json:"region"`  // 美国-州； 中国-省；其他-行政区域
	City    string `json:"city"`    // 城市
	ISP     string `json:"isp"`     // 运营商
}

func (r Result) country() string {
	if r.Country == "" {
		return "unknown"
	}
	var regx = "^[a-zA-Z]+$"
	if match, _ := regexp.MatchString(regx, r.Country); match {
		// 返回对应国家的中文名称

	}

	return r.Country
}

func (r Result) region() string {
	if r.Region == "" {
		return "unknown"
	}
	return r.Region
}

func (r Result) state() string {
	if r.State == "" {
		return "unknown"
	}
	return r.State
}

func (r Result) city() string {
	if r.City == "" {
		return "unknown"
	}
	return r.City
}

func (r Result) isp() string {
	if r.ISP == "" {
		return "unknown"
	}
	return r.ISP
}

func (r Result) String() string {
	return strings.Join([]string{r.country(), r.state(), r.region(), r.city(), r.isp()}, "|")
}
