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

// Backends is the backend connections.
type Backends []net.Conn

// Read implements the io.Reader.
func (b Backends) Read(p []byte) (n int, err error) {
	for _, conn := range b {
		n, err = conn.Read(p)
		if err != nil {
			return n, err
		}
	}
	return
}

// Write implements the io.Writer.
func (b Backends) Write(p []byte) (n int, err error) {
	for _, conn := range b {
		n, err = conn.Write(p)
		if err != nil {
			return n, err
		}
	}
	return
}

// Tee is the tee handle.
type Tee struct {
	ln    net.Listener
	laddr string   // server addr
	addrs []string // backend addrs
}

// New creates a new Tee.
func New(laddr string, addrs []string) *Tee {
	return &Tee{laddr: laddr, addrs: addrs}
}

// Listen the tee.
func (t *Tee) Listen() (err error) {
	t.ln, err = net.Listen("tcp", t.laddr)
	if err != nil {
		return err
	}
	log.Printf("tee is listening on %s\n", t.laddr)
	return nil
}

// Serve the tee.
func (t *Tee) Serve() error {
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			return err
		}
		go t.handle(conn)
	}
}

// ListenAndServe is the Listen followed by Serve.
func (t *Tee) ListenAndServe() (err error) {
	if err = t.Listen(); err != nil {
		return
	}
	return t.Serve()
}

// Handle the connection.
func (t *Tee) handle(conn net.Conn) {
	// Connect to backends
	var backends []net.Conn
	for _, addr := range t.addrs {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			log.Println(err)
		}
		backends = append(backends, c)
	}
	var err error
	_, err = io.Copy(Backends(backends), conn)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	bind := flag.String("bind", ":8000", "address to bind")
	backends := flag.String("backends", "", "backends split by comma")
	flag.Parse()
	if *backends == "" {
		log.Fatal("No backends")
	}
	tr := New(*bind, strings.Split(*backends, ","))
	log.Fatal(tr.ListenAndServe())
}
