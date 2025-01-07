package pinger

import (
	"fmt"
	"testing"
)

func TestPinger(t *testing.T) {
	pingRes, err := Ping("127.0.0.1")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pingRes)
}
