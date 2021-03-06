package types

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"strings"
	"testing"
)

func TestFlows(t *testing.T) {
	ipFlow, _ := gopacket.FlowFromEndpoints(layers.NewIPEndpoint(net.IPv4(1, 2, 3, 4)), layers.NewIPEndpoint(net.IPv4(2, 3, 4, 5)))
	tcpFlow, _ := gopacket.FlowFromEndpoints(layers.NewTCPPortEndpoint(layers.TCPPort(1)), layers.NewTCPPortEndpoint(layers.TCPPort(2)))
	tcpIpFlow := NewTcpIpFlowFromFlows(ipFlow, tcpFlow)
	ipFlow2, tcpFlow2 := tcpIpFlow.Flows()
	if ipFlow2 != ipFlow || tcpFlow2 != tcpFlow {
		t.Error("Flows method fail")
		t.Fail()
	}
}

func TestFlowString(t *testing.T) {
	ipFlow, _ := gopacket.FlowFromEndpoints(layers.NewIPEndpoint(net.IPv4(1, 2, 3, 4)), layers.NewIPEndpoint(net.IPv4(2, 3, 4, 5)))
	tcpFlow, _ := gopacket.FlowFromEndpoints(layers.NewTCPPortEndpoint(layers.TCPPort(1)), layers.NewTCPPortEndpoint(layers.TCPPort(2)))
	tcpIpFlow := NewTcpIpFlowFromFlows(ipFlow, tcpFlow)
	if !strings.EqualFold("1.2.3.4:1-2.3.4.5:2", tcpIpFlow.String()) {
		t.Error("TcpIpFlow.String() fail")
		t.Fail()
	}
}

func TestFlowEqual(t *testing.T) {
	ipFlow1, _ := gopacket.FlowFromEndpoints(layers.NewIPEndpoint(net.IPv4(1, 2, 3, 4)), layers.NewIPEndpoint(net.IPv4(2, 3, 4, 5)))
	tcpFlow1, _ := gopacket.FlowFromEndpoints(layers.NewTCPPortEndpoint(layers.TCPPort(1)), layers.NewTCPPortEndpoint(layers.TCPPort(2)))
	flow1 := NewTcpIpFlowFromFlows(ipFlow1, tcpFlow1)

	ipFlow2, _ := gopacket.FlowFromEndpoints(layers.NewIPEndpoint(net.IPv4(1, 2, 3, 4)), layers.NewIPEndpoint(net.IPv4(2, 3, 4, 5)))
	tcpFlow2, _ := gopacket.FlowFromEndpoints(layers.NewTCPPortEndpoint(layers.TCPPort(1)), layers.NewTCPPortEndpoint(layers.TCPPort(2)))
	flow2 := NewTcpIpFlowFromFlows(ipFlow2, tcpFlow2)

	if !flow1.Equal(&flow2) {
		t.Error("TcpIpFlow.Equal fail")
		t.Fail()
	}

	ipFlow3, _ := gopacket.FlowFromEndpoints(layers.NewIPEndpoint(net.IPv4(8, 2, 3, 4)), layers.NewIPEndpoint(net.IPv4(2, 3, 4, 5)))
	tcpFlow3, _ := gopacket.FlowFromEndpoints(layers.NewTCPPortEndpoint(layers.TCPPort(1)), layers.NewTCPPortEndpoint(layers.TCPPort(2)))
	flow3 := NewTcpIpFlowFromFlows(ipFlow3, tcpFlow3)

	if flow1.Equal(&flow3) {
		t.Error("TcpIpFlow.Equal fail")
		t.Fail()
	}
}

func TestNewTcpIpFlowFromPacket(t *testing.T) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	ip := layers.IPv4{
		SrcIP:    net.IP{1, 2, 3, 4},
		DstIP:    net.IP{2, 3, 4, 5},
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}
	tcp := layers.TCP{
		SYN:       true,
		SrcPort:   1,
		DstPort:   2,
		Seq:       123,
		BaseLayer: layers.BaseLayer{Payload: []byte{1, 2, 3}},
	}
	tcp.SetNetworkLayerForChecksum(&ip)
	gopacket.SerializeLayers(buf, opts, &ip, &tcp)
	packetData := buf.Bytes()
	flow1, err := NewTcpIpFlowFromPacket(packetData)

	ipFlow2, _ := gopacket.FlowFromEndpoints(layers.NewIPEndpoint(net.IPv4(1, 2, 3, 4)), layers.NewIPEndpoint(net.IPv4(2, 3, 4, 5)))
	tcpFlow2, _ := gopacket.FlowFromEndpoints(layers.NewTCPPortEndpoint(layers.TCPPort(1)), layers.NewTCPPortEndpoint(layers.TCPPort(2)))
	flow2 := NewTcpIpFlowFromFlows(ipFlow2, tcpFlow2)

	if err != nil && !flow2.Equal(flow1) {
		t.Error("NewTcpIpFlowFromPacket fail")
		t.Fail()
	}

	flow1, err = NewTcpIpFlowFromPacket([]byte{1, 2, 3, 4, 5, 6, 7})
	if err == nil || !flow1.Equal(&TcpIpFlow{}) {
		t.Error("NewTcpIpFlowFromPacket fail")
		t.Fail()
	}
}

func FlowFromPacket() *TcpIpFlow {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	ip := layers.IPv4{
		SrcIP:    net.IP{1, 2, 3, 4},
		DstIP:    net.IP{2, 3, 4, 5},
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}
	tcp := layers.TCP{
		SYN:       true,
		SrcPort:   1,
		DstPort:   2,
		Seq:       123,
		BaseLayer: layers.BaseLayer{Payload: []byte{1, 2, 3}},
	}
	tcp.SetNetworkLayerForChecksum(&ip)
	gopacket.SerializeLayers(buf, opts, &ip, &tcp)
	packetData := buf.Bytes()
	flow, _ := NewTcpIpFlowFromPacket(packetData)
	return flow
}

func TestHashedTcpIpv4Flow(t *testing.T) {
	tcpIpFlow := FlowFromPacket()
	hash1 := NewHashedTcpIpv4Flow(tcpIpFlow)
	f := tcpIpFlow.Reverse()
	hash2 := NewHashedTcpIpv4Flow(&f)
	if hash1 != hash2 {
		t.Error("hash values must be equal after sorting!")
		t.Fail()
	}
}

func TestHashedTcpIpv6Flow(t *testing.T) {
	tcpIpFlow := FlowFromPacket()
	hash1 := NewHashedTcpIpv6Flow(tcpIpFlow)
	f := tcpIpFlow.Reverse()
	hash2 := NewHashedTcpIpv6Flow(&f)
	if hash1 != hash2 {
		t.Error("hash values must be equal after sorting!")
		t.Fail()
	}
}
