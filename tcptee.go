// Copyright 2016 Chao Wang <hit9@icloud.com>

// Package main implements a tcp tee.
// Usage: ./tcptee -bind :8000 -backends :2015,:2016,:2017
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"strings"
)

// Tee is the tee handle.
type Tee struct {
	laddr    string
	backends []net.Conn
}

// New creates a new Tee.
func New(laddr string, addrs []string) (t *Tee, err error) {
	t = &Tee{laddr: laddr}
	if err = t.Connect(addrs); err != nil {
		return
	}
	return t, nil
}

// Read implements the io.Reader.
func (t *Tee) Read(p []byte) (n int, err error) {
	for _, b := range t.backends {
		n, err = b.Read(p)
		if err != nil {
			return n, err
		}
	}
	return
}

// Write implements the io.Writer.
func (t *Tee) Write(p []byte) (n int, err error) {
	for _, b := range t.backends {
		n, err = b.Write(p)
		if err != nil {
			return n, err
		}
	}
	return
}

// Connect to backends.
func (t *Tee) Connect(addrs []string) error {
	for _, addr := range addrs {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		t.backends = append(t.backends, conn)
		log.Printf("backend %s connected\n", addr)
	}
	return nil
}

// Serve the tee.
func (t *Tee) Serve() error {
	ln, err := net.Listen("tcp", t.laddr)
	if err != nil {
		return err
	}
	log.Printf("serving on %s\n", t.laddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go t.handle(conn)
	}
}

// Handle the connection.
func (t *Tee) handle(conn net.Conn) {
	var err error
	_, err = io.Copy(t, conn)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
}

func main() {
	bind := flag.String("bind", ":8000", "address to bind")
	backends := flag.String("backends", "", "backends split by comma")
	flag.Parse()
	if *backends == "" {
		log.Fatal("No backends")
	}
	tr, err := New(*bind, strings.Split(*backends, ","))
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(tr.Serve())
}
