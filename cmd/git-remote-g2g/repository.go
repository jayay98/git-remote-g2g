package main

import (
	"strings"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

type Repository struct {
	address *peer.AddrInfo
}

func NewRepo(url string) (repo Repository, err error) {
	address := strings.Replace(url, "g2g://", "", 1)
	multiaddr, err := ma.NewMultiaddr(address)
	if err != nil {
		return
	}

	info, err := peer.AddrInfoFromP2pAddr(multiaddr)
	if err != nil {
		return
	}

	repo = Repository{address: info}
	return
}

func (r Repository) AddAddressInto(node host.Host) {
	node.Peerstore().AddAddrs(r.address.ID, r.address.Addrs, peerstore.PermanentAddrTTL)
}
