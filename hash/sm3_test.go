package hash

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	var str = "cfx.f2pool.com"
	bytes, err := json.Marshal(str)
	if err != nil {
		return
	}
	fmt.Println("-- 基于json-byte/原文 --")
	fmt.Println("md51: ", Sm3Hash(bytes))
	fmt.Println("md52: ", Sm3Data(str))
	fmt.Println("-- 基于[]byte/string --")
	fmt.Println("md53: ", Sm3Hash([]byte(str)))
	fmt.Println("md54: ", Sm3String(string(str)))
}
