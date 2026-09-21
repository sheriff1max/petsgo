package main

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"strconv"
)


func startCmdServer() (addr string, stop func()) {
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
				fmt.Println(err)
				return
			}

			go func() {
				defer conn.Close()

				buf := bufio.NewScanner(conn)
				for buf.Scan() {
					s := buf.Text()

					split := strings.Fields(s)
					switch split[0] {
					case "ADD":
						v1, err := strconv.Atoi(split[1])
						if err != nil {
							fmt.Println(err)
							return
						}
						v2, err := strconv.Atoi(split[2])
						if err != nil {
							fmt.Println(err)
							return
						}

						res := strconv.Itoa(v1 + v2)
						// conn.Write([]byte(res))
						fmt.Fprintln(conn, res)

					case "ECHO":
						// conn.Write([]byte(split[1]))
						fmt.Fprintln(conn, split[1])

					default:
						// conn.Write([]byte("ERR"))
						fmt.Fprintln(conn, "ERR")
					}
				}

				err_buf := buf.Err()
				if err_buf != nil {
					fmt.Println(err_buf)
					return
				}
			}()
		}
	}()
	return ln.Addr().String(), stop
}

func ask(addr, line string) string {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// n, err := conn.Write([]byte(line))
	n, err := fmt.Fprintln(conn, line)
	if n == 0 || err != nil {
		panic(err)
	}

	var ans string
	n, err = fmt.Fscanln(conn, &ans)
	if n == 0 || err != nil {
		panic(err)
	}

	return ans
}

func askTwo(addr, l1, l2 string) (string, string) {
	var (
		ans1 string
		ans2 string
		n int
		err error
	)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	n, err = fmt.Fprintln(conn, l1)
	if n == 0 || err != nil {
		panic(err)
	}
	n, err = fmt.Fscanln(conn, &ans1)
	if n == 0 || err != nil {
		panic(err)
	}

	n, err = fmt.Fprintln(conn, l2)
	if n == 0 || err != nil {
		panic(err)
	}
	n, err = fmt.Fscanln(conn, &ans2)
	if n == 0 || err != nil {
		panic(err)
	}

	return ans1, ans2
}

func main() {
	addr, stop := startCmdServer()
	defer stop()
	a, b := askTwo(addr, "ADD 1 2", "ADD 3 4")
	fmt.Println(ask(addr, "ADD 2 3") == "5" && ask(addr, "ECHO go") == "go" &&
		ask(addr, "???") == "ERR" && a == "3" && b == "7")
}
