package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/asakher/gophercises/urlshort"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	yamlFilename := flag.String("filename", "paths_to_urls.yaml", "YAML filename to read from")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := ioutil.ReadFile(*yamlFilename)
	check(err)

	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	check(err)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
