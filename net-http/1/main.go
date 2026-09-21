package main

import (
	"fmt"
	"net"
	"bufio"
	"io"
)


func startEchoServer() (addr string, stop func()) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	stop = func() {
		listener.Close()
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			go func() {
				defer conn.Close()
				io.Copy(conn, conn)
			}()
		}
	}()

	return listener.Addr().String(), stop
}

func main() {
	addr, stop := startEchoServer()
	defer stop()

	conn, err := net.Dial("tcp", addr)
	fmt.Fprintln(conn, "ping")
	line, err2 := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(err == nil && err2 == nil && line == "ping\n")
}
