package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"golang.org/x/net/netutil"
)

const (
	LISTENING_ADDRESS = ":8080"
)

func run() int {

	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "usage: %s CMD\n", os.Args[0])
		return 1
	}

	cmd := os.Args[1]

	if cmd == "serve" {

		http_config := &http.Server{
			Handler: NewWafDevServer(),
		}

		listener, err := net.Listen("tcp", LISTENING_ADDRESS)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}

		err = http_config.Serve(netutil.LimitListener(listener, 20))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}

	} else if cmd == "fetch" {

	} else {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
