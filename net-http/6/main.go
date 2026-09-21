package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var quote string
var mut sync.Mutex

func setQuote(addr, text string) error {
	client := &http.Client{Timeout: 2 * time.Second}

	url := "http://" + addr + "/set"

	resp, err := client.Post(url, "text/plain", strings.NewReader(text))
	if err != nil {
		fmt.Println("error setQuote 1: ", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

func getQuote(addr string) string {
	client := &http.Client{Timeout: 2 * time.Second}

	url := "http://" + addr + "/get"
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("error getQuote 1: ", err)
		return ""
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		fmt.Println("error getQuote 2: ", err)
	}
	return string(bytes)
}

func setHandle(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		fmt.Println("error setHandle 1: ", err)
		return
	}
	defer r.Body.Close()

	s := string(bytes)

	mut.Lock()
	quote = s
	mut.Unlock()

	w.WriteHeader(http.StatusOK)
}

func getHandle(w http.ResponseWriter, r *http.Request) {
	mut.Lock()
	fmt.Fprint(w, quote)
	mut.Unlock()
}

func startQuoteServer() (addr string, stop func()) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	stop = func() {
		ln.Close()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/set", setHandle)
	mux.HandleFunc("/get", getHandle)

	go func() {
		err = http.Serve(ln, mux)
		if err != nil {
			fmt.Println("Server ends.")
		}
	}()

	return ln.Addr().String(), stop
}

func main() {
	addr, stop := startQuoteServer()
	defer stop()

	e1 := setQuote(addr, "Go is fun")
	g1 := getQuote(addr)
	e2 := setQuote(addr, "Channels rock")
	g2 := getQuote(addr)

	fmt.Println(e1 == nil && g1 == "Go is fun" && e2 == nil && g2 == "Channels rock")
}
