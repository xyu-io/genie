package rsyslog

import (
	"context"
	"fmt"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

// 服务端测试用例，可以配合客户端一起使用
func TestSysLogSrv(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	srv := NewSyslogServer(514, ctx)

	msg, err := srv.RunSyslogReceiver()
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case e := <-msg:
				fmt.Printf("syslog recieved >> %+v \n", e)
			default:
				// avoid panic of goroutine asleep
			}
		}
	}()
	time.Sleep(time.Second * 4)
	log.Warn("ctx will done")
	time.Sleep(time.Second * 3)
}
