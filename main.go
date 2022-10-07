package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello)
	r.HandleFunc("/headers", headers)

	http.Handle("/", r)
}
