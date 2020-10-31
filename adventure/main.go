package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
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

func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

func main() {
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

	fmt.Printf("%v\n", story)
}
