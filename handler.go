package go_url_short

import (
	"log"
	"net/http"

	je "encoding/json"

	yaml "gopkg.in/yaml.v2"
)

type UrlData struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	if pathsToUrls == nil {
		return fallback.ServeHTTP
	}
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusPermanentRedirect)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	if yml == nil {
		return fallback.ServeHTTP, nil
	}

	pathsToUrls := make(map[string]string)

	var urls []UrlData

	err := yaml.Unmarshal(yml, &urls)

	if err != nil {
		log.Fatalf("cannot unmarshal data:  %v", err)
		return fallback.ServeHTTP, err
	}

	for _, url := range urls {
		pathsToUrls[url.Path] = url.Url
	}
	if len(pathsToUrls) == 0 {
		return fallback.ServeHTTP, nil
	}
	return MapHandler(pathsToUrls, fallback), nil
}

// JSON is expected to be in the format:
// {"path": "/some-path", "url": "https://www.some-url.com/demo"}
// The only errors that can be returned all related to having
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {

	if json == nil {
		return fallback.ServeHTTP, nil
	}

	pathsToUrls := make(map[string]string)

	var urls []UrlData

	err := je.Unmarshal(json, &urls)

	if err != nil {
		log.Fatalf("cannot unmarshal data:  %v", err)
		return fallback.ServeHTTP, err
	}

	for _, url := range urls {
		pathsToUrls[url.Path] = url.Url
	}
	if len(pathsToUrls) == 0 {
		return fallback.ServeHTTP, nil
	}
	return MapHandler(pathsToUrls, fallback), nil
}
