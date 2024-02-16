package p2p

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	net "github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"log"
	"strings"
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

func (ps *PeerStream) Run(ctx context.Context, streamGroup string) {
	peerAddr := ps.getPeerFullAddr()
	log.Printf("my address: %s\n", peerAddr)

	// connect to other peers
	ps.Host.SetStreamHandler("/p2p/1.0.0", ps.handleStream)
	peerChan := initMDNS(ps.Host, streamGroup)
	go func(ctx context.Context) {
		for {
			peer := <-peerChan
			if err := ps.Host.Connect(ctx, peer); err != nil {
				fmt.Println("connection failed:", err)
				continue
			}

			//fmt.Println("connected to: ", peer)
			s, err := ps.Host.NewStream(ctx, peer.ID, "/p2p/1.0.0")
			if err != nil {
				fmt.Println("stream open failed", err)
			} else {
				rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
				ps.connections[peer.ID.String()] = rw
				go ps.readStream(s)
			}
		}
	}(ctx)
}

func (ps *PeerStream) getPeerFullAddr() ma.Multiaddr {
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", ps.Host.ID()))

	addrs := ps.Host.Addrs()
	var addr ma.Multiaddr
	for _, i := range addrs {
		if strings.HasPrefix(i.String(), "/ip4") {
			addr = i
			break
		}
	}
	return addr.Encapsulate(hostAddr)
}

func (ps *PeerStream) handleStream(s net.Stream) {
	go ps.readStream(s)
}

func (ps *PeerStream) readStream(s net.Stream) {

	for {
		buffer := make([]byte, defaultBufSize)
		n, err := s.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from stream:", err)
			}
			break
		}
		b := buffer[:n]
		b = bytes.Trim(b, "\x00")
		m := make(map[string]interface{})
		err = json.Unmarshal(b, &m)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
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
