package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)


func startSlowServer(d time.Duration) (addr string, stop func()) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	stop = func() {
		ln.Close()
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("error server 1:", err)
				return
			}

			go func() {
				defer conn.Close()

				reader := bufio.NewReader(conn)
				_, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("error server 2:", err)
					return
				}

				time.Sleep(d)

				n, err := fmt.Fprintln(conn, "ok")
				if n == 0 || err != nil {
					fmt.Println("error server 3:", err)
					return
				}
			}()
		}
	}()
	return ln.Addr().String(), stop
}

func requestWithTimeout(addr string, t time.Duration) (string, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(t))
	if err != nil {
		return "", err
	}

	n, err := fmt.Fprintln(conn, "Give me answer!")
	if n == 0 || err != nil {
		return "", err
	}

	reader := bufio.NewReader(conn)
	ans, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return ans, nil
}

func main() {
	addrSlow, stop1 := startSlowServer(300 * time.Millisecond)
	addrFast, stop2 := startSlowServer(10 * time.Millisecond)
	defer stop1()
	defer stop2()

	_, errSlow := requestWithTimeout(addrSlow, 100*time.Millisecond)
	msg, errFast := requestWithTimeout(addrFast, 100*time.Millisecond)

	fmt.Println(errors.Is(errSlow, os.ErrDeadlineExceeded) && errFast == nil &&
		strings.TrimSpace(msg) == "ok")
}
