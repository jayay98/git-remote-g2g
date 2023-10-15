package main

import (
	"fmt"
	"io"
	"log"
)

func PrintCapabilities(w io.Writer) {
	_, err := fmt.Fprintln(w, "connect")
	if err != nil {
		log.Fatalf("failed to write to stream: %v", err)
	}

	_, err = fmt.Fprintln(w, "")
	if err != nil {
		log.Fatalf("failed to write to stream: %v", err)
	}
}
