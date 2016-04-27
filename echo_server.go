// Copyright 2016 Chao Wang <hit9@icloud.com>

// +build ignore

// Simple echo server.
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	bind := flag.String("bind", ":8001", "address to bind")
	flag.Parse()
	ln, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go io.Copy(os.Stderr, conn)
	}
}
