package main

import (
	"bufio"
	"context"
	"fmt"
	"g2g/pkg/pack"
	"g2g/pkg/specs"
	"io"
	"os"
	"strings"

	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

var logger = golog.Logger("remote-helper")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	args := os.Args
	if len(args) < 3 {
		logger.Fatalln("Usage: git-remote-g2g <remoteName> <multiAddr>")
	}

	repo, err := NewRepo(args[2])
	if err != nil {
		logger.Fatalln(err)
	}

	node, err := libp2p.New()
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
		case strings.HasPrefix(command, "connect"):
			service := strings.TrimSpace(strings.TrimPrefix(command, "connect git-"))
			if err = ConnectService(node, service, ctx, repo.address.ID); err != nil {
				logger.Fatalln(err)
			}
		default:
			logger.Fatalf("Unknown command: %q", command)
		}
	}
}

func PrintCapabilities(w io.Writer) {
	fmt.Fprintln(w, "connect")
	fmt.Fprintln(w, "")
}

func ConnectService(node host.Host, service string, ctx context.Context, peerId peer.ID) (err error) {
	// Connects to given git service.
	s, err := node.NewStream(ctx, peerId, specs.HandshakeProto)
	if err != nil {
		return err
	}
	os.Stdout.WriteString("\n")

	// Inform remote the service used
	fmt.Fprintln(s, service)

	// After line feed terminating the positive (empty) response, the output of service starts.
	serviceScanner := pack.NewScanner(s)
	for serviceScanner.Scan() {
		logger.Debugf("Remote: %q\n", serviceScanner.Text())
		fmt.Fprint(os.Stdout, serviceScanner.Text())
		if serviceScanner.Text() == "0000" {
			break
		}
	}

	cmdScanner := pack.NewScanner(os.Stdin)
	for cmdScanner.Scan() {
		logger.Debugf("cmd: %q\n", cmdScanner.Text())
		s.Write(cmdScanner.Bytes())
		if cmdScanner.Text() == "0009done\n" {
			break
		}
	}

	for serviceScanner.Scan() {
		logger.Debugf("Remote: %q\n", serviceScanner.Text())
		fmt.Fprint(os.Stdout, serviceScanner.Text())
		if serviceScanner.Text() == "0000" {
			break
		}
	}

	// After the connection ends, the remote helper exits.
	s.Reset()
	os.Exit(0)
	return
}
