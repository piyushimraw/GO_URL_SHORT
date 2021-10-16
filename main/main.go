package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/piyushimraw/go_url_short"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := go_url_short.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /google
  url: https://google.com
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := go_url_short.YAMLHandler([]byte(yaml), mapHandler)
	content, _ := ioutil.ReadFile("./redirect.json")
	jsonHandler, _ := go_url_short.JSONHandler(content, yamlHandler)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
