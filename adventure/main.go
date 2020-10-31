package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Story a group of chapters
type Story map[string]Chapter

// Chapter with text content and options to continue
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option a choice for what to do next
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type handler struct {
	s Story
}

var defaultHandlerTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Choose Your Own Adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
    {{range .Paragraphs}}
        <p>{{.}}</p>
    {{end}}

    <ul>
        {{range .Options}}
            <li>
                <a href="/{{.Arc}}">{{.Text}}</a>
            </li>
        {{end}}
    </ul>
</body>
</html>`

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

func NewHandler(s Story) http.Handler {
	return handler{s}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("").Parse(defaultHandlerTmpl))
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:] // remove pre-ceeding slash "/intro" -> "intro"

	if _, ok := h.s[path]; ok {
		err := tpl.Execute(w, h.s[path])
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)

}

func main() {
	port := flag.Int("port", 3000, "port to run server on")
	file := flag.String("file", "gopher.json", "the Choose Your Own Adventure file")
	flag.Parse()
	fmt.Printf("Use file %s\n", *file)

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	story, err := JsonStory(f)
	if err != nil {
		panic(err)
	}

	h := NewHandler(story)
	fmt.Printf("Server running on port %d\n", *port)
	addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(addr, h))

}
