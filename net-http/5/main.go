package main

import (
	"fmt"
	"net/http"
	"net"
	"io"
	"strings"
	"time"
)


func echo(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("X-User")

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error echo: ", err)
		return
	}
	defer r.Body.Close()

	body := string(bytes)

	fmt.Fprintf(w, "user=%s;body=%s", user, body)
}

func slow(w http.ResponseWriter, r *http.Request) {
	time.Sleep(300 * time.Millisecond)
}

func startHTTPServer() (addr string, stop func()) {
	mux := http.NewServeMux()

	mux.HandleFunc("/echo", echo)
	mux.HandleFunc("/slow", slow)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go func() {
		http.Serve(ln, mux)
	}()

	stop = func() {
		ln.Close()
	}

	return ln.Addr().String(), stop
}

func main() {
	client := &http.Client{Timeout: 100 * time.Millisecond}

	addr, stop := startHTTPServer()
	defer stop()

	_, errSlow := client.Get("http://" + addr + "/slow")
	req, _ := http.NewRequest("POST", "http://"+addr+"/echo", strings.NewReader("hi"))
	req.Header.Set("X-User", "ann")
	resp, errEcho := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(errSlow != nil && errEcho == nil && string(body) == "user=ann;body=hi")
}