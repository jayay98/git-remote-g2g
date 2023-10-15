package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"g2g/pkg/specs"

	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

var logger = golog.Logger("git-server")

func main() {
	golog.SetAllLoggers(golog.LevelInfo)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	node, err := libp2p.New(
		libp2p.ListenAddrStrings(specs.HostAddress),
	)
	if err != nil {
		panic(err)
	}

	node.SetStreamHandler(specs.UploadPackProto, func(s network.Stream) {
		logger.Infof("%s: %s\n", s.Conn().RemotePeer().String(), s.Protocol())
		UploadPack(s, "/tmp/test_repo")
	})

	node.SetStreamHandler(specs.ReceivePackProto, func(s network.Stream) {
		logger.Infof("%s: %s\n", s.Conn().RemotePeer().String(), s.Protocol())
		ReceivePack(s, "/tmp/test_repo")
	})

	for _, addr := range node.Addrs() {
		hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", node.ID().Pretty()))
		p2pAddr := addr.Encapsulate(hostAddr).String()
		fmt.Printf("Serving on g2g://%s\n", p2pAddr)
	}

	defer node.Close()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
}
