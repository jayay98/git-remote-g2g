package main

import (
	"os/exec"

	"github.com/libp2p/go-libp2p/core/network"
)

func ReceivePack(s network.Stream, repoName string) {
	cmd := exec.Command("git", "receive-pack", repoName)
	stdin, _ := cmd.StdinPipe() // read fetch-pack, not used
	stdout, _ := cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		s.Reset()
	}

	gitChan := make(chan []byte)
	go ScannerChannel("stdout", stdout, gitChan)
	streamChan := make(chan []byte)
	go ScannerChannel("stream", s, streamChan)
	quit := make(chan int)
	go func() {
		err = cmd.Wait()
		if err != nil {
			s.Reset()
		} else {
			logger.Infoln("Command exited 0")
		}
		quit <- 0
		close(quit)
	}()

	var gitOut, streamOut []byte
comms:
	for {
		select {
		case gitOut = <-gitChan:
			s.Write(gitOut)
		case streamOut = <-streamChan:
			stdin.Write(streamOut)
		case <-quit:
			break comms
		}
	}

	s.Close()
}
