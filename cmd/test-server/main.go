package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
)

const TestProtocol = protocol.ID("/libp2p/test/data")

var testFilePath string

func main() {
	port := flag.Int("port", 4001, "server listen port")
	testFile := flag.String("file", "data", "data file to serve")

	flag.Parse()

	if _, err := os.Stat(*testFile); err != nil {
		log.Fatal(err)
	}
	testFilePath = *testFile

	host, err := libp2p.New(
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *port),
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", *port),
		),
		libp2p.DefaultTransports,
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range host.Addrs() {
		fmt.Printf("I am %s/p2p/%s\n", addr, host.ID())
	}

	host.SetStreamHandler(TestProtocol, handleStream)

	select {}
}

func handleStream(s network.Stream) {
	defer s.Close()

	log.Printf("Incoming connection from %s", s.Conn().RemoteMultiaddr())

	file, err := os.Open(testFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	start := time.Now().UnixNano()
	n, err := io.Copy(s, file)
	if err != nil {
		log.Printf("Error transmiting file: %s", err)
	}
	end := time.Now().UnixNano()
	log.Printf("Transmitted %d bytes in %s", n, time.Duration(end-start))
}
