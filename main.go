package main

import (
	"flag"
	"log"
)

func main() {
  listenAddress := flag.String("listenAddr", defaultListenAddress, "server address")
  flag.Parse()
	server := NewServer(Config{ListenAddress: *listenAddress})
	log.Fatal(server.Start())
}
