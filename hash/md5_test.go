package hash

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	var str = "cfx.f2pool.com"
	bytes, err := json.Marshal(str)
	if err != nil {
		return
	}
	fmt.Println("-- 基于json-byte --")
	fmt.Println("md51: ", MD5(bytes))
	fmt.Println("md52: ", MD5Data(str))
	fmt.Println("-- 基于[]byte --")
	fmt.Println("md53: ", MD5([]byte(str)))
}
