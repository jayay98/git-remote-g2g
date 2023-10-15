package main

import (
	"bufio"
	"context"
	"os"

	"g2g/pkg/pack"
	"g2g/pkg/specs"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func StartStreamIO(node host.Host, ctx context.Context, repo Repository) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		switch scanner.Text() {
		case "capabilities":
			PrintCapabilities(os.Stdout)
		case "connect git-upload-pack":
			s, err := node.NewStream(ctx, repo.address.ID, protocol.ID(specs.UploadPackProto))
			if err != nil {
				logger.Warnln(err)
				os.Exit(1)
			}
			connectUploadPack(s)
			s.Close()
		case "connect git-receive-pack":
			s, err := node.NewStream(ctx, repo.address.ID, protocol.ID(specs.ReceivePackProto))
			if err != nil {
				logger.Warnln(err)
				os.Exit(1)
			}
			connectReceivePack(s)
			s.Close()
		default:
			println("Sth else")
			println(scanner.Text())
		}
	}
}

func connectUploadPack(s network.Stream) {
	os.Stdout.Write([]byte("\n"))

	gitChan := make(chan []byte)
	gitEOF := make(chan int)
	go func() {
		scn := pack.NewScanner(os.Stdin)
		for scn.Scan() {
			gitChan <- scn.Bytes()
		}
		gitEOF <- 0
		close(gitChan)
		close(gitEOF)
	}()

	streamChan := make(chan []byte)
	go func() {
		scn := pack.NewScanner(s)
		for scn.Scan() {
			streamChan <- scn.Bytes()
		}
		close(streamChan)
	}()

	var gitOut, streamOut []byte
comms:
	for {
		select {
		case gitOut = <-gitChan:
			s.Write(gitOut)
		case streamOut = <-streamChan:
			os.Stdout.Write(streamOut)
		case <-gitEOF:
			break comms
		}
	}
}

func connectReceivePack(s network.Stream) {
	os.Stdout.Write([]byte("\n"))

	gitChan := make(chan []byte)
	gitEOF := make(chan int)
	go func() {
		scn := pack.NewScanner(os.Stdin)
		for scn.Scan() {
			gitChan <- scn.Bytes()
		}
		gitEOF <- 0
		close(gitChan)
		close(gitEOF)
	}()

	streamChan := make(chan []byte)
	go func() {
		scn := pack.NewScanner(s)
		for scn.Scan() {
			streamChan <- scn.Bytes()
		}
		close(streamChan)
	}()

	var gitOut, streamOut []byte
comms:
	for {
		select {
		case gitOut = <-gitChan:
			s.Write(gitOut)
		case streamOut = <-streamChan:
			os.Stdout.Write(streamOut)
		case <-gitEOF:
			break comms
		}
	}
}
