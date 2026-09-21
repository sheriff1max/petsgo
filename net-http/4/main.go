package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
)


func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello")) // v1
	// fmt.Fprint(w, "hello") // v2
}

func my_sum(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	a, err1 := strconv.Atoi(query.Get("a"))
	b, err2 := strconv.Atoi(query.Get("b"))
	if err1 != nil || err2 != nil {
		http.Error(w, "invalid parameters", 400)
		return
	}

	fmt.Fprint(w, a+b)
}

func getBody(addr string) string {
	resp, err := http.Get(addr)
	if err != nil {
		fmt.Println("error getBody 1: ", err)
		return ""
	}
	defer resp.Body.Close()

	// v1
	// body, _ := io.ReadAll(resp.Body)
	// return string(body)

	// v2
	bytes := make([]byte, 128)
	n, err := resp.Body.Read(bytes)
	if err != nil  && err != io.EOF {
		fmt.Println("error getBody 1: ", err)
		return ""
	}
	s := string(bytes[:n])
	return s
}

func startHTTPServer() (addr string, stop func()) {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/sum", my_sum)

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
	addr, stop := startHTTPServer()
	defer stop()
	resp3, _ := http.Get("http://" + addr + "/nope")
	fmt.Println(getBody("http://"+addr+"/hello") == "hello" &&
		getBody("http://"+addr+"/sum?a=2&b=3") == "5" && resp3.StatusCode == 404)
}
