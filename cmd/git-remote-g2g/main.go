package main

import (
	"context"
	"os"

	golog "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
)

var logger = golog.Logger("git-remote-helper")

func main() {
	golog.SetAllLoggers(golog.LevelInfo)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	args := os.Args
	if len(args) < 3 {
		println("Usage: git-remote-g2g <remoteName> <multiAddr>")
		os.Exit(1)
	}

	repo, err := NewRepo(args[2])
	if err != nil {
		logger.Warnln(err)
		os.Exit(1)
	}

	node, err := libp2p.New()
	if err != nil {
		logger.Warnln(err)
		os.Exit(1)
	}
	repo.AddAddressInto(node)

	StartStreamIO(node, ctx, repo)
}
