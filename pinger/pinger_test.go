package pinger

import (
	"fmt"
	"testing"
)

func TestPinger(t *testing.T) {
	pingRes, err := Ping("10.52.2.66")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pingRes)
}
