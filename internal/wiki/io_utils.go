package main

import (
	"fmt"
	"os"
	"strings"
)

func titleToInternalTitle(title string) string {
	return strings.ReplaceAll(title, " ", "_")
}

func internalTitleToTitle(filename string) string {
	return strings.ReplaceAll(filename, "_", " ")
}

// Saves a page struct to a file
// The file is named after the page's title, and the content is the just a dump of the page's body
func (p *Page) save() error {
	if _, err := os.Stat(saved_pages_path); os.IsNotExist(err) {
		os.Mkdir(saved_pages_path, 0755)
	}
	fullpath := saved_pages_path + "/" + p.InternalTitle
	return os.WriteFile(fullpath, p.Body, 0600)
}

// Loads a page from a file.
// The file is named after the page's title, and the body is the file content as a byte array
func loadPage(internalTitle string) (*Page, error) {
	fmt.Println("attempt to load page with internalTitle: ", internalTitle)
	fullpath := saved_pages_path + "/" + internalTitle
	body, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}
	return &Page{
		Title:         internalTitleToTitle(internalTitle),
		InternalTitle: internalTitle,
		Body:          body,
	}, nil
}
