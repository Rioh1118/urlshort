package urlshort

import (
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
	result := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.String()
			if url, ok := pathsToUrls[path]; ok {
				w.Header().Set("Content-Type", "text/html")
				w.Header().Set("Location", url)
				w.WriteHeader(http.StatusMovedPermanently)
			} else {
				fallback.ServeHTTP(w, r)
			}
		})
	return result
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type YAMLPathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathURLs []YAMLPathURL

	err := yaml.Unmarshal(yml, &pathURLs)
	if err != nil {
		return nil, err
	}

	pathToURLs := make(map[string]string)
	for _, pu := range pathURLs {
		pathToURLs[pu.Path] = pu.URL
	}

	return MapHandler(pathToURLs, fallback), nil
}
