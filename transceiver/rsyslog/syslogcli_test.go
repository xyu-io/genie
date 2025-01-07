package rsyslog

import (
	"github.com/sirupsen/logrus"
	"testing"
)

type FlowData struct {
	Type                string        `db:"type" json:"type"`                                   //流量类型,如NETFLOW_V9
	TimeReceivedNs      string        `db:"time_received_ns" json:"time_received_ns"`           //接收流量的时间戳,如2024-08-09T08:13:55.800190194Z
	SequenceNum         int           `db:"sequence_num" json:"sequence_num"`                   //流量的序列号,如184881
	SamplingRate        int           `db:"sampling_rate" json:"sampling_rate"`                 //采样率,如1000（表示每 1000 个数据包采样一个）
	SamplerAddress      string        `db:"sampler_address" json:"sampler_address"`             //采集流量的设备地址,如"192.168.1.100"
	TimeFlowStartNs     uint64        `db:"time_flow_start_ns" json:"time_flow_start_ns"`       //流量开始的时间戳,如1691607000100000000
	TimeFlowEndNs       uint64        `db:"time_flow_end_ns" json:"time_flow_end_ns"`           //流量结束的时间戳,如1691607000500000000
	Bytes               int           `db:"bytes" json:"bytes"`                                 //流量包含的字节总数,如5000
	Packets             int           `db:"packets" json:"packets"`                             //流量包含的数据包总数,如10
	SrcAddr             string        `db:"src_addr" json:"src_addr"`                           //源 IP 地址,如"10.0.0.2"
	DstAddr             string        `db:"dst_addr" json:"dst_addr"`                           //目标 IP 地址,如"8.8.8.8"
	Etype               string        `db:"etype" json:"etype"`                                 //以太网类型,如"IPv4"
	Proto               string        `db:"proto" json:"proto"`                                 //传输层协议类型,如"TCP"
	SourcePort          uint16        `db:"src_port" json:"src_port"`                           //源端口号,如8080
	DestinationPort     uint16        `db:"dst_port" json:"dst_port"`                           //目标端口号,如443
	InIf                int           `db:"in_if" json:"in_if"`                                 //流入接口的编号,如1
	OutIf               int           `db:"out_if" json:"out_if"`                               //流出接口的编号,如2
	SrcMac              string        `db:"src_mac" json:"src_mac"`                             //源 MAC 地址,如"00:11:22:33:44:55"
	DstMac              string        `db:"dst_mac" json:"dst_mac"`                             //目标 MAC 地址,如"55:44:33:22:11:00"
	SrcVlan             int           `db:"src_vlan" json:"src_vlan"`                           //源 VLAN ID（虚拟局域网标识）,如10
	DstVlan             int           `db:"dst_vlan" json:"dst_vlan"`                           //目标 VLAN ID,如20
	VlanID              int           `db:"vlan_id" json:"vlan_id"`                             //VLAN ID,如15
	IPTos               int           `db:"ip_tos" json:"ip_tos"`                               //IP 服务类型字段,如64
	ForwardingStatus    int           `db:"forwarding_status" json:"forwarding_status"`         //转发状态,如1
	IPTTL               int           `db:"ip_ttl" json:"ip_ttl"`                               //IP 生存时间,如64
	IPFlags             int           `db:"ip_flags" json:"ip_flags"`                           //IP 标志位,如2
	TCPFlags            int           `db:"tcp_flags" json:"tcp_flags"`                         //TCP 标志位,如2
	IcmpType            int           `db:"icmp_type" json:"icmp_type"`                         //ICMP 消息类型,如8
	IcmpCode            int           `db:"icmp_code" json:"icmp_code"`                         //ICMP 消息代码,如0
	Ipv6FlowLabel       int           `db:"ipv6_flow_label" json:"ipv6_flow_label"`             //IPv6 流标签,如12345
	FragmentID          int           `db:"fragment_id" json:"fragment_id"`                     //分片 ID，用于标识数据包的分片,如5678
	FragmentOffset      int           `db:"fragment_offset" json:"fragment_offset"`             //分片偏移量，用于确定分片在原始数据包中的位置,如1000
	SrcAs               int           `db:"src_as" json:"src_as"`                               //源自治系统编号,如65000
	DstAs               int           `db:"dst_as" json:"dst_as"`                               //目标自治系统编号,如65001
	NextHop             string        `db:"next_hop" json:"next_hop"`                           //下一跳 IP 地址,如"192.168.2.1"
	NextHopAs           int           `db:"next_hop_as" json:"next_hop_as"`                     //下一跳的自治系统编号,如65002
	SrcNet              string        `db:"src_net" json:"src_net"`                             //源网络地址,如"10.0.0.0/24"
	DstNet              string        `db:"dst_net" json:"dst_net"`                             //目标网络地址,如"8.8.8.8/32"
	BgpNextHop          string        `db:"bgp_next_hop" json:"bgp_next_hop"`                   //BGP（边界网关协议）下一跳地址,如"192.168.3.1"
	BgpCommunities      []interface{} `db:"bgp_communities" json:"bgp_communities"`             //BGP 团体属性列表,如{65000:100, 65001:200}
	AsPath              []interface{} `db:"as_path" json:"as_path"`                             //自治系统路径列表,如{65000, 65001}
	MplsTTL             []interface{} `db:"mpls_ttl" json:"mpls_ttl"`                           //MPLS（多协议标签交换）标签的生存时间,如{255, 254}
	MplsLabel           []interface{} `db:"mpls_label" json:"mpls_label"`                       //MPLS 标签列表,如{1000, 2000}
	MplsIP              []interface{} `db:"mpls_ip" json:"mpls_ip"`                             //与 MPLS 相关的 IP 地址列表,如{"192.168.4.1", "192.168.5.1"}
	ObservationDomainID int           `db:"observation_domain_id" json:"observation_domain_id"` //观测域的编号,如123
	ObservationPointID  int           `db:"observation_point_id" json:"observation_point_id"`   //观测点的编号,如456
	HttpUrl             string        `db:"http_url" json:"http_url"`                           //HTTP 请求的 URL,如"https://www.example.com/page"
	HttpRetCode         int           `db:"http_ret_code" json:"http_ret_code"`                 //HTTP 返回码,如200
	HttpHost            string        `db:"http_host" json:"http_host"`                         //HTTP 请求的主机名,如"www.example.com"
	HttpMime            string        `db:"http_mime" json:"http_mime"`                         //HTTP 的 MIME 类型,如"text/html"
	HttpUserAgent       string        `db:"http_user_agent" json:"http_user_agent"`             //HTTP 的用户代理字符串,如"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
	DnsReqDomain        string        `db:"dns_req_domain" json:"dns_req_domain"`               //DNS 请求的域名,如"example.com"
	DnsReqType          int           `db:"dns_req_type" json:"dns_req_type"`                   //DNS 请求的类型,如1（A 记录查询）
	DnsReqIp            string        `db:"dns_req_ip" json:"dns_req_ip"`                       //DNS 请求的 IP 地址,如"192.168.6.1"
	ICMPData            string        `db:"icmp_data" json:"icmp_data"`                         //ICMP 数据内容,如"Echo data"
	ICMPSeqLen          int           `db:"icmp_seq_num" json:"icmp_seq_num"`                   //ICMP 序列号长度,如128
	ICMPPayloadLen      int           `db:"icmp_payload_len" json:"icmp_payload_len"`           //ICMP 有效载荷长度,如256
	FlowUid             string        `db:"flow_uid" json:"flow_uid"`                           //流量的唯一标识符,如5b069fe19fb4ce8ad9fa3a2ec5fff70b
	FlowCount           int64         `db:"flow_count" json:"flow_count"`                       //流量计数,如100
	Collector           string        `db:"collector" json:"collector"`                         //采集流量的工具或设备名称,如"MyCollectorTool"
}

func TestSysLogCli(t *testing.T) {
	// 创建发送者实例
	sysCli := NewSyslogClient(Config{
		Network:       "udp",
		Address:       "127.0.0.1:514",
		AppName:       "syslog",
		Formatter:     RFC5424Formatter,
		Level:         logrus.InfoLevel,
		DisableOutput: false,
	})

	msg := FlowData{
		Type:                "NetFlow V9",
		TimeReceivedNs:      "2024-08-23T09:36:39+08:00",
		SequenceNum:         231,
		SamplingRate:        12,
		SamplerAddress:      "10.52.2.222",
		TimeFlowStartNs:     1713319377000000000,
		TimeFlowEndNs:       1713319377000000000,
		Bytes:               12,
		Packets:             64,
		SrcAddr:             "10.52.2.66",
		DstAddr:             "10.52.2.229",
		Etype:               "IPv4",
		Proto:               "TCP",
		SourcePort:          1032,
		DestinationPort:     443,
		InIf:                0,
		OutIf:               0,
		SrcMac:              "00:e2:69:5f:5e:00",
		DstMac:              "33:33:ff:8e:c4:64",
		SrcVlan:             0,
		DstVlan:             0,
		VlanID:              0,
		IPTos:               1,
		ForwardingStatus:    0,
		IPTTL:               0,
		IPFlags:             0,
		TCPFlags:            6,
		IcmpType:            0,
		IcmpCode:            0,
		Ipv6FlowLabel:       0,
		FragmentID:          0,
		FragmentOffset:      0,
		SrcAs:               0,
		DstAs:               0,
		NextHop:             "",
		NextHopAs:           0,
		SrcNet:              "",
		DstNet:              "",
		BgpNextHop:          "",
		BgpCommunities:      nil,
		AsPath:              nil,
		MplsTTL:             nil,
		MplsLabel:           nil,
		MplsIP:              nil,
		ObservationDomainID: 0,
		ObservationPointID:  0,
		HttpUrl:             "",
		HttpRetCode:         0,
		HttpHost:            "",
		HttpMime:            "",
		HttpUserAgent:       "",
		DnsReqDomain:        "",
		DnsReqType:          0,
		DnsReqIp:            "",
		ICMPData:            "",
		ICMPSeqLen:          0,
		ICMPPayloadLen:      0,
		FlowUid:             "526c4b40b4c019069ed5b0f44ee79f94",
		FlowCount:           1,
		Collector:           "10.52.2.222",
	}
	// 发送多条日志消息
	for i := 0; i < 2; i++ {
		sysCli.SendLog(msg)
	}
	sysCli.SendWarn("waring! waring! waring!")
}
