package rsyslog

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/influxdata/go-syslog/v3"
	"github.com/influxdata/go-syslog/v3/rfc5424"
	log "github.com/sirupsen/logrus"
)

// SyslogServer implements a UDP-based Syslog server.
type SyslogServer struct {
	udpPort uint
	ctx     context.Context
	errC    chan error
	msgC    chan Event
	conn    net.PacketConn
}

func NewSyslogServer(port uint, ctx context.Context) *SyslogServer {
	return &SyslogServer{
		udpPort: port,
		ctx:     ctx,
		errC:    make(chan error, 1),
		msgC:    make(chan Event, 1000),
		conn:    nil,
	}
}

func (svr *SyslogServer) RunSyslogReceiver() (chan Event, error) {
	if svr.ctx == nil {
		return svr.msgC, fmt.Errorf("ctx of syslog server is nil")
	}
	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", svr.udpPort))
	if err != nil {
		return svr.msgC, fmt.Errorf("listening udp syslog server err: %w", err)
	}

	svr.conn = conn

	go func() {
		maxPacketSize := 8192 // RFC5425#section-4.3.1
		buffer := make([]byte, maxPacketSize)
		syslogMachine := rfc5424.NewMachine()
		for {
			packetSize, _, err := conn.ReadFrom(buffer)
			if err != nil {
				svr.errC <- err
				return
			}
			syslogMessage, err := syslogMachine.Parse(buffer[:packetSize])
			if err != nil {
				log.Errorf("parsing syslog message: %v", err)
				continue
			}
			event, err := syslogToEvent(syslogMessage)
			if err != nil {
				log.Errorf("interpreting syslog message: %v", err)
				continue
			}
			svr.msgC <- *event
		}
	}()

	go func() {
		svr.Stop()
	}()

	return svr.msgC, nil
}

func (svr *SyslogServer) Stop() {
	select {
	case <-svr.ctx.Done():
		log.Warn("syslog server stop with ctx done")
		//svr.msgC = nil // avoid panic of goroutine asleep
		svr.conn.Close()
	case err := <-svr.errC:
		log.Warnf("syslog server stop with error: %v", err)
		//svr.msgC = nil // avoid panic of goroutine asleep
		svr.conn.Close()
	}
}

// See supervise implementation for details on Syslog field usage.
func syslogToEvent(syslogMessage syslog.Message) (*Event, error) {
	rfc5425Message, ok := syslogMessage.(*rfc5424.SyslogMessage)
	if !ok {
		panic("unexpected syslog message type")
	}
	if rfc5425Message.Appname == nil {
		return nil, errors.New("expected APP-NAME")
	}
	if rfc5425Message.MsgID == nil {
		return nil, errors.New("expected MSGID")
	}
	if rfc5425Message.Timestamp == nil {
		return nil, errors.New("expected TIMESTAMP")
	}

	appName := *rfc5425Message.Appname
	msgID := *rfc5425Message.MsgID
	procID := *rfc5425Message.ProcID
	levelID := *rfc5425Message.Severity //对应cli定义的等级
	priority := *rfc5425Message.Priority

	message := ""
	if rfc5425Message.Message != nil {
		message = strings.TrimSuffix(*rfc5425Message.Message, "\n")
	}

	return &Event{
		ProcID:    procID,
		Priority:  priority,
		Tag:       msgID,
		AppName:   appName,
		Level:     getFacility(levelID),
		Message:   message,
		Timestamp: rfc5425Message.Timestamp.Format(RFC3339MicroUTC),
	}, nil
}

func getFacility(facility uint8) string {
	switch facility {
	case 0:
		return "LOG_EMERG Priority"
	case 1:
		return "LOG_ALERT"
	case 2:
		return "LOG_CRIT"
	case 3:
		return "LOG_ERR"
	case 4:
		return "LOG_WARNING"
	case 5:
		return "LOG_NOTICE"
	case 6:
		return "LOG_INFO"
	case 7:
		return "LOG_DEBUG"
	default:
		return "LOG_Unknown"
	}
}
