package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)

const (
	listen = ":9090"
)

// type DataPoint represents a single value to be plotted on an x-y coordinate axis
// Value will be on the y axis
// Time  will be on the x axis
type DataPoint struct {
	Value float64 `json:"y"`
	Time  int64   `json:"x"`
}

func main() {

	reload := flag.Bool("reload", false, "pass this flag to have the template reloaded on each request.")

	flag.Parse()

	tmpl, err := template.ParseFiles("./template.html")
	if err != nil {
		panic(err)
	}

	m := http.NewServeMux()
	m.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if *reload {
			tmpl, err = template.ParseFiles("./template.html")
			if err != nil {
				log.Printf("problem loading template: %s", err)
				writer.WriteHeader(500)
				return
			}
		}
		err := tmpl.Execute(writer, nil)
		if err != nil {
			log.Printf("problem executing template: %s", err)
			writer.WriteHeader(500)
			return
		}
	})

	m.HandleFunc("/data", func(writer http.ResponseWriter, request *http.Request) {
		var points []DataPoint
		for n := 0; n < 5; n++ {
			points = append(points, DataPoint{
				Value: getValue(32),
				Time:  int64(n),
			})
		}
		jsonBytes, err := json.MarshalIndent(points, "", "    ")
		if err != nil {
			writer.WriteHeader(500)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		_, err = fmt.Fprint(writer, string(jsonBytes))
		if err != nil {
			writer.WriteHeader(500)
			return
		}
	})

	fmt.Printf("Server listening - %s\n", listen)
	err = http.ListenAndServe(listen, m)
	if err != nil {
		panic(err)
	}
}

// getValue produces a random 64bit floating point number on the range [-plusMinus, plusMinus)
func getValue(plusMinus float64) float64 {
	return (rand.Float64() * plusMinus * 2) - plusMinus
}
