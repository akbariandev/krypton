package p2p

import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifier struct {
	PeerChan chan peer.AddrInfo
}

func (n *discoveryNotifier) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

func initMDNS(peerHost host.Host, serviceName string) chan peer.AddrInfo {
	n := &discoveryNotifier{}
	n.PeerChan = make(chan peer.AddrInfo)

	ser := mdns.NewMdnsService(peerHost, serviceName, n)
	if err := ser.Start(); err != nil {
		panic(err)
	}
	return n.PeerChan
}
