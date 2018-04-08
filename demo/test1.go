package main

import (
	"net"
	"go-tcp-tlv"
	"fmt"
)

func main() {
	netListen, err := net.Listen("tcp", "localhost:10000")
	if err != nil {
		panic(err)
	}
	defer netListen.Close()

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		go NewConn(conn)
	}
}

func NewConn(conn net.Conn)  {
	tlv := go_tcp_tlv.TlvConn{ReadBufferSize:1024}
	tlv.NewConn(conn)
	redata := tlv.Reader()
	for v :=range redata{
		fmt.Println(string(v.Value))
		tlv.Write([]byte("wqqweasdsa"))
	}
}
