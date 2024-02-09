package p2p

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	ma "github.com/multiformats/go-multiaddr"
)

const defaultBufSize = 4096

type PeerStream struct {
	Host        host.Host
	connections map[string]*bufio.ReadWriter
}

func NewPeerStream(listenPort int) (*PeerStream, error) {
	h, err := createHost(listenPort)
	if err != nil {
		return nil, err
	}

	ps := &PeerStream{
		Host:        h,
		connections: make(map[string]*bufio.ReadWriter),
	}
	return ps, nil
}

func createHost(listenPort int) (host.Host, error) {

	r := rand.Reader
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	sourceMultiAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort))
	opts := []libp2p.Option{
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(priv),
	}

	host, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}
	return host, nil
}
