package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const yml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

func TestBuildMap(t *testing.T) {
	yml := []pathYAML{
		{Path: "test1", URL: "test1_url"},
		{Path: "test2", URL: "test2_url"},
	}

	ymlMap := buildMap(yml)
	assert.Equal(t, len(ymlMap), 2)

	url, ok := ymlMap["test1"]
	assert.True(t, ok)
	assert.Equal(t, url, "test1_url")
}

func TestParseYAML(t *testing.T) {
	yaml, err := parseYAML([]byte(yml))

	assert.Nil(t, err)
	assert.NotNil(t, yaml)
	assert.Equal(t, len(yaml), 2)

	assert.Equal(t, yaml[0].Path, "/urlshort")
	assert.Equal(t, yaml[1].URL, "https://github.com/gophercises/urlshort/tree/solution")
}

func TestYAMLHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/urlshort", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	yamlHandler, err := YAMLHandler([]byte(yml), mux)
	assert.Nil(t, err)

	handler := http.HandlerFunc(yamlHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 302)
	assert.Contains(t, rr.Body.String(), "a href")
	assert.Contains(t, rr.Body.String(), "https://github.com/gophercises/urlshort")
}
