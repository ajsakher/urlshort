package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val, ok := pathsToUrls[r.URL.Path]; ok {
			fmt.Printf("redirect '%s' to '%s'\n", r.URL.Path, val)
			http.Redirect(w, r, val, 302)
			return
		}

		fallback.ServeHTTP(w, r)
	})
}

type pathYAML struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func parseYAML(yml []byte) ([]pathYAML, error) {
	var paths []pathYAML
	err := yaml.Unmarshal(yml, &paths)
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func buildMap(yml []pathYAML) map[string]string {
	pathMap := map[string]string{}
	for _, val := range yml {
		pathMap[val.Path] = val.URL
	}
	return pathMap
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	paths, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(paths)

	return MapHandler(pathMap, fallback), nil
}
