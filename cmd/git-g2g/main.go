package main

import (
	"bufio"
	"context"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"g2g/pkg/pack"
	"g2g/pkg/specs"

	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

var logger = golog.Logger("git-server")

func main() {
	golog.SetAllLoggers(golog.LevelInfo)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	keyPath := "/tmp/key.pem"
	blob, _ := os.ReadFile(keyPath)
	block, _ := pem.Decode(blob)
	if block == nil {
		log.Fatalf("No PEM blob found")
	}
	priv, err := crypto.UnmarshalECDSAPrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(specs.HostAddress),
		libp2p.Identity(priv),
	}
	node, err := libp2p.New(opts...)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	node.SetStreamHandler(specs.HandshakeProto, func(s network.Stream) {
		defer s.Reset()
		service, err := bufio.NewReader(s).ReadString('\n')
		if err != nil {
			logger.Warnln(err)
		}
		service = strings.TrimSpace(service)

		cmd := exec.Command("git", service, "/tmp/test_repo")
		stdin, _ := cmd.StdinPipe() // read fetch-pack, not used
		stdout, _ := cmd.StdoutPipe()

		go func() {
			scn := pack.NewScanner(stdout)
			for scn.Scan() {
				logger.Warnln(scn.Text())
				s.Write(scn.Bytes())
			}
		}()
		go func() {
			scn := pack.NewScanner(s)
			for scn.Scan() {
				logger.Warnln(scn.Text())
				stdin.Write(scn.Bytes())
			}
		}()

		if err = cmd.Start(); err != nil {
			logger.Warnln(err)
			return
		}

		if err := cmd.Wait(); err != nil {
			logger.Fatal(err)
		}
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
