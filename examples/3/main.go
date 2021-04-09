package main

import (
	"fmt"
	"net/http"
)

const (
	listen = ":9090"
)

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Hello World")
	})

	fmt.Printf("Server listening - %s\n", listen)
	err := http.ListenAndServe(listen, m)
	if err != nil {
		panic(err)
	}
}
