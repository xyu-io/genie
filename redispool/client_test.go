package redispool

import (
	"context"
	"testing"
)

func TestRClient(t *testing.T) {
	redisOption := Option{
		AppName:  "hermes",
		Network:  "192.168.200.16:6379",
		Password: "",
		Timeout:  20,
		DB:       0,
		PoolSize: 50,
	}
	ctx := context.Background()
	rClient, err := MustNewClient(ctx, redisOption)
	if err != nil {
		panic(err)
	}
	pub := rClient.Subscribe("net_alert")
	select {
	case msg := <-pub.Channel():
		t.Logf("%+v", msg)
	}
}
