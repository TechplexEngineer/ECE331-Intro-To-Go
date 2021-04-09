package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	listen = ":9090"
)

func main() {
	tmpl, err := template.ParseFiles("./template.html")
	if err != nil {
		panic(err)
	}

	m := http.NewServeMux()
	m.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		err := tmpl.Execute(writer, nil)
		if err != nil {
			log.Printf("problem executing template: %s", err)
			writer.WriteHeader(500)
		}
	})
	fmt.Printf("Server listening - %s\n", listen)
	err = http.ListenAndServe(listen, m)
	if err != nil {
		panic(err)
	}
}
