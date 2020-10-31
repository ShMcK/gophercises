package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

func main() {
	file := flag.String("file", "gopher.json", "the Choose Your Own Adventure file")
	flag.Parse()
	fmt.Printf("Use file %s\n", *file)

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	d := json.NewDecoder(f)
	var story Story
	if err := d.Decode(&story); err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", story)
}
