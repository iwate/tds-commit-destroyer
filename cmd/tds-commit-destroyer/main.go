package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	// "regexp"
	// "strings"

	"github.com/iwate/tds-commit-destroyer"
)

var (
	logger  proxy.ColorLogger
	connid  uint64

	localAddr   = flag.String("l", ":9433", "local address")
	remoteAddr  = flag.String("r", "localhost:1433", "remote address")
)

func main() {
	flag.Parse()

	logger := proxy.ColorLogger{}

	logger.Info("Proxying from %v to %v", *localAddr, *remoteAddr)

	laddr, err := net.ResolveTCPAddr("tcp", *localAddr)
	if err != nil {
		logger.Warn("Failed to resolve local address: %s", err)
		os.Exit(1)
	}
	raddr, err := net.ResolveTCPAddr("tcp", *remoteAddr)
	if err != nil {
		logger.Warn("Failed to resolve remote address: %s", err)
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		logger.Warn("Failed to open local port to listen: %s", err)
		os.Exit(1)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.Warn("Failed to accept connection '%s'", err)
			continue
		}
    
		var p *proxy.Proxy
    p = proxy.New(conn, laddr, raddr)

		p.Log = proxy.ColorLogger {
			Prefix:      fmt.Sprintf("Connection #%03d ", connid),
		}

		go p.Start()
	}
}
