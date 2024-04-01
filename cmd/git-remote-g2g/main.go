package main

import (
	"bufio"
	"context"
	"os"

	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

var logger = golog.Logger("remote-helper")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	args := os.Args
	if len(args) < 3 {
		logger.Fatalln("Usage: git-remote-g2g <remoteName> <multiAddr>")
	}

	node, err := libp2p.New()
	if err != nil {
		logger.Fatalln(err)
	}
	defer node.Close()

	kdht, err := NewDHT(ctx, node)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infoln("Initialized DHT")

	repo, err := NewRepository(kdht, args[2])
	if err != nil {
		logger.Fatalln(err)
	}
	repo.AddAddressInto(node)

	stdinReader := bufio.NewReader(os.Stdin)
	for {
		command, err := stdinReader.ReadString('\n')
		if err != nil {
			logger.Fatalln(err)
		}

		switch {
		case command == "capabilities\n":
			PrintCapabilities(os.Stdout)
		case command == "connect git-upload-pack\n":
			if err = ConnectUploadPack(node, ctx, repo.addrs[0].ID, repo.id); err != nil {
				logger.Fatalln(err)
			}
		case command == "connect git-receive-pack\n":
			if err = ConnectReceivePack(node, ctx, repo.addrs[0].ID, repo.id); err != nil {
				logger.Fatalln(err)
			}
		default:
			logger.Fatalf("Unknown command: %q", command)
		}
	}
}

func NewDHT(ctx context.Context, host host.Host) (kdht *dht.IpfsDHT, err error) {
	dhtopts := []dht.Option{dht.BootstrapPeersFunc(dht.GetDefaultBootstrapPeerAddrInfos), dht.Mode(dht.ModeClient)}
	kdht, err = dht.New(ctx, host, dhtopts...)
	if err != nil {
		return
	}

	if err = kdht.Bootstrap(ctx); err != nil {
		return
	}

	for {
		if kdht.RoutingTable().Size() > 0 {
			break
		}
	}

	return
}
